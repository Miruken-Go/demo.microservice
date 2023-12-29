package subject

import (
	"slices"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	ut "github.com/go-playground/universal-translator"
	"github.com/jmoiron/sqlx"
	api3 "github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/internal/model"
	"github.com/miruken-go/demo.microservice/adb2c/azure"
	"github.com/miruken-go/miruken"
	api2 "github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/promise"
	"github.com/miruken-go/miruken/security/authorizes"
	play "github.com/miruken-go/miruken/validates/play"
	"golang.org/x/net/context"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		play.Validates1[api3.CreateSubject]
		play.Validates2[api3.AssignPrincipals]
		play.Validates3[api3.RevokePrincipals]
		play.Validates4[api3.RemoveSubject]
		play.Validates5[api3.GetSubject]
		play.Validates6[api3.FindSubjects]

		subjects *azcosmos.ContainerClient
		db       *sqlx.DB
	}
)

const (
	database  = "adb2c"
	container = "subject"
)

func (h *Handler) Constructor(
	db     *sqlx.DB,
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
	  }, create api3.CreateSubject,
	_ *struct{ args.Optional }, ctx context.Context,
) (s api3.SubjectCreated, err error) {
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
	  }, assign api3.AssignPrincipals,
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
	  }, revoke api3.RevokePrincipals,
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
	  }, remove api3.RemoveSubject,
	_ *struct{ args.Optional }, ctx context.Context,
) miruken.HandleResult {
	sid := remove.SubjectId
	pk  := azcosmos.NewPartitionKeyString(sid)
	_, found, err := azure.DeleteItem(ctx, sid, pk, h.subjects, nil)
	switch {
	case !found:
		return miruken.NotHandled
	case err != nil:
		return miruken.NotHandled.WithError(err)
	default:
		return miruken.Handled
	}
}

func (h *Handler) Get(
	_ *struct {
		handles.It
		authorizes.Required
	  }, get api3.GetSubject,
	_ *struct{ args.Optional }, ctx context.Context,
) (api3.Subject, miruken.HandleResult) {
	sid := get.SubjectId
	pk  := azcosmos.NewPartitionKeyString(sid)
	_, item, found, err := azure.ReadItem[model.Subject](ctx, sid, pk, h.subjects, nil)
	switch {
	case !found:
		return api3.Subject{}, miruken.NotHandled
	case err != nil:
		return api3.Subject{}, miruken.NotHandled.WithError(err)
	default:
		return item.ToApi(), miruken.Handled
	}
}

func (h *Handler) Find(
	_ *struct {
		handles.It
		authorizes.Required
	  }, find api3.FindSubjects,
	_ *struct{ args.Optional }, ctx context.Context,
	hc miruken.HandleContext,
) *promise.Promise[[]api3.Subject] {
	return promise.New(nil, func(
		resolve func([]api3.Subject), reject func(error), onCancel func(func())) {

		var params []any
		var sql strings.Builder
		sql.WriteString("SELECT CROSS PARTITION s.id, s.scopes FROM s")
		sql.WriteString("")

		if filter := find.Filter; filter != nil {
			if principalIds := filter.PrincipalIds; len(principalIds) > 0 {
				if !filter.Exact {
					sp, spp, err := api2.Send[[]string](hc, api3.SatisfyPrincipals{
						Scope:        filter.Scope,
						PrincipalIds: principalIds,
					})
					if spp != nil {
						sp, err = spp.Await()
					}
					if err != nil {
						reject(err)
						return
					}
					principalIds = sp
				}
				params = []any{filter.Scope, principalIds}
				sql.WriteString(" JOIN p IN s.scopes WHERE p.name = :1")
				sql.WriteString(" AND ARRAY_LENGTH(SetIntersect(p.principalIds, :2)) != 0")
			}
		}

		sql.WriteString(" WITH database=")
		sql.WriteString(database)
		sql.WriteString(" WITH collection=")
		sql.WriteString(container)

		if ctx == nil {
			ctx = context.Background()
		}
		rows, err := h.db.QueryxContext(ctx, sql.String(), params...)
		if err != nil {
			reject(err)
			return
		}
		defer func() {
			_ = rows.Close()
		}()

		results := make([]api3.Subject, 0)
		for rows.Next() {
			row := make(model.SubjectMap)
			if err := rows.MapScan(row); err != nil {
				reject(err)
				return
			}
			results = append(results, row.ToApi())
		}

		resolve(results)
	})
}

func (h *Handler) setValidationRules(
	translator ut.Translator,
) {
	_ = h.Validates1.WithRules(
		play.Rules{
			play.Type[api3.CreateSubject](play.Constraints{
				"SubjectId": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api3.AssignPrincipals](play.Constraints{
				"SubjectId":    "required",
				"Scope":        "required",
				"PrincipalIds": "gt=0,dive,required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api3.RevokePrincipals](play.Constraints{
				"SubjectId":    "required",
				"Scope":        "required",
				"PrincipalIds": "gt=0,dive,required",
			}),
		}, nil, translator)

	_ = h.Validates4.WithRules(
		play.Rules{
			play.Type[api3.RemoveSubject](play.Constraints{
				"SubjectId": "required",
			}),
		}, nil, translator)

	_ = h.Validates5.WithRules(
		play.Rules{
			play.Type[api3.GetSubject](play.Constraints{
				"SubjectId": "required",
			}),
		}, nil, translator)

	_ = h.Validates6.WithRules(
		play.Rules{
			play.Type[struct{
				Scope        string
				PrincipalIds []string
				Exact        bool
			}](play.Constraints{
				"Scope":        "required",
				"PrincipalIds": "gt=0,dive,required",
			}),
		}, nil, translator)
}
