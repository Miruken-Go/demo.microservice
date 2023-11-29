package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Factory struct {
}

func (f *Factory) Client() *mongo.Client {
	//opts := options.Client().ApplyURI("mongodb://localhost:27017")
	//client, err := mongo.Connect(context.Background(), opts)

	return nil
}
