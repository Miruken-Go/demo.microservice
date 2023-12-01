package mongo

import (
	"fmt"
	"github.com/miruken-go/miruken"
	"github.com/miruken-go/miruken/config"
	"github.com/miruken-go/miruken/provides"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"strings"
)

// Factory of mongo clients.
type Factory struct {}

func (f *Factory) DefaultClient(
	p   *provides.It,
	ctx miruken.HandleContext,
) (*mongo.Client, error) {
	key := fmt.Sprintf("%v", p.Key())
	if key == "*mongo.Client" {
		key = "Mongo"
	} else {
		key = key[strings.LastIndex(key, ".")+1:]
	}
	path := fmt.Sprintf("Databases.%s", key)
	cfg, _, ok, err := provides.Type[Config](ctx, &config.Load{Path: path})
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
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