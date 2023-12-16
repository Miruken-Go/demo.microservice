package api

type (
	// CreateSubject creates a new subject.
	CreateSubject struct {
		ObjectId     string
		PrincipalIds []string
	}
	SubjectCreated struct {
		SubjectId string
	}

	AssignPrincipals struct {
		SubjectId    string
		PrincipalIds []string
	}

	RevokePrincipals struct {
		SubjectId    string
		PrincipalIds []string
	}

	RemoveSubject struct {
		SubjectId string
	}

	// CreatePrincipal creates a new principal.
	CreatePrincipal struct {
		Type             string
		Name             string
		Domain           string
		EntitlementNames []string
	}
	PrincipalCreated struct {
		PrincipalId string
	}

	AssignEntitlements struct {
		PrincipalId      string
		Domain           string
		EntitlementNames []string
	}

	RevokeEntitlements struct {
		PrincipalId      string
		Domain           string
		EntitlementNames []string
	}

	RemovePrincipal struct {
		PrincipalId string
		Domain      string
	}

	// CreateEntitlement creates a new entitlement.
	CreateEntitlement struct {
		Name        string
		Domain      string
		Description string
	}
	EntitlementCreated struct {
		EntitlementId string
	}

	RemoveEntitlement struct {
		EntitlementId string
		Domain        string
	}
)
