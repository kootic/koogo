package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/kootic/koogo/internal/domain"
	"github.com/kootic/koogo/internal/repo"
	"github.com/kootic/koogo/internal/repo/postgres/pgmodel"
)

type petRepository struct {
	db *bun.DB
}

var _ repository.KooPetRepository = (*petRepository)(nil)

func NewPetRepository(db *bun.DB) repository.KooPetRepository {
	return &petRepository{db: db}
}

func (r *petRepository) Create(ctx context.Context, pet *domain.Pet) (*domain.Pet, error) {
	pgPet := pgmodel.KooPetFromDomain(pet)

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

func (r *petRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Pet, error) {
	var pgPet pgmodel.KooPet

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

func (r *petRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Pet, error) {
	var pgPet pgmodel.KooPet

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

func (r *petRepository) ListByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Pet, error) {
	var pgPets []*pgmodel.KooPet

	err := r.db.
		NewSelect().
		Model(&pgPets).
		Where("owner_id = ?", ownerID).
		Scan(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	pets := make([]*domain.Pet, len(pgPets))
	for i, pgPet := range pgPets {
		pets[i] = pgPet.ToDomain()
	}

	return pets, nil
}

func (r *petRepository) Update(ctx context.Context, pet *domain.Pet) error {
	pgPet := pgmodel.KooPetFromDomain(pet)

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
		Model((*pgmodel.KooPet)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return handleError(err)
}
