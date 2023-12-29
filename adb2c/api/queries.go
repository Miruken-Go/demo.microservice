package api

type (
	GetSubject struct {
		SubjectId string
	}

	FindSubjects struct {
		Filter *struct{
			Scope        string
			PrincipalIds []string
			Exact        bool
		}
	}

	GetPrincipal struct {
		PrincipalId string
		Scope       string
	}

	FindPrincipals struct {
		Type  string
		Name  string
		Scope string
	}

	ExpandPrincipals struct {
		Scope        string
		PrincipalIds []string
	}

	SatisfyPrincipals struct {
		Scope        string
		PrincipalIds []string
	}
)
