package handler

import (
	"github.com/kootic/koogo/internal/service"
)

type Handler struct {
	KooUserHandler KooUserHandler
}

func NewHandler(services *service.Services) *Handler {
	kooUserHandler := NewKooUserHandler(services.KooUserService)

	return &Handler{
		KooUserHandler: kooUserHandler,
	}
}
