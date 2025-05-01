package middleware

import (
	"fmt"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/kootic/koogo/pkg/kooctx"
)

var ignorePaths = []string{
	"/api/v1/health",
}

func LogRequestResponse(c *fiber.Ctx) error {
	startTime := time.Now()

	err := c.Next()

	url := c.OriginalURL()
	if slices.Contains(ignorePaths, url) {
		return err
	}

	statusCode := c.Response().StatusCode()
	params := c.AllParams()
	queries := c.Queries()
	headers := c.GetReqHeaders()
	requestBody := string(c.Body())
	responseBody := string(c.Response().Body())
	latencyMs := float64(time.Since(startTime).Microseconds()) / 1000.0

	logger := kooctx.GetContextLogger(c.UserContext()).With(
		zap.String("method", c.Method()),
		zap.Int("status", statusCode),
		zap.Any("params", params),
		zap.Any("queries", queries),
		zap.Float64("latency_ms", latencyMs),
		zap.Any("headers", headers),
		zap.Any("request_body", requestBody),
		zap.Any("response_body", responseBody),
	)

	if statusCode >= 500 || err != nil {
		logger.Error("Internal server error", zap.Error(err))
	} else {
		logger.Debug(fmt.Sprintf("%s %s", c.Method(), url))
	}

	return err
}
