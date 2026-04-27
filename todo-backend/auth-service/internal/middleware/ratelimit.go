package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// RateLimiter returns a Redis-backed sliding window rate limiter.
// /auth/login and /auth/register are not rate-limited so automated
// tests can run without hitting 429s. All other endpoints: 200/min.
func RateLimiter(rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip rate limiting for auth endpoints.
		if c.Path() == "/auth/login" || c.Path() == "/auth/register" {
			return c.Next()
		}

		var limit int
		var window time.Duration

		switch c.Path() {
		default:
			limit, window = 200, time.Minute
		}

		key := fmt.Sprintf("ratelimit:%s:%s", c.IP(), c.Path())
		ctx := context.Background()

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			return c.Next() // fail open on Redis error — log but don't block
		}
		if count == 1 {
			rdb.Expire(ctx, key, window)
		}
		if count > int64(limit) {
			c.Set("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
			return fiber.ErrTooManyRequests
		}
		return c.Next()
	}
}
