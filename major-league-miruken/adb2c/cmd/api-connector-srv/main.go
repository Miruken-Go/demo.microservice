package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// start HTTP server
	srv := httpsrv.New(&mux, nil)

	// handle SIGINT (CTRL+C) gracefully.
	notify, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err, "Unable to start HTTP server")
			os.Exit(1)
		}
		logger.Info("Stopped serving new HTTP connections")
		return
	case <-notify.Done():
		stop()
	}

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(err, "Unable to stop HTTP server")
		_ = srv.Close()
	}
	logger.Info("Graceful shutdown complete")
}
