// Code generated by https://github.com/Miruken-Go/miruken/tools/cmd/miruken; DO NOT EDIT.

package handles

import "github.com/miruken-go/miruken"

var Feature miruken.Feature = miruken.FeatureFunc(func(setup *miruken.SetupBuilder) error {
	setup.Specs(
		&CreatePersonIntegrity{},
		&PersonHandler{},
	)
	return nil
})
