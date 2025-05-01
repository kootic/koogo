package pgmodel

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/kootic/koogo/internal/domain"
)

type KooPet struct {
	bun.BaseModel `bun:"table:koo_pets,alias:p"`

	ID      uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	OwnerID uuid.UUID `bun:"owner_id,notnull,type:uuid"`

	// Relations
	Owner *KooUser `bun:"rel:belongs-to,join:owner_id=id"`
}

// ToDomain converts the database model to a domain model.
func (p *KooPet) ToDomain() *domain.Pet {
	if p == nil {
		return nil
	}

	pet := &domain.Pet{
		ID:      p.ID,
		OwnerID: p.OwnerID,
	}
	if p.Owner != nil {
		pet.Owner = p.Owner.ToDomain()
	}

	return pet
}

// KooPetFromDomain converts a domain model to a database model.
func KooPetFromDomain(pet *domain.Pet) *KooPet {
	if pet == nil {
		return nil
	}

	pgPet := &KooPet{
		ID:      pet.ID,
		OwnerID: pet.OwnerID,
	}
	if pet.Owner != nil {
		pgPet.Owner = KooUserFromDomain(pet.Owner)
	}

	return pgPet
}
