package storage

type (
	User struct {
		ID    string `bson:"_id,omitempty"`
		Email string `bson:"email"`
	}

	Principal struct {
		ID   string `bson:"_id,omitempty"`
		Name string `bson:"name"`
		Type string `bson:"type"`
	}

	Entitlement struct {
		ID   string `bson:"_id,omitempty"`
		Name string `bson:"name"`
	}

	UserPrincipal struct {
		UserID      string `bson:"user_id"`
		PrincipalID string `bson:"principal_id"`
		Scope       string `bson:"scope"`
	}

	PrincipalEntitlement struct {
		PrincipalID   string `bson:"principal_id"`
		EntitlementID string `bson:"entitlement_id"`
	}
)
