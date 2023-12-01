package handle

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	TagHandler struct {
		database *mongo.Database
	}
)


func (h *TagHandler) Constructor(
	client *mongo.Client,
) {
	h.database = client.Database("adb2c")
}

func (h *TagHandler) Create(
	_*struct{
		handles.It
		authorizes.Required
	}, create api.CreateTag,
) (api.Tag, error) {
	var tag api.Tag
	return tag, nil
}

func (h *TagHandler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	  }, remove api.RemoveTags,
) error {
	return nil
}

func (h *TagHandler) Get(
	_ *handles.It, get api.GetTag,
) (api.Tag, error) {
	return api.Tag{
	}, nil
}

func (h *TagHandler) Find(
	_ *handles.It, find api.FindTags,
) ([]api.Tag, error) {
	return []api.Tag{
	}, nil
}

