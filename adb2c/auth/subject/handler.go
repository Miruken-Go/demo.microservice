package subject

import (
	"slices"
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
	h.db       = db
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
	id := create.SubjectId
	if id == "" {
		id = model.NewId()
	}
	subject := model.Subject{
		Id:        id,
		CreatedAt: time.Now().UTC(),
	}
	pk := azcosmos.NewPartitionKeyString(id)
	if _, err = azure.CreateItem(ctx, &subject, pk, h.subjects, nil); err == nil {
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
) miruken.HandleResult {
	sid := assign.SubjectId
	pk  := azcosmos.NewPartitionKeyString(sid)
	_, found, err := azure.ReplaceItem(ctx, func(subject *model.Subject) (bool, error) {
		var scope *model.Scope
		idx := slices.IndexFunc(subject.Scopes, func(s model.Scope) bool {
			return s.Name == assign.Scope
		})
		if idx < 0 {
			subject.Scopes = append(subject.Scopes, model.Scope{Name: assign.Scope})
			idx = len(subject.Scopes)-1
		}
		var changed bool
		scope = &subject.Scopes[idx]
		if newPrincipalIds, ok := model.Union(scope.PrincipalIds, assign.PrincipalIds...); ok {
			scope.PrincipalIds = newPrincipalIds
			changed = true
		}
		return changed, nil
	}, sid, pk, h.subjects, nil)

	switch {
	case !found:
		return miruken.NotHandled
	case err != nil:
		return miruken.NotHandled.WithError(err)
	default:
		return miruken.Handled
	}
}

func (h *Handler) Revoke(
	_ *struct {
		handles.It
		authorizes.Required
	  }, revoke api.RevokePrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) miruken.HandleResult {
	sid := revoke.SubjectId
	pk := azcosmos.NewPartitionKeyString(sid)
	_, found, err := azure.ReplaceItem(ctx, func(subject *model.Subject) (bool, error) {
		var changed bool
		idx := slices.IndexFunc(subject.Scopes, func(scope model.Scope) bool {
			return scope.Name == revoke.Scope
		})
		if idx >= 0 {
			if newPrincipalIds, ok := model.Difference(
				subject.Scopes[idx].PrincipalIds,
				revoke.PrincipalIds...,
			); ok {
				if len(newPrincipalIds) == 0 {
					subject.Scopes = slices.Delete(subject.Scopes, idx, idx+1)
				} else {
					subject.Scopes[idx].PrincipalIds = newPrincipalIds
				}
				changed = true
			}
		}
		return changed, nil
	}, sid, pk, h.subjects, nil)

	switch {
	case !found:
		return miruken.NotHandled
	case err != nil:
		return miruken.NotHandled.WithError(err)
	default:
		return miruken.Handled
	}
}

func (h *Handler) Remove(
	_ *struct {
		handles.It
		authorizes.Required
	  }, remove api.RemoveSubject,
	_ *struct{ args.Optional }, ctx context.Context,
) error {
	sid := remove.SubjectId
	pk  := azcosmos.NewPartitionKeyString(sid)
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
	sid := get.SubjectId
	pk  := azcosmos.NewPartitionKeyString(sid)
	item, found, err := azure.ReadItem[model.Subject](ctx, sid, pk, h.subjects, nil)
	switch {
	case !found:
		return api.Subject{}, miruken.NotHandled
	case err != nil:
		return api.Subject{}, miruken.NotHandled.WithError(err)
	default:
		return item.ToApi(), miruken.Handled
	}
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
	sql.WriteString("SELECT CROSS PARTITION s.id, s.scopes FROM s")
	sql.WriteString("")

	if filter := find.Principals; filter != nil {
		if principalIds := filter.Ids; len(principalIds) > 0 {
			params = []any{filter.Scope, principalIds}
			sql.WriteString(" JOIN p IN s.scopes WHERE p.name = :1")
			if all := find.Principals.All; all {
				sql.WriteString(" AND ARRAY_LENGTH(SetIntersect(p.principalIds, :2)) = :3")
				params = append(params, len(principalIds))
			} else {
				sql.WriteString(" AND ARRAY_LENGTH(SetIntersect(p.principalIds, :2)) != 0")
			}
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
				"SubjectId": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.AssignPrincipals](play.Constraints{
				"SubjectId":    "required",
				"Scope":        "required",
				"PrincipalIds": "gt=0",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.RevokePrincipals](play.Constraints{
				"SubjectId":    "required",
				"Scope":        "required",
				"PrincipalIds": "gt=0",
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

	_ = h.Validates6.WithRules(
		play.Rules{
			play.Type[struct{
				Scope  string
				Ids    []string
				All    bool
			}](play.Constraints{
				"Scope": "required",
				"Ids":   "gt=0",
			}),
		}, nil, translator)
}
