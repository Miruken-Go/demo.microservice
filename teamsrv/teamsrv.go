package teamsrv

import (
	"github.com/miruken-go/demo.microservice/teamapi"
	"github.com/miruken-go/demo.microservice/teamsrv/commands"
	"github.com/miruken-go/miruken"
)

//go:generate $GOPATH/bin/miruken -tests

var Feature = miruken.GroupFeatures(teamapi.Feature, commands.Feature)
