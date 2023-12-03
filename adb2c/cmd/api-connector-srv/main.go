package main

import (
	"errors"
	"github.com/go-logr/zerologr"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/miruken-go/demo.microservice/adb2c/mongo"
	"github.com/miruken-go/demo.microservice/adb2c/token"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/api/http/httpsrv"
	"github.com/miruken-go/miruken/api/http/httpsrv/auth"
	"github.com/miruken-go/miruken/config"
	koanfp "github.com/miruken-go/miruken/config/koanf"
	"github.com/miruken-go/miruken/context"
	"github.com/miruken-go/miruken/logs"
	"github.com/rs/zerolog"
	"net/http"
	"os"
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

	handler, err := miruken.Setup(
		token.Feature(),
		config.Feature(koanfp.P(k)),
		logs.Feature(logger),
		mongo.Feature(),
	).Handler()

	if err != nil {
		logger.Error(err, "setup failed")
		os.Exit(1)
	}

	ctx := context.New(handler)
	defer ctx.End(nil)

	http.Handle("/enrich", httpsrv.Use(ctx,
		httpsrv.H[*token.EnrichHandler](),
		auth.WithFlowRef("Login.Adb2c").Basic().Required()))

	// start http server
	port := k.String("App.Port")
	if port == "" {
		port = "8080"
	}
	err = http.ListenAndServe(":"+port, nil)

	if errors.Is(err, http.ErrServerClosed) {
		logger.Info("server closed")
	} else if err != nil {
		logger.Error(err, "error starting server")
		os.Exit(1)
	}
}
