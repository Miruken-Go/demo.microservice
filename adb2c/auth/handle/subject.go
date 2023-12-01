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
}

func (h *SubjectHandler) Create(
	_*struct{
		handles.It
		authorizes.Required
	  }, create api.CreateSubject,
	_*struct{args.Optional}, ctx context.Context,
) (api.Subject, error) {
	now := time.Now().UTC()
	subject := internal.Subject{
		ID:         uuid.New(),
		ObjectID:   create.ObjectId,
		CreatedAt:  now,
	}
	_, err := h.database.Collection("subject").InsertOne(ctx, subject)
	if err != nil {
		return api.Subject{}, err
	}
	return api.Subject{
		Id: subject.ID,
	}, nil
}

func (h *SubjectHandler) Assign(
	_*struct{
		handles.It
		authorizes.Required
	}, assign api.AssignPrincipals,
) error {
	return nil
}

func (h *SubjectHandler) Revoke(
	_*struct{
		handles.It
		authorizes.Required
	  }, revoke api.RevokePrincipals,
) error {
	return nil
}

func (h *SubjectHandler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	  }, remove api.RemoveSubjects,
) error {
	return nil
}

func (h *SubjectHandler) Get(
	_*struct{
		handles.It
		authorizes.Required
	  }, get api.GetSubject,
	_*struct{args.Optional}, ctx context.Context,
) (api.Subject, miruken.HandleResult) {
	filter := bson.M{"_id": get.SubjectId}
	var result internal.Subject
	err := h.database.Collection("subject").FindOne(ctx, filter).Decode(&result)
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
) ([]api.Subject, error) {
	return []api.Subject{
	}, nil
}
