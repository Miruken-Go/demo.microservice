package entitlement

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
		play.Validates1[api.CreateEntitlement]
		play.Validates2[api.RemoveEntitlement]
		play.Validates3[api.GetEntitlement]
		play.Validates4[api.FindEntitlements]

		entitlements *azcosmos.ContainerClient
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
	h.db           = db
	h.entitlements = azure.Container(client, database, container)
	h.setValidationRules(translator)
}

func (h *Handler) Create(
	_*struct{
		handles.It
		authorizes.Required
	  }, create api.CreateEntitlement,
	_*struct{args.Optional}, ctx context.Context,
) (e api.EntitlementCreated, err error) {
	id := model.NewId()
	entitlement := model.Entitlement{
		Id:             id.String(),
		Type:           "Entitlement",
		Name:           create.Name,
		Scope:          create.Domain,
	}
	pk := azcosmos.NewPartitionKeyString(entitlement.Scope)
	_, err = azure.CreateItem(&entitlement, ctx, pk, h.entitlements, nil)
	if err == nil {
		e.EntitlementId = id
	}
	return
}


func (h *Handler) Remove(
	_*struct{
		handles.It
		authorizes.Required
	}, remove api.RemoveEntitlement,
	_*struct{args.Optional}, ctx context.Context,
) error {
	eid := remove.EntitlementId.String()
	pk  := azcosmos.NewPartitionKeyString(remove.Domain)
	_, err := h.entitlements.DeleteItem(ctx, pk, eid, nil)
	return err
}

func (h *Handler) Get(
	_ *handles.It, get api.GetEntitlement,
	_*struct{args.Optional}, ctx context.Context,
) (api.Entitlement, miruken.HandleResult) {
	eid := get.EntitlementId.String()
	pk  := azcosmos.NewPartitionKeyString(get.Domain)
	item, found, err := azure.ReadItem[model.Entitlement](ctx, eid, pk, h.entitlements, nil)
	if !found {
		return api.Entitlement{}, miruken.NotHandled
	} else if err != nil {
		return api.Entitlement{}, miruken.NotHandled.WithError(err)
	} else {
		return item.ToApi(), miruken.Handled
	}
}

func (h *Handler) Find(
	_ *handles.It, find api.FindEntitlements,
	_*struct{args.Optional}, ctx context.Context,
) ([]api.Entitlement, error) {
	rows, err := h.db.QueryxContext(ctx, fmt.Sprintf(
		`SELECT CROSS PARTITION * FROM e
 			WITH database=%s WITH collection=%s`, database, container),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	results := make([]api.Entitlement, 0)
	for rows.Next() {
		row := make(model.EntitlementMap)
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
			play.Type[api.CreateEntitlement](map[string]string{
				"Name":   "required",
				"Domain": "required",
			}),
		}, nil, translator)

	_ = h.Validates2.WithRules(
		play.Rules{
			play.Type[api.RemoveEntitlement](map[string]string{
				"EntitlementId": "required",
			}),
		}, nil, translator)

	_ = h.Validates3.WithRules(
		play.Rules{
			play.Type[api.GetEntitlement](map[string]string{
				"EntitlementId": "required",
			}),
		}, nil, translator)
}