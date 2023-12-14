package team

import (
	teamapi "github.com/miruken-go/demo.microservice/team-api"
	"github.com/miruken-go/demo.microservice/team/person"
	"github.com/miruken-go/demo.microservice/team/team"
	"github.com/miruken-go/miruken/setup"
)

var Feature = setup.FeatureSet(teamapi.Feature, person.Feature, team.Feature)
