package api

type (
	// CreateSubject creates a new subject.
	CreateSubject struct {
		Principals []Principal
	}

	AssignPrincipals struct {
		SubjectId    string
		PrincipalIds []string
	}

	RevokePrincipals struct {
		SubjectId    string
		PrincipalIds []string
	}

	RemoveSubjects struct {
		SubjectIds []string
	}


	// CreatePrincipal creates a new principal.
	CreatePrincipal struct {
		Name         string
		Tags         []Tag
		Entitlements []Entitlement
	}

	TagPrincipal struct {
		PrincipalId string
		TagIds      []string
	}

	UntagPrincipal struct {
		PrincipalId string
		TagIds      []string
	}

	AssignEntitlements struct {
		PrincipalId    string
		EntitlementIds []string
	}

	RevokeEntitlements struct {
		PrincipalId    string
		EntitlementIds []string
	}

	RemovePrincipals struct {
		PrincipalIds []string
	}


	// CreateEntitlement creates a new entitlement.
	CreateEntitlement struct {
		Name string
		Tags []Tag
	}

	TagEntitlement struct {
		EntitlementId string
		TagIds        []string
	}

	UntagEntitlement struct {
		EntitlementId string
		TagIds        []string
	}

	RemoveEntitlements struct {
		EntitlementIds []string
	}


	CreateTag struct {
		Name        string
		Description string
	}

	RemoveTags struct {
		TagIds []string
	}
)

