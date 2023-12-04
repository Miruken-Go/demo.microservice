package handle

import (
	"errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/google/uuid"
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/internal"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	play "github.com/miruken-go/miruken/validates/play"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"time"
)

import "go.mongodb.org/mongo-driver/mongo"

//go:generate $GOPATH/bin/miruken -tests

type (
	SubjectHandler struct {
		play.Validates1[api.CreateSubject]
		play.Validates2[api.AssignPrincipals]
		play.Validates3[api.RevokePrincipals]
		play.Validates4[api.RemoveSubjects]
		play.Validates5[api.GetSubject]
		play.Validates6[api.FindSubjects]
		database *mongo.Database
	}
)


func (h *SubjectHandler) Constructor(
	client *mongo.Client,
	_*struct{args.Optional}, translator ut.Translator,
) {
	h.database = client.Database("adb2c")

	_ = h.Validates1.WithRules(
		play.Rules{
			play.Type[api.CreateSubject](map[string]string{
				"ObjectId": "required",
			}),
		}, nil, translator)

	_ = h.Validates5.WithRules(
		play.Rules{
			play.Type[api.GetSubject](map[string]string{
				"SubjectId": "required",
			}),
		}, nil, translator)

	_ = h.Validates6.WithRules(
		play.Rules{
			play.Type[api.FindSubjects](map[string]string{
				"PrincipalIds": "gt=0,required",
			}),
		}, nil, translator)
}

func (h *SubjectHandler) Create(
	_*struct{
		handles.It
		authorizes.Required
	  }, create api.CreateSubject,
	_*struct{args.Optional}, ctx context.Context,
) (api.SubjectCreated, error) {
	now := time.Now().UTC()
	subject := internal.Subject{
		ID:         uuid.New(),
		ObjectID:   create.ObjectId,
		CreatedAt:  now,
	}
	subjects := h.database.Collection("subject")
	if _, err := subjects.InsertOne(ctx, subject); err != nil {
		return api.SubjectCreated{}, err
	}
	return api.SubjectCreated{
		SubjectId: subject.ID,
	}, nil
}

func (h *SubjectHandler) Assign(
	_*struct{
		handles.It
		authorizes.Required
	}, assign api.AssignPrincipals,
	_*struct{args.Optional}, ctx context.Context,
) error {
	subjectId  := assign.SubjectId
	operations := make([]mongo.WriteModel, len(assign.PrincipalIds))
	for i, principalId := range assign.PrincipalIds {
		subPrincipal := internal.SubjectPrincipal{
			SubjectID:   subjectId,
			PrincipalID: principalId,
		}
		filter := bson.M{"subject_id": subjectId, "principal_id": principalId}
		update := bson.M{"$setOnInsert": subPrincipal}
		operations[i] = mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
	}
	subjectPrincipals := h.database.Collection("subject_principal")
	_, err := subjectPrincipals.BulkWrite(ctx, operations)
	return err
}

func (h *SubjectHandler) Revoke(
	_*struct{
		handles.It
		authorizes.Required
	  }, revoke api.RevokePrincipals,
	_*struct{args.Optional}, ctx context.Context,
) error {
	filter := bson.M{
		"subject_id":   revoke.SubjectId,
		"principal_id": bson.M{"$in": revoke.PrincipalIds},
	}
	subjectPrincipals := h.database.Collection("subject_principal")
	_, err := subjectPrincipals.DeleteMany(ctx, filter)
	return err
}

func (h *SubjectHandler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	  }, remove api.RemoveSubjects,
	_*struct{args.Optional}, ctx context.Context,
) error {
	if subjectIds := remove.SubjectIds; len(subjectIds) > 0 {
		subjects := h.database.Collection("subject")
		filter := bson.M{"_id": bson.M{"$in": subjectIds}}
		_, err := subjects.DeleteMany(ctx, filter)
		return err
	}
	return nil
}

func (h *SubjectHandler) Get(
	_*struct{
		handles.It
		authorizes.Required
	  }, get api.GetSubject,
	_*struct{args.Optional}, ctx context.Context,
) (api.Subject, miruken.HandleResult) {
	var result internal.Subject
	filter := bson.M{"_id": get.SubjectId}
	subjects := h.database.Collection("subject")
	err := subjects.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return api.Subject{}, miruken.NotHandled
	} else if err != nil {
		return api.Subject{}, miruken.NotHandled.WithError(err)
	}
	return api.Subject{
		Id:       result.ID,
		ObjectId: result.ObjectID,
	}, miruken.Handled
}

func (h *SubjectHandler) Find(
	_*struct{
		handles.It
		authorizes.Required
	  }, find api.FindSubjects,
	_*struct{args.Optional}, ctx context.Context,
) ([]api.Subject, error) {
	principalIds := find.PrincipalIds
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"principal_id": bson.M{"$in": principalIds},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "subject",
				"localField":   "subject_id",
				"foreignField": "_id",
				"as":           "subject",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "principal",
				"localField":   "principal_id",
				"foreignField": "_id",
				"as":           "matched_principals",
			},
		},
		{
			"$unwind": "$subject",
		},
		{
			"$unwind": "$matched_principals",
		},
		{
			"$group": bson.M{
				"_id":   "$subject._id",
				"subject": bson.M{"$first": "$subject"},
				"matched_principals": bson.M{"$push": "$matched_principals"},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$match": bson.M{
				"count": len(principalIds),
			},
		},
		{
			"$project": bson.M{
				"_id":                0,
				"subject":            1,
				"matched_principals": 1,
			},
		},
	}

	subjectPrincipals := h.database.Collection("subject_principal")
	cursor, err := subjectPrincipals.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	type Result struct {
		Subject           internal.Subject     `bson:"subject"`
		MatchedPrincipals []internal.Principal `bson:"matched_principals"`
	}

	var results []Result
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	response := make([]api.Subject, len(results))
	for i, result := range results {
		principals := make([]api.Principal, len(result.MatchedPrincipals))
		for j, principal := range result.MatchedPrincipals {
			principals[j] = api.Principal{
				Id:   principal.ID,
				Name: principal.Name,
			}
		}
		response[i] = api.Subject{
			Id:        result.Subject.ID,
			ObjectId:  result.Subject.ObjectID,
			Principals: principals,
		}
	}
	return response, nil
}
