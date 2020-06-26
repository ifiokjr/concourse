module Data exposing
    ( check
    , dashboardPipeline
    , elementPosition
    , httpInternalServerError
    , httpNotFound
    , httpNotImplemented
    , httpUnauthorized
    , input
    , job
    , jobBuild
    , jobBuildId
    , jobId
    , jobName
    , leftClickEvent
    , longJobBuildId
    , pipeline
    , pipelineId
    , pipelineName
    , resource
    , resourceId
    , resourceName
    , resourceVersionId
    , shortJobId
    , shortPipelineId
    , shortResourceId
    , teamName
    , version
    , versionedResource
    , withArchived
    , withBuildName
    , withDisableManualTrigger
    , withFinishedBuild
    , withGroups
    , withInputs
    , withJobName
    , withName
    , withNextBuild
    , withPaused
    , withPipelineName
    , withPublic
    , withResourceName
    , withShortJobId
    , withShortPipelineId
    , withShortResourceId
    , withTeamName
    , withTransitionBuild
    )

import Browser.Dom
import Concourse
import Concourse.BuildStatus as BuildStatus
import Dashboard.Group.Models
import Dict exposing (Dict)
import Http
import Json.Encode
import Test.Html.Event as Event
import Time


httpUnauthorized : Result Http.Error a
httpUnauthorized =
    Err <|
        Http.BadStatus
            { url = "http://example.com"
            , status =
                { code = 401
                , message = ""
                }
            , headers = Dict.empty
            , body = ""
            }


httpNotFound : Result Http.Error a
httpNotFound =
    Err <|
        Http.BadStatus
            { url = "http://example.com"
            , status =
                { code = 404
                , message = "not found"
                }
            , headers = Dict.empty
            , body = ""
            }


httpNotImplemented : Result Http.Error a
httpNotImplemented =
    Err <|
        Http.BadStatus
            { url = "http://example.com"
            , status =
                { code = 501
                , message = "not implemented"
                }
            , headers = Dict.empty
            , body = ""
            }


httpInternalServerError : Result Http.Error a
httpInternalServerError =
    Err <|
        Http.BadStatus
            { url = "http://example.com"
            , status =
                { code = 500
                , message = "internal server error"
                }
            , headers = Dict.empty
            , body = ""
            }


check : Concourse.CheckStatus -> Concourse.Check
check status =
    case status of
        Concourse.Started ->
            { id = 0
            , status = Concourse.Started
            , createTime = Just <| Time.millisToPosix 0
            , startTime = Just <| Time.millisToPosix 0
            , endTime = Nothing
            , checkError = Nothing
            }

        Concourse.Succeeded ->
            { id = 0
            , status = Concourse.Succeeded
            , createTime = Just <| Time.millisToPosix 0
            , startTime = Just <| Time.millisToPosix 0
            , endTime = Just <| Time.millisToPosix 1000
            , checkError = Nothing
            }

        Concourse.Errored ->
            { id = 0
            , status = Concourse.Errored
            , createTime = Just <| Time.millisToPosix 0
            , startTime = Just <| Time.millisToPosix 0
            , endTime = Just <| Time.millisToPosix 1000
            , checkError = Just "something broke"
            }


resource : String -> Concourse.Resource
resource pinnedVersion =
    { teamName = teamName
    , pipelineName = pipelineName
    , name = resourceName
    , failingToCheck = False
    , checkError = ""
    , checkSetupError = ""
    , lastChecked = Nothing
    , pinnedVersion = Just <| version pinnedVersion
    , pinnedInConfig = False
    , pinComment = Nothing
    , icon = Nothing
    }


pipeline : String -> Int -> Concourse.Pipeline
pipeline team id =
    { id = id
    , name = "pipeline-" ++ String.fromInt id
    , paused = False
    , archived = False
    , public = True
    , teamName = team
    , groups = []
    }


dashboardPipeline : Int -> Bool -> Dashboard.Group.Models.Pipeline
dashboardPipeline id public =
    { id = id
    , name = pipelineName
    , teamName = teamName
    , public = public
    , isToggleLoading = False
    , isVisibilityLoading = False
    , paused = False
    , archived = False
    , stale = False
    , jobsDisabled = False
    }


withPaused : Bool -> { r | paused : Bool } -> { r | paused : Bool }
withPaused paused p =
    { p | paused = paused }


withArchived : Bool -> { r | archived : Bool } -> { r | archived : Bool }
withArchived archived p =
    { p | archived = archived }


withPublic : Bool -> { r | public : Bool } -> { r | public : Bool }
withPublic public p =
    { p | public = public }


withName : String -> { r | name : String } -> { r | name : String }
withName name p =
    { p | name = name }


withGroups : List Concourse.PipelineGroup -> { r | groups : List Concourse.PipelineGroup } -> { r | groups : List Concourse.PipelineGroup }
withGroups groups p =
    { p | groups = groups }


job : Int -> Int -> Concourse.Job
job jobID pipelineID =
    { id = jobID
    , name = jobName
    , pipelineName = "pipeline-" ++ String.fromInt pipelineID
    , teamName = teamName
    , nextBuild = Nothing
    , finishedBuild = Nothing
    , transitionBuild = Nothing
    , paused = False
    , disableManualTrigger = False
    , inputs = []
    , outputs = []
    , groups = []
    }


withDisableManualTrigger : Bool -> { r | disableManualTrigger : Bool } -> { r | disableManualTrigger : Bool }
withDisableManualTrigger disableManualTrigger p =
    { p | disableManualTrigger = disableManualTrigger }


withTeamName : String -> { r | teamName : String } -> { r | teamName : String }
withTeamName name p =
    { p | teamName = name }


