package principal

//go:generate $GOPATH/bin/miruken -tests

import (
	"encoding/json"
	"fmt"
	"github.com/miruken-go/miruken/promise"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	ut "github.com/go-playground/universal-translator"
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

type (
	Handler struct {
		play.Validates1[api.CreatePrincipal]
		play.Validates2[api.IncludePrincipals]
		play.Validates3[api.ExcludePrincipals]
		play.Validates4[api.RemovePrincipal]
		play.Validates5[api.GetPrincipal]
		play.Validates6[api.FindPrincipals]
		play.Validates7[api.FlattenPrincipals]

		principals *azcosmos.ContainerClient
	}
)

const (
	database  = "adb2c"
	container = "principal"
)

func (h *Handler) Constructor(
	client *azcosmos.Client,
	_ *struct{ args.Optional }, translator ut.Translator,
) {
	h.principals = azure.Container(client, database, container)
	h.setValidationRules(translator)
}

func (h *Handler) Create(
	_ *struct {
		handles.It
		authorizes.Required
	  }, create api.CreatePrincipal,
	_ *struct{ args.Optional }, ctx context.Context,
) (p api.PrincipalCreated, err error) {
	id := model.NewId()
	principal := model.Principal{
		Id:    id,
		Type:  create.Type,
		Name:  create.Name,
		Scope: create.Scope,
	}
	if included, ok := model.Union([]string(nil), create.IncludedIds...); ok {
		principal.IncludedIds = included
	}
	pk := azcosmos.NewPartitionKeyString(principal.Scope)
	_, err = azure.CreateItem(ctx, &principal, pk, h.principals, nil)
	if err == nil {
		p.PrincipalId = id
	}
	return
}

func (h *Handler) Include(
	_ *struct {
		handles.It
		authorizes.Required
	  }, assign api.IncludePrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) miruken.HandleResult {
	pid := assign.PrincipalId
	pk  := azcosmos.NewPartitionKeyString(assign.Scope)
	_, found, err := azure.ReplaceItem(ctx, func(principal *model.Principal) (bool, error) {
		updated, changed := model.Union(principal.IncludedIds, assign.IncludedIds...)
		if changed {
			principal.IncludedIds = updated
		}
		return changed, nil
	}, pid, pk, h.principals, nil)
	switch {
	case !found:
		return miruken.NotHandled
	case err != nil:
		return miruken.NotHandled.WithError(err)
	default:
		return miruken.Handled
	}
}

func (h *Handler) Exclude(
	_ *struct {
		handles.It
		authorizes.Required
	  }, revoke api.ExcludePrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) miruken.HandleResult {
	pid := revoke.PrincipalId
	pk := azcosmos.NewPartitionKeyString(revoke.Scope)
	_, found, err := azure.ReplaceItem(ctx, func(principal *model.Principal) (bool, error) {
		updated, changed := model.Difference(principal.IncludedIds, revoke.ExcludedIds...)
		if changed {
			if len(updated) == 0 {
				principal.IncludedIds = nil
			} else {
				principal.IncludedIds = updated
			}
		}
		return changed, nil
	}, pid, pk, h.principals, nil)
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
	  }, remove api.RemovePrincipal,
	_ *struct{ args.Optional }, ctx context.Context,
) miruken.HandleResult {
	pid := remove.PrincipalId
	pk := azcosmos.NewPartitionKeyString(remove.Scope)
	_, found, err := azure.DeleteItem(ctx, pid, pk, h.principals, nil)
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
	_ *handles.It, get api.GetPrincipal,
	_ *struct{ args.Optional }, ctx context.Context,
) (api.Principal, miruken.HandleResult) {
	pid := get.PrincipalId
	pk := azcosmos.NewPartitionKeyString(get.Scope)
	_, item, found, err := azure.ReadItem[model.Principal](ctx, pid, pk, h.principals, nil)
	switch {
	case !found:
		return api.Principal{}, miruken.NotHandled
	case err != nil:
		return api.Principal{}, miruken.NotHandled.WithError(err)
	default:
		return item.ToApi(), miruken.Handled
	}
}

