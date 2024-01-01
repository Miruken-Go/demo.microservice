package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/miruken-go/miruken/promise"
)

type Bootstrap struct {
	db *sqlx.DB
}


func (b *Bootstrap) Constructor(
	db *sqlx.DB,
) {
	b.db = db
}

func (b *Bootstrap) Startup(
	ctx context.Context,
) *promise.Promise[struct{}] {
	return promise.Empty()
}


func (b *Bootstrap) Shutdown(
	ctx context.Context,
) *promise.Promise[struct{}] {
	return promise.Empty()
}
