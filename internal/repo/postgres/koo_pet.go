package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/kootic/koogo/internal/domain"
	"github.com/kootic/koogo/internal/repo"
	bun1 "github.com/kootic/koogo/internal/repo/postgres/bun"
)

type petRepository struct {
	db *bun.DB
}

var _ repo.KooPetRepository = (*petRepository)(nil)

func NewKooPetRepository(db *bun.DB) repo.KooPetRepository {
	return &petRepository{db: db}
}

func (r *petRepository) Create(ctx context.Context, pet *domain.KooPet) (*domain.KooPet, error) {
	pgPet := bun1.KooPetFromDomain(pet)

	_, err := r.db.
		NewInsert().
		Model(pgPet).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return pgPet.ToDomain(), nil
}

func (r *petRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.KooPet, error) {
	var pgPet bun1.KooPet

	err := r.db.
		NewSelect().
		Model(&pgPet).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return pgPet.ToDomain(), nil
}

func (r *petRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.KooPet, error) {
	var pgPet bun1.KooPet

	err := r.db.
		NewSelect().
		Model(&pgPet).
		Where("owner_id = ?", ownerID).
		Scan(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return pgPet.ToDomain(), nil
}

func (r *petRepository) ListByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.KooPet, error) {
	var pgPets []*bun1.KooPet

	err := r.db.
		NewSelect().
		Model(&pgPets).
		Where("owner_id = ?", ownerID).
		Scan(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	pets := make([]*domain.KooPet, len(pgPets))
	for i, pgPet := range pgPets {
		pets[i] = pgPet.ToDomain()
	}

	return pets, nil
}

func (r *petRepository) Update(ctx context.Context, pet *domain.KooPet) error {
	pgPet := bun1.KooPetFromDomain(pet)

	_, err := r.db.
		NewUpdate().
		Model(pgPet).
		Where("id = ?", pgPet.ID).
		Exec(ctx)

	return handleError(err)
}

func (r *petRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.
		NewDelete().
		Model((*bun1.KooPet)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return handleError(err)
}
