package service

import (
	"context"

	"github.com/kootic/koogo/internal/repository/dbrepo"
)

type HealthService interface {
	HealthCheck(ctx context.Context) error
}

type healthService struct {
	dbRepo dbrepo.DatabaseRepository
}

func NewHealthService(dbRepo dbrepo.DatabaseRepository) HealthService {
	return &healthService{dbRepo: dbRepo}
}

func (s *healthService) HealthCheck(ctx context.Context) error {
	return s.dbRepo.Ping(ctx)
}
