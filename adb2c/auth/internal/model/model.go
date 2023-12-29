package model

import (
	"time"

	"github.com/miruken-go/demo.microservice/adb2c/api"
)

type (
	Scope struct {
		Name         string   `json:"name"`
		PrincipalIds []string `json:"principalIds"`
	}

	Subject struct {
		Id        string    `json:"id"`
		Scopes    []Scope   `json:"scopes"`
		CreatedAt time.Time `json:"createdAt"`
	}
	SubjectMap map[string]any

	Principal struct {
		Id          string   `json:"id"`
		Type        string   `json:"type"`
		Name        string   `json:"name"`
		Scope       string   `json:"scope"`
		IncludedIds []string `json:"includedIds"`
	}
)

// Subject

func (m *Subject) ToApi() api.Subject {
	scopes   := m.Scopes
	scopedPs := make([]api.ScopedPrincipals, len(scopes))
	for i, scope := range scopes {
		principals := make([]api.Principal, len(scope.PrincipalIds))
		for j, pid := range scope.PrincipalIds {
			principals[j] = api.Principal{Id: pid}
		}
		scopedPs[i] = api.ScopedPrincipals{Scope: scope.Name, Principals: principals}
	}
	return api.Subject{
		Id:     m.Id,
		Scopes: scopedPs,
	}
}

func (m SubjectMap) ToApi() api.Subject {
	var scopedPs []api.ScopedPrincipals
	if val, ok := m["scopes"]; ok && val != nil {
		scopes := val.([]any)
		scopedPs = make([]api.ScopedPrincipals, len(scopes))
		for i, scope := range scopes {
			scopeMap     := scope.(map[string]any)
			scopeName    := scopeMap["name"].(string)
			principalIds := scopeMap["principalIds"].([]any)
			ps  := make([]api.Principal, len(principalIds))
			for j, pid := range principalIds {
				ps[j] = api.Principal{Id: pid.(string)}
			}
			scopedPs[i] = api.ScopedPrincipals{Scope: scopeName, Principals: ps}
		}
	}

	return api.Subject{
		Id:     m["id"].(string),
		Scopes: scopedPs,
	}
}

// Principal

func (m *Principal) ToApi() api.Principal {
	ids      := m.IncludedIds
	included := make([]api.Principal, len(ids))
	for i, id := range ids {
		included[i] = api.Principal{Id: id}
	}
	return api.Principal{
		Id:       m.Id,
		Type:     m.Type,
		Name:     m.Name,
		Includes: included,
	}
}

