package mongo

import (
	"fmt"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/args"
	"github.com/miruken-go/miruken/config"
	"github.com/miruken-go/miruken/provides"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"reflect"
)

type (
	// Options for mongo clients.
	Options struct {
		Aliases map[reflect.Type]string
		Clients map[reflect.Type]Config
	}

	// Factory of mongo clients.
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

func (f *Factory) DefaultClient(
	_*struct{
		provides.Single `mode:"covariant"`
	  }, p *provides.It,
	ctx miruken.HandleContext,
) (client *mongo.Client, err error) {
	typ := p.Key().(reflect.Type)
	cfg, ok := f.opts.Clients[typ]
	if !ok {
		var key string
		if key, ok = f.opts.Aliases[typ]; !ok {
			if typ == ClientType {
				key = "Mongo"
			} else {
				key = typ.Name()
			}
		}
		path := fmt.Sprintf("Databases.%s", key)
		cfg, _, ok, err = provides.Type[Config](ctx, &config.Load{Path: path})
		if !ok || err != nil {
			return
		}
	}
	return newClient(cfg)
}


func newClient(cfg Config) (*mongo.Client, error) {
	if uri := cfg.ConnectionUri; uri == "" {
		return nil, nil
	} else {
		opts := options.Client()
		if timeout := cfg.Timeout; timeout > 0 {
			opts.SetTimeout(timeout)
		}
		opts.ApplyURI(uri)
		return mongo.Connect(context.Background(), opts)
	}
}

var (
	ClientType = reflect.TypeOf((*mongo.Client)(nil))
)