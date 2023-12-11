package model

import (
	"github.com/google/uuid"
	"time"
)

type (
	Subject struct {
		Id           string    `json:"id"`
		ObjectId     string    `json:"objectId"`
		PrincipalIds []string  `json:"principalIds"`
		CreatedAt    time.Time `json:"createdAt"`
	}

	Principal struct {
		Id             string   `json:"id"`
		Type           string   `json:"type"`
		Name           string   `json:"name"`
		Scope          string   `json:"scope"`
		EntitlementIds []string `json:"entitlementIds"`
	}

	Entitlement struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	SubjectM struct {
		ID        uuid.UUID `bson:"_id,omitempty"`
		ObjectID  string    `bson:"object_id,omitempty"`
		CreatedAt time.Time `bson:"created_at"`
	}

	PrincipalM struct {
		ID     uuid.UUID   `bson:"_id,omitempty"`
		Name   string      `bson:"name"`
		TagIDs []uuid.UUID `bson:"tags"`
	}

	EntitlementM struct {
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
