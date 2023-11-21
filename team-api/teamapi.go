package teamapi

import (
	"github.com/miruken-go/demo.microservice/team-api/commands"
	"github.com/miruken-go/demo.microservice/team-api/data"
	"github.com/miruken-go/demo.microservice/team-api/queries"
	"github.com/miruken-go/miruken"
)

//go:generate $GOPATH/bin/miruken -tests

var Feature = miruken.FeatureSet(commands.Feature, queries.Feature, data.Feature)
