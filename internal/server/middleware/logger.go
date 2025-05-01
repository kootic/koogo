package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/kootic/koogo/pkg/kooctx"
)

func LogRequestResponse(c *fiber.Ctx) error {
	startTime := time.Now()

	err := c.Next()

	statusCode := c.Response().StatusCode()
	url := c.OriginalURL()
	params := c.AllParams()
	queries := c.Queries()
	headers := c.GetReqHeaders()
	requestBody := string(c.Body())
	responseBody := string(c.Response().Body())
	latencyMs := float64(time.Since(startTime).Microseconds()) / 1000.0

	logger := kooctx.GetContextLogger(c.UserContext())
	logger.Debug(
		fmt.Sprintf("%s %s", c.Method(), url),
		zap.String("method", c.Method()),
		zap.Int("status", statusCode),
		zap.Any("params", params),
		zap.Any("queries", queries),
		zap.Float64("latency_ms", latencyMs),
		zap.Any("headers", headers),
		zap.Any("request_body", requestBody),
		zap.Any("response_body", responseBody),
	)

	return err
}
