package person

import (
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"github.com/miruken-go/demo.microservice/teamapi/queries"
	"github.com/miruken-go/miruken/handles"
	play "github.com/miruken-go/miruken/validates/play"
	"time"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		play.Validates1[*commands.CreatePerson]
		play.Validates2[*commands.UpdatePerson]
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

