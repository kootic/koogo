package service

import "github.com/kootic/koogo/internal/repo"

type Services struct {
	HealthService  HealthService
	KooUserService KooUserService
}

func NewServices(repos *repo.Repositories) *Services {
	return &Services{
		HealthService:  NewHealthService(repos.Health),
		KooUserService: NewKooUserService(repos.User, repos.Pet),
	}
}
