package data

type (
	Scope string

	User struct {
		Id         string
		Principals map[Scope][]Principal
	}

	Principal struct {
		Id           string
		Type         string
		Name         string
		Entitlements []Entitlement
	}

	Entitlement struct {
		Id   string
		Name string
	}
)
