package dbrepo

import (
	"context"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/repository/dbrepo/dbschema"
)

type KooUser interface {
	KooCreateUser(ctx context.Context, kooUser *dbschema.KooUser) error
	KooGetUserByID(ctx context.Context, id uuid.UUID) (*dbschema.KooUser, error)
	KooGetPetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*dbschema.KooPet, error)
}

var _ KooUser = (*databaseRepository)(nil)

func (d *databaseRepository) KooCreateUser(ctx context.Context, kooUser *dbschema.KooUser) error {
	_, err := d.db.
		NewInsert().
		Model(kooUser).
		Exec(ctx)
	if err != nil {
		return bunErrorHandler(err)
	}

	return nil
}

func (d *databaseRepository) KooGetUserByID(ctx context.Context, id uuid.UUID) (*dbschema.KooUser, error) {
	var kooUser dbschema.KooUser

	err := d.db.
		NewSelect().
		Model(&kooUser).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, bunErrorHandler(err)
	}

	return &kooUser, nil
}

func (d *databaseRepository) KooGetPetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*dbschema.KooPet, error) {
	var kooPet dbschema.KooPet

	err := d.db.
		NewSelect().
		Model(&kooPet).
		Where("owner_id = ?", ownerID).
		Scan(ctx)
	if err != nil {
		return nil, bunErrorHandler(err)
	}

	return &kooPet, nil
}
