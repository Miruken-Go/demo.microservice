package token

import (
	"encoding/json"
	"github.com/go-logr/logr"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"net/http"
)

type (
	// EnrichHandler is an Azure ADB2C Api Connector that provides
	// an external source of claims to enrich tokens during user flows.
	// https://learn.microsoft.com/en-us/azure/active-directory-b2c/add-api-connector-token-enrichment?pivots=b2c-custom-policy
	EnrichHandler struct {
		logger logr.Logger
	}

	// EnrichRequest is the request body sent to the Api Connector.
	// It receives a set et of InputClaims and returns a set of OutputClaims.
	// This enrichment support groups, roles and entitlements.
	EnrichRequest struct {
		Email    string
		ObjectId string
		Scope    string
	}

	// EnrichResponse is the response body returned from the Api Connector.
	// It contains the enriched groups, roles and entitlement OutputClaims.
	EnrichResponse struct {
		Groups       []string
		Roles        []string
		Entitlements []string
	}
)

func (e EnrichHandler) Constructor(
	_ *struct{ args.Optional }, logger logr.Logger,
) {
	if logger == e.logger {
		e.logger = logr.Discard()
	} else {
		e.logger = logger
	}
}

func (e EnrichHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
	h miruken.Handler,
) {
	var request EnrichRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	e.logger.Info("Enrich token",
		"Email", request.Email,
		"ObjectId", request.ObjectId,
		"Scope", request.Scope)

	response := EnrichResponse{
		Groups:       []string{"oncall"},
		Roles:        []string{"admin", "coach", "player"},
		Entitlements: []string{"createTeam", "updateTeam", "createPerson", "updatePerson"},
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
