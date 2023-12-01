package token

import (
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/security/password"
)

// Installer enables Azure ADB2C token enrichment.
type Installer struct{}

func (i *Installer) DependsOn() []miruken.Feature {
	return []miruken.Feature{password.Feature()}
}

func (i *Installer) Install(setup *miruken.SetupBuilder) error {
	if setup.Tag(&featureTag) {
		setup.Specs(&EnrichHandler{})
	}
	return nil
}

func Feature(config ...func(*Installer)) miruken.Feature {
	installer := &Installer{}
	for _, configure := range config {
		if configure != nil {
			configure(installer)
		}
	}
	return installer
}

var featureTag byte
