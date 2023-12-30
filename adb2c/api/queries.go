package api

type (
	// Users

	// ListUsers returns all users.
	ListUsers struct {

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
	ExpandPrincipals struct {
		Scope        string
		PrincipalIds []string
	}

	// SatisfyPrincipals returns all Principal ids that are ancestors
	// of the provided Principal ids.
	SatisfyPrincipals struct {
		Scope        string
		PrincipalIds []string
	}
)
