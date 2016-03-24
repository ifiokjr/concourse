package radar

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
	"github.com/tedsuo/ifrit"

	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
)

type ResourceNotConfiguredError struct {
	ResourceName string
}

func (err ResourceNotConfiguredError) Error() string {
	return fmt.Sprintf("resource '%s' was not found in config", err.ResourceName)
}

//go:generate counterfeiter . RadarDB

type RadarDB interface {
	GetPipelineName() string
	ScopedName(string) string

	IsPaused() (bool, error)

	GetConfig() (atc.Config, db.ConfigVersion, bool, error)

	GetLatestVersionedResource(resource db.SavedResource) (db.SavedVersionedResource, bool, error)
	GetResource(resourceName string) (db.SavedResource, error)
	GetResourceType(resourceTypeName string) (db.SavedResourceType, bool, error)
	PauseResource(resourceName string) error
	UnpauseResource(resourceName string) error

	SaveResourceVersions(atc.ResourceConfig, []atc.Version) error
	SaveResourceTypeVersion(atc.ResourceType, atc.Version) error
	SetResourceCheckError(resource db.SavedResource, err error) error
	LeaseResourceChecking(logger lager.Logger, resource string, interval time.Duration, immediate bool) (db.Lease, bool, error)
}

type Radar struct {
	logger          lager.Logger
	tracker         resource.Tracker
	defaultInterval time.Duration
	db              RadarDB
	clock           clock.Clock
	externalURL     string
}

func NewRadar(
	tracker resource.Tracker,
	defaultInterval time.Duration,
	db RadarDB,
	clock clock.Clock,
	externalURL string,
) *Radar {
	return &Radar{
		tracker:         tracker,
		defaultInterval: defaultInterval,
		db:              db,
		clock:           clock,
		externalURL:     externalURL,
	}
}

func (radar *Radar) Scanner(logger lager.Logger, resourceName string) ifrit.Runner {
	return ifrit.RunFunc(func(signals <-chan os.Signal, ready chan<- struct{}) error {
		// do an immediate initial check
		var interval time.Duration = 0

		close(ready)

		for {
			timer := radar.clock.NewTimer(interval)

			var resourceConfig atc.ResourceConfig
			var resourceTypes atc.ResourceTypes

			select {
			case <-signals:
				timer.Stop()
				return nil

			case <-timer.C():
				var err error

				resourceConfig, resourceTypes, err = radar.getResourceConfig(logger, resourceName)
				if err != nil {
					return err
				}

				savedResource, err := radar.db.GetResource(resourceConfig.Name)
				if err != nil {
					return err
				}

				interval, err = radar.checkInterval(resourceConfig)
				if err != nil {
					setErr := radar.db.SetResourceCheckError(savedResource, err)
					if setErr != nil {
						logger.Error("failed-to-set-check-error", err)
					}

					return err
				}

				leaseLogger := logger.Session("lease", lager.Data{
					"resource": resourceName,
				})

				lease, leased, err := radar.db.LeaseResourceChecking(logger, resourceName, interval, false)

				if err != nil {
					leaseLogger.Error("failed-to-get-lease", err, lager.Data{
						"resource": resourceName,
					})
					break
				}

				if !leased {
					leaseLogger.Debug("did-not-get-lease")
					break
				}

				err = radar.scan(logger.Session("tick"), resourceConfig, resourceTypes, savedResource)

				lease.Break()

				if err != nil {
					return err
				}
			}
		}
	})
}

func (radar *Radar) ResourceTypeScanner(logger lager.Logger, resourceTypeName string) ifrit.Runner {
	return ifrit.RunFunc(func(signals <-chan os.Signal, ready chan<- struct{}) error {
		// do an immediate initial check
		interval := radar.defaultInterval

		close(ready)

		for {
			timer := radar.clock.NewTimer(interval)

			var resourceType atc.ResourceType

			select {
			case <-signals:
				timer.Stop()
				return nil

			case <-timer.C():
				var err error

				resourceType, err = radar.getResourceTypeConfig(logger, resourceTypeName)
				if err != nil {
					return err
				}

				err = radar.resourceTypeScan(logger.Session("tick"), resourceType)
				if err != nil {
					return err
				}
			}
		}
	})
}

