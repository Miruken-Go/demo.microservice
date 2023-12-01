package handle

import "go.mongodb.org/mongo-driver/mongo"

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		database *mongo.Database
	}
)

func (h *Handler) Constructor(
	client *mongo.Client,
) {
	h.database = client.Database("adb2c")
}