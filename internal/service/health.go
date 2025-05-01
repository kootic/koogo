package service

import (
	"context"

	"github.com/kootic/koogo/internal/repo"
)

type HealthService interface {
	HealthCheck(ctx context.Context) error
}

type healthService struct {
	healthRepo repo.HealthRepository
}

func NewHealthService(healthRepo repo.HealthRepository) HealthService {
	return &healthService{healthRepo: healthRepo}
}

func (s *healthService) HealthCheck(ctx context.Context) error {
	return s.healthRepo.Ping(ctx)
}
