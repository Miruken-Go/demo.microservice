package person

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security"
	"github.com/miruken-go/miruken/security/authorizes"
	"github.com/miruken-go/miruken/security/jwt"
	"github.com/miruken-go/miruken/security/principal"
	play "github.com/miruken-go/miruken/validates/play"
	"sync/atomic"
	"time"
)

type (
	CreateIntegrity struct {
		play.ValidatorT[*commands.CreatePerson]
	}
)

func (h *Handler) AuthorizeCreate(
	_ *authorizes.It, _ *commands.CreatePerson,
	subject security.Subject,
) bool {
	return principal.All(subject, jwt.Scope("Person.Create"))
}

func (h *Handler) Create(
	_*struct {
		handles.It
		authorizes.Required
	  }, create *commands.CreatePerson,
	_*struct{ args.Optional }, parts api.PartContainer,
) (any, error) {
	person := data.Person{
		Id:        atomic.AddInt32(&h.nextId, 1),
		FirstName: create.FirstName,
		LastName:  create.LastName,
		BirthDate: create.BirthDate,
	}
	if parts == nil {
		return person, nil
	}
	var pb api.WritePartsBuilder
	pb.MainPart(pb.NewPart().
		MediaType("application/json").
		Metadata(map[string]any{
			"role": []string{"coach", "manager"},
		}).
		Body(person).
		Build()).
		AddParts(parts.Parts())
	return pb.Build(), nil
}


func (i *CreateIntegrity) Constructor(
	_*struct{args.Optional}, translator ut.Translator,
) error {
	return i.InitWithRules(
		play.Rules{
			play.Type[commands.CreatePerson](map[string]string{
				"FirstName": "required",
				"LastName":  "required",
				"BirthDate": "notfuture",
			}),
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
