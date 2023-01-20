package queries

import "github.com/miruken-go/miruken"

// Factory creates queries from a type id.
type Factory struct{}

func (F *Factory) New(
	_*struct {
		fp miruken.Creates `key:"queries.FindPeople"`
		ft miruken.Creates `key:"queries.FindTeams"`
	  }, create *miruken.Creates,
) any {
	switch create.Key() {
	case "queries.FindPeople":
		return new(FindPeople)
	case "queries.FindTeams":
		return new(FindTeams)
	}
	return nil
}
