package entitlement

//go:generate $GOPATH/bin/miruken -tests

import (
	"errors"
	"github.com/google/uuid"
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/internal/model"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type (
	Handler struct {
		database *mongo.Database
	}
)


func (h *Handler) Constructor(
	client *mongo.Client,
) {
	h.database = client.Database("adb2c")
}

func (h *Handler) Create(
	_*struct{
		handles.It
		authorizes.Required
	  }, create api.CreateEntitlement,
	_*struct{args.Optional}, ctx context.Context,
) (api.EntitlementCreated, error) {
	entitlement := model.EntitlementM{
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

func (h *Handler) Tag(
	_*struct{
		handles.It
		authorizes.Required
	}, tag api.TagEntitlement,
) error {
	return nil
}

func (h *Handler) Untag(
	_*struct{
		handles.It
		authorizes.Required
	}, untag api.UntagEntitlement,
) error {
	return nil
}

func (h *Handler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	}, remove api.RemoveEntitlement,
	_*struct{args.Optional}, ctx context.Context,
) error {
	entitlements := h.database.Collection("entitlement")
	_, err := entitlements.DeleteOne(ctx, bson.M{"_id": remove.EntitlementId})
	return err
}

func (h *Handler) Get(
	_ *handles.It, get api.GetEntitlement,
	_*struct{args.Optional}, ctx context.Context,
) (api.Entitlement, miruken.HandleResult) {
	var result model.EntitlementM
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

func (h *Handler) Find(
	_ *handles.It, find api.FindEntitlements,
) ([]api.Entitlement, error) {
	return []api.Entitlement{
	}, nil
}
