package commands

import (
	"github.com/miruken-go/miruken/creates"
)

//go:generate $GOPATH/bin/miruken -tests

// Factory creates commands from a type id.
type Factory struct{}

func (f *Factory) New(
	_ *struct {
		cp creates.It `key:"commands.CreatePerson"`
		up creates.It `key:"commands.UpdatePerson"`
		dp creates.It `key:"commands.DeletePeople"`
		ct creates.It `key:"commands.CreateTeam"`
		ut creates.It `key:"commands.UpdateTeam"`
	  }, create *creates.It,
) any {
	switch create.Key() {
	case "commands.CreatePerson":
		return new(CreatePerson)
	case "commands.UpdatePerson":
		return new(UpdatePerson)
	case "commands.DeletePeople":
		return new(DeletePeople)
	case "commands.CreateTeam":
		return new(CreateTeam)
	case "commands.UpdateTeam":
		return new(UpdateTeam)
	}
	return nil
}
