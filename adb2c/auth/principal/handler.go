package principal

//go:generate $GOPATH/bin/miruken -tests

import (
	"fmt"
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

type (
	Handler struct {
		play.Validates1[api.CreatePrincipal]
		play.Validates2[api.AssignEntitlements]
		play.Validates3[api.RevokeEntitlements]
		play.Validates4[api.RemovePrincipal]
		play.Validates5[api.GetPrincipal]
		play.Validates6[api.FindPrincipals]

		principals *azcosmos.ContainerClient
		db *sqlx.DB
	}
)


const (
	database  = "adb2c"
	container = "principal"
)


func (h *Handler) Constructor(
	db     *sqlx.DB,
	client *azcosmos.Client,
	_*struct{args.Optional}, translator ut.Translator,
) {
	h.db         = db
	h.principals = azure.Container(client, database, container)
	h.setValidationRules(translator)
}

func (h *Handler) Create(
	_*struct{
		handles.It
		authorizes.Required
	  }, create api.CreatePrincipal,
	_*struct{args.Optional}, ctx context.Context,
) (p api.PrincipalCreated, err error) {
	id := model.NewId()
	principal := model.Principal{
		Id:             id.String(),
		Type:           create.Type,
		Name:           create.Name,
		Scope:          create.Domain,
		EntitlementIds: model.Strings(create.EntitlementIds),
	}
	pk := azcosmos.NewPartitionKeyString(principal.Scope)
	_, err = azure.CreateItem(&principal, ctx, pk, h.principals, nil)
	if err == nil {
		p.PrincipalId = id
	}
	return
}

func (h *Handler) Assign(
	_*struct{
		handles.It
		authorizes.Required
	}, assign api.AssignEntitlements,
	_*struct{args.Optional}, ctx context.Context,
) error {
	pid := assign.PrincipalId.String()
	pk  := azcosmos.NewPartitionKeyString(assign.Domain)
	_, _, err := azure.ReplaceItem(func(principal *model.Principal) (bool, error) {
		add := model.Strings(assign.EntitlementIds)
		updated, changed := model.Union(principal.EntitlementIds, add...)
		if changed {
			principal.EntitlementIds = updated
		}
		return changed, nil
	}, ctx, pid, pk, h.principals, nil)
	return err
}

func (h *Handler) Revoke(
	_*struct{
		handles.It
		authorizes.Required
	  }, revoke api.RevokeEntitlements,
	_*struct{args.Optional}, ctx context.Context,
) error {
	pid := revoke.PrincipalId.String()
	pk  := azcosmos.NewPartitionKeyString(revoke.Domain)
	_, _, err := azure.ReplaceItem(func(principal *model.Principal) (bool, error) {
		remove := model.Strings(revoke.EntitlementIds)
		updated, changed := model.Difference(principal.EntitlementIds, remove...)
		if changed {
			principal.EntitlementIds = updated
		}
		return changed, nil
	}, ctx, pid, pk, h.principals, nil)
	return err
}

func (h *Handler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	}, remove api.RemovePrincipal,
	_*struct{args.Optional}, ctx context.Context,
) error {
	pid := remove.PrincipalId.String()
	pk  := azcosmos.NewPartitionKeyString(remove.Domain)
	_, err := h.principals.DeleteItem(ctx, pk, pid, nil)
	return err
}

func (h *Handler) Get(
	_ *handles.It, get api.GetPrincipal,
	_*struct{args.Optional}, ctx context.Context,
) (api.Principal, miruken.HandleResult) {
	pid := get.PrincipalId.String()
	pk  := azcosmos.NewPartitionKeyString(get.Domain)
	item, found, err := azure.ReadItem[model.Principal](ctx, pid, pk, h.principals, nil)
	if !found || item.Type == "Entitlement" {
		return api.Principal{}, miruken.NotHandled
	} else if err != nil {
		return api.Principal{}, miruken.NotHandled.WithError(err)
	} else {
		return item.ToApi(), miruken.Handled
	}
}

func (h *Handler) Find(
	_ *handles.It, find api.FindPrincipals,
	_*struct{args.Optional}, ctx context.Context,
) ([]api.Principal, error) {
	rows, err := h.db.QueryxContext(ctx, fmt.Sprintf(
		`SELECT * FROM p WHERE p.type != 'Entitlement' AND p.scope = :1
 			WITH database=%s WITH collection=%s`,
			 database, container),
			 find.Domain,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	results := make([]api.Principal, 0)
	for rows.Next() {
		row := make(model.PrincipalMap)
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
			play.Type[api.CreatePrincipal](map[string]string{
				"Type":   "required",
				"Name":   "required",
				"Domain": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.AssignEntitlements](map[string]string{
				"PrincipalId":    "required",
				"EntitlementIds": "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.RevokeEntitlements](map[string]string{
				"PrincipalId":    "required",
				"EntitlementIds": "gt=0,required",
			}),
		}, nil, translator)

	_ = h.Validates4.WithRules(
		play.Rules{
			play.Type[api.RemovePrincipal](map[string]string{
				"PrincipalId": "required",
			}),
		}, nil, translator)

	_ = h.Validates5.WithRules(
		play.Rules{
			play.Type[api.GetPrincipal](map[string]string{
				"PrincipalId": "required",
			}),
		}, nil, translator)

	_ = h.Validates6.WithRules(
		play.Rules{
			play.Type[api.FindPrincipals](map[string]string{
				"Domain": "required",
			}),
		}, nil, translator)
}