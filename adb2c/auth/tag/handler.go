package tag

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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type (
	Handler struct {
		play.Validates1[api.CreateTag]
		play.Validates2[api.RemoveTag]
		play.Validates3[api.GetTag]
		play.Validates4[api.FindTags]
		database *mongo.Database
	}
)


func (h *Handler) Constructor(
	client *mongo.Client,
	_*struct{args.Optional}, translator ut.Translator,
) {
	h.database = client.Database("adb2c")

	_ = h.Validates1.WithRules(
		play.Rules{
			play.Type[api.CreateTag](map[string]string{
				"Name": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.RemoveTag](map[string]string{
				"TagId": "required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.GetTag](map[string]string{
				"TagId": "required",
			}),
		}, nil, translator)
}

func (h *Handler) Create(
	_*struct{
		handles.It
		authorizes.Required
	}, create api.CreateTag,
	_*struct{args.Optional}, ctx context.Context,
) (api.TagCreated, error) {
	tag := model.Tag{
		ID:   uuid.New(),
		Name: create.Name,
	}
	tags := h.database.Collection("tag")
	if _, err := tags.InsertOne(ctx, tag); err != nil {
		return api.TagCreated{}, err
	}
	return api.TagCreated{
		TagId: tag.ID,
	}, nil
}

func (h *Handler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	  }, remove api.RemoveTag,
	_*struct{args.Optional}, ctx context.Context,
) error {
	tags := h.database.Collection("tag")
	_, err := tags.DeleteOne(ctx, bson.M{"_id": remove.TagId})
	return err
}

func (h *Handler) Get(
	_ *handles.It, get api.GetTag,
	_*struct{args.Optional}, ctx context.Context,
) (api.Tag, miruken.HandleResult) {
	var result model.Tag
	filter := bson.M{"_id": get.TagId}
	tags   := h.database.Collection("tag")
	err    := tags.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return api.Tag{}, miruken.NotHandled
	} else if err != nil {
		return api.Tag{}, miruken.NotHandled.WithError(err)
	}
	return api.Tag{
		Id:   result.ID,
		Name: result.Name,
	}, miruken.Handled
}

func (h *Handler) Find(
	_ *handles.It, find api.FindTags,
	_*struct{args.Optional}, ctx context.Context,
) ([]api.Tag, error) {
	var filter bson.M
	if name := find.Name; name != "" {
		regex := bson.M{"$regex": primitive.Regex{Pattern: find.Name, Options: "i"}}
		filter = bson.M{"name": regex}
	} else {
		filter = bson.M{}
	}

	tags := h.database.Collection("tag")
	cursor, err := tags.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var results []model.Tag
	err = cursor.All(ctx, &results)

	tagResults := make([]api.Tag, len(results), len(results))
	for i, result := range results {
		tagResults[i] = api.Tag{
			Id:   result.ID,
			Name: result.Name,
		}
	}
	return tagResults, err
}
