package teamsrv

import (
	"github.com/miruken-go/demo.microservice/teamapi"
	"github.com/miruken-go/demo.microservice/teamsrv/person"
	"github.com/miruken-go/demo.microservice/teamsrv/team"
	"github.com/miruken-go/miruken"
)

//go:generate $GOPATH/bin/miruken -tests

var Feature = miruken.FeatureSet(teamapi.Feature, person.Feature, team.Feature)
