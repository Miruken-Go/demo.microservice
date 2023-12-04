package api

import "github.com/google/uuid"

type (
	GetSubject struct {
		SubjectId uuid.UUID
	}

	FindSubjects struct {
		PrincipalIds []uuid.UUID
	}

	GetPrincipal struct {
		PrincipalId uuid.UUID
	}

	FindPrincipals struct {
		Name string
		Tags []Tag
	}

	GetEntitlement struct {
		EntitlementId uuid.UUID
	}

	FindEntitlements struct {
		Name   string
		TagIds []uuid.UUID
	}

	GetTag struct {
		TagId uuid.UUID
	}

	FindTags struct {
		Name string
	}
)
