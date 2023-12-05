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
		Name           string
		TagIds         []uuid.UUID
		EntitlementIds []uuid.UUID
	}
	PrincipalCreated struct {
		PrincipalId uuid.UUID
	}

	TagPrincipal struct {
		PrincipalId uuid.UUID
		TagIds      []uuid.UUID
	}

	UntagPrincipal struct {
		PrincipalId uuid.UUID
		TagIds      []uuid.UUID
	}

	AssignEntitlements struct {
		PrincipalId    uuid.UUID
		EntitlementIds []uuid.UUID
	}

	RevokeEntitlements struct {
		PrincipalId    uuid.UUID
		EntitlementIds []uuid.UUID
	}

	RemovePrincipal struct {
		PrincipalId uuid.UUID
	}


	// CreateEntitlement creates a new entitlement.
	CreateEntitlement struct {
		Name   string
		TagIds []uuid.UUID
	}
	EntitlementCreated struct {
		EntitlementId uuid.UUID
	}

	TagEntitlement struct {
		EntitlementId uuid.UUID
		TagIds        []uuid.UUID
	}

	UntagEntitlement struct {
		EntitlementId uuid.UUID
		TagIds        []uuid.UUID
	}

	RemoveEntitlement struct {
		EntitlementId uuid.UUID
	}


	CreateTag struct {
		Name        string
		Description string
	}
	TagCreated struct {
		TagId uuid.UUID
	}

	RemoveTag struct {
		TagId uuid.UUID
	}
)

