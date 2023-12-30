package user

import (
	"context"

	"github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
	}
)

func (h *Handler) Constructor(
) {
}

func (h *Handler) List(
	_ *struct {
		handles.It
		authorizes.Required
	  }, list api.ListUsers,
	_ *struct{ args.Optional }, ctx context.Context,
) ([]api.User, error) {
	return []api.User{}, nil
}


