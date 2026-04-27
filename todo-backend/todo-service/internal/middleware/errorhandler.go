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
		switch code {
		case fiber.StatusUnauthorized:
			msg = "unauthorized"
		case fiber.StatusForbidden:
			msg = "forbidden"
		case fiber.StatusNotFound:
			msg = "not found"
		case fiber.StatusUnprocessableEntity:
			msg = "validation error"
		case fiber.StatusTooManyRequests:
			msg = "too many requests"
		}
	}
	log.Error().Err(err).Int("status", code).Str("path", c.Path()).Msg("request error")
	return c.Status(code).JSON(fiber.Map{"error": msg})
}

func Recover() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Str("path", c.Path()).Msg("panic recovered")
				err = fiber.ErrInternalServerError
			}
		}()
		return c.Next()
	}
}
