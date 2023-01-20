package commands

import "github.com/miruken-go/miruken"

// Factory creates commands from a type id.
type Factory struct{}

func (F *Factory) New(
	_*struct {
		cp miruken.Creates `key:"commands.CreatePerson"`
		up miruken.Creates `key:"commands.UpdatePerson"`
		ct miruken.Creates `key:"commands.CreateTeam"`
		ut miruken.Creates `key:"commands.UpdateTeam"`
	  }, create *miruken.Creates,
) any {
	switch create.Key() {
	case "commands.CreatePerson":
		return new(CreatePerson)
	case "commands.UpdatePerson":
		return new(UpdatePerson)
	case "commands.CreateTeam":
		return new(CreateTeam)
	case "commands.UpdateTeam":
		return new(UpdateTeam)
	}
	return nil
}