func (radar *Radar) Scan(logger lager.Logger, resourceName string) error {
	leaseLogger := logger.Session("lease", lager.Data{
		"resource": resourceName,
	})

	resourceConfig, resourceTypes, err := radar.getResourceConfig(logger, resourceName)
	if err != nil {
		return err
	}

	savedResource, err := radar.db.GetResource(resourceConfig.Name)
	if err != nil {
		return err
	}

	interval, err := radar.checkInterval(resourceConfig)
	if err != nil {
		setErr := radar.db.SetResourceCheckError(savedResource, err)
		if setErr != nil {
			logger.Error("failed-to-set-check-error", err)
		}

		return err
	}

	for {
		lease, leased, err := radar.db.LeaseResourceChecking(logger, resourceName, interval, true)
		if err != nil {
			leaseLogger.Error("failed-to-get-lease", err, lager.Data{
				"resource": resourceName,
			})

			return err
		}

		if !leased {
			leaseLogger.Debug("did-not-get-lease")
			radar.clock.Sleep(time.Second)
			continue
		}

		defer lease.Break()

		break
	}

	return radar.scan(logger, resourceConfig, resourceTypes, savedResource)
}

func (radar *Radar) scan(logger lager.Logger, resourceConfig atc.ResourceConfig, resourceTypes atc.ResourceTypes, savedResource db.SavedResource) error {
	pipelinePaused, err := radar.db.IsPaused()
	if err != nil {
		logger.Error("failed-to-check-if-pipeline-paused", err)
		return err
	}

	if pipelinePaused {
		logger.Debug("pipeline-paused")
		return nil
	}

	if savedResource.Paused {
		logger.Debug("resource-paused")
		return nil
	}

	vr, found, err := radar.db.GetLatestVersionedResource(savedResource)
	if err != nil {
		logger.Error("failed-to-get-current-version", err)
		return err
	}

	var from db.Version
	if found {
		from = vr.Version
	}

	pipelineName := radar.db.GetPipelineName()

	var resourceTypeVersion atc.Version
	_, found = resourceTypes.Lookup(resourceConfig.Type)
	if found {
		savedResourceType, resourceTypeFound, err := radar.db.GetResourceType(resourceConfig.Type)
		if err != nil {
			logger.Error("failed-to-find-resource-type", err)
			return err
		}
		if resourceTypeFound {
			resourceTypeVersion = atc.Version(savedResourceType.Version)
		}
	}

	session := resource.Session{
		ID: worker.Identifier{
			ResourceTypeVersion: resourceTypeVersion,
			ResourceID:          savedResource.ID,
			Stage:               db.ContainerStageRun,
			CheckType:           resourceConfig.Type,
			CheckSource:         resourceConfig.Source,
		},
		Metadata: worker.Metadata{
			Type:         db.ContainerTypeCheck,
			PipelineName: pipelineName,
		},
		Ephemeral: true,
	}

	res, err := radar.tracker.Init(
		logger,
		resource.TrackerMetadata{
			ResourceName: resourceConfig.Name,
			PipelineName: pipelineName,
			ExternalURL:  radar.externalURL,
		},
		session,
		resource.ResourceType(resourceConfig.Type),
		[]string{},
		resourceTypes,
		worker.NoopImageFetchingDelegate{},
	)
	if err != nil {
		logger.Error("failed-to-initialize-new-resource", err)
		return err
	}

	defer res.Release(nil)

	logger.Debug("checking", lager.Data{
		"from": from,
	})

	newVersions, err := res.Check(resourceConfig.Source, atc.Version(from))

	setErr := radar.db.SetResourceCheckError(savedResource, err)
	if setErr != nil {
		logger.Error("failed-to-set-check-error", err)
	}

	if err != nil {
		if rErr, ok := err.(resource.ErrResourceScriptFailed); ok {
			logger.Info("check-failed", lager.Data{"exit-status": rErr.ExitStatus})
			return nil
		}

		logger.Error("failed-to-check", err)
		return err
	}

	if len(newVersions) == 0 {
		logger.Debug("no-new-versions")
		return nil
	}

	logger.Info("versions-found", lager.Data{
		"versions": newVersions,
		"total":    len(newVersions),
	})

	err = radar.db.SaveResourceVersions(resourceConfig, newVersions)
	if err != nil {
		logger.Error("failed-to-save-versions", err, lager.Data{
			"versions": newVersions,
		})
	}

	return nil
}

