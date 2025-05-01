package service

import "github.com/kootic/koogo/internal/repository/dbrepo"

type Services struct {
	KooUserService KooUserService
}

func NewServices(dbRepo dbrepo.DatabaseRepository) *Services {
	return &Services{
		KooUserService: NewKooUserService(dbRepo),
	}
}
