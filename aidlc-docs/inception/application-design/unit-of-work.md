# Unit of Work — Todo App

## Decomposition Summary

| Decision | Choice |
|---|---|
| Unit mapping | 2 units: Backend + Frontend |
| Development sequence | Parallel (all backend services simultaneously, then frontend) |
| Repository structure | Polyrepo (each service in its own repository) |
| Database isolation | Shared database, shared schema |

---

## Unit 1: backend

**Description**: All 5 Go/Fiber microservices developed in parallel, each in its own repository.

**Services**:

| Service | Repository | Purpose |
|---|---|---|
| auth-service | `todo-auth-service` | JWT auth, MFA, brute-force protection |
| todo-service | `todo-todo-service` | Todo CRUD, tags, search |
| scheduler-service | `todo-scheduler-service` | Reminders, recurring tasks |
| file-service | `todo-file-service` | File attachment storage |
| notification-service | `todo-notification-service` | WebSocket real-time notifications |

**Shared infrastructure** (defined in a separate `todo-infra` repo):
- `docker-compose.yml` — orchestrates all services + PostgreSQL
- Shared PostgreSQL instance (single database, shared schema)
- Internal Docker network for inter-service communication

**Code organisation per service**:
```
<service-name>/
├── cmd/
│   └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── model/
│   └── middleware/
├── Dockerfile
├── go.mod
└── go.sum
```

---

## Unit 2: frontend

**Description**: Vue 3 SPA in its own repository, developed after backend services are stable.

**Repository**: `todo-vue-app`

**Code organisation**:
```
todo-vue-app/
├── src/
│   ├── components/
│   ├── views/
│   ├── stores/          # Pinia stores
│   ├── composables/
│   ├── router/
│   ├── api/             # API client modules
│   └── assets/
├── cypress/             # Cypress E2E tests
├── Dockerfile
├── package.json
└── vite.config.ts
```

---

## Development Sequence

```
Phase 1 (Parallel):
  todo-auth-service        ──┐
  todo-todo-service        ──┤
  todo-scheduler-service   ──┼──> All backend services developed simultaneously
  todo-file-service        ──┤
  todo-notification-service──┘

Phase 2 (Sequential, after Phase 1):
  todo-vue-app             ──> Frontend developed against stable backend APIs
```
