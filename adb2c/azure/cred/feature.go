package cred

import (
	"github.com/miruken-go/miruken/setup"
)

// Installer enables azure credentials support.
type Installer struct {}

func (i *Installer) Install(b *setup.Builder) error {
	if b.Tag(&featureTag) {
		b.Specs(&Factory{})
	}
	return nil
}

// Feature creates and configures configuration support
// using the supplied configuration Provider.
func Feature(
	config ...func(*Installer),
) setup.Feature {
	installer := &Installer{}
	for _, configure := range config {
		if configure != nil {
			configure(installer)
		}
	}
	return installer
}

var featureTag byte
