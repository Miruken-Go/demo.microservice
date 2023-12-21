package person

import (
	"sync/atomic"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/miruken-go/demo.microservice/team-api/commands"
	"github.com/miruken-go/demo.microservice/team-api/data"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security"
	"github.com/miruken-go/miruken/security/authorizes"
	"github.com/miruken-go/miruken/security/principal"
	play "github.com/miruken-go/miruken/validates/play"
)

func (h *Handler) InitCreate(
	_ *struct{ args.Optional }, translator ut.Translator,
) error {
	return h.Validates1.WithRules(
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

func (h *Handler) AuthorizeCreate(
	_ *authorizes.It, _ *commands.CreatePerson,
	subject security.Subject,
) bool {
	return principal.Any(subject,
		principal.Role("Manager"),
		principal.Entitlement("Player.Add"),
	)
}

func (h *Handler) Create(
	_ *struct {
		handles.It
		authorizes.Required
	  }, create *commands.CreatePerson,
	_ *struct{ args.Optional }, parts api.PartContainer,
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

func notfuture(fl validator.FieldLevel) bool {
	if t, ok := fl.Field().Interface().(time.Time); ok {
		return t.Before(time.Now())
	}
	return false
}
