package service

import "github.com/kootic/koogo/internal/repo/postgres"

type Services struct {
	HealthService  HealthService
	KooUserService KooUserService
}

func NewServices(repos *postgres.Repositories) *Services {
	return &Services{
		HealthService:  NewHealthService(repos.Health),
		KooUserService: NewKooUserService(repos.User, repos.Pet),
	}
}