withPipelineName : String -> { r | pipelineName : String } -> { r | pipelineName : String }
withPipelineName name p =
    { p | pipelineName = name }


withJobName : String -> { r | jobName : String } -> { r | jobName : String }
withJobName name p =
    { p | jobName = name }


withResourceName : String -> { r | resourceName : String } -> { r | resourceName : String }
withResourceName name p =
    { p | resourceName = name }


withBuildName : String -> { r | buildName : String } -> { r | buildName : String }
withBuildName name p =
    { p | buildName = name }


withNextBuild :
    Maybe Concourse.Build
    -> { r | nextBuild : Maybe Concourse.Build, name : String, pipelineName : String, teamName : String }
    -> { r | nextBuild : Maybe Concourse.Build, name : String, pipelineName : String, teamName : String }
withNextBuild build j =
    { j | nextBuild = build |> Maybe.map (updateJobIdentifier j) }


withFinishedBuild :
    Maybe Concourse.Build
    -> { r | finishedBuild : Maybe Concourse.Build, name : String, pipelineName : String, teamName : String }
    -> { r | finishedBuild : Maybe Concourse.Build, name : String, pipelineName : String, teamName : String }
withFinishedBuild build j =
    { j | finishedBuild = build |> Maybe.map (updateJobIdentifier j) }


withTransitionBuild :
    Maybe Concourse.Build
    -> { r | transitionBuild : Maybe Concourse.Build, name : String, pipelineName : String, teamName : String }
    -> { r | transitionBuild : Maybe Concourse.Build, name : String, pipelineName : String, teamName : String }
withTransitionBuild build j =
    { j | transitionBuild = build |> Maybe.map (updateJobIdentifier j) }


updateJobIdentifier :
    { r | name : String, pipelineName : String, teamName : String }
    -> Concourse.Build
    -> Concourse.Build
updateJobIdentifier j b =
    { b
        | job =
            Just
                { jobName = j.name
                , pipelineName = j.pipelineName
                , teamName = j.teamName
                }
    }


withInputs : List Concourse.JobInput -> { r | inputs : List Concourse.JobInput } -> { r | inputs : List Concourse.JobInput }
withInputs inputs j =
    { j | inputs = inputs }


input : List String -> Concourse.JobInput
input passed =
    { name = "input"
    , resource = "res0"
    , passed = passed
    , trigger = True
    }


jobName =
    "job"


teamName =
    "team"


pipelineName =
    "pipeline"


resourceName =
    "resource"


buildName =
    "1"


withShortPipelineId =
    withPipelineName "p"
        >> withTeamName "t"


withShortJobId =
    withShortPipelineId >> withJobName "j"


withShortResourceId =
    withShortPipelineId >> withResourceName "r"


versionedResource : String -> Int -> Concourse.VersionedResource
versionedResource v id =
    { id = id
    , version = version v
    , metadata = []
    , enabled = True
    }


version : String -> Dict String String
version v =
    Dict.fromList [ ( "version", v ) ]


pipelineId : Concourse.PipelineIdentifier
pipelineId =
    { teamName = teamName
    , pipelineName = pipelineName
    }


shortPipelineId : Concourse.PipelineIdentifier
shortPipelineId =
    pipelineId |> withShortPipelineId


jobId : Concourse.JobIdentifier
jobId =
    { teamName = teamName
    , pipelineName = pipelineName
    , jobName = jobName
    }


shortJobId : Concourse.JobIdentifier
shortJobId =
    jobId |> withShortJobId


resourceId : Concourse.ResourceIdentifier
resourceId =
    { teamName = teamName
    , pipelineName = pipelineName
    , resourceName = resourceName
    }


shortResourceId : Concourse.ResourceIdentifier
shortResourceId =
    resourceId |> withShortResourceId


resourceVersionId : Int -> Concourse.VersionedResourceIdentifier
resourceVersionId v =
    { teamName = teamName
    , pipelineName = pipelineName
    , resourceName = resourceName
    , versionID = v
    }



-- jobBuildId is really shortJobBuildId, but since jobBuild returns a short jobId,
-- it would be weird for jobBuildId to not represent jobBuild


jobBuildId : Concourse.JobBuildIdentifier
jobBuildId =
    longJobBuildId |> withShortJobId


longJobBuildId : Concourse.JobBuildIdentifier
longJobBuildId =
    { teamName = teamName
    , pipelineName = pipelineName
    , jobName = jobName
    , buildName = buildName
    }


jobBuild : BuildStatus.BuildStatus -> Concourse.Build
jobBuild status =
    { id = 1
    , name = buildName
    , job = Just (jobId |> withShortJobId)
    , status = status
    , duration =
        { startedAt =
            case status of
                BuildStatus.BuildStatusPending ->
                    Nothing

                _ ->
                    Just <| Time.millisToPosix 0
        , finishedAt =
            if BuildStatus.isRunning status then
                Nothing

            else
                Just <| Time.millisToPosix 0
        }
    , reapTime = Nothing
    }


elementPosition : Browser.Dom.Element
elementPosition =
    { scene =
        { width = 0
        , height = 0
        }
    , viewport =
        { width = 0
        , height = 0
        , x = 0
        , y = 0
        }
    , element =
        { x = 0
        , y = 0
        , width = 1
        , height = 1
        }
    }


leftClickEvent : ( String, Json.Encode.Value )
leftClickEvent =
    Event.custom "click" <|
        Json.Encode.object
            [ ( "ctrlKey", Json.Encode.bool False )
            , ( "altKey", Json.Encode.bool False )
            , ( "metaKey", Json.Encode.bool False )
            , ( "shiftKey", Json.Encode.bool False )
            , ( "button", Json.Encode.int 0 )
            ]
