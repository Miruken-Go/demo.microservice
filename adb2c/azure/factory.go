package azure

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/config"
	"github.com/miruken-go/miruken/provides"
	"golang.org/x/net/context"
	"reflect"
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
	_*struct{
		args.Optional
		args.FromOptions
	  }, options Options,
) {
	f.opts = options
}

func (f *Factory) NewClient(
	_*struct{
		provides.Single `mode:"covariant"`
	  }, p *provides.It,
	_*struct{args.Optional}, ctx context.Context,
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
	return newClient(cfg, ctx)
}


func newClient(
	cfg Config,
	ctx context.Context,
) (*azcosmos.Client, error) {
	if uri := cfg.ConnectionUri; uri == "" {
		return nil, nil
	} else {
		return azcosmos.NewClientFromConnectionString(uri, nil)
	}
}

var (
	ClientType = reflect.TypeOf((*azcosmos.Client)(nil))
)