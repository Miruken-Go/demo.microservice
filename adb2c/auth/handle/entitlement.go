package handle

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
)

func (h *Handler) CreateEntitlement(
	_*struct {
		handles.It
		authorizes.Required
	  }, create api.CreateEntitlement,
) (api.Entitlement, error) {
	var entitlement api.Entitlement
	return entitlement, nil
}

func (h *Handler) RemoveEntitlements(
	_*struct {
		handles.It
		authorizes.Required
	}, remove api.RemoveEntitlements,
) error {
	return nil
}

func (h *Handler) GetEntitlement(
	_ *handles.It, get api.GetEntitlement,
) (api.Entitlement, error) {
	return api.Entitlement{
	}, nil
}

func (h *Handler) FindEntitlements(
	_ *handles.It, find api.FindEntitlements,
) ([]api.Entitlement, error) {
	return []api.Entitlement{
	}, nil
}
