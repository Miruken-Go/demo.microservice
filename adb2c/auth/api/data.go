package api

import "github.com/google/uuid"

type (
	Tag struct {
		Id          uuid.UUID
		Name        string
		Description string
	}

	Subject struct {
		Id         uuid.UUID
		Name       string
		Principals []Principal
	}

	Principal struct {
		Id           uuid.UUID
		Name         string
		Tags         []Tag
		Entitlements []Entitlement
	}

	Entitlement struct {
		Id   uuid.UUID
		Name string
		Tags []Tag
	}
)
