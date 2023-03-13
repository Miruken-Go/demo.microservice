package teamapi

import (
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"github.com/miruken-go/demo.microservice/teamapi/queries"
	"github.com/miruken-go/miruken"
)

//go:generate $GOPATH/bin/miruken -tests

var Feature = miruken.FeatureSet(commands.Feature, queries.Feature, data.Feature)
