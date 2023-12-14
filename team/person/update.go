package person

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/miruken-go/demo.microservice/team-api/commands"
	"github.com/miruken-go/demo.microservice/team-api/data"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	play "github.com/miruken-go/miruken/validates/play"
)

func (h *Handler) InitUpdate(
	_ *struct{ args.Optional }, translator ut.Translator,
) error {
	return h.Validates2.WithRules(
		play.Rules{
			play.Type[commands.UpdatePerson](map[string]string{
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
	_ *handles.It, update *commands.UpdatePerson,
) (data.Person, error) {
	return data.Person{
		Id:        update.Id,
		FirstName: update.FirstName,
		LastName:  update.LastName,
		BirthDate: update.BirthDate,
	}, nil
}
