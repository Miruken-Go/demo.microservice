package data

import (
	"github.com/miruken-go/miruken/creates"
)

//go:generate $GOPATH/bin/miruken -tests

// Factory creates queries from a type id.
type Factory struct{}

func (f *Factory) New(
	_ *struct {
		p creates.It `key:"data.Person"`
	}, create *creates.It,
) any {
	switch create.Key() {
	case "data.Person":
		return new(Person)
	}
	return nil
}
