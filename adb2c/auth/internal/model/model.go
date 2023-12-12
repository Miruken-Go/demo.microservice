package model

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"time"
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
		Type        string `json:"type"`  // always Entitlement
		Name        string `json:"name"`
		Scope       string `json:"scope"`
		Description string `json:"description"`
	}
)


// Subject

func (m *Subject) ToApi() api.Subject {
	ps := ParseIds(m.PrincipalIds)
	principals := make([]api.Principal, len(ps))
	for i, pid := range ps {
		principals[i] = api.Principal{Id: pid}
	}
	return api.Subject{
		Id:         ParseId(m.Id),
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
				Id: ParseId(pid.(string)),
			}
		}
	}
	return api.Subject{
		Id:         ParseId(m["id"].(string)),
		ObjectId:   m["objectId"].(string),
		Principals: principals,
	}
}


// Principal

func (m *Principal) ToApi() api.Principal {
	es := m.EntitlementNames
	entitlements := make([]api.Entitlement, len(es))
	for i, e := range es {
		entitlements[i] = api.Entitlement{Name:e}
	}
	return api.Principal{
		Id:           ParseId(m.Id),
		Type:         m.Type,
		Name:         m.Name,
		Domain:       m.Scope,
		Entitlements: entitlements,
	}
}


// Entitlement

func (m *Entitlement) ToApi() api.Entitlement {
	return api.Entitlement{
		Id:          ParseId(m.Id),
		Name:        m.Name,
		Domain:      m.Scope,
		Description: m.Description,
	}
}


const EntitlementType = "$E"