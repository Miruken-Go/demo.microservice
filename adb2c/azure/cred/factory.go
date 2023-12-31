package cred

import (
	"maps"
	"sync"
	"sync/atomic"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/miruken-go/miruken/config"
	"github.com/miruken-go/miruken/provides"
	"github.com/miruken-go/miruken/security"
)

// Factory of azure token credentials.
type Factory struct {
	cfg   Config
	cache atomic.Pointer[map[string]azcore.TokenCredential]
	lock  sync.Mutex
}


func (f *Factory) Constructor(
	_ *struct {
		config.Load `path:"Credentials"`
	  }, cfg Config,
) {
	f.cfg = cfg
}

func (f *Factory) ClientSecret(
	_ *provides.It,
	subject security.Subject,
) (*azidentity.ClientSecretCredential, error) {
	var token *jwt.Token

	// Find jwt in credentials
	for _, c := range subject.Credentials() {
		if t, ok := c.(*jwt.Token); ok {
			token = t
		}
	}
	if token == nil {
		return nil, nil
	}

	// Extract claims from jwt
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil
	}

	// Extract aud (client id) claim
	var clientId string
	if aud, ok := claims["aud"]; !ok {
		return nil, nil
	} else {
		clientId = aud.(string)
	}

	// Use existing client secret credential
	if cache := f.cache.Load(); cache != nil {
		if cred, ok := (*cache)[clientId]; ok {
			if csc, ok := cred.(*azidentity.ClientSecretCredential); ok {
				return csc, nil
			}
			return nil, nil
		}
	}

	// Find secret for client id
	secret := f.secret(clientId)
	if secret == "" {
		return nil, nil
	}

	// Extract tenant id claim
	var tenantId string
	if tid, ok := claims["tid"]; !ok {
		return nil, nil
	} else {
		tenantId = tid.(string)
	}

	f.lock.Lock()
	defer f.lock.Unlock()

	cache := f.cache.Load()
	if cache != nil {
		if cred, ok := (*cache)[clientId]; ok {
			if csc, ok := cred.(*azidentity.ClientSecretCredential); ok {
				return csc, nil
			}
			return nil, nil
		}
		cc := maps.Clone(*cache)
		cache = &cc
	} else {
		cache = &map[string]azcore.TokenCredential{}
	}

	// Create new client secret credential
	cred, err := azidentity.NewClientSecretCredential(tenantId, clientId, secret, nil)
	if err != nil {
		return nil, err
	}

	(*cache)[clientId] = cred
	f.cache.Store(cache)
	return cred, nil
}

func(f *Factory) secret(clientId string) string {
	for _, client := range f.cfg.ClientSecret {
		if client.Id == clientId {
			return client.Secret
		}
	}
	return ""
}