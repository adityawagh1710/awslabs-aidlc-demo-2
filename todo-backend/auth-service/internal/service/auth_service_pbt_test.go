package service_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/todo-app/auth-service/internal/middleware"
	"github.com/todo-app/auth-service/internal/service"
)

// PBT: for any valid userID string, a JWT signed and immediately parsed
// must yield the same userID and not be expired.
func TestPBT_JWTRoundTrip(t *testing.T) {
	secret := "test-secret-32-chars-minimum-ok!"
	properties := gopter.NewProperties(nil)

	properties.Property("JWT round-trip preserves userID", prop.ForAll(
		func(userID string) bool {
			if userID == "" {
				return true // skip empty — not a valid UUID but not the property under test
			}
			claims := &middleware.Claims{
				UserID: userID,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}
			tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
			if err != nil {
				return false
			}
			parsed := &middleware.Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, parsed, func(_ *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			return err == nil && token.Valid && parsed.UserID == userID
		},
		gen.AlphaString(),
	))

	properties.TestingRun(t)
}

// PBT: hashToken must be deterministic — same input always yields same output.
func TestPBT_HashTokenDeterministic(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("hashToken is deterministic", prop.ForAll(
		func(s string) bool {
			return service.HashToken(s) == service.HashToken(s)
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}
