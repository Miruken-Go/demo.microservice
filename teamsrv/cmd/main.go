package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/miruken-go/miruken/security/password"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-logr/zerologr"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/miruken-go/demo.microservice/team"
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

func authzHandler(w http.ResponseWriter, r *http.Request) {
	username := "ooYymDzee5!V&v8gk7*s"
	password := "i**72R#PLWbx8&#$I$ok"
	u, p, ok := r.BasicAuth()
	if !ok || u != username || p != password {
		w.WriteHeader(401)
		return
	}

	type Request struct {
		Type     string `json:"@type"`
		Email    string
		ObjectId string
		Scope    string
	}

	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(422)
		return
	}

	fmt.Println("")
	fmt.Println("New Request")
	fmt.Println("@type: ", request.Type)
	fmt.Println("email: ", request.Email)
	fmt.Println("objectId: ", request.ObjectId)
	fmt.Println("scope: ", request.Scope)

	type Response struct {
		Groups       []string
		Roles        []string
		Entitlements []string
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Response{
		Groups:       []string{"oncall"},
		Roles:        []string{"admin", "coach", "player"},
		Entitlements: []string{"createTeam", "updateTeam", "createPerson", "updatePerson"},
	})
	if err != nil {
		w.WriteHeader(500)
	}
}

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
			Title:       "Team Api",
			Description: "REST Api for managing Teams ",
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
			Description: "teamsrv/" + appConfig.App.Version,
			URL:         appConfig.App.Source.Url,
		},
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"team_auth": &openapi3.SecuritySchemeRef{
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
				"team_auth": appConfig.OpenApi.ScopeNames(),
			},
		},
	})

	// initialize miruken
	handler, err := miruken.Setup(
		team.Feature, jwt.Feature(),
		password.Feature(), play.Feature(),
		config.Feature(koanfp.P(k)), stdjson.Feature(),
		logs.Feature(logger), openapiGen).
		Specs(&api.GoPolymorphism{}).
		Options(stdjson.CamelCase).
		Handler()

	if err != nil {
		logger.Error(err, "setup failed")
		os.Exit(1)
	}

	docs := openapiGen.Docs()

	h := httpsrv.Api(handler,
		auth.WithFlowRef("Login.OAuth").Bearer(),
	)

	// configure routes
	var mux http.ServeMux
	mux.Handle("/process", h)
	mux.Handle("/process/", h)
	mux.Handle("/publish", h)
	mux.Handle("/publish/", h)
	mux.Handle("/openapi", openapi.Handler(docs, true))
	mux.Handle("/", ui.Handler("", docs, appConfig.OpenApi))
	mux.HandleFunc("/authz/", authzHandler)

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
