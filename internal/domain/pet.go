package domain

import "github.com/google/uuid"

// Pet represents a pet in the domain layer.
// This is the business entity, free from database implementation details.
type Pet struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
	Owner   *User // Optional relation
}
