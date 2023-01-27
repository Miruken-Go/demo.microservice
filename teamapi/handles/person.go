package handles

import (
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"github.com/miruken-go/miruken"
	"sync/atomic"
)

type (
	PersonHandler struct {
		nextId int32
	}
)

func (h *PersonHandler) Create(
	_*struct {
		miruken.Handles
	  }, create *commands.CreatePerson,
) (data.Person, error) {
	atomic.AddInt32(&h.nextId, 1)
	person := data.Person{
		Id:        h.nextId,
		FirstName: create.FirstName,
		LastName:  create.LastName,
		BirthDate: create.BirthDate,
	}
	return person, nil
}
