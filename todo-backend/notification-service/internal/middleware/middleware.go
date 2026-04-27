package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	}
	log.Error().Err(err).Int("status", code).Msg("error")
	return c.Status(code).JSON(fiber.Map{"error": msg})
}

func Recover() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Msg("panic recovered")
				err = fiber.ErrInternalServerError
			}
		}()
		return c.Next()
	}
}

func InternalAPIKey(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("X-Internal-API-Key") != key {
			return fiber.ErrUnauthorized
		}
		return c.Next()
	}
}
