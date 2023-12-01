package handle

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	EntitlementHandler struct {
		database *mongo.Database
	}
)


func (h *EntitlementHandler) Constructor(
	client *mongo.Client,
) {
	h.database = client.Database("adb2c")
}

func (h *EntitlementHandler) Create(
	_*struct{
		handles.It
		authorizes.Required
	  }, create api.CreateEntitlement,
) (api.Entitlement, error) {
	var entitlement api.Entitlement
	return entitlement, nil
}

func (h *EntitlementHandler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	}, remove api.RemoveEntitlements,
) error {
	return nil
}

func (h *EntitlementHandler) Get(
	_ *handles.It, get api.GetEntitlement,
) (api.Entitlement, error) {
	return api.Entitlement{
	}, nil
}

func (h *EntitlementHandler) Find(
	_ *handles.It, find api.FindEntitlements,
) ([]api.Entitlement, error) {
	return []api.Entitlement{
	}, nil
}
