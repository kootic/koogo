package service

import (
	"context"

	"github.com/kootic/koogo/internal/repo"
)

type HealthService interface {
	HealthCheck(ctx context.Context) error
}

type healthService struct {
	healthRepo repository.HealthRepository
}

func NewHealthService(healthRepo repository.HealthRepository) HealthService {
	return &healthService{healthRepo: healthRepo}
}

func (s *healthService) HealthCheck(ctx context.Context) error {
	return s.healthRepo.Ping(ctx)
}
