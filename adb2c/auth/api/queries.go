package api

import "github.com/google/uuid"

type (
	GetSubject struct {
		SubjectId uuid.UUID
	}

	FindSubjects struct {
		ObjectId   string
		Principals struct {
			All bool
			Ids []uuid.UUID
		}
	}

	GetPrincipal struct {
		PrincipalId uuid.UUID
		Domain      string
	}

	FindPrincipals struct {
		Type   string
		Name   string
		Domain string
	}

	GetEntitlement struct {
		EntitlementId uuid.UUID
		Domain        string
	}

	FindEntitlements struct {
		Name   string
		Domain string
	}
)
