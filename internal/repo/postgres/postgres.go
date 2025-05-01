package postgres

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/kootic/koogo/internal/repo"
)

// Repositories holds all repository implementations.
type Repositories struct {
	User   repository.KooUserRepository
	Pet    repository.KooPetRepository
	Health repository.HealthRepository

	db *bun.DB
}

// NewRepositories creates all PostgreSQL repository implementations.
func NewRepositories(sqlDB *sql.DB) (*Repositories, error) {
	db := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	return &Repositories{
		User:   NewUserRepository(db),
		Pet:    NewPetRepository(db),
		Health: NewHealthRepository(db),
		db:     db,
	}, nil
}

// Close closes the database connection.
func (r *Repositories) Close() error {
	return r.db.Close()
}
