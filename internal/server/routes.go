package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const (
	APIBasePath = "/api"
)

type route struct {
	Method     string
	Version    int
	Path       string
	Middleware []fiber.Handler
	Handler    fiber.Handler
}

func (s *server) allRoutes() []route {
	return []route{
		{
			Method:  http.MethodPost,
			Version: 1,
			Path:    "/koo/users",
			Handler: s.handler.KooUserHandler.KooCreateUser,
		},
		{
			Method:  http.MethodGet,
			Version: 1,
			Path:    "/koo/users/:userId",
			Handler: s.handler.KooUserHandler.KooGetUserByID,
		},
		{
			Method:  http.MethodGet,
			Version: 1,
			Path:    "/koo/users/:userId/pet",
			Handler: s.handler.KooUserHandler.KooGetUserPet,
		},
	}
}
