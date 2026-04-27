package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// RequestLogger logs each request with correlation ID, method, path, status, latency.
// No PII or tokens are logged (SECURITY-03).
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("X-Request-ID", requestID)
		c.Locals("requestID", requestID)

		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		log.Info().
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("latency", latency).
			Msg("request")

		return err
	}
}
