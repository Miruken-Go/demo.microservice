package subject

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	ut "github.com/go-playground/universal-translator"
	"github.com/jmoiron/sqlx"
	"github.com/miruken-go/demo.microservice/adb2c/auth/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/internal/model"
	"github.com/miruken-go/demo.microservice/adb2c/azure"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	play "github.com/miruken-go/miruken/validates/play"
	"golang.org/x/net/context"
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

		subjects *azcosmos.ContainerClient
		db       *sqlx.DB
	}
)

const (
	database  = "adb2c"
	container = "subject"
)

func (h *Handler) Constructor(
	db *sqlx.DB,
	client *azcosmos.Client,
	_ *struct{ args.Optional }, translator ut.Translator,
) {
	h.db = db
	h.subjects = azure.Container(client, database, container)
	h.setValidationRules(translator)
}

func (h *Handler) Create(
	_ *struct {
		handles.It
		authorizes.Required
	}, create api.CreateSubject,
	_ *struct{ args.Optional }, ctx context.Context,
) (s api.SubjectCreated, err error) {
	id := model.NewId()
	subject := model.Subject{
		Id:           id.String(),
		ObjectId:     create.ObjectId,
		PrincipalIds: model.Strings(create.PrincipalIds),
		CreatedAt:    time.Now().UTC(),
	}
	pk := azcosmos.NewPartitionKeyString(subject.Id)
	_, err = azure.CreateItem(&subject, ctx, pk, h.subjects, nil)
	if err == nil {
		s.SubjectId = id
	}
	return
}

func (h *Handler) Assign(
	_ *struct {
		handles.It
		authorizes.Required
	}, assign api.AssignPrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) error {
	sid := assign.SubjectId.String()
	pk := azcosmos.NewPartitionKeyString(sid)
	_, _, err := azure.ReplaceItem(func(subject *model.Subject) (bool, error) {
		add := model.Strings(assign.PrincipalIds)
		updated, changed := model.Union(subject.PrincipalIds, add...)
		if changed {
			subject.PrincipalIds = updated
		}
		return changed, nil
	}, ctx, sid, pk, h.subjects, nil)
	return err
}

func (h *Handler) Revoke(
	_ *struct {
		handles.It
		authorizes.Required
	}, revoke api.RevokePrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) error {
	sid := revoke.SubjectId.String()
	pk := azcosmos.NewPartitionKeyString(sid)
	_, _, err := azure.ReplaceItem(func(subject *model.Subject) (bool, error) {
		remove := model.Strings(revoke.PrincipalIds)
		updated, changed := model.Difference(subject.PrincipalIds, remove...)
		if changed {
			subject.PrincipalIds = updated
		}
		return changed, nil
	}, ctx, sid, pk, h.subjects, nil)
	return err
}

func (h *Handler) Remove(
	_ *struct {
		handles.It
		authorizes.Required
	}, remove api.RemoveSubject,
	_ *struct{ args.Optional }, ctx context.Context,
) error {
	sid := remove.SubjectId.String()
	pk := azcosmos.NewPartitionKeyString(sid)
	_, err := h.subjects.DeleteItem(ctx, pk, sid, nil)
	return err
}

func (h *Handler) Get(
	_ *struct {
		handles.It
		authorizes.Required
	}, get api.GetSubject,
	_ *struct{ args.Optional }, ctx context.Context,
) (api.Subject, miruken.HandleResult) {
	sid := get.SubjectId.String()
	pk := azcosmos.NewPartitionKeyString(sid)
	item, found, err := azure.ReadItem[model.Subject](ctx, sid, pk, h.subjects, nil)
	if !found {
		return api.Subject{}, miruken.NotHandled
	} else if err != nil {
		return api.Subject{}, miruken.NotHandled.WithError(err)
	}
	return item.ToApi(), miruken.Handled
}

func (h *Handler) Find(
	_ *struct {
		handles.It
		authorizes.Required
	}, find api.FindSubjects,
	_ *struct{ args.Optional }, ctx context.Context,
) ([]api.Subject, error) {
	var params []any
	var sql strings.Builder
	sql.WriteString("SELECT CROSS PARTITION * FROM s")
	sql.WriteString("")

	cond := " WHERE"

	if objectId := find.ObjectId; objectId != "" {
		sql.WriteString(" WHERE s.objectId = :1")
		params = append(params, objectId)
		cond = " AND"
	}

	if principalIds := find.Principals.Ids; len(principalIds) > 0 {
		sql.WriteString(cond)
		if all := find.Principals.All; all {
			sql.WriteString(fmt.Sprintf(
				"%s ARRAY_LENGTH(SetIntersect(s.principalIds, :%d)) = :%d",
				cond, len(params)+1, len(params)+2))
			params = append(params, principalIds, len(principalIds))
		} else {
			sql.WriteString(fmt.Sprintf(
				"%s ARRAY_LENGTH(SetIntersect(s.principalIds, :%d)) != 0",
				cond, len(params)+1))
			params = append(params, principalIds)
		}
	}

	sql.WriteString(" WITH database=")
	sql.WriteString(database)
	sql.WriteString(" WITH collection=")
	sql.WriteString(container)

	rows, err := h.db.QueryxContext(ctx, sql.String(), params...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	results := make([]api.Subject, 0)
	for rows.Next() {
		row := make(model.SubjectMap)
		if err := rows.MapScan(row); err != nil {
			return nil, err
		}
		results = append(results, row.ToApi())
	}

	return results, nil
}

func (h *Handler) setValidationRules(
	translator ut.Translator,
) {
	_ = h.Validates1.WithRules(
		play.Rules{
			play.Type[api.CreateSubject](play.Constraints{
				"ObjectId": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.AssignPrincipals](play.Constraints{
				"SubjectId":    "required",
				"PrincipalIds": "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.RevokePrincipals](play.Constraints{
				"SubjectId":    "required",
				"PrincipalIds": "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates4.WithRules(
		play.Rules{
			play.Type[api.RemoveSubject](play.Constraints{
				"SubjectId": "required",
			}),
		}, nil, translator)

	_ = h.Validates5.WithRules(
		play.Rules{
			play.Type[api.GetSubject](play.Constraints{
				"SubjectId": "required",
			}),
		}, nil, translator)
}
