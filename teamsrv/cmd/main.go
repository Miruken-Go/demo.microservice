package main

import (
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-logr/zerologr"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/miruken-go/demo.microservice/team"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/api/http/httpsrv"
	"github.com/miruken-go/miruken/api/http/httpsrv/middleware"
	"github.com/miruken-go/miruken/api/http/httpsrv/openapi"
	"github.com/miruken-go/miruken/api/http/httpsrv/openapi/ui"
	"github.com/miruken-go/miruken/api/json/stdjson"
	"github.com/miruken-go/miruken/config"
	koanfp "github.com/miruken-go/miruken/config/koanf"
	"github.com/miruken-go/miruken/logs"
	"github.com/miruken-go/miruken/security/jwt"
	"github.com/miruken-go/miruken/security/jwt/jwks"
	play "github.com/miruken-go/miruken/validates/play"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

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

	// openapi generator
	openapiGen := openapi.Feature(openapi3.T{
		Info: &openapi3.Info{
			Title:       "Team Api",
			Description: "REST Api for managing Teams " + k.String("App.Version"),
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				URL: "https://github.com/Miruken-Go/demo.microservice",
			},
		},
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"team_auth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type: "oauth2",
						Flows: &openapi3.OAuthFlows{
							Implicit: &openapi3.OAuthFlow{
								AuthorizationURL: "https://teamsrvdevcraig.b2clogin.com/teamsrvdevcraig.onmicrosoft.com/b2c_1_signIn/oauth2/v2.0/authorize",
								TokenURL:         "https://teamsrvdevcraig.b2clogin.com/teamsrvdevcraig.onmicrosoft.com/b2c_1_signIn/oauth2/v2.0/token",
								Scopes: map[string]string{
									"https://teamsrvdevcraig.onmicrosoft.com/60f123ab-de4d-4d2f-bb93-b54fddc38ee1/Person.Create": "create a person",
									"https://teamsrvdevcraig.onmicrosoft.com/60f123ab-de4d-4d2f-bb93-b54fddc38ee1/Team.Create":   "create a team",
								},
							},
						},
						OpenIdConnectUrl: "https://login.microsoftonline.com/048cf208-778f-496b-b892-9d03d15652cd/v2.0/.well-known/openid-configuration",
					},
				},
			},
		},
		Security: openapi3.SecurityRequirements{
			{"team_auth": []string{
				"https://teamsrvdevcraig.onmicrosoft.com/60f123ab-de4d-4d2f-bb93-b54fddc38ee1/Person.Create",
				"https://teamsrvdevcraig.onmicrosoft.com/60f123ab-de4d-4d2f-bb93-b54fddc38ee1/Team.Create",
			}},
		},
	})

	// initialize miruken
	handler, err := miruken.Setup(
		team.Feature, jwt.Feature(), jwks.Feature(),
		play.Feature(), config.Feature(koanfp.P(k)),
		stdjson.Feature(), logs.Feature(logger), openapiGen).
		Specs(&api.GoPolymorphism{}).
		Options(stdjson.CamelCase).
		Handler()

	if err != nil {
		logger.Error(err, "setup failed")
		os.Exit(1)
	}

	docs := openapiGen.Docs()

	h := httpsrv.Pipeline(handler, &middleware.Login{Flows: []string{"login.oauth"}})

	// configure routes
	var mux http.ServeMux
	mux.Handle("/process", h)
	mux.Handle("/process/", h)
	mux.Handle("/publish", h)
	mux.Handle("/publish/", h)
	mux.Handle("/openapi", openapi.Handler(docs, true))
	mux.Handle("/", ui.Handler("", docs))

	// start http server
	err = http.ListenAndServe(":8080", &mux)

	if errors.Is(err, http.ErrServerClosed) {
		logger.Info("server closed")
	} else if err != nil {
		logger.Error(err, "error starting server")
		os.Exit(1)
	}
}
