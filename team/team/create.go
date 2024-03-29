package team

import (
	"sync/atomic"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/miruken-go/demo.microservice/team-api/commands"
	"github.com/miruken-go/demo.microservice/team-api/data"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/promise"
	"github.com/miruken-go/miruken/security/authorizes"
	"github.com/miruken-go/miruken/security/jwt"
	play "github.com/miruken-go/miruken/validates/play"
)

type (
	CreateIntegrity struct {
		play.Validates[*commands.CreateTeam]
	}
)

func (i *CreateIntegrity) Constructor(
	_ *struct{ args.Optional }, translator ut.Translator,
) error {
	return i.WithRules(
		play.Rules{
			play.Type[commands.CreateTeam](play.Constraints{
				"Name": "required",
			}),
			play.Type[data.Coach](play.Constraints{
				"Person":  "required",
				"License": "required,len=10",
			}),
			play.Type[data.Manager](play.Constraints{
				"Person": "required",
			}),
			play.Type[data.Person](play.Constraints{
				"Id":        "eq=0",
				"FirstName": "required",
				"LastName":  "required",
				"BirthDate": "notfuture",
			}),
		}, func(v *validator.Validate) error {
			return v.RegisterValidation("notfuture", notfuture)
		}, translator)
}

func (h *Handler) Create(
	_ *struct {
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

func notfuture(fl validator.FieldLevel) bool {
	if t, ok := fl.Field().Interface().(time.Time); ok {
		return t.Before(time.Now())
	}
	return false
}
