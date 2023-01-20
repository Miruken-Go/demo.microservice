package handles

import (
	"github.com/miruken-go/demo-microservice/client/api/commands"
	"github.com/miruken-go/demo-microservice/client/api/data"
	"github.com/miruken-go/miruken"
	"sync/atomic"
)

type PersonHandler struct {
	nextId int32
}

func (h *PersonHandler) Create(
	_*struct{
		miruken.Handles
	  }, create *commands.CreatePerson,
) (data.Person, error) {
	atomic.AddInt32(&h.nextId, 1)
	person := create.Person
	person.Id = h.nextId
	return person, nil
}