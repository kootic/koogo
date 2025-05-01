package service

import "github.com/kootic/koogo/internal/repository/dbrepo"

type Services struct {
	HealthService  HealthService
	KooUserService KooUserService
}

func NewServices(dbRepo dbrepo.DatabaseRepository) *Services {
	return &Services{
		HealthService:  NewHealthService(dbRepo),
		KooUserService: NewKooUserService(dbRepo),
	}
}
