package api

import "github.com/google/uuid"

type (
	GetSubject struct {
		SubjectId uuid.UUID
	}

	FindSubjects struct {
		Principals []Principal
	}

	GetPrincipal struct {
		PrincipalId uuid.UUID
	}

	FindPrincipals struct {
		Tags []Tag
	}

	GetEntitlement struct {
		EntitlementId uuid.UUID
	}

	FindEntitlements struct {
		Tags []Tag
	}

	GetTag struct {
		TagId uuid.UUID
	}

	FindTags struct {
		Names []string
	}
)
