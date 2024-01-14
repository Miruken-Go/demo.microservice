package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-logr/zerologr"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/miruken-go/demo.microservice/adb2c/azure"
	"github.com/miruken-go/demo.microservice/adb2c/azure/db"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/api/http/httpsrv"
	"github.com/miruken-go/miruken/api/http/httpsrv/auth"
	"github.com/miruken-go/miruken/api/http/httpsrv/openapi"
	"github.com/miruken-go/miruken/api/http/httpsrv/openapi/ui"
	"github.com/miruken-go/miruken/api/json/stdjson"
	"github.com/miruken-go/miruken/config"
	koanfp "github.com/miruken-go/miruken/config/koanf"
	"github.com/miruken-go/miruken/logs"
	"github.com/miruken-go/miruken/security/jwt"
	"github.com/miruken-go/miruken/setup"
	"github.com/rs/zerolog"
)

type Config struct {
	App struct {
		Version string
		Source  struct {
			Url string
		}
	}
	Server  httpsrv.Config
	OpenApi openapi.Config
}

func main() {
	// logging
	zl := zerolog.New(os.Stdout)
	zl = zl.With().Timestamp().Logger()
	logger := zerologr.New(&zl)

	// configuration
	var k = koanf.New(".")
	err := k.Load(file.Provider("./app.yml"), yaml.Parser(),
		koanf.WithMergeFunc(koanfp.Merge))
	if err != nil {
		logger.Error(err, "error loading app.yml configuration")
		os.Exit(1)
	}
	err = k.Load(env.Provider("", "__", nil), nil,
		koanf.WithMergeFunc(koanfp.Merge))
	if err != nil {
		logger.Error(err, "error loading env configuration")
		os.Exit(1)
	}

	var appConfig Config
	if err := k.Unmarshal("", &appConfig); err != nil {
		logger.Error(err, "error unmarshalling configuration")
		os.Exit(1)
	}

	// openapi generator
	openapiGen := openapi.Feature(&openapi3.T{
		Info: &openapi3.Info{
			Title:       "Authorization Api",
			Description: "REST Api for managing User Authorization",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				Name: "Miruken",
				URL:  "https://github.com/Miruken-Go/demo.microservice",
			},
		},
		ExternalDocs: &openapi3.ExternalDocs{
			Description: "adb2c-auth-srv/" + appConfig.App.Version,
			URL:         appConfig.App.Source.Url,
		},
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"adb2c_auth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type: "oauth2",
						Flows: &openapi3.OAuthFlows{
							Implicit: &openapi3.OAuthFlow{
								AuthorizationURL: appConfig.OpenApi.AuthorizationUrl,
								TokenURL:         appConfig.OpenApi.TokenUrl,
								Scopes:           appConfig.OpenApi.ScopeMap(),
							},
						},
						OpenIdConnectUrl: appConfig.OpenApi.OpenIdConnectUrl,
					},
				},
			},
		},
		Security: openapi3.SecurityRequirements{
			{"adb2c_auth": appConfig.OpenApi.ScopeNames()},
		},
	})

	// initialize context
	ctx, err := setup.New(
		azure.Feature, jwt.Feature(),
		db.Feature(db.Provision[*azcosmos.Client]),
		config.Feature(koanfp.P(k)), stdjson.Feature(),
		logs.Feature(logger), openapiGen).
		Specs(&api.GoPolymorphism{}).
		Options(stdjson.CamelCase).
		Context()

	if err != nil {
		logger.Error(err, "setup failed")
		os.Exit(1)
	}

	defer ctx.End(nil)

	// Polymorphic api endpoints
	poly := httpsrv.Api(ctx,
		auth.WithFlowAlias("Login.OAuth").Bearer().Required(),
	)

	// configure routes
	var mux http.ServeMux
	mux.Handle("/process", poly)
	mux.Handle("/process/", poly)
	mux.Handle("/publish", poly)
	mux.Handle("/publish/", poly)

	// OpenAPI document and swagger endpoints
	docs := openapiGen.Docs()
	mux.Handle("/openapi", openapi.Handler(docs, true))
	mux.Handle("/", ui.Handler("", docs, &appConfig.OpenApi))

	// start http server
	if err := httpsrv.ListenAndServe(&mux, &appConfig.Server); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("server closed")
		} else if err != nil {
			logger.Error(err, "error starting server")
		}
	}
}
