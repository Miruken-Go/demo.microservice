package api

type (
	// CreateSubject creates a new subject.
	CreateSubject struct {
		SubjectId string
	}
	SubjectCreated struct {
		SubjectId string
	}

	AssignPrincipals struct {
		SubjectId    string
		Scope        string
		PrincipalIds []string
	}

	RevokePrincipals struct {
		SubjectId    string
		Scope        string
		PrincipalIds []string
	}

	RemoveSubject struct {
		SubjectId string
	}

	// CreatePrincipal creates a new principal.
	CreatePrincipal struct {
		Type        string
		Name        string
		Scope       string
		IncludedIds []string
	}
	PrincipalCreated struct {
		PrincipalId string
	}

	IncludePrincipals struct {
		PrincipalId string
		Scope       string
		IncludedIds []string
	}

	ExcludePrincipals struct {
		PrincipalId string
		Scope       string
		ExcludedIds []string
	}

	RemovePrincipal struct {
		PrincipalId string
		Scope       string
	}
)