func (radar *Radar) resourceTypeScan(logger lager.Logger, resourceType atc.ResourceType) error {
	pipelineName := radar.db.GetPipelineName()

	session := resource.Session{
		ID: worker.Identifier{
			Stage:               db.ContainerStageCheck,
			CheckType:           resourceType.Type,
			CheckSource:         resourceType.Source,
			ImageResourceType:   resourceType.Type,
			ImageResourceSource: resourceType.Source,
		},
		Metadata: worker.Metadata{
			Type:                 db.ContainerTypeCheck,
			PipelineName:         pipelineName,
			WorkingDirectory:     "",
			EnvironmentVariables: nil,
		},
		Ephemeral: true,
	}

	res, err := radar.tracker.Init(
		logger.Session("check-image"),
		resource.EmptyMetadata{},
		session,
		resource.ResourceType(resourceType.Type),
		[]string{},
		atc.ResourceTypes{},
		worker.NoopImageFetchingDelegate{},
	)
	if err != nil {
		return err
	}

	defer res.Release(nil)

	logger.Debug("checking")

	newVersions, err := res.Check(resourceType.Source, atc.Version{})
	if err != nil {
		if rErr, ok := err.(resource.ErrResourceScriptFailed); ok {
			logger.Info("check-failed", lager.Data{"exit-status": rErr.ExitStatus})
			return nil
		}

		logger.Error("failed-to-check", err)
		return err
	}

	if len(newVersions) == 0 {
		logger.Debug("no-new-versions")
		return nil
	}

	logger.Info("versions-found", lager.Data{
		"versions": newVersions,
		"total":    len(newVersions),
	})

	version := newVersions[len(newVersions)-1]
	err = radar.db.SaveResourceTypeVersion(resourceType, version)
	if err != nil {
		logger.Error("failed-to-save-resource-type-version", err, lager.Data{
			"version": version,
		})
		return err
	}

	return nil
}

func (radar *Radar) checkInterval(resourceConfig atc.ResourceConfig) (time.Duration, error) {
	interval := radar.defaultInterval
	if resourceConfig.CheckEvery != "" {
		configuredInterval, err := time.ParseDuration(resourceConfig.CheckEvery)
		if err != nil {
			return 0, err
		}

		interval = configuredInterval
	}

	return interval, nil
}

var errPipelineRemoved = errors.New("pipeline removed")

func (radar *Radar) getResourceConfig(logger lager.Logger, resourceName string) (atc.ResourceConfig, atc.ResourceTypes, error) {
	config, _, found, err := radar.db.GetConfig()
	if err != nil {
		logger.Error("failed-to-get-config", err)
		return atc.ResourceConfig{}, nil, err
	}

	if !found {
		logger.Info("pipeline-removed")
		return atc.ResourceConfig{}, nil, errPipelineRemoved
	}

	resourceConfig, found := config.Resources.Lookup(resourceName)
	if !found {
		logger.Info("resource-removed-from-configuration")
		return resourceConfig, nil, ResourceNotConfiguredError{ResourceName: resourceName}
	}

	return resourceConfig, config.ResourceTypes, nil
}

func (radar *Radar) getResourceTypeConfig(logger lager.Logger, resourceTypeName string) (atc.ResourceType, error) {
	config, _, found, err := radar.db.GetConfig()
	if err != nil {
		logger.Error("failed-to-get-config", err)
		return atc.ResourceType{}, err
	}

	if !found {
		logger.Info("pipeline-removed")
		return atc.ResourceType{}, errPipelineRemoved
	}

	resourceType, found := config.ResourceTypes.Lookup(resourceTypeName)
	if !found {
		logger.Info("resource-type-removed-from-configuration")
		return resourceType, ResourceNotConfiguredError{ResourceName: resourceTypeName}
	}

	return resourceType, nil
}
