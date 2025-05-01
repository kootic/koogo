package dto

// BOILERPLATE: This file demonstrates DTOs and validation.
// Delete this file when bootstrapping a new project.
// See docs/BOOTSTRAPPING.md for details.

import (
	"errors"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/domain"
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

func (r *KooCreateUserRequest) ToModel() *domain.KooUser {
	return &domain.KooUser{
		ID:        uuid.New(),
		FirstName: r.FirstName,
	}
}

type KooUserResponse struct {
	ID           uuid.UUID `json:"id"`
	IsSubscribed bool      `json:"isSubscribed"`
	FirstName    string    `json:"firstName"`
}

func (k *KooUserResponse) FromModel(m *domain.KooUser) {
	k.ID = m.ID
	k.IsSubscribed = m.IsSubscribed
	k.FirstName = m.FirstName
}

type KooPetResponse struct {
	ID      uuid.UUID `json:"id"`
	OwnerID uuid.UUID `json:"ownerId"`
}

func (k *KooPetResponse) FromModel(m *domain.KooPet) {
	k.ID = m.ID
	k.OwnerID = m.OwnerID
}