func (h *Handler) Find(
	_ *handles.It, find api.FindPrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) ([]api.Principal, error) {
	var sql strings.Builder
	sql.WriteString("SELECT * FROM principal p")

	var params []azcosmos.QueryParameter

	cond := " WHERE"

	if typ := find.Type; typ != "" {
		sql.WriteString(" WHERE CONTAINS(p.type, @type, true)")
		params = append(params, azcosmos.QueryParameter{Name: "@type", Value: typ})
		cond = " AND"
	}

	if name := find.Name; name != "" {
		sql.WriteString(cond)
		sql.WriteString(" CONTAINS(p.name, @name, true)")
		params = append(params, azcosmos.QueryParameter{Name: "@name", Value: name})
	}

	if ctx == nil {
		ctx = context.Background()
	}

	pk := azcosmos.NewPartitionKeyString(find.Scope)
	queryPager := h.principals.NewQueryItemsPager(sql.String(),
		pk, &azcosmos.QueryOptions{QueryParameters: params})

	principals := make([]api.Principal, 0)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range queryResponse.Items {
			var principal model.Principal
			if err := json.Unmarshal(item, &principal); err != nil {
				return nil, err
			}
			principals = append(principals, principal.ToApi())
		}
	}

	return principals, nil
}

func (h *Handler) Flatten(
	_ *handles.It, flatten api.FlattenPrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) *promise.Promise[[]api.Principal] {
	return promise.New(nil, func(resolve func([]api.Principal), reject func(error), onCancel func(func())) {
		queue := make(map[string][]string, len(flatten.PrincipalIds))
		for _, id := range flatten.PrincipalIds {
			if _, ok := queue[id]; !ok {
				queue[id] = []string{id}
			}
		}

		if ctx == nil {
			ctx = context.Background()
		}

		pk := azcosmos.NewPartitionKeyString(flatten.Scope)
		principals := make(map[string]api.Principal, len(flatten.PrincipalIds))

		for len(queue) > 0 {
			ids := make([]string, 0, len(queue))
			for pid := range queue {
				if _, found := principals[pid]; found {
					continue
				}
				ids = append(ids, pid)
			}

			if len(ids) == 0 {
				break
			}

			next := make(map[string][]string)
			queryPager := h.principals.NewQueryItemsPager(
				"SELECT * FROM p WHERE ARRAY_CONTAINS(@ids, p.id)",
				pk, &azcosmos.QueryOptions{QueryParameters: []azcosmos.QueryParameter{
					{Name: "@ids", Value: ids}}})

			for queryPager.More() {
				queryResponse, err := queryPager.NextPage(ctx)
				if err != nil {
					reject(err)
					return
				}
				for _, item := range queryResponse.Items {
					var principal model.Principal
					if err := json.Unmarshal(item, &principal); err != nil {
						reject(err)
						return
					}
					path := queue[principal.Id]
					for _, childId := range principal.IncludedIds {
						for _, ancestorId := range path {
							if childId == ancestorId {
								reject(fmt.Errorf("circular reference detected: %s", childId))
								return
							}
						}
						next[childId] = append(path, childId)
					}
					principals[principal.Id] = principal.ToApi()
				}
			}

			queue = next
		}

		result := make([]api.Principal, 0, len(principals))
		for _, p := range principals {
			result = append(result, p)
		}
		resolve(result)
	})
}

func (h *Handler) setValidationRules(
	translator ut.Translator,
) {
	_ = h.Validates1.WithRules(
		play.Rules{
			play.Type[api.CreatePrincipal](play.Constraints{
				"Type":  "required,alphanum",
				"Name":  "required",
				"Scope": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.IncludePrincipals](play.Constraints{
				"PrincipalId": "required",
				"Scope":       "required",
				"IncludedIds": "gt=0,dive,required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.ExcludePrincipals](play.Constraints{
				"PrincipalId": "required",
				"Scope":       "required",
				"ExcludedIds": "gt=0,dive,required",
			}),
		}, nil, translator)

	_ = h.Validates4.WithRules(
		play.Rules{
			play.Type[api.RemovePrincipal](play.Constraints{
				"PrincipalId": "required",
				"Scope":       "required",
			}),
		}, nil, translator)

	_ = h.Validates5.WithRules(
		play.Rules{
			play.Type[api.GetPrincipal](play.Constraints{
				"PrincipalId": "required",
				"Scope":       "required",
			}),
		}, nil, translator)

	_ = h.Validates6.WithRules(
		play.Rules{
			play.Type[api.FindPrincipals](play.Constraints{
				"Type":  "omitempty,alphanum",
				"Scope": "required",
			}),
		}, nil, translator)

	_ = h.Validates7.WithRules(
		play.Rules{
			play.Type[api.FlattenPrincipals](play.Constraints{
				"Scope":        "required",
				"PrincipalIds": "gt=0,dive,required",
			}),
		}, nil, translator)
}
