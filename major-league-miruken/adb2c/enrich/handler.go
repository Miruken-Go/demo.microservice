package enrich

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
	api2 "github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/args"
)

type (
	// Handler is an Azure ADB2C Api Connector that provides
	// an external source of claims to enrich tokens during user flows.
	// https://learn.microsoft.com/en-us/azure/active-directory-b2c/add-api-connector-token-enrichment?pivots=b2c-custom-policy
	Handler struct {
		logger logr.Logger
	}

	// Request is the request body sent to the Api Connector.
	// It receives a set of InputClaims and returns a set of OutputClaims.
	//   objectId - the objectId of the user in ADB2C
	//   scope - the scope of the token to enrich
	Request struct {
		ObjectId string `json:"objectId"`
		Scope    string `json:"scope"`
	}
)

func (e *Handler) Constructor(
	_ *struct{ args.Optional }, logger logr.Logger,
) {
	if logger == e.logger {
		e.logger = logr.Discard()
	} else {
		e.logger = logger
	}
}

func (e *Handler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
	h miruken.Handler,
) {
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	objectId := request.ObjectId
	if objectId == "" {
		http.Error(w, "objectId is required", http.StatusUnprocessableEntity)
		return
	}

	scope := request.Scope
	if scope == "" {
		http.Error(w, "scope is required", http.StatusUnprocessableEntity)
		return
	}

	e.logger.Info("Enrich enrich", "ObjectId", objectId, "Scope", scope)

	domain, principals, err := e.parseScope(scope)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	e.logger.Info("Parsed scope", "Domain", domain, "Principals", principals)

	var claims map[string]any

	s, ps, err := api.Send[api2.Subject](h, api2.GetSubject{SubjectId: objectId})
	if ps != nil {
		s, err = ps.Await()
	}

	var nh *miruken.NotHandledError
	if errors.As(err, &nh) {
		// return nothing if subject not found
		w.WriteHeader(http.StatusOK)
		return
	} else if err == nil {
		claims, err = e.getClaims(domain, principals, s.Scopes, h)
	}

	if err == nil {
		e.logger.Info("Enriched enrich", "Claims", claims)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(claims)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (e *Handler) getClaims(
	domain           string,
	principalTypes   []string,
	scopedPrincipals []api2.ScopedPrincipals,
	h                miruken.Handler,
) (map[string]any, error) {
	claims := make(map[string]any)

	var principalIds []string
	for _, sp := range scopedPrincipals {
		if sp.Scope == domain {
			for _, p := range sp.Principals {
				principalIds = append(principalIds, p.Id)
			}
			break
		}
	}

	if len(principalIds) == 0 {
		return claims, nil
	}

	ep, epp, err := api.Send[[]api2.Principal](h, api2.ExpandPrincipals{
		Scope:        domain,
		PrincipalIds: principalIds,
		Squash:       true,
	})
	if epp != nil {
		ep, err = epp.Await()
	}
	if err != nil {
		return nil, err
	}

	principalTypeMap := make(map[string]string, len(principalTypes))
	for _, typ := range principalTypes {
		principalTypeMap[strings.ToLower(typ)] = typ
	}

	for _, principal := range ep {
		pt := strings.ToLower(principal.Type)
		if principalType, ok := principalTypeMap[pt]; ok {
			if claim, ok := claims[principalType]; ok {
				claims[principalType] = append(claim.([]any), principal.Name)
			} else {
				claims[principalType] = []any{principal.Name}
			}
		}
	}

	return claims, nil
}

func (e *Handler) parseScope(
	scope string,
) (string, []string, error) {
	var domain string
	var types []string
	principals := strings.Split(scope, " ")
	for _, principal := range principals {
		d, typ := filepath.Split(principal)
		d = strings.TrimSuffix(d, "/")
		if d == "" {
			return "", nil, fmt.Errorf("scope %q is missing domain", principal)
		} else if typ == "" {
			return "", nil, fmt.Errorf("scope %q is missing principal type", principal)
		} else if domain != "" {
			if d != domain {
				return "", nil, fmt.Errorf(
					"scope %q does not match expected domain %q", principal, domain)
			}
		} else {
			domain = d
		}
		types = append(types, typ)
	}
	return domain, types, nil
}