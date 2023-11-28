package api

type (
	GetSubject struct {
		SubjectId string
	}

	FindSubjects struct {
		Principals []Principal
	}

	GetPrincipal struct {
		PrincipalId string
	}

	FindPrincipals struct {
		Tags []Tag
	}

	GetEntitlement struct {
		EntitlementId string
	}

	FindEntitlements struct {
		Tags []Tag
	}

	GetTag struct {
		TagId string
	}

	FindTags struct {
		Names []string
	}
)
