package domain

import "github.com/google/uuid"

// KooPet represents a pet in the domain layer.
// This is the business entity, free from database implementation details.
type KooPet struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
	Owner   *KooUser // Optional relation
}
