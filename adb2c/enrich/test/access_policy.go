package test

import (
	api2 "github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/miruken/security/authorizes"
)

// EnrichAccessPolicy holds the authorization policy for the test's
// local GetSubject stub, kept separate from EnrichTestSuite.Get - the
// same cross-cutting-concern separation used in production (e.g.
// subject.AccessPolicy). Keeps the Enrich prefix here since this type
// lives in package "test", not "enrich" - the prefix isn't redundant
// the way it would be inside the subject/user/principal/team packages.
type EnrichAccessPolicy struct{}

// AuthorizeGet is a placeholder authorization policy for the
// authorizes.Required marker on EnrichTestSuite.Get, needed now that
// authorizes.Required fails closed instead of allowing by default when
// no policy answers the check.
func (p *EnrichAccessPolicy) AuthorizeGet(
	_ *authorizes.It, _ api2.GetSubject,
) bool {
	return true
}
