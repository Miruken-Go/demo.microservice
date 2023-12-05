package auth

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/entitlement"
	"github.com/miruken-go/demo.microservice/adb2c/auth/principal"
	"github.com/miruken-go/demo.microservice/adb2c/auth/subject"
	"github.com/miruken-go/demo.microservice/adb2c/auth/tag"
	"github.com/miruken-go/demo.microservice/adb2c/mongo"
	"github.com/miruken-go/miruken"
)

var Feature = miruken.FeatureSet(
	api.Feature,
	subject.Feature,
	principal.Feature,
	entitlement.Feature,
	tag.Feature,
	mongo.Feature())


