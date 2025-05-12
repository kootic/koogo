package handler

import (
	"github.com/kootic/koogo/internal/service"
)

type Handler struct {
	HealthHandler  HealthHandler
	KooUserHandler KooUserHandler
}

func NewHandler(services *service.Services) *Handler {
	healthHandler := NewHealthHandler(services)
	kooUserHandler := NewKooUserHandler(services)

	return &Handler{
		HealthHandler:  healthHandler,
		KooUserHandler: kooUserHandler,
	}
}
