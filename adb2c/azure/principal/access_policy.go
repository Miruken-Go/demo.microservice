package principal

import (
	"github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/miruken/security/authorizes"
)

// AccessPolicy holds the authorization policies for actions in
// this package, kept separate from Handler - the same cross-cutting-
// concern separation used for validation (Handler.setValidationRules).
type AccessPolicy struct{}

// The methods below are placeholder authorization policies for the
// authorizes.Required markers on Handler's methods. Miruken's
// authorizes.Required fails closed (denies) when no policy answers the
// check, rather than allowing by default; before this, these methods
// had no policy at all and were implicitly allowed for everyone. These
// placeholders preserve that prior de facto behavior explicitly rather
// than leaving it accidental - replace with real authorization logic
// as needed.

func (p *AccessPolicy) AuthorizeCreate(
	_ *authorizes.It, _ api.CreatePrincipal,
) bool {
	return true
}

func (p *AccessPolicy) AuthorizeInclude(
	_ *authorizes.It, _ api.IncludePrincipals,
) bool {
	return true
}

func (p *AccessPolicy) AuthorizeExclude(
	_ *authorizes.It, _ api.ExcludePrincipals,
) bool {
	return true
}

func (p *AccessPolicy) AuthorizeRemove(
	_ *authorizes.It, _ api.RemovePrincipal,
) bool {
	return true
}
