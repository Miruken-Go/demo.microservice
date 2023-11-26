package data

type (
	User struct {
		Id         string
		Principals map[string][]Principal
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
