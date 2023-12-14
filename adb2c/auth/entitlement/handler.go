package entitlement

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
		play.Validates1[api.CreateEntitlement]
		play.Validates2[api.RemoveEntitlement]
		play.Validates3[api.GetEntitlement]
		play.Validates4[api.FindEntitlements]

		entitlements *azcosmos.ContainerClient
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
	h.entitlements = azure.Container(client, database, container)
	h.setValidationRules(translator)
}

func (h *Handler) Create(
	_ *struct {
		handles.It
		authorizes.Required
	}, create api.CreateEntitlement,
	_ *struct{ args.Optional }, ctx context.Context,
) (e api.EntitlementCreated, err error) {
	id := model.NewId()
	entitlement := model.Entitlement{
		Id:          id.String(),
		Type:        model.EntitlementType,
		Name:        create.Name,
		Scope:       create.Domain,
		Description: create.Description,
	}
	pk := azcosmos.NewPartitionKeyString(entitlement.Scope)
	_, err = azure.CreateItem(ctx, &entitlement, pk, h.entitlements, nil)
	if err == nil {
		e.EntitlementId = id
	}
	return
}

func (h *Handler) Remove(
	_ *struct {
		handles.It
		authorizes.Required
	}, remove api.RemoveEntitlement,
	_ *struct{ args.Optional }, ctx context.Context,
) error {
	eid := remove.EntitlementId.String()
	pk := azcosmos.NewPartitionKeyString(remove.Domain)
	_, err := h.entitlements.DeleteItem(ctx, pk, eid, nil)
	return err
}

func (h *Handler) Get(
	_ *handles.It, get api.GetEntitlement,
	_ *struct{ args.Optional }, ctx context.Context,
) (api.Entitlement, miruken.HandleResult) {
	eid := get.EntitlementId.String()
	pk := azcosmos.NewPartitionKeyString(get.Domain)
	item, found, err := azure.ReadItem[model.Entitlement](ctx, eid, pk, h.entitlements, nil)
	switch {
	case !found || item.Type != model.EntitlementType:
		return api.Entitlement{}, miruken.NotHandled
	case err != nil:
		return api.Entitlement{}, miruken.NotHandled.WithError(err)
	default:
		return item.ToApi(), miruken.Handled
	}
}

func (h *Handler) Find(
	_ *handles.It, find api.FindEntitlements,
	_ *struct{ args.Optional }, ctx context.Context,
) ([]api.Entitlement, error) {
	params := []azcosmos.QueryParameter{
		{"@entitlement", model.EntitlementType},
	}
	var sql strings.Builder
	sql.WriteString("SELECT * FROM principal p WHERE p.type = @entitlement")

	if name := find.Name; name != "" {
		sql.WriteString(" AND CONTAINS(p.name, @name, true)")
		params = append(params, azcosmos.QueryParameter{Name: "@name", Value: name})
	}

	pk := azcosmos.NewPartitionKeyString(find.Domain)
	queryPager := h.entitlements.NewQueryItemsPager(sql.String(), pk, &azcosmos.QueryOptions{
		QueryParameters: params,
	})

	entitlements := make([]api.Entitlement, 0)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range queryResponse.Items {
			var entitlement model.Entitlement
			if err := json.Unmarshal(item, &entitlement); err != nil {
				return nil, err
			}
			entitlements = append(entitlements, entitlement.ToApi())
		}
	}

	return entitlements, nil
}

func (h *Handler) setValidationRules(
	translator ut.Translator,
) {
	_ = h.Validates1.WithRules(
		play.Rules{
			play.Type[api.CreateEntitlement](play.Constraints{
				"Name":   "required",
				"Domain": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.RemoveEntitlement](play.Constraints{
				"EntitlementId": "required",
				"Domain":        "required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.GetEntitlement](play.Constraints{
				"EntitlementId": "required",
				"Domain":        "required",
			}),
		}, nil, translator)

	_ = h.Validates4.WithRules(
		play.Rules{
			play.Type[api.FindEntitlements](play.Constraints{
				"Domain": "required",
			}),
		}, nil, translator)
}
