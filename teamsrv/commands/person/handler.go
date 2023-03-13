package person

import (
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"github.com/miruken-go/demo.microservice/teamapi/queries"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"sync/atomic"
	"time"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		nextId int32
	}
)


func (h *Handler) Find(
	_ *handles.It, find queries.FindPeople,
) ([]data.Person, error) {
	return []data.Person{
		{1, "John", "Smith", time.Now()},
	}, nil
}

func (h *Handler) Create(
	_ *handles.It, create *commands.CreatePerson,
	_*struct{args.Optional}, parts api.PartContainer,
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
