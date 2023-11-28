package api

type (
	Tag struct {
		Id          string
		Name        string
		Description string
	}

	Subject struct {
		Id         string
		Principals []Principal
	}

	Principal struct {
		Id           string
		Name         string
		Tags         []Tag
		Entitlements []Entitlement
	}

	Entitlement struct {
		Id   string
		Name string
		Tags []Tag
	}
)
