package postgres

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/kootic/koogo/internal/repo"
)

type healthRepository struct {
	db *bun.DB
}

var _ repository.HealthRepository = (*healthRepository)(nil)

func NewHealthRepository(db *bun.DB) repository.HealthRepository {
	return &healthRepository{db: db}
}

func (r *healthRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
