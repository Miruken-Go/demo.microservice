package team

import (
	"github.com/miruken-go/demo.microservice/team-api/commands"
	"github.com/miruken-go/miruken/security/authorizes"
)

// AccessPolicy holds the authorization policies for actions in this
// package, kept separate from the handlers that process those actions -
// the same cross-cutting-concern separation used for validation
// (CreateIntegrity/UpdateIntegrity).
type AccessPolicy struct{}

// AuthorizeCreate is a placeholder authorization policy for the
// authorizes.Required marker on Handler.Create. Miruken's
// authorizes.Required fails closed (denies) when no policy answers the
// check, rather than allowing by default; before this, Create had no
// policy at all (only the jwt.Scope "Team.Create" constraint) and was
// implicitly allowed for everyone once that scope check passed. This
// placeholder preserves that prior de facto behavior explicitly -
// replace with real authorization logic (e.g. role or entitlement
// checks like team/person's AuthorizeCreate) as needed.
func (p *AccessPolicy) AuthorizeCreate(
	_ *authorizes.It, _ *commands.CreateTeam,
) bool {
	return true
}
