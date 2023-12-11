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
		Id             string   `json:"id"`
		Type           string   `json:"type"`
		Name           string   `json:"name"`
		Scope          string   `json:"scope"`
		EntitlementIds []string `json:"entitlementIds"`
	}
	PrincipalMap map[string]any

	Entitlement struct {
		Id    string `json:"id"`
		Type  string `json:"type"`  // always Entitlement
		Name  string `json:"name"`
		Scope string `json:"scope"`
	}
	EntitlementMap map[string]any
)


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


func (m *Principal) ToApi() api.Principal {
	es := ParseIds(m.EntitlementIds)
	entitlements := make([]api.Entitlement, len(es))
	for i, eid := range es {
		entitlements[i] = api.Entitlement{Id: eid}
	}
	return api.Principal{
		Id:           ParseId(m.Id),
		Type:         m.Type,
		Name:         m.Name,
		Domain:       m.Scope,
		Entitlements: entitlements,
	}
}

func (m PrincipalMap) ToApi() api.Principal {
	var entitlements []api.Entitlement
	if val, ok := m["entitlementIds"]; ok && val != nil {
		es := val.([]any)
		entitlements = make([]api.Entitlement, len(es))
		for i, eid := range es {
			entitlements[i] = api.Entitlement{
				Id: ParseId(eid.(string)),
			}
		}
	}
	return api.Principal{
		Id:           ParseId(m["id"].(string)),
		Type:         m["type"].(string),
		Name:         m["name"].(string),
		Domain:       m["scope"].(string),
		Entitlements: entitlements,
	}
}


func (m *Entitlement) ToApi() api.Entitlement {
	return api.Entitlement{
		Id:     ParseId(m.Id),
		Name:   m.Name,
		Domain: m.Scope,
	}
}

func (m EntitlementMap) ToApi() api.Entitlement {
	return api.Entitlement{
		Id:     ParseId(m["id"].(string)),
		Name:   m["name"].(string),
		Domain: m["scope"].(string),
	}
}