package api

type (
	GetSubject struct {
		SubjectId string
	}

	FindSubjects struct {
		Filter *struct{
			Scope        string
			PrincipalIds []string
			All          bool
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

	FlattenPrincipals struct {
		Scope        string
		PrincipalIds []string
	}
)
