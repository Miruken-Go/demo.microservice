package team

import (
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"github.com/miruken-go/demo.microservice/teamapi/queries"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/promise"
	"github.com/miruken-go/miruken/security/authorizes"
	"github.com/miruken-go/miruken/security/jwt"
	"sync/atomic"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		nextId int32
	}
)

func (h *Handler) Find(
	_ *handles.It, find queries.FindTeams,
) ([]data.Team, error) {
	return []data.Team{
		{Id: 1, Name: "Breakaway", Color: data. ColorOrange},
	}, nil
}

func (h *Handler) Create(
	_*struct {
		handles.It
		authorizes.Required
		jwt.Scope `name:"Team.Create"`
	  }, create *commands.CreateTeam,
	ctx miruken.HandleContext,
) *promise.Promise[data.Team] {
	composer := ctx.Composer
	cp, _, err := api.Send[data.Person](composer,
		&commands.CreatePerson{
			FirstName: create.Coach.Person.FirstName,
			LastName:  create.Coach.Person.LastName,
			BirthDate: create.Coach.Person.BirthDate,
		})
	if err != nil {
		return promise.Reject[data.Team](err)
	}
	create.Coach.Person = &cp
	mp, _, err := api.Send[data.Person](composer,
		&commands.CreatePerson{
			FirstName: create.Manager.Person.FirstName,
			LastName:  create.Manager.Person.LastName,
			BirthDate: create.Manager.Person.BirthDate,
		})
	if err != nil {
		return promise.Reject[data.Team](err)
	}
	create.Manager.Person = &mp
	return promise.Resolve(data.Team{
		Id:      atomic.AddInt32(&h.nextId, 1),
		Name:    create.Name,
		Coach:   &create.Coach,
		Manager: &create.Manager,
	})
}

func (h *Handler) Update(
	_ *handles.It, update *commands.UpdateTeam,
	ctx miruken.HandleContext,
) *promise.Promise[data.Team] {
	composer := ctx.Composer
	if coach := update.Coach; coach != nil {
		cp, _, err := api.Send[data.Person](composer,
			&commands.UpdatePerson{
				Id:        coach.Person.Id,
				FirstName: coach.Person.FirstName,
				LastName:  coach.Person.LastName,
				BirthDate: coach.Person.BirthDate,
			})
		if err != nil {
			return promise.Reject[data.Team](err)
		}
		coach.Person = &cp
	}
	if manager := update.Manager; manager != nil {
		mp, _, err := api.Send[data.Person](composer,
			&commands.UpdatePerson{
				Id:        manager.Person.Id,
				FirstName: manager.Person.FirstName,
				LastName:  manager.Person.LastName,
				BirthDate: manager.Person.BirthDate,
			})
		if err != nil {
			return promise.Reject[data.Team](err)
		}
		manager.Person = &mp
	}
	return promise.Resolve(data.Team{
		Id:      update.Id,
		Name:    update.Name,
		Coach:   update.Coach,
		Manager: update.Manager,
	})
}
