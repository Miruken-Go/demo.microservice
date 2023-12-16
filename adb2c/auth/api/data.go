package api

type (
	Subject struct {
		Id         string
		ObjectId   string
		Principals []Principal
	}

	Principal struct {
		Id           string
		Type         string
		Name         string
		Domain       string
		Entitlements []Entitlement
	}

	Entitlement struct {
		Id          string
		Name        string
		Domain      string
		Description string
	}
)
