package user

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	a "github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
	}
)

func (h *Handler) Constructor(
) {
}

func (h *Handler) List(
	_ *struct {
		handles.It
		authorizes.Required
	  }, list api.ListUsers,
	cred *azidentity.ClientSecretCredential,
	_ *struct{ args.Optional }, ctx context.Context,
) ([]api.User, error) {
	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(
		cred, []string{"https://graph.microsoft.com/.default"})

	if err != nil {
		fmt.Println(err)
	}

	requestAdapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		fmt.Println(err)
	}

	graphClient := msgraphsdk.NewGraphServiceClient(requestAdapter)

	requestParameters := &graphusers.UsersRequestBuilderGetQueryParameters{
		//Select: [] string {"displayName","jobTitle"},
	}
	configuration := &graphusers.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	users, err := graphClient.Users().Get(context.Background(), configuration)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(users)

	return []api.User{}, nil
}


