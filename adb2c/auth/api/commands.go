package api

import "github.com/google/uuid"

type (
	// CreateSubject creates a new subject.
	CreateSubject struct {
		ObjectId     string
		PrincipalIds []uuid.UUID
	}
	SubjectCreated struct {
		SubjectId uuid.UUID
	}

	AssignPrincipals struct {
		SubjectId    uuid.UUID
		PrincipalIds []uuid.UUID
	}

	RevokePrincipals struct {
		SubjectId    uuid.UUID
		PrincipalIds []uuid.UUID
	}

	RemoveSubject struct {
		SubjectId uuid.UUID
	}


	// CreatePrincipal creates a new principal.
	CreatePrincipal struct {
		Type           string
		Name           string
		Domain         string
		EntitlementIds []uuid.UUID
	}
	PrincipalCreated struct {
		PrincipalId uuid.UUID
	}

	AssignEntitlements struct {
		PrincipalId    uuid.UUID
		Domain         string
		EntitlementIds []uuid.UUID
	}

	RevokeEntitlements struct {
		PrincipalId    uuid.UUID
		Domain         string
		EntitlementIds []uuid.UUID
	}

	RemovePrincipal struct {
		PrincipalId uuid.UUID
		Domain      string
	}


	// CreateEntitlement creates a new entitlement.
	CreateEntitlement struct {
		Name   string
		Domain string
	}
	EntitlementCreated struct {
		EntitlementId uuid.UUID
	}

	RemoveEntitlement struct {
		EntitlementId uuid.UUID
		Domain        string
	}
)

