package user

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/demo.microservice/adb2c/azure/graph"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
)

//go:generate $GOPATH/bin/miruken -tests

type Handler struct {}


func (h *Handler) List(
	_ *struct {
		handles.It
		authorizes.Required
	  }, list api.ListUsers,
	client *graph.Client[*azidentity.ClientSecretCredential],
	_ *struct{ args.Optional }, ctx context.Context,
) ([]api.User, error) {
	query := graphusers.UsersRequestBuilderGetQueryParameters{
		Select: userFields,
	}
	if filter := list.Filter; filter != "" {
		// The Users resource does not support the 'contains' operator.
		criteria := fmt.Sprintf(`
			startsWith(displayName,'%s') or
			startsWith(givenName,'%s') or
			startsWith(surname,'%s') or
			startsWith(mail,'%s')`,
			filter, filter, filter, filter)
		query.Filter = &criteria
	}
	configuration := &graphusers.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: &query,
	}
	result, err := client.Users().Get(ctx, configuration)
	if err != nil {
		return nil, err
	}

	col := result.GetValue()
	users := make([]api.User, len(col))
	for i, user := range col {
		ToApi(user, &users[i])
	}

	return users, nil
}


// AuthorizeList is a placeholder authorization policy for the
// authorizes.Required marker above. Miruken's authorizes.Required now
// fails closed (denies) when no policy answers the check, rather than
// allowing by default; before this change List had no policy at all
// and was implicitly allowed for everyone. This placeholder preserves
// that prior de facto behavior explicitly - replace with real
// authorization logic as needed.
func (h *Handler) AuthorizeList(
	_ *authorizes.It, _ api.ListUsers,
) bool {
	return true
}

var (
	userFields = []string {"id","givenName","surname","displayName","mail"}
)