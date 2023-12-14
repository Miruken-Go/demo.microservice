package token

import (
	"github.com/miruken-go/miruken/security/password"
	"github.com/miruken-go/miruken/setup"
)

// Installer enables Azure ADB2C token enrichment.
type Installer struct{}

func (i *Installer) DependsOn() []setup.Feature {
	return []setup.Feature{password.Feature()}
}

func (i *Installer) Install(b *setup.Builder) error {
	if b.Tag(&featureTag) {
		b.Specs(&EnrichHandler{})
	}
	return nil
}

func Feature(config ...func(*Installer)) setup.Feature {
	installer := &Installer{}
	for _, configure := range config {
		if configure != nil {
			configure(installer)
		}
	}
	return installer
}

var featureTag byte
