package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/auth-service/internal/handler"
	"github.com/todo-app/auth-service/internal/middleware"
	"github.com/todo-app/auth-service/internal/repository"
	"github.com/todo-app/auth-service/internal/service"
	"github.com/todo-app/auth-service/internal/testutil"
)

const testSecret = "test-secret-32-chars-minimum-ok!"

func setupApp(t *testing.T) *fiber.App {
	t.Helper()
	ctx := context.Background()
	pgc, db := testutil.StartPostgres(t, ctx)
	t.Cleanup(func() { pgc.Terminate(ctx) })

	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	svc := service.NewAuthService(userRepo, tokenRepo, rdb, testSecret)
	h := handler.NewAuthHandler(svc)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Use(middleware.Recover())
	app.Post("/auth/register", h.Register)
	app.Post("/auth/login", h.Login)
	app.Post("/auth/refresh", h.Refresh)
	app.Post("/auth/logout", middleware.JWTAuth(testSecret, rdb), h.Logout)
	return app
}

func post(t *testing.T, app *fiber.App, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 10000)
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	rec.WriteHeader(resp.StatusCode)
	return rec
}

func TestE2E_RegisterLoginRefreshLogout(t *testing.T) {
	app := setupApp(t)

	// Register
	b, _ := json.Marshal(map[string]string{"email": "e2e@example.com", "password": "password123"})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var pair service.TokenPair
	json.NewDecoder(resp.Body).Decode(&pair)
	require.NotEmpty(t, pair.AccessToken)
	require.NotEmpty(t, pair.RefreshToken)

	// Login
	b, _ = json.Marshal(map[string]string{"email": "e2e@example.com", "password": "password123"})
	req = httptest.NewRequest("POST", "/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var loginPair service.TokenPair
	json.NewDecoder(resp.Body).Decode(&loginPair)

	// Refresh
	b, _ = json.Marshal(map[string]string{"refresh_token": loginPair.RefreshToken})
	req = httptest.NewRequest("POST", "/auth/refresh", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var refreshedPair service.TokenPair
	json.NewDecoder(resp.Body).Decode(&refreshedPair)

	// Logout
	b, _ = json.Marshal(map[string]string{"refresh_token": refreshedPair.RefreshToken})
	req = httptest.NewRequest("POST", "/auth/logout", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+refreshedPair.AccessToken)
	resp, err = app.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestE2E_DuplicateRegister(t *testing.T) {
	app := setupApp(t)
	body := map[string]string{"email": "dup@example.com", "password": "password123"}
	b, _ := json.Marshal(body)

	for i, wantStatus := range []int{201, 409} {
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, 10000)
		require.NoError(t, err)
		assert.Equal(t, wantStatus, resp.StatusCode, "attempt %d", i+1)
	}
}
