package dbschema

import "github.com/google/uuid"

type KooPet struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
}
