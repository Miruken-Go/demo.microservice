package subject

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	ut "github.com/go-playground/universal-translator"
	"github.com/google/uuid"
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/internal/model"
	"github.com/miruken-go/demo.microservice/adb2c/azure"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	play "github.com/miruken-go/miruken/validates/play"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"time"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		play.Validates1[api.CreateSubject]
		play.Validates2[api.AssignPrincipals]
		play.Validates3[api.RevokePrincipals]
		play.Validates4[api.RemoveSubject]
		play.Validates5[api.GetSubject]
		play.Validates6[api.FindSubjects]

		subjects     *azcosmos.ContainerClient
		principals   *azcosmos.ContainerClient
		entitlements *azcosmos.ContainerClient
	}
)


func (h *Handler) Constructor(
	client *azcosmos.Client,
	_*struct{args.Optional}, translator ut.Translator,
) {
	h.setupValidationRules(translator)
	if err := h.createContainers(client); err != nil {
		panic(err)
	}
}

func (h *Handler) Create(
	_*struct{
		handles.It
		authorizes.Required
	  }, create api.CreateSubject,
	_*struct{args.Optional}, ctx context.Context,
) (api.SubjectCreated, error) {
	id  := model.NewId()
	subject := model.Subject{
		Id:           id.String(),
		ObjectId:     create.ObjectId,
		PrincipalIds: model.Strings(create.PrincipalIds),
		CreatedAt:    time.Now().UTC(),
	}
	pk := azcosmos.NewPartitionKeyString(subject.Id)
	_, err := azure.CreateItem(&subject, ctx, pk, h.subjects, nil)
	if err != nil {
		return api.SubjectCreated{}, err
	}
	return api.SubjectCreated{
		SubjectId: id,
	}, nil
}

func (h *Handler) Assign(
	_*struct{
		handles.It
		authorizes.Required
	}, assign api.AssignPrincipals,
	_*struct{args.Optional}, ctx context.Context,
) error {
	subjectId := assign.SubjectId.String()
	pk := azcosmos.NewPartitionKeyString(subjectId)
	_, _, err := azure.ReplaceItem(func(subject *model.Subject) (bool, error) {
		add := model.Strings(assign.PrincipalIds)
		updated, changed := model.Union(subject.PrincipalIds, add...)
		if changed {
			subject.PrincipalIds = updated
		}
		return changed, nil
	}, ctx, subjectId, pk, h.subjects, nil)
	return err
}

func (h *Handler) Revoke(
	_*struct{
		handles.It
		authorizes.Required
	  }, revoke api.RevokePrincipals,
	_*struct{args.Optional}, ctx context.Context,
) error {
	subjectId := revoke.SubjectId.String()
	pk := azcosmos.NewPartitionKeyString(subjectId)
	_, _, err := azure.ReplaceItem(func(subject *model.Subject) (bool, error) {
		remove := model.Strings(revoke.PrincipalIds)
		updated, changed := model.Difference(subject.PrincipalIds, remove...)
		if changed {
			subject.PrincipalIds = updated
		}
		return changed, nil
	}, ctx, subjectId, pk, h.subjects, nil)
	return err
}

func (h *Handler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	  }, remove api.RemoveSubject,
	_*struct{args.Optional}, ctx context.Context,
) error {
	subjectId := remove.SubjectId.String()
	pk := azcosmos.NewPartitionKeyString(subjectId)
	_, err := h.subjects.DeleteItem(ctx, pk, subjectId, nil)
	return err
}

func (h *Handler) Get(
	_*struct{
		handles.It
		authorizes.Required
	  }, get api.GetSubject,
	_*struct{args.Optional}, ctx context.Context,
) (api.Subject, miruken.HandleResult) {
	subjectId := get.SubjectId.String()
	pk := azcosmos.NewPartitionKeyString(subjectId)
	item, found, err := azure.ReadItem[model.Subject](ctx, subjectId, pk, h.subjects, nil)
	if !found {
		return api.Subject{}, miruken.NotHandled
	} else if err != nil {
		return api.Subject{}, miruken.NotHandled.WithError(err)
	} else {
		id, err := uuid.Parse(item.Id)
		if err != nil {
			return api.Subject{}, miruken.NotHandled.WithError(err)
		} else {
			return api.Subject{
				Id:       id,
				ObjectId: item.ObjectId,
			}, miruken.Handled
		}
	}
}

func (h *Handler) Find(
	_*struct{
		handles.It
		authorizes.Required
	  }, find api.FindSubjects,
	_*struct{args.Optional}, ctx context.Context,
) ([]api.Subject, error) {
	return nil, nil
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

func (h *Handler) setupValidationRules(
	translator ut.Translator,
) {
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

func (h *Handler) createContainers(
	azure *azcosmos.Client,
) error {
	database, err := azure.NewDatabase("adb2c")
	if err != nil {
		return fmt.Errorf("error creating \"adb2c\" database client: %w", err)
	}
	if h.subjects, err = database.NewContainer("subject"); err != nil {
		return fmt.Errorf("error creating \"subject\" container client: %w", err)
	}
	if h.principals, err = database.NewContainer("principal"); err != nil {
		return fmt.Errorf("error creating \"principal\" container client: %w", err)
	}
	if h.entitlements, err = database.NewContainer("entitlement"); err != nil {
		return fmt.Errorf("error creating \"entitlement\" container client: %w", err)
	}
	return nil
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
