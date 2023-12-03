package mongo

import (
	"github.com/miruken-go/miruken"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type (
	// Installer enables configuration support.
	Installer struct {
		aliases map[reflect.Type]string
		clients map[reflect.Type]Config
	}
)

func (i *Installer) Install(setup *miruken.SetupBuilder) error {
	if setup.Tag(&featureTag) {
		setup.Specs(&Factory{})
		if i.aliases != nil {
			setup.Options(Options{Aliases: i.aliases})
		}
	}
	return nil
}

func Client[T ~*mongo.Client](cfg Config) func(*Installer) {
	return func(installer *Installer) {
		if installer.clients == nil {
			installer.clients = make(map[reflect.Type]Config, 1)
		}
		installer.clients[reflect.TypeOf((*T)(nil)).Elem()] = cfg
	}
}

func ClientAlias[T ~*mongo.Client](path string) func(*Installer) {
	if path == "" {
		panic("path is required")
	}
	return func(installer *Installer) {
		if installer.aliases == nil {
			installer.aliases = make(map[reflect.Type]string, 1)
		}
		installer.aliases[reflect.TypeOf((*T)(nil)).Elem()] = path
	}
}

// Feature creates and configures configuration support
// using the supplied configuration Provider.
func Feature(
	config ...func(*Installer),
) miruken.Feature {
	installer := &Installer{}
	for _, configure := range config {
		if configure != nil {
			configure(installer)
		}
	}
	return installer
}

var featureTag byte
