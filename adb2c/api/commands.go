package api

type (
	// Subjects

	// CreateSubject creates a new Subject.
	CreateSubject struct {
		SubjectId string
	}
	SubjectCreated struct {
		SubjectId string
	}


	// Principals

	// AssignPrincipals assigns principals to a Subject.
	AssignPrincipals struct {
		SubjectId    string
		Scope        string
		PrincipalIds []string
	}

	// RevokePrincipals removes principals from a Subject.
	RevokePrincipals struct {
		SubjectId    string
		Scope        string
		PrincipalIds []string
	}

	// RemoveSubject removes a Subject.
	RemoveSubject struct {
		SubjectId string
	}

	// CreatePrincipal creates a new Principal.
	CreatePrincipal struct {
		Type        string
		Name        string
		Scope       string
		IncludedIds []string
	}
	PrincipalCreated struct {
		PrincipalId string
	}

	// IncludePrincipals add principals to a Principal hierarchy.
	// Principal hierarchies simplify the management of principals.
	IncludePrincipals struct {
		PrincipalId string
		Scope       string
		IncludedIds []string
	}

	// ExcludePrincipals removes principals from a Principal hierarchy.
	ExcludePrincipals struct {
		PrincipalId string
		Scope       string
		ExcludedIds []string
	}

	// RemovePrincipal removes a Principal.
	RemovePrincipal struct {
		PrincipalId string
		Scope       string
	}
)
