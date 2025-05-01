package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/kootic/koogo/pkg/kooctx"
)

func InjectContext(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.SetUserContext(kooctx.SetContextLogger(c.UserContext(), logger))

		return c.Next()
	}
}
