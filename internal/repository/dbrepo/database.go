package dbrepo

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type DatabaseRepository interface {
	KooUser
}

type databaseRepository struct {
	db *bun.DB
}

func NewDatabaseRepository(sqlDB *sql.DB) (DatabaseRepository, error) {
	db := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	return &databaseRepository{db: db}, nil
}

func (d *databaseRepository) Close() error {
	return nil
}
