package handle

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
)

func (h *Handler) CreateTag(
	_*struct {
		handles.It
		authorizes.Required
	}, create api.CreateTag,
) (api.Tag, error) {
	var tag api.Tag
	return tag, nil
}

func (h *Handler) RemoveTags(
	_*struct {
		handles.It
		authorizes.Required
	  }, remove api.RemoveTags,
) error {
	return nil
}

func (h *Handler) GetTag(
	_ *handles.It, get api.GetTag,
) (api.Tag, error) {
	return api.Tag{
	}, nil
}

func (h *Handler) FindTags(
	_ *handles.It, find api.FindTags,
) ([]api.Tag, error) {
	return []api.Tag{
	}, nil
}

