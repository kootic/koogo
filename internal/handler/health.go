package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/kootic/koogo/internal/service"
	"github.com/kootic/koogo/pkg/kooctx"
	"github.com/kootic/koogo/pkg/koohttp"
)

type HealthHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type healthHandler struct {
	healthService service.HealthService
}

func NewHealthHandler(healthService service.HealthService) HealthHandler {
	return &healthHandler{healthService: healthService}
}

// HealthCheck godoc
//
//	@tags			Health
//	@Summary		Health check endpoint
//	@Description	Returns the health status of the application
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		503
//	@Router			/v1/health [get]
func (h *healthHandler) HealthCheck(c *fiber.Ctx) error {
	ctx := c.UserContext()
	logger := kooctx.GetContextLogger(ctx)

	// Check database connection
	if err := h.healthService.HealthCheck(ctx); err != nil {
		logger.Error("Health check failed", zap.Error(err))
		return koohttp.ServiceUnavailable(c)
	}

	return koohttp.Success(c, nil)
}
