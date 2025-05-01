package koohttp

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	ErrInvalidParamUUID = NewAPIError(http.StatusBadRequest, "invalid_param_uuid")
)

type WithValidate interface {
	Validate() error
}

func GetBodyAndValidate[T any](c *fiber.Ctx) (*T, error) {
	var body T

	if err := c.BodyParser(&body); err != nil {
		return nil, err
	}

	if dto, ok := any(body).(WithValidate); ok {
		if err := dto.Validate(); err != nil {
			return nil, err
		}
	}

	return &body, nil
}

func GetQueryAndValidate[T any](c *fiber.Ctx) (*T, error) {
	var query T

	if err := c.QueryParser(&query); err != nil {
		return nil, err
	}

	if dto, ok := any(query).(WithValidate); ok {
		if err := dto.Validate(); err != nil {
			return nil, err
		}
	}

	return &query, nil
}

func GetParamUUID(c *fiber.Ctx, paramKey string) (uuid.UUID, error) {
	param := c.Params(paramKey)
	if param == "" {
		return uuid.Nil, ErrInvalidParamUUID
	}

	parsedUUID, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil, ErrInvalidParamUUID
	}

	return parsedUUID, nil
}
