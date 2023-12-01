package internal

import (
	"github.com/google/uuid"
	"time"
)

type (
	Subject struct {
		ID         uuid.UUID `bson:"_id,omitempty"`
		Name       string    `bson:"name"`
		CreatedAt  time.Time `bson:"created_at"`
		ModifiedAt time.Time `bson:"modified_at"`
	}

	Principal struct {
		ID      uuid.UUID  `bson:"_id,omitempty"`
		Name    string     `bson:"name"`
		TagIDs []uuid.UUID `bson:"tags"`
	}

	Entitlement struct {
		ID     uuid.UUID  `bson:"_id,omitempty"`
		Name   string      `bson:"name"`
		TagIDs []uuid.UUID `bson:"tags"`
	}

	SubjectPrincipal struct {
		SubjectID   uuid.UUID `bson:"subject_id"`
		PrincipalID uuid.UUID `bson:"principal_id"`
	}

	PrincipalEntitlement struct {
		PrincipalID   uuid.UUID `bson:"principal_id"`
		EntitlementID uuid.UUID `bson:"entitlement_id"`
	}

	Tag struct {
		ID          uuid.UUID `bson:"_id"`
		Name        string    `bson:"name"`
		Description string    `bson:"description"`
	}
)
