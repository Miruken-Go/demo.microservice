package azure

import (
	"github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/demo.microservice/adb2c/azure/db"
	"github.com/miruken-go/demo.microservice/adb2c/azure/principal"
	"github.com/miruken-go/demo.microservice/adb2c/azure/subject"
	"github.com/miruken-go/miruken/setup"
	play "github.com/miruken-go/miruken/validates/play"
)

var Feature = setup.FeatureSet(
	api.Feature, subject.Feature, principal.Feature,
	play.Feature(), db.Feature())
