package team

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"github.com/miruken-go/miruken/args"
	play "github.com/miruken-go/miruken/validates/play"
	"time"
)

type (
	CreateTeamIntegrity struct {
		play.ValidatorT[*commands.CreateTeam]
	}

	UpdateTeamIntegrity struct {
		play.ValidatorT[*commands.UpdateTeam]
	}
)


// CreateTeamIntegrity

func (i *CreateTeamIntegrity) Constructor(
	_*struct{args.Optional}, translator ut.Translator,
) error {
	return i.ConstructWithRules(
		play.Rules{
			{commands.CreateTeam{}, map[string]string{
				"Name": "required",
			}},
			{data.Coach{}, map[string]string{
				"Person": "required",
				"License": "required,len=10",
			}},
			{data.Manager{}, map[string]string{
				"Person": "required",
			}},
			{data.Person{}, map[string]string{
				"Id":        "eq=0",
				"FirstName": "required",
				"LastName":  "required",
				"BirthDate": "notfuture",
			}},
		}, func(v *validator.Validate) error {
			return v.RegisterValidation("notfuture", notfuture)
		}, translator)
}


// UpdateTeamIntegrity

func (i *UpdateTeamIntegrity) Constructor(
	_*struct{args.Optional}, translator ut.Translator,
) error {
	return i.ConstructWithRules(
		play.Rules{
			{commands.UpdateTeam{}, map[string]string{
				"Id": "required,gt=0",
				"Name": "omitempty,min=1",
			}},
			{data.Coach{}, map[string]string{
				"License": "omitempty,len=10",
			}},
			{data.Person{}, map[string]string{
				"Id":        "required,gt=0",
				"FirstName": "omitempty,min=1",
				"LastName":  "omitempty,min=1",
				"BirthDate": "notfuture",
			}},
		}, func(v *validator.Validate) error {
			return v.RegisterValidation("notfuture", notfuture)
		}, translator)
}


func notfuture(fl validator.FieldLevel) bool {
	if t, ok := fl.Field().Interface().(time.Time); ok {
		return t.Before(time.Now())
	}
	return false
}
