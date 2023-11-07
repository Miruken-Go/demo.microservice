package main

import (
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-logr/zerologr"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
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
	"github.com/miruken-go/miruken/security/jwt/jwks"
	play "github.com/miruken-go/miruken/validates/play"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// logging
	zl := zerolog.New(os.Stderr)
	zl = zl.With().Timestamp().Logger()
	logger := zerologr.New(&zl)

	// configuration
	var k = koanf.New(".")
	err := k.Load(file.Provider("cmd/app.yml"), yaml.Parser())
	if err != nil {
		logger.Error(err, "error loading configuration")
		os.Exit(1)
	}

	// openapi generator
	openapiGen := openapi.Feature(openapi3.T{
		Info: &openapi3.Info{
			Title:       "Team Api",
			Description: "REST Api for managing Teams",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				Name: "Miruken",
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
								AuthorizationURL: "https://teamsrvidentitydev.b2clogin.com/teamsrvidentitydev.onmicrosoft.com/oauth2/v2.0/authorize?p=b2c_1a_signup_signin",
								TokenURL:         "https://teamsrvidentitydev.b2clogin.com/teamsrvidentitydev.onmicrosoft.com/oauth2/v2.0/token?p=b2c_1a_signup_signin",
								Scopes: map[string]string{
									"https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Groups": "Groups to which the user belongs.",
									"https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Roles":  "Roles to which the user belongs.",
									"https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Entitlements":  "Entitlements the user has.",
								},
							},
						},
						OpenIdConnectUrl: "https://teamsrvidentitydev.b2clogin.com/teamsrvidentitydev.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=B2C_1A_SIGNUP_SIGNIN",
					},
				},
			},
		},
		Security: openapi3.SecurityRequirements{
			{"team_auth": []string{
				"https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Groups",
				"https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Roles",
				"https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Roles",
			}},
		},
	})

	// initialize miruken
	handler, err := miruken.Setup(
		team.Feature, jwt.Feature(), jwks.Feature(),
		auth.Feature(), play.Feature(),
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

	h := httpsrv.Pipeline(handler,
		auth.WithFlowRef("login.oauth").Bearer(),
	)

	// configure routes
	var mux http.ServeMux
	mux.Handle("/process", h)
	mux.Handle("/process/", h)
	mux.Handle("/publish", h)
	mux.Handle("/publish/", h)
	mux.Handle("/openapi", openapi.Handler(docs, true))
	mux.Handle("/", ui.Handler("", docs))

	// Register pprof handlers
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// start http server
	server := &http.Server{
		Addr:   ":8080",
		Handler: &mux,
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err, "HTTP server error")
			os.Exit(1)
		}
		logger.Info("Stopped serving new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(err, "HTTP shutdown error")
		_ = server.Close()
		os.Exit(1)
	}
	logger.Info("Graceful shutdown complete")
}
