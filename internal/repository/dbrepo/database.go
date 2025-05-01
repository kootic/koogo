package dbrepo

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type DatabaseRepository interface {
	KooUser
	Close() error
	Ping(ctx context.Context) error
}

type databaseRepository struct {
	db *bun.DB
}

func NewDatabaseRepository(sqlDB *sql.DB) (DatabaseRepository, error) {
	db := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	return &databaseRepository{db: db}, nil
}

func (d *databaseRepository) Close() error {
	return d.db.Close()
}

func (d *databaseRepository) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
