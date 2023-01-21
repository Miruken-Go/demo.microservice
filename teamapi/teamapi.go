package teamapi

import (
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/handles"
	"github.com/miruken-go/demo.microservice/teamapi/queries"
	"github.com/miruken-go/miruken"
)

//go:generate $GOPATH/bin/miruken -tests

var Feature = miruken.GroupFeatures(
	commands.Feature,
	queries.Feature,
	handles.Feature)
