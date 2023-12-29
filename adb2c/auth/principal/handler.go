package principal

//go:generate $GOPATH/bin/miruken -tests

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	ut "github.com/go-playground/universal-translator"
	api2 "github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/demo.microservice/adb2c/auth/internal/model"
	"github.com/miruken-go/demo.microservice/adb2c/azure"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/promise"
	"github.com/miruken-go/miruken/security/authorizes"
	play "github.com/miruken-go/miruken/validates/play"
	"golang.org/x/net/context"
)

type (
	Handler struct {
		play.Validates1[api2.CreatePrincipal]
		play.Validates2[api2.IncludePrincipals]
		play.Validates3[api2.ExcludePrincipals]
		play.Validates4[api2.RemovePrincipal]
		play.Validates5[api2.GetPrincipal]
		play.Validates6[api2.FindPrincipals]
		play.Validates7[api2.ExpandPrincipals]
		play.Validates8[api2.SatisfyPrincipals]

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
	  }, create api2.CreatePrincipal,
	_ *struct{ args.Optional }, ctx context.Context,
) (p api2.PrincipalCreated, err error) {
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
	  }, assign api2.IncludePrincipals,
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
	  }, revoke api2.ExcludePrincipals,
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
	  }, remove api2.RemovePrincipal,
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
	_ *handles.It, get api2.GetPrincipal,
	_ *struct{ args.Optional }, ctx context.Context,
) (api2.Principal, miruken.HandleResult) {
	pid := get.PrincipalId
	pk := azcosmos.NewPartitionKeyString(get.Scope)
	_, item, found, err := azure.ReadItem[model.Principal](ctx, pid, pk, h.principals, nil)
	switch {
	case !found:
		return api2.Principal{}, miruken.NotHandled
	case err != nil:
		return api2.Principal{}, miruken.NotHandled.WithError(err)
	default:
		return item.ToApi(), miruken.Handled
	}
}

func (h *Handler) Find(
	_ *handles.It, find api2.FindPrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) ([]api2.Principal, error) {
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

	principals := make([]api2.Principal, 0)
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

func (h *Handler) Expand(
	_ *handles.It, expand api2.ExpandPrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) *promise.Promise[[]api2.Principal] {
	return promise.New(nil, func(resolve func([]api2.Principal), reject func(error), onCancel func(func())) {
		queue := make(map[string][]string, len(expand.PrincipalIds))
		for _, pid := range expand.PrincipalIds {
			if _, ok := queue[pid]; !ok {
				queue[pid] = []string{pid}
			}
		}

		if ctx == nil {
			ctx = context.Background()
		}

		pk := azcosmos.NewPartitionKeyString(expand.Scope)
		principals := make(map[string]api2.Principal, len(expand.PrincipalIds))

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
					{Name: "@ids", Value: ids}},
				})

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

		result := make([]api2.Principal, 0, len(principals))
		for _, p := range principals {
			result = append(result, p)
		}
		resolve(result)
	})
}

func (h *Handler) Satisfy(
	_ *handles.It, satisfy api2.SatisfyPrincipals,
	_ *struct{ args.Optional }, ctx context.Context,
) *promise.Promise[[]string] {
	return promise.New(nil, func(
		resolve func([]string), reject func(error), onCancel func(func())) {

		if ctx == nil {
			ctx = context.Background()
		}

		queue := satisfy.PrincipalIds
		pk := azcosmos.NewPartitionKeyString(satisfy.Scope)
		principals := make(map[string]struct{}, len(queue))

		for len(queue) > 0 {
			ids := make([]string, 0, len(queue))
			for _, pid := range queue {
				if _, found := principals[pid]; found {
					continue
				}
				principals[pid] = struct{}{}
				ids = append(ids, pid)
			}

			if len(ids) == 0 {
				break
			}

			queryPager := h.principals.NewQueryItemsPager(
				"SELECT * FROM p WHERE ARRAY_LENGTH(SetIntersect(p.includedIds, @ids)) != 0",
				pk, &azcosmos.QueryOptions{QueryParameters: []azcosmos.QueryParameter{
					{Name: "@ids", Value: ids}},
				})

			for queryPager.More() {
				queryResponse, err := queryPager.NextPage(ctx)
				if err != nil {
					reject(err)
					return
				}
				queue := queue[:0]
				for _, item := range queryResponse.Items {
					var principal model.Principal
					if err := json.Unmarshal(item, &principal); err != nil {
						reject(err)
						return
					}
					queue = append(queue, principal.Id)
				}
			}
		}

		result := make([]string, 0, len(principals))
		for pid := range principals {
			result = append(result, pid)
		}
		resolve(result)
	})
}

func (h *Handler) setValidationRules(
	translator ut.Translator,
) {
	_ = h.Validates1.WithRules(
		play.Rules{
			play.Type[api2.CreatePrincipal](play.Constraints{
				"Type":  "required,alphanum",
				"Name":  "required",
				"Scope": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api2.IncludePrincipals](play.Constraints{
				"PrincipalId": "required",
				"Scope":       "required",
				"IncludedIds": "gt=0,dive,required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api2.ExcludePrincipals](play.Constraints{
				"PrincipalId": "required",
				"Scope":       "required",
				"ExcludedIds": "gt=0,dive,required",
			}),
		}, nil, translator)

	_ = h.Validates4.WithRules(
		play.Rules{
			play.Type[api2.RemovePrincipal](play.Constraints{
				"PrincipalId": "required",
				"Scope":       "required",
			}),
		}, nil, translator)

	_ = h.Validates5.WithRules(
		play.Rules{
			play.Type[api2.GetPrincipal](play.Constraints{
				"PrincipalId": "required",
				"Scope":       "required",
			}),
		}, nil, translator)

	_ = h.Validates6.WithRules(
		play.Rules{
			play.Type[api2.FindPrincipals](play.Constraints{
				"Type":  "omitempty,alphanum",
				"Scope": "required",
			}),
		}, nil, translator)

	_ = h.Validates7.WithRules(
		play.Rules{
			play.Type[api2.ExpandPrincipals](play.Constraints{
				"Scope":        "required",
				"PrincipalIds": "gt=0,dive,required",
			}),
		}, nil, translator)

	_ = h.Validates8.WithRules(
		play.Rules{
			play.Type[api2.SatisfyPrincipals](play.Constraints{
				"Scope":        "required",
				"PrincipalIds": "gt=0,dive,required",
			}),
		}, nil, translator)
}
