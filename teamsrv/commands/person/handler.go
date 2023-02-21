package person

import (
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"github.com/miruken-go/miruken/handles"
	"sync/atomic"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		nextId int32
	}
)

func (h *Handler) Create(
	_ *handles.It, create *commands.CreatePerson,
) (data.Person, error) {
	return data.Person{
		Id:        atomic.AddInt32(&h.nextId, 1),
		FirstName: create.FirstName,
		LastName:  create.LastName,
		BirthDate: create.BirthDate,
	}, nil
}

func (h *Handler) Update(
	_ *handles.It, update *commands.UpdatePerson,
) (data.Person, error) {
	return data.Person{
		Id:        update.Id,
		FirstName: update.FirstName,
		LastName:  update.LastName,
		BirthDate: update.BirthDate,
	}, nil
}

func (h *Handler) Delete(
	_ *handles.It, delete *commands.DeletePeople,
) error {
	return nil
}