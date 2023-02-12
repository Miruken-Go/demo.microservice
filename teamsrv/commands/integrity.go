package commands

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/validates"
	play "github.com/miruken-go/miruken/validates/play"
	"time"
)

type (
	CreatePersonIntegrity struct {
		play.Validator
	}
)

func (i *CreatePersonIntegrity) Constructor(
	_ *struct{ miruken.Optional }, translator ut.Translator,
) error {
	val := validator.New()
	if err := val.RegisterValidation("notfuture", notfuture); err != nil {
		return err
	}

	i.ConstructWithRules(
		play.Rules{
			{commands.CreatePerson{}, map[string]string{
				"FirstName": "required",
				"LastName":  "required",
				"BirthDate": "notfuture",
			}},
		}, val, translator)

	return nil
}

func (i *CreatePersonIntegrity) Validate(
	v *validates.It, create *commands.CreatePerson,
) miruken.HandleResult {
	return i.ValidateAndStop(create, v.Outcome())
}

func notfuture(fl validator.FieldLevel) bool {
	if t, ok := fl.Field().Interface().(time.Time); ok {
		return t.Before(time.Now())
	}
	return false
}
