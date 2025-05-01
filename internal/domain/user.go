package domain

import "github.com/google/uuid"

// User represents a user in the domain layer.
// This is the business entity, free from database implementation details.
type User struct {
	ID           uuid.UUID
	IsSubscribed bool
	FirstName    string
}
