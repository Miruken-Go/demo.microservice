package api

import "github.com/google/uuid"

type (
	// CreateSubject creates a new subject.
	CreateSubject struct {
		ObjectId   string
		Principals []Principal
	}

	AssignPrincipals struct {
		SubjectId  uuid.UUID
		Principals []Principal
	}

	RevokePrincipals struct {
		SubjectId  uuid.UUID
		Principals []Principal
	}

	RemoveSubjects struct {
		Subjects []Subject
	}


	// CreatePrincipal creates a new principal.
	CreatePrincipal struct {
		Name         string
		Tags         []Tag
		Entitlements []Entitlement
	}

	TagPrincipal struct {
		PrincipalId uuid.UUID
		Tags        []Tag
	}

	UntagPrincipal struct {
		PrincipalId uuid.UUID
		Tags        []Tag
	}

	AssignEntitlements struct {
		PrincipalId  uuid.UUID
		Entitlements []Entitlement
	}

	RevokeEntitlements struct {
		PrincipalId  uuid.UUID
		Entitlements []Entitlement
	}

	RemovePrincipals struct {
		Principals []Principal
	}


	// CreateEntitlement creates a new entitlement.
	CreateEntitlement struct {
		Name string
		Tags []Tag
	}

	TagEntitlement struct {
		EntitlementId uuid.UUID
		Tags          []Tag
	}

	UntagEntitlement struct {
		EntitlementId uuid.UUID
		Tags          []Tag
	}

	RemoveEntitlements struct {
		Entitlements []Entitlement
	}


	CreateTag struct {
		Name        string
		Description string
	}

	RemoveTags struct {
		Tags []Tag
	}
)

