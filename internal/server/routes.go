package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const (
	APIBasePath = "/api"
)

type route struct {
	Version    int
	Method     string
	Path       string
	Middleware []fiber.Handler
	Handler    fiber.Handler
}

func (s *server) allRoutes() []route {
	return []route{
		{
			Version: 1,
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: s.handler.HealthHandler.HealthCheck,
		},
		{
			Version: 1,
			Method:  http.MethodPost,
			Path:    "/koo/users",
			Handler: s.handler.KooUserHandler.CreateUser,
		},
		{
			Version: 1,
			Method:  http.MethodGet,
			Path:    "/koo/users/:userId",
			Handler: s.handler.KooUserHandler.GetUserByID,
		},
		{
			Version: 1,
			Method:  http.MethodGet,
			Path:    "/koo/users/:userId/pet",
			Handler: s.handler.KooUserHandler.GetUserPet,
		},
	}
}
