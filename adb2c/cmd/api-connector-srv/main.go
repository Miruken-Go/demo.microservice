package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/go-logr/zerologr"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/miruken-go/demo.microservice/adb2c/enrich"
	"github.com/miruken-go/miruken/api/http/httpsrv"
	"github.com/miruken-go/miruken/api/http/httpsrv/auth"
	"github.com/miruken-go/miruken/config"
	koanfp "github.com/miruken-go/miruken/config/koanf"
	"github.com/miruken-go/miruken/logs"
	"github.com/miruken-go/miruken/setup"
	"github.com/rs/zerolog"
)

func main() {
	// logging
	zl := zerolog.New(os.Stdout)
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

	// initialize context
	ctx, err := setup.New(
		enrich.Feature(),
		config.Feature(koanfp.P(k)),
		logs.Feature(logger),
	).Context()

	if err != nil {
		logger.Error(err, "setup failed")
		os.Exit(1)
	}

	defer ctx.End(nil)

	var mux http.ServeMux
	mux.Handle("/enrich", httpsrv.Use(ctx,
		httpsrv.H[*enrich.Handler](),
		auth.WithFlowAlias("Login.Adb2c").Basic().Required()))

	if err := httpsrv.ListenAndServe(&mux, nil); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("server closed")
		} else if err != nil {
			logger.Error(err, "error starting server")
		}
	}
}
