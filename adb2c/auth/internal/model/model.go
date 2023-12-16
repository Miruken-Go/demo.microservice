package model

import (
	"time"

	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
)

type (
	Subject struct {
		Id           string    `json:"id"`
		ObjectId     string    `json:"objectId"`
		PrincipalIds []string  `json:"principalIds"`
		CreatedAt    time.Time `json:"createdAt"`
	}
	SubjectMap map[string]any

	Principal struct {
		Id               string   `json:"id"`
		Type             string   `json:"type"`
		Name             string   `json:"name"`
		Scope            string   `json:"scope"`
		EntitlementNames []string `json:"entitlementNames"`
	}

	Entitlement struct {
		Id          string `json:"id"`
		Type        string `json:"type"` // always Entitlement
		Name        string `json:"name"`
		Scope       string `json:"scope"`
		Description string `json:"description"`
	}
)

// Subject

func (m *Subject) ToApi() api.Subject {
	ps := m.PrincipalIds
	principals := make([]api.Principal, len(ps))
	for i, pid := range ps {
		principals[i] = api.Principal{Id: pid}
	}
	return api.Subject{
		Id:         m.Id,
		ObjectId:   m.ObjectId,
		Principals: principals,
	}
}

func (m SubjectMap) ToApi() api.Subject {
	var principals []api.Principal
	if val, ok := m["principalIds"]; ok && val != nil {
		ps := val.([]any)
		principals = make([]api.Principal, len(ps))
		for i, pid := range ps {
			principals[i] = api.Principal{
				Id: pid.(string),
			}
		}
	}
	return api.Subject{
		Id:         m["id"].(string),
		ObjectId:   m["objectId"].(string),
		Principals: principals,
	}
}

// Principal

func (m *Principal) ToApi() api.Principal {
	es := m.EntitlementNames
	entitlements := make([]api.Entitlement, len(es))
	for i, e := range es {
		entitlements[i] = api.Entitlement{Name: e}
	}
	return api.Principal{
		Id:           m.Id,
		Type:         m.Type,
		Name:         m.Name,
		Domain:       m.Scope,
		Entitlements: entitlements,
	}
}

// Entitlement

func (m *Entitlement) ToApi() api.Entitlement {
	return api.Entitlement{
		Id:          m.Id,
		Name:        m.Name,
		Domain:      m.Scope,
		Description: m.Description,
	}
}

const EntitlementType = "$E"
