package main

import (
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-logr/zerologr"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/miruken-go/demo.microservice/teamsrv"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/api/http/httpsrv"
	"github.com/miruken-go/miruken/api/http/httpsrv/openapi"
	"github.com/miruken-go/miruken/api/http/httpsrv/openapi/ui"
	"github.com/miruken-go/miruken/api/json/jsonstd"
	"github.com/miruken-go/miruken/config"
	koanfp "github.com/miruken-go/miruken/config/koanf"
	"github.com/miruken-go/miruken/context"
	"github.com/miruken-go/miruken/log"
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
	err := k.Load(env.Provider("", ".", nil), nil)
	if err != nil {
		logger.Error(err, "error loading configuration")
		os.Exit(1)
	}

	// openapi generator
	openapiGen := openapi.Feature()

	// initialize miruken
	handler, err := miruken.Setup(
		teamsrv.Feature, jsonstd.Feature(),
		play.Feature(), config.Feature(koanfp.P(k)),
		log.Feature(logger), openapiGen).
		Specs(&api.GoPolymorphismMapper{}).
		Options(jsonstd.CamelCase).
		Handler()

	if err != nil {
		logger.Error(err, "setup failed")
		os.Exit(1)
	}

	ctl := httpsrv.Api(context.New(handler))

	// configure routes
	var mux http.ServeMux
	mux.Handle("/process", ctl)
	mux.Handle("/process/", ctl)
	mux.Handle("/publish", ctl)
	mux.Handle("/publish/", ctl)

	// swagger ui
	api := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Team Api",
			Description: "REST Api used for managing Teams",
			Version:     "0.0.0",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				URL: "https://github.com/craig/team-microservice",
			},
		},
	}
	openapiGen.Merge(&api)
	mux.Handle("/swagger", openapi.Handler(&api, true))
	mux.Handle("/swagger_ui/", ui.Handler("/swagger_ui/", &api))

	// start http server
	err = http.ListenAndServe(":8080", &mux)

	if errors.Is(err, http.ErrServerClosed) {
		logger.Info("server closed")
	} else if err != nil {
		logger.Error(err, "error starting server")
		os.Exit(1)
	}
}
