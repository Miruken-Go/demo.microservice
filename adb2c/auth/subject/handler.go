package subject

import (
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
	Handler struct {
		play.Validates1[api.CreateSubject]
		play.Validates2[api.AssignPrincipals]
		play.Validates3[api.RevokePrincipals]
		play.Validates4[api.RemoveSubject]
		play.Validates5[api.GetSubject]
		play.Validates6[api.FindSubjects]
		database *mongo.Database
	}

 	subjectResult struct {
		Subject           internal.Subject    `bson:"subject"`
		RelatedPrincipals []internal.Principal `bson:"related_principals"`
 	}
)


func (h *Handler) Constructor(
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

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.AssignPrincipals](map[string]string{
				"SubjectId":    "required",
				"PrincipalIds": "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.RevokePrincipals](map[string]string{
				"SubjectId":    "required",
				"PrincipalIds": "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates4.WithRules(
		play.Rules{
			play.Type[api.RemoveSubject](map[string]string{
				"SubjectId": "required",
			}),
		}, nil, translator)

	_ = h.Validates5.WithRules(
		play.Rules{
			play.Type[api.GetSubject](map[string]string{
				"SubjectId": "required",
			}),
		}, nil, translator)
}

func (h *Handler) Create(
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

func (h *Handler) Assign(
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

func (h *Handler) Revoke(
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

func (h *Handler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	  }, remove api.RemoveSubject,
	_*struct{args.Optional}, ctx context.Context,
) error {
	subjects := h.database.Collection("subject")
	_, err := subjects.DeleteOne(ctx, bson.M{"_id": remove.SubjectId})
	return err
}

func (h *Handler) Get(
	_*struct{
		handles.It
		authorizes.Required
	  }, get api.GetSubject,
	_*struct{args.Optional}, ctx context.Context,
) (api.Subject, miruken.HandleResult) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"subject_id": get.SubjectId,
			},
		},
		joinPrincipals, unwindPrincipals,
		{
			"$group": bson.M{
				"_id":                "$_id",
				"subject_id":         bson.M{"$first": "$subject_id"},
				"related_principals": bson.M{"$push": "$related_principals"},
			},
		},
		joinSubject, unwindSubject, projectSubject,
	}

	subjectPrincipals := h.database.Collection("subject_principal")
	cursor, err := subjectPrincipals.Aggregate(ctx, pipeline)
	if err != nil {
		return api.Subject{}, miruken.NotHandled.WithError(err)
	}

	var result subjectResult
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return api.Subject{}, miruken.NotHandled.WithError(err)
		}
		return result.mapSubject(), miruken.Handled
	}
	return api.Subject{}, miruken.NotHandled
}

func (h *Handler) Find(
	_*struct{
		handles.It
		authorizes.Required
	  }, find api.FindSubjects,
	_*struct{args.Optional}, ctx context.Context,
) ([]api.Subject, error) {
	principalIds := find.PrincipalIds

	var pipeline []bson.M
	if len(principalIds) > 0 {
		pipeline = append(pipeline,
			bson.M{
				"$match": bson.M{
					"principal_id": bson.M{"$in": principalIds},
				},
			},
		)
	}

	pipeline = append(pipeline,
		joinSubject, unwindSubject,
		joinPrincipals, unwindPrincipals,
		bson.M{
			"$group": bson.M{
				"_id":               "$subject._id",
				"subject":            bson.M{"$first": "$subject"},
				"related_principals": bson.M{"$push": "$related_principals"},
				"count":              bson.M{"$sum": 1},
			},
		},
		bson.M{
			"$match": bson.M{"count": bson.M{"$gte": len(principalIds)}},
		},
		projectSubject,
	)

	subjectPrincipals := h.database.Collection("subject_principal")
	cursor, err := subjectPrincipals.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var results []subjectResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	response := make([]api.Subject, len(results))
	for i, result := range results {
		response[i] = result.mapSubject()
	}
	return response, nil
}


func (s subjectResult) mapSubject() api.Subject {
	principals := make([]api.Principal, len(s.RelatedPrincipals))
	for i, principal := range s.RelatedPrincipals {
		tags := make([]api.Tag, len(principal.TagIDs))
		for j, tagId := range principal.TagIDs {
			tags[j] = api.Tag{
				Id: tagId,
			}
		}
		principals[i] = api.Principal{
			Id:   principal.ID,
			Name: principal.Name,
			Tags: tags,
		}
	}

	return api.Subject{
		Id:         s.Subject.ID,
		ObjectId:   s.Subject.ObjectID,
		Principals: principals,
	}
}


var (
	joinSubject = bson.M{
		"$lookup": bson.M{
			"from":         "subject",
			"localField":   "subject_id",
			"foreignField": "_id",
			"as":           "subject",
		},
	}

	unwindSubject = bson.M{
		"$unwind": "$subject",
	}

	joinPrincipals = bson.M{
		"$lookup": bson.M{
			"from":         "principal",
			"localField":   "principal_id",
			"foreignField": "_id",
			"as":           "related_principals",
		},
	}

	unwindPrincipals = bson.M{
		"$unwind": "$related_principals",
	}

	projectSubject = bson.M{
		"$project": bson.M{
			"_id":                0,
			"subject":            1,
			"related_principals": 1,
		},
	}
)