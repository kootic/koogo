package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/kootic/koogo/internal/domain"
	"github.com/kootic/koogo/internal/repo"
	"github.com/kootic/koogo/internal/repo/postgres/pgmodel"
)

type userRepository struct {
	db *bun.DB
}

// Ensure interface compliance at compile time.
var _ repository.KooUserRepository = (*userRepository)(nil)

func NewUserRepository(db *bun.DB) repository.KooUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	pgUser := pgmodel.KooUserFromDomain(user)

	_, err := r.db.
		NewInsert().
		Model(pgUser).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return pgUser.ToDomain(), nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var pgUser pgmodel.KooUser

	err := r.db.
		NewSelect().
		Model(&pgUser).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return pgUser.ToDomain(), nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	pgUser := pgmodel.KooUserFromDomain(user)

	_, err := r.db.
		NewUpdate().
		Model(pgUser).
		Where("id = ?", pgUser.ID).
		Exec(ctx)

	return handleError(err)
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.
		NewDelete().
		Model((*pgmodel.KooUser)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return handleError(err)
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var pgUsers []*pgmodel.KooUser

	err := r.db.
		NewSelect().
		Model(&pgUsers).
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	users := make([]*domain.User, len(pgUsers))
	for i, pgUser := range pgUsers {
		users[i] = pgUser.ToDomain()
	}

	return users, nil
}
