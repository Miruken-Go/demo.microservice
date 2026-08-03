package user

import (
	"github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/miruken/security/authorizes"
)

// UserAccessPolicy holds the authorization policies for actions in this
// package, kept separate from Handler.
type UserAccessPolicy struct{}

// AuthorizeList is a placeholder authorization policy for the
// authorizes.Required marker on Handler.List. Miruken's
// authorizes.Required fails closed (denies) when no policy answers the
// check, rather than allowing by default; before this, List had no
// policy at all and was implicitly allowed for everyone. This
// placeholder preserves that prior de facto behavior explicitly -
// replace with real authorization logic as needed.
func (p *UserAccessPolicy) AuthorizeList(
	_ *authorizes.It, _ api.ListUsers,
) bool {
	return true
}
