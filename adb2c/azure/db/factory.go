package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	_ "github.com/btnguyen2k/gocosmos"
	"github.com/go-logr/logr"
	"github.com/jmoiron/sqlx"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/config"
	"github.com/miruken-go/miruken/promise"
	"github.com/miruken-go/miruken/provides"
)

type (
	// Options for azure cosmosdb resources.
	Options struct {
		Aliases   map[reflect.Type]string
		Clients   map[reflect.Type]*Config
		Provision []reflect.Type
	}

	// Factory of azure cosmosdb resources.
	Factory struct {
		opts   Options
		logger logr.Logger
	}
)


func (f *Factory) Constructor(
	_ *struct {
		args.Optional
		args.FromOptions
	  }, options Options,
	_ *struct{ args.Optional }, logger logr.Logger,
) {
	f.opts = options
	if logger == f.logger {
		f.logger = logr.Discard()
	} else {
		f.logger = logger
	}
}

func (f *Factory) NewAzClient(
	_ *struct {
		provides.Single `mode:"covariant"`
	  }, p *provides.It,
	hc miruken.HandleContext,
) (*azcosmos.Client, error) {
	typ := p.Key().(reflect.Type)
	if cfg, err := f.config(typ, clientType, hc); err != nil {
		return nil, err
	} else {
		return newAzClient(cfg, false)
	}
}

func (f *Factory) NewSqlxClient(
	_ *struct {
		provides.Single `mode:"covariant"`
	  }, p *provides.It,
	hc miruken.HandleContext,
) (*sqlx.DB, error) {
	typ := p.Key().(reflect.Type)
	if cfg, err := f.config(typ, sqlxDbType, hc); err != nil {
		return nil, err
	} else {
		return newSqlxClient(cfg, false)
	}
}

func (f *Factory) Startup(
	ctx context.Context,
	h   miruken.Handler,
) *promise.Promise[struct{}] {
	provisions := f.opts.Provision
	if len(provisions) == 0 {
		return promise.Empty()
	}
	promises := make([]*promise.Promise[struct{}], len(provisions))
	for i, provision := range provisions {
		promises[i] = f.provision(ctx, provision, h)
	}
	return promise.Erase(promise.All(ctx, promises...))
}


func (f *Factory) Shutdown(
	ctx context.Context,
) *promise.Promise[struct{}] {
	return promise.Empty()
}

func (f *Factory) config(
	typ, defTyp reflect.Type,
	h    miruken.Handler,
) (*Config, error) {
	cfg, ok := f.opts.Clients[typ]
	if ok {
		return cfg, nil
	}
	var key string
	if key, ok = f.opts.Aliases[typ]; !ok {
		if typ == defTyp {
			key = "Azure"
		} else {
			key = typ.Name()
		}
	}
	path := fmt.Sprintf("Databases.%s", key)
	cfg, _, ok, err := provides.Type[*Config](h, &config.Load{Path: path})
	if !ok || err != nil {
		return nil, err
	}
	return cfg, nil
}

func (f *Factory) provision(
	ctx context.Context,
	typ reflect.Type,
	h   miruken.Handler,
) *promise.Promise[struct{}] {
	cfg, err := f.config(typ, clientType, h)
	if err != nil {
		return promise.RejectEmpty(err)
	}
	client, err := newAzClient(cfg, true)
	if err != nil {
		return promise.RejectEmpty(err)
	}
	return ProvisionDatabase(ctx, client, cfg, f.logger)
}


func newAzClient(
	cfg     *Config,
	require bool,
) (*azcosmos.Client, error) {
	if uri := cfg.ConnectionUri; uri == "" {
		if require {
			return nil, errors.New("missing connection uri")
		}
		return nil, nil
	} else {
		return azcosmos.NewClientFromConnectionString(uri, nil)
	}
}

func newSqlxClient(
	cfg     *Config,
	require bool,
) (*sqlx.DB, error) {
	if uri := cfg.ConnectionUri; uri == "" {
		if require {
			return nil, errors.New("missing connection uri")
		}
		return nil, nil
	} else {
		return sqlx.Open("gocosmos", uri)
	}
}


var (
	clientType = reflect.TypeOf((*azcosmos.Client)(nil))
	sqlxDbType = reflect.TypeOf((*sqlx.DB)(nil))
)
