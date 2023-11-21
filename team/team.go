package team

import (
	"github.com/miruken-go/demo.microservice/team-api"
	"github.com/miruken-go/demo.microservice/team/person"
	"github.com/miruken-go/demo.microservice/team/team"
	"github.com/miruken-go/miruken"
)

var Feature = miruken.FeatureSet(teamapi.Feature, person.Feature, team.Feature)
