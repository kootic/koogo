package domain

import "github.com/google/uuid"

// KooUser represents a user in the domain layer.
// This is the business entity, free from database implementation details.
type KooUser struct {
	ID           uuid.UUID
	IsSubscribed bool
	FirstName    string
}
