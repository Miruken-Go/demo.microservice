// Code generated by https://github.com/Miruken-Go/miruken/tools/cmd/miruken; DO NOT EDIT.

package data

import "github.com/miruken-go/miruken"

var Feature miruken.Feature = miruken.FeatureFunc(func(setup *miruken.SetupBuilder) error {
	setup.Specs(
		&Factory{},
	)
	return nil
})
