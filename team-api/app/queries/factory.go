package queries

import (
	"github.com/miruken-go/miruken/creates"
)

//go:generate $GOPATH/bin/miruken -tests

// Factory creates queries from a type id.
type Factory struct{}

func (f *Factory) New(
	_ *struct {
		fp creates.It `key:"queries.FindPeople"`
		ft creates.It `key:"queries.FindTeams"`
	  }, create *creates.It,
) any {
	switch create.Key() {
	case "queries.FindPeople":
		return new(FindPeople)
	case "queries.FindTeams":
		return new(FindTeams)
	}
	return nil
}
