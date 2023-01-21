package main

import (
	"errors"
	"github.com/Rican7/conjson/transform"
	"github.com/go-logr/zerologr"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/miruken-go/demo.microservice/teamapi"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/api/http/httpsrv"
	"github.com/miruken-go/miruken/api/json"
	"github.com/miruken-go/miruken/config"
	koanfp "github.com/miruken-go/miruken/config/koanf"
	"github.com/miruken-go/miruken/log"
	"github.com/miruken-go/miruken/validate/go"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

func main() {
	// logging
	zl := zerolog.New(os.Stderr)
	zl = zl.With().Caller().Timestamp().Logger()
	logger := zerologr.New(&zl)

	// configuration
	var k = koanf.New(".")
	err := k.Load(env.Provider("", ".", nil), nil)
	if err != nil {
		logger.Error(err, "error loading configuration")
		os.Exit(1)
	}

	// initialize miruken
	ctx, err := miruken.SetupContext(
		teamapi.Feature,
		httpsrv.Feature(),
		govalidator.Feature(),
		config.Feature(koanfp.P(k)),
		log.Feature(logger),
		miruken.Builders(
			json.StdTransform(transform.OnlyForDirection(
				transform.Marshal,
				transform.CamelCaseKeys(false))),
		))

	if err != nil {
		logger.Error(err, "setup failed")
		os.Exit(1)
	}

	// start http server
	err = http.ListenAndServe(":8080", httpsrv.NewController(ctx))

	if errors.Is(err, http.ErrServerClosed) {
		logger.Info("server closed")
	} else if err != nil {
		logger.Error(err, "error starting server")
		os.Exit(1)
	}
}
