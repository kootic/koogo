package postgres

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/kootic/koogo/internal/repo"
)

// NewRepositories creates all PostgreSQL repository implementations.
func NewRepositories(sqlDB *sql.DB) (*repo.Repositories, error) {
	db := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	return &repo.Repositories{
		DB:     db,
		User:   NewKooUserRepository(db),
		Pet:    NewKooPetRepository(db),
		Health: NewHealthRepository(db),
	}, nil
}
