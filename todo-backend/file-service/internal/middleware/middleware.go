package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func JWTAuth(secret string, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return fiber.ErrUnauthorized
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(_ *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}
		blacklisted, _ := rdb.Exists(context.Background(), "blacklist:"+tokenStr).Result()
		if blacklisted > 0 {
			return fiber.ErrUnauthorized
		}
		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}

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
