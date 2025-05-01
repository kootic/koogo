package dbschema

import "github.com/google/uuid"

type KooUser struct {
	ID           uuid.UUID
	IsSubscribed bool
	FirstName    string
}
