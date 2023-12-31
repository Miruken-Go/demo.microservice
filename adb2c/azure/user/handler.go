package user

import (
	"context"
	"github.com/miruken-go/miruken/args"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	a "github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
	}
)


func (h *Handler) List(
	_ *struct {
		handles.It
		authorizes.Required
	  }, _ api.ListUsers,
	cred *azidentity.ClientSecretCredential,
	_ *struct{ args.Optional }, ctx context.Context,
) ([]api.User, error) {
	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(
		cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil, err
	}

	requestAdapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		ctx = context.Background()
	}

	graphClient := msgraphsdk.NewGraphServiceClient(requestAdapter)
	configuration := &graphusers.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: &graphusers.UsersRequestBuilderGetQueryParameters{
			Select: [] string {"id", "displayName","givenName","surname"},
		},
	}
	result, err := graphClient.Users().Get(ctx, configuration)
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


