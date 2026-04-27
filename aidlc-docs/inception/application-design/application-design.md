# Application Design — Todo App

## Architecture Summary

| Decision | Choice |
|---|---|
| Backend architecture | Separate microservices (Go + Fiber) |
| Authentication | Built-in JWT (access + refresh tokens, MFA via TOTP) |
| File storage | Local filesystem on Docker volume |
| Scheduling | In-process goroutine scheduler (no extra infra) |
| Frontend ↔ backend | REST + WebSockets (notifications) |
| Frontend state | Pinia |
| API versioning | None (no prefix) |

---

## Services

| Service | Language | Purpose |
|---|---|---|
| auth-service | Go + Fiber | User registration, login, JWT, MFA |
| todo-service | Go + Fiber | Todo CRUD, tags, search |
| scheduler-service | Go + Fiber | Reminders, recurring tasks |
| file-service | Go + Fiber | File attachment storage |
| notification-service | Go + Fiber | Real-time WebSocket notifications |
| vue-app | Vue 3 + Pinia | Web frontend |

---

## Layers per Service (Go)

```
handler/       HTTP handlers (Fiber routes, request parsing, response formatting)
service/       Business logic orchestration
repository/    Database access (PostgreSQL via pgx)
model/         Domain structs
middleware/    JWT validation, rate limiting, logging, error handling
```

---

## Inter-Service Communication

```
vue-app ──REST──> auth-service
vue-app ──REST──> todo-service
vue-app ──REST──> file-service
vue-app ──WS───> notification-service

todo-service ──internal REST──> file-service        (cascade delete)
todo-service ──internal REST──> scheduler-service   (reminders, completion)

scheduler-service ──internal REST──> notification-service  (fire events)
scheduler-service ──internal REST──> todo-service          (create recurrence)
```

All inter-service calls travel over the internal Docker network and are never exposed publicly.

---

## Key Design Decisions

### JWT Validation (Stateless)
All services validate JWT tokens independently using a shared secret. No runtime calls to auth-service for validation — keeps services decoupled.

### Goroutine Scheduler
scheduler-service runs a background goroutine that polls PostgreSQL for due reminders at a configurable interval. Simple, no extra infrastructure (no Redis, no message broker).

### Circular Dependency (scheduler ↔ todo)
scheduler-service calls todo-service to create the next recurrence. This is a deliberate, managed circular dependency via internal REST — acceptable at this scale.

### File Storage
Files stored on a named Docker volume mounted into file-service. Metadata (filename, path, size, MIME type, ownerID, todoID) stored in PostgreSQL.

---

## Detailed Artifacts

- Components: `application-design/components.md`
- Methods: `application-design/component-methods.md`
- Services: `application-design/services.md`
- Dependencies: `application-design/component-dependency.md`
