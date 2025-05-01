package koohttp

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// There is intentionally not a function to return an InternalServerError,
// the intended way to handle unexpected errors is to simply return the error
// and let the middleware handle it.

func Success(c *fiber.Ctx, data any) error {
	return c.Status(http.StatusOK).JSON(data)
}

func SuccessCreated(c *fiber.Ctx, data any) error {
	return c.Status(http.StatusCreated).JSON(data)
}

func SuccessNoContent(c *fiber.Ctx) error {
	return c.Status(http.StatusNoContent).JSON(nil)
}

func BadRequest(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(NewAPIError(http.StatusBadRequest, APIErrorCodeBadRequest))
}

func Unauthorized(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(NewAPIError(http.StatusUnauthorized, APIErrorCodeUnauthorized))
}

func Forbidden(c *fiber.Ctx) error {
	return c.Status(http.StatusForbidden).JSON(NewAPIError(http.StatusForbidden, APIErrorCodeForbidden))
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(NewAPIError(http.StatusNotFound, APIErrorCodeNotFound))
}

func RequestTimeout(c *fiber.Ctx) error {
	return c.Status(http.StatusRequestTimeout).JSON(NewAPIError(http.StatusRequestTimeout, APIErrorCodeRequestTimeout))
}

func Conflict(c *fiber.Ctx) error {
	return c.Status(http.StatusConflict).JSON(NewAPIError(http.StatusConflict, APIErrorCodeConflict))
}

func UnprocessableEntity(c *fiber.Ctx) error {
	return c.Status(http.StatusUnprocessableEntity).JSON(NewAPIError(http.StatusUnprocessableEntity, APIErrorCodeUnprocessableEntity))
}

func ServiceUnavailable(c *fiber.Ctx) error {
	return c.Status(http.StatusServiceUnavailable).JSON(NewAPIError(http.StatusServiceUnavailable, APIErrorCodeServiceUnavailable))
}
