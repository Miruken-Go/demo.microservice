package person

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

	UpdatePersonIntegrity struct {
		play.Validator
	}
)


// CreatePersonIntegrity

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


// UpdatePersonIntegrity

func (i *UpdatePersonIntegrity) Constructor(
	_ *struct{ miruken.Optional }, translator ut.Translator,
) error {
	val := validator.New()
	if err := val.RegisterValidation("notfuture", notfuture); err != nil {
		return err
	}

	i.ConstructWithRules(
		play.Rules{
			{commands.UpdatePerson{}, map[string]string{
				"Id":        "required,gt=0",
				"FirstName": "omitempty,min=1",
				"LastName":  "omitempty,min=1",
				"BirthDate": "notfuture",
			}},
		}, val, translator)

	return nil
}

func (i *UpdatePersonIntegrity) Validate(
	v *validates.It, update *commands.UpdatePerson,
) miruken.HandleResult {
	return i.ValidateAndStop(update, v.Outcome())
}


func notfuture(fl validator.FieldLevel) bool {
	if t, ok := fl.Field().Interface().(time.Time); ok {
		return t.Before(time.Now())
	}
	return false
}
