package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	"github.com/kootic/koogo/internal/config"
	"github.com/kootic/koogo/internal/handler"
	"github.com/kootic/koogo/internal/repo/postgres"
	"github.com/kootic/koogo/internal/server/middleware"
	"github.com/kootic/koogo/internal/service"
)

// Server represents the HTTP server interface.
type Server interface {
	Initialize() error
	Start() error
	Shutdown(ctx context.Context) error
}

type server struct {
	config        *config.Config
	logger        *zap.Logger
	handler       *handler.Handler
	fiberApp      *fiber.App
	isInitialized bool
}

func NewServer(config *config.Config, logger *zap.Logger, sqlDB *sql.DB, fiberApp *fiber.App) (*server, error) {
	// Create repositories
	repos, err := postgres.NewRepositories(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("failed to create repositories: %w", err)
	}

	// Create services
	services := service.NewServices(repos)

	// Create handlers
	handler := handler.NewHandler(services)

	return &server{
		config:   config,
		logger:   logger,
		handler:  handler,
		fiberApp: fiberApp,
	}, nil
}

func (s *server) RegisterMiddleware() {
	s.fiberApp.Use(
		otelfiber.Middleware(),
		middleware.InjectContext(s.logger),
		middleware.LogRequestResponse,
		middleware.CaptureError,
	)
}

func (s *server) RegisterRoutes() {
	for _, route := range s.allRoutes() {
		s.fiberApp.Add(
			route.Method,
			fmt.Sprintf("%s/v%d%s", APIBasePath, route.Version, route.Path),
			append(route.Middleware, route.Handler)...,
		)
	}
}

func (s *server) RegisterSwagger() {
	if !s.config.Swagger.Enabled {
		return
	}

	swaggerHandler := swagger.New(swagger.Config{
		Title: "koogo API Docs",
	})

	s.fiberApp.Add(
		http.MethodGet,
		"/swagger/*",
		basicauth.New(basicauth.Config{
			Users: map[string]string{
				s.config.Swagger.Username: s.config.Swagger.Password,
			},
		}),
		swaggerHandler,
	)
}

func (s *server) Initialize() error {
	if s.fiberApp == nil {
		s.logger.Error("Fiber app has not been created yet")
		return fmt.Errorf("fiber app has not been created yet")
	}

	s.RegisterMiddleware()
	s.RegisterRoutes()
	s.RegisterSwagger()

	s.isInitialized = true

	return nil
}

func (s *server) Start() error {
	if s.fiberApp == nil || !s.isInitialized {
		s.logger.Error("Fiber app is not initialized")
		return fmt.Errorf("fiber app is not initialized")
	}

	if err := s.fiberApp.Listen(fmt.Sprintf(":%d", s.config.App.Port)); err != nil {
		s.logger.Error("Failed to start server", zap.Error(err))
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	if s.fiberApp == nil {
		s.logger.Warn("Fiber app is not initialized")
		return nil
	}

	return s.fiberApp.ShutdownWithContext(ctx)
}
