// Code generated by https://github.com/Miruken-Go/miruken/tools/cmd/miruken; DO NOT EDIT.

package person

import "github.com/miruken-go/miruken/setup"

var Feature setup.Feature = setup.FeatureFunc(func(setup *setup.Builder) error {
	setup.Specs(
		&Handler{},
	)
	return nil
})
