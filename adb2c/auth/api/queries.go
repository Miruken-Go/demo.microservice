package api

type (
	GetSubject struct {
		SubjectId string
	}

	FindSubjects struct {
		ObjectId   string
		Principals struct {
			All bool
			Ids []string
		}
	}

	GetPrincipal struct {
		PrincipalId string
		Domain      string
	}

	FindPrincipals struct {
		Type   string
		Name   string
		Domain string
	}

	GetEntitlement struct {
		EntitlementId string
		Domain        string
	}

	FindEntitlements struct {
		Name   string
		Domain string
	}
)
