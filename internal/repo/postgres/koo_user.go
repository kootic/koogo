package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/kootic/koogo/internal/domain"
	"github.com/kootic/koogo/internal/repo"
	bun1 "github.com/kootic/koogo/internal/repo/postgres/bun"
)

type userRepository struct {
	db *bun.DB
}

// Ensure interface compliance at compile time.
var _ repo.KooUserRepository = (*userRepository)(nil)

func NewKooUserRepository(db *bun.DB) repo.KooUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.KooUser) (*domain.KooUser, error) {
	pgUser := bun1.KooUserFromDomain(user)

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

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.KooUser, error) {
	var pgUser bun1.KooUser

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

func (r *userRepository) Update(ctx context.Context, user *domain.KooUser) error {
	pgUser := bun1.KooUserFromDomain(user)

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
		Model((*bun1.KooUser)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return handleError(err)
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*domain.KooUser, error) {
	var pgUsers []*bun1.KooUser

	err := r.db.
		NewSelect().
		Model(&pgUsers).
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	users := make([]*domain.KooUser, len(pgUsers))
	for i, pgUser := range pgUsers {
		users[i] = pgUser.ToDomain()
	}

	return users, nil
}
