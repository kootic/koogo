package service

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/dto"
	"github.com/kootic/koogo/internal/repository/dbrepo"
	"github.com/kootic/koogo/pkg/koohttp"
)

var (
	ErrUserNotFound        = koohttp.NewAPIError(http.StatusNotFound, "user_not_found")
	ErrUserIsNotSubscribed = koohttp.NewAPIError(http.StatusForbidden, "user_is_not_subscribed")
)

type KooUserService interface {
	KooCreateUser(ctx context.Context, kooUser *dto.KooCreateUserRequest) (*dto.KooUserResponse, error)
	KooGetUserByID(ctx context.Context, id uuid.UUID) (*dto.KooUserResponse, error)
	KooGetPetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*dto.KooPetResponse, error)
}

type kooUserService struct {
	dbRepo dbrepo.DatabaseRepository
}

func NewKooUserService(dbRepo dbrepo.DatabaseRepository) KooUserService {
	return &kooUserService{dbRepo: dbRepo}
}

func (s *kooUserService) KooCreateUser(ctx context.Context, kooUser *dto.KooCreateUserRequest) (*dto.KooUserResponse, error) {
	newKooUser := kooUser.ToSchema()

	err := s.dbRepo.KooCreateUser(ctx, newKooUser)
	if err != nil {
		return nil, err
	}

	return s.KooGetUserByID(ctx, newKooUser.ID)
}

func (s *kooUserService) KooGetUserByID(ctx context.Context, id uuid.UUID) (*dto.KooUserResponse, error) {
	kooUser, err := s.dbRepo.KooGetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if kooUser == nil {
		return nil, ErrUserNotFound
	}

	var dtoKooUser dto.KooUserResponse

	dtoKooUser.FromSchema(kooUser)

	return &dtoKooUser, nil
}

func (s *kooUserService) KooGetPetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*dto.KooPetResponse, error) {
	user, err := s.KooGetUserByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	if !user.IsSubscribed {
		return nil, ErrUserIsNotSubscribed
	}

	userPet, err := s.dbRepo.KooGetPetByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	var dtoKooPet dto.KooPetResponse

	dtoKooPet.FromSchema(userPet)

	return &dtoKooPet, nil
}
