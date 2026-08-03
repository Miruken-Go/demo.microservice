package subject

import (
	"github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/miruken/security/authorizes"
)

// SubjectAccessPolicy holds the authorization policies for actions in
// this package, kept separate from Handler - the same cross-cutting-
// concern separation used for validation (Handler.setValidationRules).
type SubjectAccessPolicy struct{}

// The methods below are placeholder authorization policies for the
// authorizes.Required markers on Handler's methods. Miruken's
// authorizes.Required fails closed (denies) when no policy answers the
// check, rather than allowing by default; before this, these methods
// had no policy at all and were implicitly allowed for everyone. These
// placeholders preserve that prior de facto behavior explicitly rather
// than leaving it accidental - replace with real authorization logic
// (e.g. role or entitlement checks like team/person's AuthorizeCreate)
// as needed.

func (p *SubjectAccessPolicy) AuthorizeCreate(
	_ *authorizes.It, _ api.CreateSubject,
) bool {
	return true
}

func (p *SubjectAccessPolicy) AuthorizeAssign(
	_ *authorizes.It, _ api.AssignPrincipals,
) bool {
	return true
}

func (p *SubjectAccessPolicy) AuthorizeRevoke(
	_ *authorizes.It, _ api.RevokePrincipals,
) bool {
	return true
}

func (p *SubjectAccessPolicy) AuthorizeRemove(
	_ *authorizes.It, _ api.RemoveSubject,
) bool {
	return true
}

func (p *SubjectAccessPolicy) AuthorizeGet(
	_ *authorizes.It, _ api.GetSubject,
) bool {
	return true
}

func (p *SubjectAccessPolicy) AuthorizeFind(
	_ *authorizes.It, _ api.FindSubjects,
) bool {
	return true
}
