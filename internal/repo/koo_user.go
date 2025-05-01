package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/domain"
)

type KooUserRepository interface {
	Create(ctx context.Context, user *domain.KooUser) (*domain.KooUser, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.KooUser, error)
	Update(ctx context.Context, user *domain.KooUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.KooUser, error)
}
