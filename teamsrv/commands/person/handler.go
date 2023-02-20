package person

import (
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"github.com/miruken-go/miruken"
	"sync/atomic"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		nextId int32
	}
)

func (h *Handler) Create(
	_ *struct {
		miruken.Handles
	}, create *commands.CreatePerson,
) (data.Person, error) {
	return data.Person{
		Id:        atomic.AddInt32(&h.nextId, 1),
		FirstName: create.FirstName,
		LastName:  create.LastName,
		BirthDate: create.BirthDate,
	}, nil
}

func (h *Handler) Update(
	_ *struct {
		miruken.Handles
	}, update *commands.UpdatePerson,
) (data.Person, error) {
	return data.Person{
		Id:        update.Id,
		FirstName: update.FirstName,
		LastName:  update.LastName,
		BirthDate: update.BirthDate,
	}, nil
}

func (h *Handler) Delete(
	_ *struct {
	miruken.Handles
}, delete *commands.DeletePeople,
) error {
	return nil
}