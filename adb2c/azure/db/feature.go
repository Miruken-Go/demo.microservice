package db

import (
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/jmoiron/sqlx"
	"github.com/miruken-go/miruken/setup"
)

// Installer enables azure cosmos support.
type Installer struct {
	aliases map[reflect.Type]string
	clients map[reflect.Type]Config
}


func (i *Installer) Install(b *setup.Builder) error {
	if b.Tag(&featureTag) {
		b.Specs(&Factory{})
		if i.aliases != nil {
			b.Options(Options{Aliases: i.aliases, Clients: i.clients})
		}
	}
	return nil
}

func Client[T ~*azcosmos.Client](cfg Config) func(*Installer) {
	return func(installer *Installer) {
		if installer.clients == nil {
			installer.clients = make(map[reflect.Type]Config, 1)
		}
		installer.clients[reflect.TypeOf((*T)(nil)).Elem()] = cfg
	}
}

func ClientAlias[T ~*azcosmos.Client](path string) func(*Installer) {
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

func SqlClient[T ~*sqlx.DB](cfg Config) func(*Installer) {
	return func(installer *Installer) {
		if installer.clients == nil {
			installer.clients = make(map[reflect.Type]Config, 1)
		}
		installer.clients[reflect.TypeOf((*T)(nil)).Elem()] = cfg
	}
}

func SqlClientAlias[T ~*sqlx.DB](path string) func(*Installer) {
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
