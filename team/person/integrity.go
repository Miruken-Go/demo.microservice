package person

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/miruken/args"
	play "github.com/miruken-go/miruken/validates/play"
	"time"
)

type (
	CreatePersonIntegrity struct {
		play.ValidatorT[*commands.CreatePerson]
	}

	UpdatePersonIntegrity struct {
		play.ValidatorT[*commands.UpdatePerson]
	}
)


// CreatePersonIntegrity

func (i *CreatePersonIntegrity) Constructor(
	_*struct{args.Optional}, translator ut.Translator,
) error {
	return i.ConstructWithRules(
		play.Rules{
			{commands.CreatePerson{}, map[string]string{
				"FirstName": "required",
				"LastName":  "required",
				"BirthDate": "notfuture",
			}},
		}, func(v *validator.Validate) error {
			return v.RegisterValidation("notfuture", notfuture)
		}, translator)
}


// UpdatePersonIntegrity

func (i *UpdatePersonIntegrity) Constructor(
	_*struct{args.Optional}, translator ut.Translator,
) error {
	return i.ConstructWithRules(
		play.Rules{
			{commands.UpdatePerson{}, map[string]string{
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
