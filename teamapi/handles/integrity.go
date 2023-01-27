package handles

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/validate"
	play "github.com/miruken-go/miruken/validate/play"
	"time"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	CreatePersonIntegrity struct {
		play.Base
	}
)

func (v *CreatePersonIntegrity) Constructor(
	_ *struct{ miruken.Optional }, translator ut.Translator,
) error {
	val := validator.New()
	if err := val.RegisterValidation("notfuture", notfuture); err != nil {
		return err
	}

	v.ConstructWithRules(
		play.Rules{
			{commands.CreatePerson{}, map[string]string{
				"FirstName": "required",
				"LastName":  "required",
				"BirthDate": "notfuture",
			}},
		}, val, translator)

	return nil
}

func (v *CreatePersonIntegrity) Validate(
	validates *validate.Validates, create *commands.CreatePerson,
) miruken.HandleResult {
	return v.Base.ValidateAndStop(create, validates.Outcome())
}

func notfuture(fl validator.FieldLevel) bool {
	if t, ok := fl.Field().Interface().(time.Time); ok {
		return t.Before(time.Now())
	}
	return false
}
