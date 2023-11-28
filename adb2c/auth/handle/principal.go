package handle

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
)

func (h *Handler) CreatePrincipal(
	_*struct {
		handles.It
		authorizes.Required
	  }, create api.CreatePrincipal,
) (api.Principal, error) {
	var principal api.Principal
	return principal, nil
}

func (h *Handler) AssignEntitlements(
	_*struct {
		handles.It
		authorizes.Required
	}, assign api.AssignEntitlements,
) error {
	return nil
}

func (h *Handler) RevokeEntitlements(
	_*struct {
		handles.It
		authorizes.Required
	  }, revoke api.RevokeEntitlements,
) error {
	return nil
}

func (h *Handler) RemovePrincipals(
	_*struct {
		handles.It
		authorizes.Required
	}, remove api.RemovePrincipals,
) error {
	return nil
}

func (h *Handler) GetPrincipal(
	_ *handles.It, get api.GetPrincipal,
) (api.Principal, error) {
	return api.Principal{
	}, nil
}

func (h *Handler) FindPrincipals(
	_ *handles.It, find api.FindPrincipals,
) ([]api.Principal, error) {
	return []api.Principal{
	}, nil
}
