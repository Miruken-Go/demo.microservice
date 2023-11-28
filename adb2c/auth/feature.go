package auth

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/handle"
	"github.com/miruken-go/miruken"
)

var Feature = miruken.FeatureSet(api.Feature, handle.Feature)


