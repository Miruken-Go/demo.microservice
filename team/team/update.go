package team

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/miruken-go/demo.microservice/team-api/commands"
	"github.com/miruken-go/demo.microservice/team-api/data"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/promise"
	play "github.com/miruken-go/miruken/validates/play"
)

type (
	UpdateIntegrity struct {
		play.Validates[*commands.UpdateTeam]
	}
)

func (i *UpdateIntegrity) Constructor(
	_*struct{args.Optional}, translator ut.Translator,
) error {
	return i.WithRules(
		play.Rules{
			play.Type[commands.UpdateTeam](map[string]string{
				"Id": "required,gt=0",
				"Name": "omitempty,min=1",
			}),
			play.Type[data.Coach](map[string]string{
				"License": "omitempty,len=10",
			}),
			play.Type[data.Person](map[string]string{
				"Id":        "required,gt=0",
				"FirstName": "omitempty,min=1",
				"LastName":  "omitempty,min=1",
				"BirthDate": "notfuture",
			}),
		}, func(v *validator.Validate) error {
			return v.RegisterValidation("notfuture", notfuture)
		}, translator)
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
