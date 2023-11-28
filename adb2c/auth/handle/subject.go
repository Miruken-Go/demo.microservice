package handle

import (
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
)

func (h *Handler) CreateSubject(
	_*struct {
		handles.It
		authorizes.Required
	  }, create api.CreateSubject,
) (api.Subject, error) {
	var subject api.Subject
	return subject, nil
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
