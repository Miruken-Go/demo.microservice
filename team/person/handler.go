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
		{Id: 1, FirstName: "John", LastName: "Smith", BirthDate: time.Now()},
	}, nil
}

func (h *Handler) Create(
	_ *handles.It, create *commands.CreatePerson,
	_*struct{args.Optional}, parts api.PartContainer,
) (any, error) {
	person := data.Person{
		Id:        atomic.AddInt32(&h.nextId, 1),
		FirstName: create.FirstName,
		LastName:  create.LastName,
		BirthDate: create.BirthDate,
	}
	if parts == nil {
		return person, nil
	}
	var pb api.WritePartsBuilder
	pb.MainPart(pb.NewPart().
		 MediaType("application/json").
		 Metadata(map[string]any {
		 	"role": []string{"coach", "manager"},
	     }).
		 Body(person).
		 Build()).
	   AddParts(parts.Parts())
	return pb.Build(), nil
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
