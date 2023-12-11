package api

import "github.com/google/uuid"

type (
	Subject struct {
		Id         uuid.UUID
		ObjectId   string
		Principals []Principal
	}

	Principal struct {
		Id           uuid.UUID
		Type         string
		Name         string
		Domain       string
		Entitlements []Entitlement
	}

	Entitlement struct {
		Id     uuid.UUID
		Name   string
		Domain string
	}
)
