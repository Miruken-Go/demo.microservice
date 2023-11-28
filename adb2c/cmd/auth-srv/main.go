package main

import (
	"errors"
	auth2 "github.com/miruken-go/demo.microservice/adb2c/auth"
	"github.com/miruken-go/miruken/context"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-logr/zerologr"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/miruken-go/miruken"
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
	play "github.com/miruken-go/miruken/validates/play"
	"github.com/rs/zerolog"
)

type Config struct {
	App struct {
		Version string
		Source  struct {
			Url string
		}
		Port    string
	}
	OpenApi openapi.Config
}

func main() {
	// logging
	zl := zerolog.New(os.Stderr)
	zl = zl.With().Timestamp().Logger()
	logger := zerologr.New(&zl)

	// configuration
	var k = koanf.New(".")
	err := k.Load(env.Provider("", "__", nil), nil,
		koanf.WithMergeFunc(koanfp.Merge))
	if err != nil {
		logger.Error(err, "error loading configuration")
		os.Exit(1)
	}

	var appConfig Config
	if err := k.Unmarshal("", &appConfig); err != nil {
		logger.Error(err, "error unmarshalling configuration")
		os.Exit(1)
	}

	// openapi generator
	openapiGen := openapi.Feature(openapi3.T{
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
			{
				"adb2c_auth": appConfig.OpenApi.ScopeNames(),
			},
		},
	})

	// initialize miruken
	handler, err := miruken.Setup(
		auth2.Feature, jwt.Feature(), play.Feature(),
		config.Feature(koanfp.P(k)), stdjson.Feature(),
		logs.Feature(logger),
		openapiGen).
		Specs(&api.GoPolymorphism{}).
		Options(stdjson.CamelCase).
		Handler()

	if err != nil {
		logger.Error(err, "setup failed")
		os.Exit(1)
	}

	ctx := context.New(handler)
	defer ctx.End(nil)

	// configure routes
	var mux http.ServeMux

	// Polymorphic miruken endpoints
	poly := httpsrv.Api(ctx,
		auth.WithFlowRef("Login.OAuth").Bearer(),
	)
	mux.Handle("/process", poly)
	mux.Handle("/process/", poly)
	mux.Handle("/publish", poly)
	mux.Handle("/publish/", poly)

	// OpenAPI document and swagger endpoints
	docs := openapiGen.Docs()
	mux.Handle("/openapi", openapi.Handler(docs, true))
	mux.Handle("/", ui.Handler("", docs, appConfig.OpenApi))

	// start http server
	port := appConfig.App.Port
	if port == "" {
		port = "8080"
	}
	err = http.ListenAndServe(":"+port, &mux)

	if errors.Is(err, http.ErrServerClosed) {
		logger.Info("server closed")
	} else if err != nil {
		logger.Error(err, "error starting server")
		os.Exit(1)
	}
}

