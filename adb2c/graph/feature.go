package graph

import (
	"github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/demo.microservice/adb2c/graph/user"
	"github.com/miruken-go/miruken/setup"
	play "github.com/miruken-go/miruken/validates/play"
)

var Feature = setup.FeatureSet(
	api.Feature, user.Feature,
	play.Feature())
