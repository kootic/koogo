package bun

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/kootic/koogo/internal/domain"
)

type KooUser struct {
	bun.BaseModel `bun:"table:koo_users,alias:u"`

	ID           uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	IsSubscribed bool      `bun:"is_subscribed,notnull,default:false"`
	FirstName    string    `bun:"first_name,notnull"`
}

// ToDomain converts the database model to a domain model.
func (u *KooUser) ToDomain() *domain.KooUser {
	if u == nil {
		return nil
	}

	return &domain.KooUser{
		ID:           u.ID,
		IsSubscribed: u.IsSubscribed,
		FirstName:    u.FirstName,
	}
}

// KooUserFromDomain converts a domain model to a database model.
func KooUserFromDomain(user *domain.KooUser) *KooUser {
	if user == nil {
		return nil
	}

	return &KooUser{
		ID:           user.ID,
		IsSubscribed: user.IsSubscribed,
		FirstName:    user.FirstName,
	}
}
