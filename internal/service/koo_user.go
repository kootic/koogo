package service

// BOILERPLATE: This file demonstrates the service pattern.
// Delete this file when bootstrapping a new project.
// See docs/BOOTSTRAPPING.md for details.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/kootic/koogo/internal/dto"
	"github.com/kootic/koogo/internal/repo"
	"github.com/kootic/koogo/pkg/koohttp"
)

var (
	ErrUserNotFound        = koohttp.NewAPIError(http.StatusNotFound, "user_not_found")
	ErrUserIsNotSubscribed = koohttp.NewAPIError(http.StatusForbidden, "user_is_not_subscribed")
)

type KooUserService interface {
	KooCreateUser(ctx context.Context, req *dto.KooCreateUserRequest) (*dto.KooUserResponse, error)
	KooGetUserByID(ctx context.Context, id uuid.UUID) (*dto.KooUserResponse, error)
	KooGetPetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*dto.KooPetResponse, error)
}

type userService struct {
	userRepo repo.KooUserRepository
	petRepo  repo.KooPetRepository
}

func NewKooUserService(userRepo repo.KooUserRepository, petRepo repo.KooPetRepository) KooUserService {
	return &userService{
		userRepo: userRepo,
		petRepo:  petRepo,
	}
}

func (s *userService) KooCreateUser(ctx context.Context, req *dto.KooCreateUserRequest) (*dto.KooUserResponse, error) {
	newUser := req.ToModel()

	createdUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	var response dto.KooUserResponse
	response.FromModel(createdUser)

	return &response, nil
}

func (s *userService) KooGetUserByID(ctx context.Context, id uuid.UUID) (*dto.KooUserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	var response dto.KooUserResponse
	response.FromModel(user)

	return &response, nil
}

func (s *userService) KooGetPetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*dto.KooPetResponse, error) {
	// Check if user exists and is subscribed
	user, err := s.userRepo.GetByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	if !user.IsSubscribed {
		return nil, ErrUserIsNotSubscribed
	}

	// Get pet
	pet, err := s.petRepo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	var response dto.KooPetResponse
	response.FromModel(pet)

	return &response, nil
}
