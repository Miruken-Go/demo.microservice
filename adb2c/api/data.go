package api

type (
	User struct {

	}

	ScopedPrincipals struct {
		Scope      string
		Principals []Principal
	}

	Subject struct {
		Id     string
		Scopes []ScopedPrincipals
	}

	Principal struct {
		Id       string
		Type     string
		Name     string
		Includes []Principal
	}
)
