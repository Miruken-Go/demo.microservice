package graph

import (
	a "github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
	"maps"
	"sync"
	"sync/atomic"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang-jwt/jwt/v5"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/miruken-go/miruken/config"
	"github.com/miruken-go/miruken/provides"
	"github.com/miruken-go/miruken/security"
)

type (
	// Factory of GraphServiceClient.
	Factory struct {
		cfg     Config
		clients atomic.Pointer[map[string]any]
		lock    sync.Mutex
	}

	// Client is a typed wrapper for a GraphServiceClient.
	Client[T azcore.TokenCredential] struct {
		*msgraphsdkgo.GraphServiceClient
		cred T
	}
)


// Client

func (c *Client[T]) Credential() azcore.TokenCredential {
	return c.cred
}


// Factory

func (f *Factory) Constructor(
	_ *struct {
		config.Load `path:"Graph"`
	  }, cfg Config,
) {
	f.cfg = cfg
}

func (f *Factory) ClientSecret(
	_ *provides.It,
	subject security.Subject,
) (*Client[*azidentity.ClientSecretCredential], error) {
	var token *jwt.Token

	// Find jwt credential
	for _, c := range subject.Credentials() {
		if t, ok := c.(*jwt.Token); ok {
			token = t
		}
	}
	if token == nil {
		return nil, nil
	}

	// Expect jwt map claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil
	}

	// Require aud (client id) claim
	var clientId string
	if aud, ok := claims["aud"]; !ok {
		return nil, nil
	} else {
		clientId = aud.(string)
	}

	// Use existing client secret credential
	if cache := f.clients.Load(); cache != nil {
		if cred, ok := (*cache)[clientId]; ok {
			if csc, ok := cred.(*Client[*azidentity.ClientSecretCredential]); ok {
				return csc, nil
			}
			return nil, nil
		}
	}

	// Lookup secret for client id
	secret := f.secret(clientId)
	if secret == "" {
		return nil, nil
	}

	// Require tenant id claim
	var tenantId string
	if tid, ok := claims["tid"]; !ok {
		return nil, nil
	} else {
		tenantId = tid.(string)
	}

	f.lock.Lock()
	defer f.lock.Unlock()

	cache := f.clients.Load()
	if cache != nil {
		if cred, ok := (*cache)[clientId]; ok {
			if csc, ok := cred.(*Client[*azidentity.ClientSecretCredential]); ok {
				return csc, nil
			}
			return nil, nil
		}
		cc := maps.Clone(*cache)
		cache = &cc
	} else {
		cache = &map[string]any{}
	}

	// Create client secret credential
	cred, err := azidentity.NewClientSecretCredential(tenantId, clientId, secret, nil)
	if err != nil {
		return nil, err
	}

	// Create request adapter with authentication
	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(
		cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil, err
	}

	requestAdapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		return nil, err
	}

	client := &Client[*azidentity.ClientSecretCredential]{
		msgraphsdk.NewGraphServiceClient(requestAdapter),
		cred,
	}

	(*cache)[clientId] = client
	f.clients.Store(cache)
	return client, nil
}

func (f *Factory) secret(clientId string) string {
	for _, client := range f.cfg.ClientSecret {
		if client.Id == clientId {
			return client.Secret
		}
	}
	return ""
}
