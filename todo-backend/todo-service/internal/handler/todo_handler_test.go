package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo-app/todo-service/internal/handler"
	"github.com/todo-app/todo-service/internal/middleware"
	"github.com/todo-app/todo-service/internal/model"
	"github.com/todo-app/todo-service/internal/repository"
	"github.com/todo-app/todo-service/internal/service"
	"github.com/todo-app/todo-service/internal/testutil"
)

const testSecret = "test-secret-32-chars-minimum-ok!"

func setupApp(t *testing.T) (*fiber.App, string) {
	t.Helper()
	ctx := context.Background()
	pgc, db := testutil.StartPostgres(t, ctx)
	t.Cleanup(func() { pgc.Terminate(ctx) })

	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	todoRepo := repository.NewTodoRepository(db)
	tagRepo := repository.NewTagRepository(db)
	outboxRepo := repository.NewOutboxRepository(db)
	todoSvc := service.NewTodoService(todoRepo, tagRepo, outboxRepo, nil, "", "")
	tagSvc := service.NewTagService(tagRepo)

	todoH := handler.NewTodoHandler(todoSvc)
	tagH := handler.NewTagHandler(tagSvc)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Use(middleware.Recover())
	jwtMW := middleware.JWTAuth(testSecret, rdb)

	app.Get("/health", handler.NewHealthHandler().Health)
	todos := app.Group("/todos", jwtMW)
	todos.Post("/", todoH.Create)
	todos.Get("/", todoH.List)
	todos.Get("/search", todoH.Search)
	todos.Get("/:id", todoH.Get)
	todos.Put("/:id", todoH.Update)
	todos.Delete("/:id", todoH.Delete)
	tags := app.Group("/tags", jwtMW)
	tags.Post("/", tagH.Create)
	tags.Get("/", tagH.List)
	tags.Delete("/:id", tagH.Delete)

	// Generate a test JWT
	userID := "00000000-0000-0000-0000-000000000001"
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}).SignedString([]byte(testSecret))

	return app, token
}

func doRequest(t *testing.T, app *fiber.App, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	var b []byte
	if body != nil {
		b, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := app.Test(req, 10000)
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	rec.WriteHeader(resp.StatusCode)
	rec.Body = new(bytes.Buffer)
	rec.Body.ReadFrom(resp.Body)
	return rec
}

func TestE2E_TodoCRUD(t *testing.T) {
	app, token := setupApp(t)

	// Create
	rec := doRequest(t, app, "POST", "/todos", map[string]any{"title": "Buy milk", "priority": "high"}, token)
	assert.Equal(t, 201, rec.Code)
	var created model.Todo
	json.Unmarshal(rec.Body.Bytes(), &created)
	require.NotEmpty(t, created.ID)

	// List
	rec = doRequest(t, app, "GET", "/todos", nil, token)
	assert.Equal(t, 200, rec.Code)
	var todos []model.Todo
	json.Unmarshal(rec.Body.Bytes(), &todos)
	assert.Len(t, todos, 1)

	// Get
	rec = doRequest(t, app, "GET", fmt.Sprintf("/todos/%s", created.ID), nil, token)
	assert.Equal(t, 200, rec.Code)

	// Update status: pending → in_progress
	inProgress := string(model.StatusInProgress)
	rec = doRequest(t, app, "PUT", fmt.Sprintf("/todos/%s", created.ID), map[string]any{"status": inProgress}, token)
	assert.Equal(t, 200, rec.Code)

	// Invalid transition: in_progress → pending (should fail)
	rec = doRequest(t, app, "PUT", fmt.Sprintf("/todos/%s", created.ID), map[string]any{"status": "pending"}, token)
	assert.Equal(t, 422, rec.Code)

	// Delete
	rec = doRequest(t, app, "DELETE", fmt.Sprintf("/todos/%s", created.ID), nil, token)
	assert.Equal(t, 204, rec.Code)

	// Get after delete should 404
	rec = doRequest(t, app, "GET", fmt.Sprintf("/todos/%s", created.ID), nil, token)
	assert.Equal(t, 404, rec.Code)
}

func TestE2E_Unauthorized(t *testing.T) {
	app, _ := setupApp(t)
	rec := doRequest(t, app, "GET", "/todos", nil, "")
	assert.Equal(t, 401, rec.Code)
}
