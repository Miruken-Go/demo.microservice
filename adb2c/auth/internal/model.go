package internal

import "time"

type (
	Subject struct {
		ID         string    `bson:"_id,omitempty"`
		CreatedAt  time.Time `bson:"created_at"`
		ModifiedAt time.Time `bson:"modified_at"`
	}

	Principal struct {
		ID      string  `bson:"_id,omitempty"`
		Name    string  `bson:"name"`
		TagIDs []string `bson:"tags"`
	}

	Entitlement struct {
		ID     string   `bson:"_id,omitempty"`
		Name   string   `bson:"name"`
		TagIDs []string `bson:"tags"`
	}

	SubjectPrincipal struct {
		SubjectID   string `bson:"subject_id"`
		PrincipalID string `bson:"principal_id"`
	}

	PrincipalEntitlement struct {
		PrincipalID   string `bson:"principal_id"`
		EntitlementID string `bson:"entitlement_id"`
	}

	Tag struct {
		ID          string `bson:"_id"`
		Name        string `bson:"name"`
		Description string `bson:"description"`
	}
)
