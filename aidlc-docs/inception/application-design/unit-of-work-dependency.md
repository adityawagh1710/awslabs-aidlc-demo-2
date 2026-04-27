# Unit of Work Dependencies — Todo App

## Unit Dependency Matrix

| Unit | Depends On | Type |
|---|---|---|
| backend (all services) | Shared PostgreSQL | Runtime |
| backend (scheduler-service) | backend (todo-service) | Internal REST (recurrence) |
| backend (scheduler-service) | backend (notification-service) | Internal REST (events) |
| backend (todo-service) | backend (file-service) | Internal REST (cascade delete) |
| backend (todo-service) | backend (scheduler-service) | Internal REST (reminders) |
| frontend | backend (all services) | REST + WebSocket |

## Inter-Service Dependencies Within Unit 1 (Backend)

```
auth-service          (no dependencies — standalone)
file-service          (no dependencies — standalone)

todo-service
  --> file-service        (cascade delete on todo deletion)
  --> scheduler-service   (create/cancel reminders, complete todo)

scheduler-service
  --> todo-service        (create next recurrence todo)
  --> notification-service (fire reminder events)

notification-service
  <-- scheduler-service   (receives events)
  <-- vue-app             (WebSocket connections)
```

## Development Dependency Order (Within Parallel Phase)

Although all backend services are developed in parallel, integration testing requires this order:

1. `auth-service` + `file-service` — no dependencies, can be fully tested standalone
2. `notification-service` — depends only on incoming events, can be tested with mock events
3. `todo-service` — depends on file-service and scheduler-service (use mocks during unit dev)
4. `scheduler-service` — depends on todo-service and notification-service (use mocks during unit dev)

Integration testing (Testcontainers) runs all services together after individual development.

## Shared Infrastructure Dependencies

All services depend on:
- PostgreSQL (shared database, shared schema)
- Internal Docker network (`todo-network`)
- JWT shared secret (environment variable, injected via Docker Compose)
