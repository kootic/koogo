package dto

import (
	"errors"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/repository/dbrepo/dbschema"
)

type KooCreateUserRequest struct {
	FirstName string `json:"firstName"`
}

func (r *KooCreateUserRequest) Validate() error {
	if r.FirstName == "" {
		return errors.New("firstName is required")
	}

	return nil
}

func (r *KooCreateUserRequest) ToSchema() *dbschema.KooUser {
	return &dbschema.KooUser{
		ID:        uuid.New(),
		FirstName: r.FirstName,
	}
}

type KooUserResponse struct {
	ID           uuid.UUID `json:"id"`
	IsSubscribed bool      `json:"isSubscribed"`
	FirstName    string    `json:"firstName"`
}

func (k *KooUserResponse) FromSchema(schema *dbschema.KooUser) {
	k.ID = schema.ID
	k.IsSubscribed = schema.IsSubscribed
	k.FirstName = schema.FirstName
}

type KooPetResponse struct {
	ID      uuid.UUID `json:"id"`
	OwnerID uuid.UUID `json:"ownerId"`
}

func (k *KooPetResponse) FromSchema(schema *dbschema.KooPet) {
	k.ID = schema.ID
	k.OwnerID = schema.OwnerID
}
