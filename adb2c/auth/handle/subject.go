package handle

import (
	"github.com/google/uuid"
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/internal"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	"golang.org/x/net/context"
	"time"
)

func (h *Handler) CreateSubject(
	_*struct {
		handles.It
		authorizes.Required
	  }, create api.CreateSubject,
	_*struct{args.Optional}, ctx context.Context,
) (api.Subject, error) {
	now := time.Now().UTC()
	subject := internal.Subject{
		ID:         uuid.New(),
		Name:       create.Name,
		CreatedAt:  now,
		ModifiedAt: now,
	}
	if _, err := h.database.Collection("subject").InsertOne(ctx, subject); err != nil {
		return api.Subject{}, err
	} else {
		return api.Subject{
			Id:   subject.ID,
			Name: subject.Name,
		}, nil
	}
}

func (h *Handler) AssignPrincipals(
	_*struct {
		handles.It
		authorizes.Required
	}, assign api.AssignPrincipals,
) error {
	return nil
}

func (h *Handler) RevokePrincipals(
	_*struct {
		handles.It
		authorizes.Required
	  }, revoke api.RevokePrincipals,
) error {
	return nil
}

func (h *Handler) RemoveSubjects(
	_*struct {
		handles.It
		authorizes.Required
	  }, remove api.RemoveSubjects,
) error {
	return nil
}

func (h *Handler) GetSubject(
	_ *handles.It, get api.GetSubject,
) (api.Subject, error) {
	return api.Subject{
	}, nil
}

func (h *Handler) FindSubjects(
	_ *handles.It, find api.FindSubjects,
) ([]api.Subject, error) {
	return []api.Subject{
	}, nil
}
