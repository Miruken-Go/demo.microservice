package principal

//go:generate $GOPATH/bin/miruken -tests

import (
	"encoding/json"
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
		play.Validates2[api.AssignEntitlements]
		play.Validates3[api.RevokeEntitlements]
		play.Validates4[api.RemovePrincipal]
		play.Validates5[api.GetPrincipal]
		play.Validates6[api.FindPrincipals]

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
		Id:               id.String(),
		Type:             create.Type,
		Name:             create.Name,
		Scope:            create.Domain,
		EntitlementNames: create.EntitlementNames,
	}
	pk := azcosmos.NewPartitionKeyString(principal.Scope)
	_, err = azure.CreateItem(ctx, &principal, pk, h.principals, nil)
	if err == nil {
		p.PrincipalId = id
	}
	return
}

func (h *Handler) Assign(
	_ *struct {
		handles.It
		authorizes.Required
	}, assign api.AssignEntitlements,
	_ *struct{ args.Optional }, ctx context.Context,
) error {
	pid := assign.PrincipalId.String()
	pk := azcosmos.NewPartitionKeyString(assign.Domain)
	_, _, err := azure.ReplaceItem(ctx, func(principal *model.Principal) (bool, error) {
		add := assign.EntitlementNames
		updated, changed := model.Union(principal.EntitlementNames, add...)
		if changed {
			principal.EntitlementNames = updated
		}
		return changed, nil
	}, pid, pk, h.principals, nil)
	return err
}

func (h *Handler) Revoke(
	_ *struct {
		handles.It
		authorizes.Required
	}, revoke api.RevokeEntitlements,
	_ *struct{ args.Optional }, ctx context.Context,
) error {
	pid := revoke.PrincipalId.String()
	pk := azcosmos.NewPartitionKeyString(revoke.Domain)
	_, _, err := azure.ReplaceItem(ctx, func(principal *model.Principal) (bool, error) {
		remove := revoke.EntitlementNames
		updated, changed := model.Difference(principal.EntitlementNames, remove...)
		if changed {
			principal.EntitlementNames = updated
		}
		return changed, nil
	}, pid, pk, h.principals, nil)
	return err
}

func (h *Handler) Remove(
	_ *struct {
		handles.It
		authorizes.Required
	}, remove api.RemovePrincipal,
	_ *struct{ args.Optional }, ctx context.Context,
) error {
	pid := remove.PrincipalId.String()
	pk := azcosmos.NewPartitionKeyString(remove.Domain)
	_, err := h.principals.DeleteItem(ctx, pk, pid, nil)
	return err
}

func (h *Handler) Get(
	_ *handles.It, get api.GetPrincipal,
	_ *struct{ args.Optional }, ctx context.Context,
) (api.Principal, miruken.HandleResult) {
	pid := get.PrincipalId.String()
	pk := azcosmos.NewPartitionKeyString(get.Domain)
	item, found, err := azure.ReadItem[model.Principal](ctx, pid, pk, h.principals, nil)
	switch {
	case !found || item.Type == model.EntitlementType:
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
	params := []azcosmos.QueryParameter{
		{"@entitlement", model.EntitlementType},
	}
	var sql strings.Builder
	sql.WriteString("SELECT * FROM principal p WHERE p.type != @entitlement")

	if typ := find.Type; typ != "" {
		sql.WriteString(" AND CONTAINS(p.type, @type, true)")
		params = append(params, azcosmos.QueryParameter{Name: "@type", Value: typ})
	}

	if name := find.Name; name != "" {
		sql.WriteString(" AND CONTAINS(p.name, @name, true)")
		params = append(params, azcosmos.QueryParameter{Name: "@name", Value: name})
	}

	pk := azcosmos.NewPartitionKeyString(find.Domain)
	queryPager := h.principals.NewQueryItemsPager(sql.String(), pk, &azcosmos.QueryOptions{
		QueryParameters: params,
	})

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

func (h *Handler) setValidationRules(
	translator ut.Translator,
) {
	_ = h.Validates1.WithRules(
		play.Rules{
			play.Type[api.CreatePrincipal](play.Constraints{
				"Type":   "required,alphanum",
				"Name":   "required",
				"Domain": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.AssignEntitlements](play.Constraints{
				"PrincipalId":      "required",
				"Domain":           "required",
				"EntitlementNames": "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.RevokeEntitlements](play.Constraints{
				"PrincipalId":      "required",
				"Domain":           "required",
				"EntitlementNames": "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates4.WithRules(
		play.Rules{
			play.Type[api.RemovePrincipal](play.Constraints{
				"PrincipalId": "required",
				"Domain":      "required",
			}),
		}, nil, translator)

	_ = h.Validates5.WithRules(
		play.Rules{
			play.Type[api.GetPrincipal](play.Constraints{
				"PrincipalId": "required",
				"Domain":      "required",
			}),
		}, nil, translator)

	_ = h.Validates6.WithRules(
		play.Rules{
			play.Type[api.FindPrincipals](play.Constraints{
				"Type":   "omitempty,alphanum",
				"Domain": "required",
			}),
		}, nil, translator)
}
