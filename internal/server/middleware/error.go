package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/kootic/koogo/pkg/kooctx"
	"github.com/kootic/koogo/pkg/koohttp"
)

func CaptureError(c *fiber.Ctx) error {
	originalErr := c.Next()
	if originalErr == nil {
		return nil
	}

	var fiberErr *fiber.Error

	// Let fiber handle its own errors
	ok := errors.As(originalErr, &fiberErr)
	if ok {
		return originalErr
	}

	var apiErr koohttp.APIError

	// Handle our own API errors
	ok = errors.As(originalErr, &apiErr)
	if ok {
		return c.Status(apiErr.HTTPStatus()).JSON(apiErr)
	}

	// Anything else we wrap the original error in our own internal server error
	apiErr = koohttp.NewAPIError(http.StatusInternalServerError, koohttp.APIErrorCodeInternalServerError)

	respErr := c.Status(http.StatusInternalServerError).JSON(apiErr)
	if respErr != nil {
		return fmt.Errorf("failed to send internal server error response: %w: %w", respErr, apiErr)
	}

	// Log the error
	logger := kooctx.GetContextLogger(c.UserContext())
	logger.Error("Unexpected error", zap.Error(originalErr))

	return nil
}
