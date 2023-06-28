package person

import (
	"github.com/miruken-go/demo.microservice/teamapi/commands"
	"github.com/miruken-go/miruken/security"
	"github.com/miruken-go/miruken/security/authorizes"
	"github.com/miruken-go/miruken/security/jwt"
	"github.com/miruken-go/miruken/security/principal"
)

type (
	AccessPolicy struct {}
)


func (p *AccessPolicy) AuthorizeCreate(
	_ *authorizes.It, _ commands.CreatePerson,
	subject security.Subject,
) bool {
	return principal.All(subject, jwt.Scope("Person.Create"))
}