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
	_*struct{args.Optional}, ctx context.Context,
) (api.EntitlementCreated, error) {
	entitlement := internal.Entitlement{
		ID:      uuid.New(),
		Name:    create.Name,
		TagIDs:  create.TagIds,
	}
	entitlements := h.database.Collection("entitlement")
	if _, err := entitlements.InsertOne(ctx, entitlement); err != nil {
		return api.EntitlementCreated{}, err
	}
	return api.EntitlementCreated{
		EntitlementId: entitlement.ID,
	}, nil
}

func (h *EntitlementHandler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	}, remove api.RemoveEntitlements,
	_*struct{args.Optional}, ctx context.Context,
) error {
	if entitlementIds := remove.EntitlementIds; len(entitlementIds) > 0 {
		entitlements := h.database.Collection("entitlement")
		filter := bson.M{"_id": bson.M{"$in": entitlementIds}}
		_, err := entitlements.DeleteMany(ctx, filter)
		return err
	}
	return nil
}

func (h *EntitlementHandler) Get(
	_ *handles.It, get api.GetEntitlement,
	_*struct{args.Optional}, ctx context.Context,
) (api.Entitlement, miruken.HandleResult) {
	var result internal.Entitlement
	filter := bson.M{"_id": get.EntitlementId}
	entitlements := h.database.Collection("entitlement")
	err := entitlements.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return api.Entitlement{}, miruken.NotHandled
	} else if err != nil {
		return api.Entitlement{}, miruken.NotHandled.WithError(err)
	}
	return api.Entitlement{
		Id:   result.ID,
		Name: result.Name,
	}, miruken.Handled
}

func (h *EntitlementHandler) Find(
	_ *handles.It, find api.FindEntitlements,
) ([]api.Entitlement, error) {
	return []api.Entitlement{
	}, nil
}
