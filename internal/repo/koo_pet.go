package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/domain"
)

type KooPetRepository interface {
	Create(ctx context.Context, pet *domain.KooPet) (*domain.KooPet, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.KooPet, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.KooPet, error)
	ListByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.KooPet, error)
	Update(ctx context.Context, pet *domain.KooPet) error
	Delete(ctx context.Context, id uuid.UUID) error
}
