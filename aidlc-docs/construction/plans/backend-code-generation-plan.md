# Code Generation Plan — Unit 1: backend

## Unit Context
- **Workspace root**: `/home/adityawagh/awslabs-aidlc-demo-2`
- **Structure pattern**: Greenfield multi-unit (polyrepo microservices)
- **Code location**: Each service in its own subdirectory at workspace root
  - `auth-service/`, `todo-service/`, `scheduler-service/`, `file-service/`, `notification-service/`
- **Shared infra**: `infra/` at workspace root (Docker Compose, Nginx/Traefik config, migrations)

## Dependencies
- PostgreSQL 16, Redis 7, Elasticsearch 8.13 (via Docker Compose)
- Shared JWT secret across all services
- Internal Docker network `todo-network`

## Stories / Requirements Covered
- FR-01: Auth (register, login, logout, MFA, password reset)
- FR-02: Todo CRUD
- FR-03: Due dates, priorities
- FR-04: Categories/tags
- FR-05: Search (Elasticsearch)
- FR-06: Reminders
- FR-07: Recurring tasks
- FR-08: File attachments
- NFR: Security baseline (SECURITY-01–15), Testcontainers E2E, partial PBT

---

## Steps

### Step 1: Shared infrastructure setup
- [x] Create `infra/docker-compose.yml` with all 13 services
- [x] Create `infra/.env.example` with all required variables
- [x] Create `infra/traefik/traefik.yml` (static config: entrypoints, Let's Encrypt, dashboard)
- [x] Create `infra/traefik/dynamic.yml` (middleware: HSTS, rate limit, security headers)
- [x] Create `infra/prometheus/prometheus.yml` (scrape configs for all services)
- [x] Create `infra/grafana/provisioning/datasources/prometheus.yml`

### Step 2: auth-service — project structure
- [ ] Create `auth-service/go.mod` (module: `github.com/<org>/todo-auth-service`)
- [ ] Create `auth-service/cmd/main.go` (Fiber app bootstrap, middleware chain, graceful shutdown)
- [ ] Create `auth-service/internal/model/user.go` (User, RefreshToken structs + GORM tags)
- [ ] Create `auth-service/internal/middleware/jwt.go` (JWT validation + Redis blacklist check)
- [ ] Create `auth-service/internal/middleware/ratelimit.go` (Redis-backed rate limiter)
- [ ] Create `auth-service/internal/middleware/errorhandler.go` (global error handler, fail-closed)
- [ ] Create `auth-service/internal/middleware/logger.go` (zerolog request logger with request_id)
- [ ] Create `auth-service/migrations/000001_create_users.up.sql`
- [ ] Create `auth-service/migrations/000001_create_users.down.sql`
- [ ] Create `auth-service/migrations/000002_create_refresh_tokens.up.sql`
- [ ] Create `auth-service/migrations/000002_create_refresh_tokens.down.sql`
- [ ] Create `auth-service/Dockerfile` (multi-stage: golang:1.22-alpine → alpine:3.19)

### Step 3: auth-service — repository layer
- [ ] Create `auth-service/internal/repository/user_repository.go` (CreateUser, FindUserByEmail, UpdateMFASecret)
- [ ] Create `auth-service/internal/repository/token_repository.go` (SaveRefreshToken, DeleteRefreshToken, FindRefreshToken)
- [ ] Create `auth-service/internal/repository/user_repository_test.go` (Testcontainers, table-driven)
- [ ] Create `auth-service/internal/repository/token_repository_test.go` (Testcontainers)

### Step 4: auth-service — service layer
- [ ] Create `auth-service/internal/service/auth_service.go` (RegisterUser, AuthenticateUser, RefreshTokens, RevokeToken, EnrollMFA, VerifyMFA)
- [ ] Create `auth-service/internal/service/auth_service_test.go` (unit tests with mocks)
- [ ] Create `auth-service/internal/service/auth_service_pbt_test.go` (PBT: token validation pure functions)

### Step 5: auth-service — handler layer
- [ ] Create `auth-service/internal/handler/auth_handler.go` (RegisterHandler, LoginHandler, RefreshHandler, LogoutHandler, MFAEnrollHandler, MFAVerifyHandler)
- [ ] Create `auth-service/internal/handler/health_handler.go` (GET /health)
- [ ] Create `auth-service/internal/handler/auth_handler_test.go` (Testcontainers E2E: register→login→refresh→logout flow)
- [ ] Create `auth-service/docs/api.md` (endpoint reference)

### Step 6: todo-service — project structure + models + migrations
- [ ] Create `todo-service/go.mod`
- [ ] Create `todo-service/cmd/main.go`
- [ ] Create `todo-service/internal/model/todo.go` (Todo, Tag, TodoTag, SearchOutbox structs)
- [ ] Create `todo-service/internal/middleware/` (jwt.go, ratelimit.go, errorhandler.go, logger.go — same pattern as auth-service)
- [ ] Create `todo-service/migrations/` (todos, tags, todo_tags, search_outbox tables)
- [ ] Create `todo-service/Dockerfile`

### Step 7: todo-service — repository layer
- [ ] Create `todo-service/internal/repository/todo_repository.go` (InsertTodo, FindTodosByUser, FindTodoByID, UpdateTodo, DeleteTodo)
- [ ] Create `todo-service/internal/repository/tag_repository.go` (InsertTag, FindTagsByUser, DeleteTag)
- [ ] Create `todo-service/internal/repository/outbox_repository.go` (InsertOutboxEvent, FindUnprocessed, MarkProcessed)
- [ ] Create `todo-service/internal/repository/todo_repository_test.go` (Testcontainers)
- [ ] Create `todo-service/internal/repository/tag_repository_test.go` (Testcontainers)

### Step 8: todo-service — service layer
- [ ] Create `todo-service/internal/service/todo_service.go` (CreateTodo, ListTodos, GetTodo, UpdateTodo, DeleteTodo, SearchTodos — with status transition enforcement, ownership checks, outbox writes)
- [ ] Create `todo-service/internal/service/tag_service.go` (CreateTag, ListTags, DeleteTag)
- [ ] Create `todo-service/internal/service/outbox_worker.go` (background goroutine: poll outbox → sync to Elasticsearch)
- [ ] Create `todo-service/internal/service/circuit_breaker.go` (gobreaker wrappers for file-service and scheduler-service calls)
- [ ] Create `todo-service/internal/service/todo_service_test.go` (unit tests with mocks)
- [ ] Create `todo-service/internal/service/todo_service_pbt_test.go` (PBT: status transition logic)

### Step 9: todo-service — handler layer
- [ ] Create `todo-service/internal/handler/todo_handler.go` (CRUD + search handlers with data-testid-friendly JSON field names)
- [ ] Create `todo-service/internal/handler/tag_handler.go`
- [ ] Create `todo-service/internal/handler/health_handler.go`
- [ ] Create `todo-service/internal/handler/todo_handler_test.go` (Testcontainers E2E: create→list→update→delete→search)
- [ ] Create `todo-service/docs/api.md`

### Step 10: scheduler-service — full implementation
- [ ] Create `scheduler-service/go.mod`
- [ ] Create `scheduler-service/cmd/main.go`
- [ ] Create `scheduler-service/internal/model/` (Reminder, RecurrenceConfig structs)
- [ ] Create `scheduler-service/internal/middleware/` (jwt.go, errorhandler.go, logger.go)
- [ ] Create `scheduler-service/migrations/` (reminders, recurrence_configs tables)
- [ ] Create `scheduler-service/internal/repository/reminder_repository.go` + test
- [ ] Create `scheduler-service/internal/repository/recurrence_repository.go` + test
- [ ] Create `scheduler-service/internal/service/scheduler_service.go` (ScheduleReminder, CancelReminder, SetRecurrence, HandleTodoCompletion, RunScheduler goroutine)
- [ ] Create `scheduler-service/internal/service/recurrence_engine.go` (cron parsing, next occurrence computation)
- [ ] Create `scheduler-service/internal/service/recurrence_engine_pbt_test.go` (PBT: cron next-occurrence pure function)
- [ ] Create `scheduler-service/internal/service/circuit_breaker.go`
- [ ] Create `scheduler-service/internal/handler/scheduler_handler.go` + health + tests
- [ ] Create `scheduler-service/Dockerfile`

### Step 11: file-service — full implementation
- [ ] Create `file-service/go.mod`
- [ ] Create `file-service/cmd/main.go`
- [ ] Create `file-service/internal/model/file.go` (FileAttachment struct)
- [ ] Create `file-service/internal/middleware/` (jwt.go, errorhandler.go, logger.go)
- [ ] Create `file-service/migrations/` (file_attachments table)
- [ ] Create `file-service/internal/repository/file_repository.go` + test
- [ ] Create `file-service/internal/service/file_service.go` (UploadFile, GetFile, DeleteFile — MIME validation, size check, disk I/O)
- [ ] Create `file-service/internal/service/file_service_test.go`
- [ ] Create `file-service/internal/handler/file_handler.go` (upload, download, delete) + health + tests
- [ ] Create `file-service/Dockerfile`

### Step 12: notification-service — full implementation
- [ ] Create `notification-service/go.mod`
- [ ] Create `notification-service/cmd/main.go`
- [ ] Create `notification-service/internal/model/notification.go`
- [ ] Create `notification-service/internal/middleware/` (jwt.go, errorhandler.go, logger.go)
- [ ] Create `notification-service/migrations/` (notifications table)
- [ ] Create `notification-service/internal/repository/notification_repository.go` + test
- [ ] Create `notification-service/internal/service/hub.go` (WebSocket connection registry)
- [ ] Create `notification-service/internal/service/notification_service.go` (RegisterConnection, DeliverNotification, StoreUndelivered, GetPending)
- [ ] Create `notification-service/internal/service/notification_service_test.go`
- [ ] Create `notification-service/internal/handler/ws_handler.go` (WebSocket upgrade, auth, push pending on connect)
- [ ] Create `notification-service/internal/handler/event_handler.go` (internal event ingestion)
- [ ] Create `notification-service/internal/handler/health_handler.go`
- [ ] Create `notification-service/Dockerfile`

### Step 13: Code documentation summaries
- [ ] Create `aidlc-docs/construction/backend/code/auth-service-summary.md`
- [ ] Create `aidlc-docs/construction/backend/code/todo-service-summary.md`
- [ ] Create `aidlc-docs/construction/backend/code/scheduler-service-summary.md`
- [ ] Create `aidlc-docs/construction/backend/code/file-service-summary.md`
- [ ] Create `aidlc-docs/construction/backend/code/notification-service-summary.md`
