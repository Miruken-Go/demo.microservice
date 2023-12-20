package auth

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/principal"
	"github.com/miruken-go/demo.microservice/adb2c/auth/subject"
	"github.com/miruken-go/demo.microservice/adb2c/azure"
	"github.com/miruken-go/miruken/setup"
	play "github.com/miruken-go/miruken/validates/play"
)

var Feature = setup.FeatureSet(
	api.Feature, subject.Feature, principal.Feature,
	play.Feature(), azure.Feature())
