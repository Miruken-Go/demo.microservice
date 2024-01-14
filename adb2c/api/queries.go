package api

type (
	// Users

	// ListUsers returns all users.
	ListUsers struct {
		Filter string
	}

	// Subjects

	// GetSubject retrieves a Subject by id.
	GetSubject struct {
		SubjectId string
	}

	// FindSubjects returns all subjects matching an optional filter.
	// The filter matches subjects containing ANY of the supplied
	// principals.
	// Set Exact to true to suppress Principal hierarchy evaluation.
	FindSubjects struct {
		Filter *struct{
			Scope        string
			PrincipalIds []string
			Exact        bool
		}
	}


	// Principals

	// GetPrincipal retrieves a Principal by id and scope.
	GetPrincipal struct {
		PrincipalId string
		Scope       string
	}

	// FindPrincipals returns all principals in a scope matching
	// optional criteria.
	// The Type and Name are used to constrain the search.
	FindPrincipals struct {
		Type  string
		Name  string
		Scope string
	}

	// ExpandPrincipals flattens a Principal hierarchy in a scope.
	// Set Squash true to squash related Principal hierarchies.
	// e.g. Role.A (Role.B, Role.C) => Role.B, Role.C
	//      otherwise Role.A, Role.B, Role.C
	ExpandPrincipals struct {
		Scope        string
		PrincipalIds []string
		Squash       bool
	}

	// ImpliedPrincipals returns all Principal ids that directly
	// or indirectly include the provided Principal ids.
	ImpliedPrincipals struct {
		Scope        string
		PrincipalIds []string
	}
)
