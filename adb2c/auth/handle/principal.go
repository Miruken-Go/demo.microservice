package handle

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	PrincipalHandler struct {
		database *mongo.Database
	}
)


func (h *PrincipalHandler) Constructor(
	client *mongo.Client,
) {
	h.database = client.Database("adb2c")
}

func (h *PrincipalHandler) Create(
	_*struct{
		handles.It
		authorizes.Required
	  }, create api.CreatePrincipal,
) (api.Principal, error) {
	var principal api.Principal
	return principal, nil
}

func (h *PrincipalHandler) Assign(
	_*struct{
		handles.It
		authorizes.Required
	}, assign api.AssignEntitlements,
) error {
	return nil
}

func (h *PrincipalHandler) Revoke(
	_*struct{
		handles.It
		authorizes.Required
	  }, revoke api.RevokeEntitlements,
) error {
	return nil
}

func (h *PrincipalHandler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	}, remove api.RemovePrincipals,
) error {
	return nil
}

func (h *PrincipalHandler) Get(
	_ *handles.It, get api.GetPrincipal,
) (api.Principal, error) {
	return api.Principal{
	}, nil
}

func (h *PrincipalHandler) Find(
	_ *handles.It, find api.FindPrincipals,
) ([]api.Principal, error) {
	return []api.Principal{
	}, nil
}
