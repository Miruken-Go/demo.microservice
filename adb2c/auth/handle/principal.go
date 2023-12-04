package handle

import (
	"errors"
	"github.com/google/uuid"
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/internal"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
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
	_*struct{args.Optional}, ctx context.Context,
) (api.PrincipalCreated, error) {
	principal := internal.Principal{
		ID:     uuid.New(),
		Name:   create.Name,
		TagIDs: create.TagIds,
	}
	principals := h.database.Collection("principal")
	if _, err := principals.InsertOne(ctx, principal); err != nil {
		return api.PrincipalCreated{}, err
	}
	return api.PrincipalCreated{
		PrincipalId: principal.ID,
	}, nil
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
	_*struct{args.Optional}, ctx context.Context,
) error {
	if principalIds := remove.PrincipalIds; len(principalIds) > 0 {
		principals := h.database.Collection("principal")
		filter := bson.M{"_id": bson.M{"$in": principalIds}}
		_, err := principals.DeleteMany(ctx, filter)
		return err
	}
	return nil
}

func (h *PrincipalHandler) Get(
	_ *handles.It, get api.GetPrincipal,
	_*struct{args.Optional}, ctx context.Context,
) (api.Principal, miruken.HandleResult) {
	var result internal.Principal
	filter := bson.M{"_id": get.PrincipalId}
	principals := h.database.Collection("principal")
	err := principals.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return api.Principal{}, miruken.NotHandled
	} else if err != nil {
		return api.Principal{}, miruken.NotHandled.WithError(err)
	}
	return api.Principal{
		Id:   result.ID,
		Name: result.Name,
	}, miruken.Handled
}

func (h *PrincipalHandler) Find(
	_ *handles.It, find api.FindPrincipals,
) ([]api.Principal, error) {
	return []api.Principal{
	}, nil
}
