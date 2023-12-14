package auth

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/entitlement"
	"github.com/miruken-go/demo.microservice/adb2c/auth/principal"
	"github.com/miruken-go/demo.microservice/adb2c/auth/subject"
	"github.com/miruken-go/demo.microservice/adb2c/azure"
	"github.com/miruken-go/miruken/setup"
)

var Feature = setup.FeatureSet(
	api.Feature,
	subject.Feature,
	principal.Feature,
	entitlement.Feature,
	azure.Feature())
