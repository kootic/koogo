package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/domain"
)

type KooPetRepository interface {
	Create(ctx context.Context, pet *domain.Pet) (*domain.Pet, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Pet, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Pet, error)
	ListByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Pet, error)
	Update(ctx context.Context, pet *domain.Pet) error
	Delete(ctx context.Context, id uuid.UUID) error
}
