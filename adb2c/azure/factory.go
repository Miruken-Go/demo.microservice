package azure

import (
	"fmt"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	_ "github.com/btnguyen2k/gocosmos"
	"github.com/jmoiron/sqlx"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/config"
	"github.com/miruken-go/miruken/provides"
)

type (
	// Options for azure cosmosdb resources.
	Options struct {
		Aliases map[reflect.Type]string
		Clients map[reflect.Type]Config
	}

	// Factory of azure cosmosdb resources.
	Factory struct {
		opts Options
	}
)

func (f *Factory) Constructor(
	_ *struct {
		args.Optional
		args.FromOptions
	  }, options Options,
) {
	f.opts = options
}

func (f *Factory) NewClient(
	_ *struct {
		provides.Single `mode:"covariant"`
	  }, p *provides.It,
	hc miruken.HandleContext,
) (client *azcosmos.Client, err error) {
	typ := p.Key().(reflect.Type)
	cfg, ok := f.opts.Clients[typ]
	if !ok {
		var key string
		if key, ok = f.opts.Aliases[typ]; !ok {
			if typ == ClientType {
				key = "Azure"
			} else {
				key = typ.Name()
			}
		}
		path := fmt.Sprintf("Databases.%s", key)
		cfg, _, ok, err = provides.Type[Config](hc, &config.Load{Path: path})
		if !ok || err != nil {
			return
		}
	}
	return newClient(cfg)
}

func (f *Factory) NewSqlClient(
	_ *struct {
		provides.Single `mode:"covariant"`
	  }, p *provides.It,
	hc miruken.HandleContext,
) (db *sqlx.DB, err error) {
	typ := p.Key().(reflect.Type)
	cfg, ok := f.opts.Clients[typ]
	if !ok {
		var key string
		if key, ok = f.opts.Aliases[typ]; !ok {
			if typ == SqlxDbType {
				key = "Azure"
			} else {
				key = typ.Name()
			}
		}
		path := fmt.Sprintf("Databases.%s", key)
		cfg, _, ok, err = provides.Type[Config](hc, &config.Load{Path: path})
		if !ok || err != nil {
			return
		}
	}
	return newSqlxClient(cfg)
}

func newClient(
	cfg Config,
) (*azcosmos.Client, error) {
	if uri := cfg.ConnectionUri; uri == "" {
		return nil, nil
	} else {
		return azcosmos.NewClientFromConnectionString(uri, nil)
	}
}

func newSqlxClient(
	cfg Config,
) (*sqlx.DB, error) {
	if uri := cfg.ConnectionUri; uri == "" {
		return nil, nil
	} else {
		return sqlx.Open("gocosmos", uri)
	}
}

var (
	ClientType = reflect.TypeOf((*azcosmos.Client)(nil))
	SqlxDbType = reflect.TypeOf((*sqlx.DB)(nil))
)
