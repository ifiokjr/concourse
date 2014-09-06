// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/concourse/atc/builds"
	"github.com/concourse/atc/callbacks/handler"
)

type FakeBuildDB struct {
	GetBuildStub        func(job string, build int) (builds.Build, error)
	getBuildMutex       sync.RWMutex
	getBuildArgsForCall []struct {
		job   string
		build int
	}
	getBuildReturns struct {
		result1 builds.Build
		result2 error
	}
	SaveBuildStatusStub        func(job string, build int, status builds.Status) error
	saveBuildStatusMutex       sync.RWMutex
	saveBuildStatusArgsForCall []struct {
		job    string
		build  int
		status builds.Status
	}
	saveBuildStatusReturns struct {
		result1 error
	}
	SaveBuildInputStub        func(job string, build int, input builds.VersionedResource) error
	saveBuildInputMutex       sync.RWMutex
	saveBuildInputArgsForCall []struct {
		job   string
		build int
		input builds.VersionedResource
	}
	saveBuildInputReturns struct {
		result1 error
	}
	SaveBuildOutputStub        func(job string, build int, output builds.VersionedResource) error
	saveBuildOutputMutex       sync.RWMutex
	saveBuildOutputArgsForCall []struct {
		job    string
		build  int
		output builds.VersionedResource
	}
	saveBuildOutputReturns struct {
		result1 error
	}
}

func (fake *FakeBuildDB) GetBuild(job string, build int) (builds.Build, error) {
	fake.getBuildMutex.Lock()
	fake.getBuildArgsForCall = append(fake.getBuildArgsForCall, struct {
		job   string
		build int
	}{job, build})
	fake.getBuildMutex.Unlock()
	if fake.GetBuildStub != nil {
		return fake.GetBuildStub(job, build)
	} else {
		return fake.getBuildReturns.result1, fake.getBuildReturns.result2
	}
}

func (fake *FakeBuildDB) GetBuildCallCount() int {
	fake.getBuildMutex.RLock()
	defer fake.getBuildMutex.RUnlock()
	return len(fake.getBuildArgsForCall)
}

func (fake *FakeBuildDB) GetBuildArgsForCall(i int) (string, int) {
	fake.getBuildMutex.RLock()
	defer fake.getBuildMutex.RUnlock()
	return fake.getBuildArgsForCall[i].job, fake.getBuildArgsForCall[i].build
}

func (fake *FakeBuildDB) GetBuildReturns(result1 builds.Build, result2 error) {
	fake.GetBuildStub = nil
	fake.getBuildReturns = struct {
		result1 builds.Build
		result2 error
	}{result1, result2}
}

func (fake *FakeBuildDB) SaveBuildStatus(job string, build int, status builds.Status) error {
	fake.saveBuildStatusMutex.Lock()
	fake.saveBuildStatusArgsForCall = append(fake.saveBuildStatusArgsForCall, struct {
		job    string
		build  int
		status builds.Status
	}{job, build, status})
	fake.saveBuildStatusMutex.Unlock()
	if fake.SaveBuildStatusStub != nil {
		return fake.SaveBuildStatusStub(job, build, status)
	} else {
		return fake.saveBuildStatusReturns.result1
	}
}

func (fake *FakeBuildDB) SaveBuildStatusCallCount() int {
	fake.saveBuildStatusMutex.RLock()
	defer fake.saveBuildStatusMutex.RUnlock()
	return len(fake.saveBuildStatusArgsForCall)
}

func (fake *FakeBuildDB) SaveBuildStatusArgsForCall(i int) (string, int, builds.Status) {
	fake.saveBuildStatusMutex.RLock()
	defer fake.saveBuildStatusMutex.RUnlock()
	return fake.saveBuildStatusArgsForCall[i].job, fake.saveBuildStatusArgsForCall[i].build, fake.saveBuildStatusArgsForCall[i].status
}

func (fake *FakeBuildDB) SaveBuildStatusReturns(result1 error) {
	fake.SaveBuildStatusStub = nil
	fake.saveBuildStatusReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildDB) SaveBuildInput(job string, build int, input builds.VersionedResource) error {
	fake.saveBuildInputMutex.Lock()
	fake.saveBuildInputArgsForCall = append(fake.saveBuildInputArgsForCall, struct {
		job   string
		build int
		input builds.VersionedResource
	}{job, build, input})
	fake.saveBuildInputMutex.Unlock()
	if fake.SaveBuildInputStub != nil {
		return fake.SaveBuildInputStub(job, build, input)
	} else {
		return fake.saveBuildInputReturns.result1
	}
}

func (fake *FakeBuildDB) SaveBuildInputCallCount() int {
	fake.saveBuildInputMutex.RLock()
	defer fake.saveBuildInputMutex.RUnlock()
	return len(fake.saveBuildInputArgsForCall)
}

func (fake *FakeBuildDB) SaveBuildInputArgsForCall(i int) (string, int, builds.VersionedResource) {
	fake.saveBuildInputMutex.RLock()
	defer fake.saveBuildInputMutex.RUnlock()
	return fake.saveBuildInputArgsForCall[i].job, fake.saveBuildInputArgsForCall[i].build, fake.saveBuildInputArgsForCall[i].input
}

func (fake *FakeBuildDB) SaveBuildInputReturns(result1 error) {
	fake.SaveBuildInputStub = nil
	fake.saveBuildInputReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildDB) SaveBuildOutput(job string, build int, output builds.VersionedResource) error {
	fake.saveBuildOutputMutex.Lock()
	fake.saveBuildOutputArgsForCall = append(fake.saveBuildOutputArgsForCall, struct {
		job    string
		build  int
		output builds.VersionedResource
	}{job, build, output})
	fake.saveBuildOutputMutex.Unlock()
	if fake.SaveBuildOutputStub != nil {
		return fake.SaveBuildOutputStub(job, build, output)
	} else {
		return fake.saveBuildOutputReturns.result1
	}
}

func (fake *FakeBuildDB) SaveBuildOutputCallCount() int {
	fake.saveBuildOutputMutex.RLock()
	defer fake.saveBuildOutputMutex.RUnlock()
	return len(fake.saveBuildOutputArgsForCall)
}

func (fake *FakeBuildDB) SaveBuildOutputArgsForCall(i int) (string, int, builds.VersionedResource) {
	fake.saveBuildOutputMutex.RLock()
	defer fake.saveBuildOutputMutex.RUnlock()
	return fake.saveBuildOutputArgsForCall[i].job, fake.saveBuildOutputArgsForCall[i].build, fake.saveBuildOutputArgsForCall[i].output
}

func (fake *FakeBuildDB) SaveBuildOutputReturns(result1 error) {
	fake.SaveBuildOutputStub = nil
	fake.saveBuildOutputReturns = struct {
		result1 error
	}{result1}
}

var _ handler.BuildDB = new(FakeBuildDB)
