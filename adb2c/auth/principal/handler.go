package principal

//go:generate $GOPATH/bin/miruken -tests

import (
	"errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/google/uuid"
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/internal/model"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	play "github.com/miruken-go/miruken/validates/play"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type (
	Handler struct {
		play.Validates1[api.CreatePrincipal]
		play.Validates2[api.TagPrincipal]
		play.Validates3[api.UntagPrincipal]
		play.Validates4[api.AssignEntitlements]
		play.Validates5[api.RemoveEntitlement]
		play.Validates6[api.RemovePrincipal]
		play.Validates7[api.GetPrincipal]
		play.Validates8[api.FindPrincipals]
		database *mongo.Database
	}

	principalResult struct {
		Principal           model.PrincipalM     `bson:"principal"`
		RelatedEntitlements []model.EntitlementM `bson:"related_entitlements"`
	}
)


func (h *Handler) Constructor(
	client *mongo.Client,
	_*struct{args.Optional}, translator ut.Translator,
) {
	h.database = client.Database("adb2c")

	_ = h.Validates1.WithRules(
		play.Rules{
			play.Type[api.CreatePrincipal](map[string]string{
				"Name":           "required",
				"TagIds":         "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.TagPrincipal](map[string]string{
				"PrincipalId": "required",
				"TagIds":      "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.UntagPrincipal](map[string]string{
				"PrincipalId": "required",
				"TagIds":      "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates4.WithRules(
		play.Rules{
			play.Type[api.AssignEntitlements](map[string]string{
				"PrincipalId":    "required",
				"EntitlementIds": "required",
			}),
		}, nil, translator)

	_ = h.Validates5.WithRules(
		play.Rules{
			play.Type[api.RevokeEntitlements](map[string]string{
				"PrincipalId":    "required",
				"EntitlementIds": "required",
			}),
		}, nil, translator)

	_ = h.Validates6.WithRules(
		play.Rules{
			play.Type[api.RemovePrincipal](map[string]string{
				"PrincipalId": "required",
			}),
		}, nil, translator)

	_ = h.Validates7.WithRules(
		play.Rules{
			play.Type[api.GetPrincipal](map[string]string{
				"PrincipalId": "required",
			}),
		}, nil, translator)
}

func (h *Handler) Create(
	_*struct{
		handles.It
		authorizes.Required
	  }, create api.CreatePrincipal,
	_*struct{args.Optional}, ctx context.Context,
) (api.PrincipalCreated, error) {
	principal := model.PrincipalM{
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

func (h *Handler) Tag(
	_*struct{
		handles.It
		authorizes.Required
	}, tag api.TagPrincipal,
	_*struct{args.Optional}, ctx context.Context,
) error {
	principals := h.database.Collection("principal")
	_, err := principals.UpdateOne(
		ctx,
		bson.M{"_id": tag.PrincipalId},
		bson.M{"$addToSet": bson.M{"tags": bson.M{"$each": tag.TagIds}}},
	)
	return err
}

func (h *Handler) Untag(
	_*struct{
		handles.It
		authorizes.Required
	  }, untag api.UntagPrincipal,
	_*struct{args.Optional}, ctx context.Context,
) error {
	principals := h.database.Collection("principal")
	_, err := principals.UpdateOne(
		ctx,
		bson.M{"_id": untag.PrincipalId},
		bson.M{"$pull": bson.M{"tags": bson.M{"$in": untag.TagIds}}},
	)
	return err
}

func (h *Handler) Assign(
	_*struct{
		handles.It
		authorizes.Required
	}, assign api.AssignEntitlements,
) error {
	return nil
}

func (h *Handler) Revoke(
	_*struct{
		handles.It
		authorizes.Required
	  }, revoke api.RevokeEntitlements,
) error {
	return nil
}

func (h *Handler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	}, remove api.RemovePrincipal,
	_*struct{args.Optional}, ctx context.Context,
) error {
	principals := h.database.Collection("principal")
	_, err := principals.DeleteOne(ctx, bson.M{"_id": remove.PrincipalId})
	return err
}

func (h *Handler) Get(
	_ *handles.It, get api.GetPrincipal,
	_*struct{args.Optional}, ctx context.Context,
) (api.Principal, miruken.HandleResult) {
	var result model.PrincipalM
	filter := bson.M{"_id": get.PrincipalId}
	principals := h.database.Collection("principal")
	err := principals.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return api.Principal{}, miruken.NotHandled
	} else if err != nil {
		return api.Principal{}, miruken.NotHandled.WithError(err)
	}
	tags := make([]api.Tag, len(result.TagIDs))
	for i, tagId := range result.TagIDs {
		tags[i] = api.Tag{Id: tagId}
	}
	return api.Principal{
		Id:   result.ID,
		Name: result.Name,
		Tags: tags,
	}, miruken.Handled
}

func (h *Handler) Find(
	_ *handles.It, find api.FindPrincipals,
	_*struct{args.Optional}, ctx context.Context,
) ([]api.Principal, error) {
	var pipeline []bson.M
	if name := find.Name; name != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"name": bson.M{"$regex": name, "$options": "i"},
			},
		})
	}

	pipeline = append(pipeline,
		bson.M{
			"$lookup": bson.M{
				"from":         "principal_entitlement",
				"localField":   "_id",
				"foreignField": "principal_id",
				"as":           "entitlements",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "entitlement",
				"localField":   "entitlements.entitlement_id",
				"foreignField": "_id",
				"as":           "entitlement_details",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$entitlement_details",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":      "$_id",
				"name":     bson.M{"$first": "$name"},
				"tags":     bson.M{"$first": "$tags"},
				"related_entitlements": bson.M{
					"$addToSet": "$entitlement_details._id",
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":                  1,
				"name":                 1,
				"tags":                 1,
				"related_entitlements": 1,
			},
		},
	)

	principals := h.database.Collection("principal")
	cursor, err := principals.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var results []principalResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	response := make([]api.Principal, len(results))
	for i, result := range results {
		response[i] = result.mapPrincipal()
	}
	return response, nil
}


func (p principalResult) mapPrincipal() api.Principal {
	entitlements := make([]api.Entitlement, len(p.RelatedEntitlements))
	for i, entitlement := range p.RelatedEntitlements {
		tags := make([]api.Tag, len(entitlement.TagIDs))
		for j, tagId := range entitlement.TagIDs {
			tags[j] = api.Tag{
				Id: tagId,
			}
		}
		entitlements[i] = api.Entitlement{
			Id:   entitlement.ID,
			Name: entitlement.Name,
			Tags: tags,
		}
	}

	tags := make([]api.Tag, len(p.Principal.TagIDs))
	for j, tagId := range p.Principal.TagIDs {
		tags[j] = api.Tag{
			Id: tagId,
		}
	}

	return api.Principal{
		Id:           p.Principal.ID,
		Name:         p.Principal.Name,
		Tags:         tags,
		Entitlements: entitlements,
	}
}