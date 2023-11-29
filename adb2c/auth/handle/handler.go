package handle

import "go.mongodb.org/mongo-driver/mongo"

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		Database *mongo.Database
	}
)
