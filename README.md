# Taskly — Production-Grade Todo Application

A full-stack, multi-user todo application built with Go microservices and a Vue 3 frontend. Designed for production: JWT authentication with MFA, full-text search, real-time notifications, file attachments, recurring tasks, and a complete test suite.

---

## Architecture

```
                         ┌─────────────────────┐
                         │    Browser / Client │
                         └──────────┬──────────┘
                                    │ HTTPS / WSS
                         ┌──────────▼──────────┐
                         │     Traefik v3      │  ← API Gateway, TLS termination
                         │  (reverse proxy)    │    Let's Encrypt auto-renew
                         └──────────┬──────────┘
                                    │
                  ┌─────────────────┼─────────────────┐
                  │ /, /assets/...  │ /auth /todos    │
                  │                 │ /tags /files /ws│
        ┌─────────▼─────────┐  ┌────▼─────────────────────────────────┐
        │  frontend         │  │  backend (single image)              │
        │  nginx + Vue SPA  │  │  ┌──────────────────────────────┐    │
        │  :80              │  │  │ supervisord                   │   │
        │  (also reverse-   │  │  │  ├ auth-service          :3000│   │
        │   proxies API to  │  │  │  ├ todo-service          :3001│   │
        │   backend)        │  │  │  ├ file-service          :3002│   │
        └───────────────────┘  │  │  ├ notification-service  :3003│   │
                               │  │  └ scheduler-service     :3004│   │
                               │  └──────────────────────────────┘    │
                               └──┬────────────┬─────────────────┬───-┘
                                  │            │                 │
                       ┌──────────▼──┐  ┌──────▼──────┐  ┌──────▼──────-┐
                       │ PostgreSQL  │  │   Redis     │  │ Elasticsearch│
                       │ (shared,    │  │ (cache,     │  │  (search)    │
                       │  per-svc    │  │  sessions)  │  └──────────────┘
                       │  schema)    │  └─────────────┘
                       └─────────────┘

              ┌────────────────────────────────────────┐
              │  Prometheus + Grafana  (observability) │
              └────────────────────────────────────────┘
```

Two long-running images: **`todo-infra-go-backend`** (one container running all five Go services under supervisord) and **`todo-infra-go-frontend`** (nginx serving the Vue SPA + reverse-proxying API/WebSocket paths to the backend). Inter-service calls travel over `localhost` *inside* the backend container; ports are not published to the host. The `scheduler-service` calls `todo-service` (create recurrence) and `notification-service` (fire reminder events). `todo-service` calls `file-service` (cascade delete) and `scheduler-service` (create/cancel reminders).

---

## Tech Stack

| Layer | Technology | Version |
|---|---|---|
| Backend language | Go | 1.23 |
| HTTP framework | Fiber v2 | 2.52.12 |
| Database | PostgreSQL | 16 |
| Migrations | golang-migrate | 4.17 |
| ORM | GORM + pgx | latest |
| Cache / sessions | Redis | 7 |
| Search | Elasticsearch | 8.13 |
| Password hashing | argon2id | — |
| Auth tokens | JWT (golang-jwt/jwt v5) | 5.2.2 |
| MFA | TOTP via pquerna/otp | — |
| Circuit breaker | sony/gobreaker | — |
| HTTP client | go-resty/resty | — |
| Logging | zerolog (JSON, structured) | — |
| Frontend | Vue 3 + Pinia + Vue Router | 3.5+ |
| Build tool | Vite | 8 |
| CSS | Tailwind CSS | 3 |
| Icons | Heroicons v2 | — |
| API gateway | Traefik v3 | — |
| Observability | Prometheus + Grafana | — |
| Unit tests | testify + gopter (PBT) | — |
| Integration tests | Testcontainers (real PostgreSQL) | 0.31 |
| Container | Docker + Docker Compose v2 | — |
| Process supervisor (backend) | supervisord | — |
| SPA host (frontend) | nginx | 1.27 |

---

## Repository Structure

```
.
├── todo-backend/                      # Single image, all 5 Go services
│   ├── go.work                        # Workspace tying the 5 modules together
│   ├── Dockerfile                     # Multi-stage build → all 5 binaries in one image
│   ├── supervisord.conf               # Supervises the 5 services at runtime
│   ├── auth-service/                  # Identity — registration, login, JWT, MFA
│   ├── todo-service/                  # Core domain — todo CRUD, tags, full-text search
│   ├── scheduler-service/             # Time-based — reminders, recurring tasks
│   ├── file-service/                  # File attachments — upload, download, delete
│   └── notification-service/          # Real-time — WebSocket hub, event delivery
├── todo-frontend/                     # Vue 3 SPA, served by nginx
│   ├── Dockerfile                     # Build SPA, package into nginx:alpine
│   ├── nginx.conf                     # SPA fallback + reverse proxy to backend:300x
│   └── vue-app/                       # Vue source
├── infra/                             # Docker Compose stack + config
│   ├── docker-compose.yml             # Project name: todo-infra-go
│   ├── traefik/
│   ├── prometheus/
│   └── grafana/
└── aidlc-docs/                        # AI-DLC lifecycle documentation
```

Each Go service follows the same internal layout:

```
todo-backend/<service>/
├── cmd/main.go             # Entry point: wiring, migrations, server start
├── internal/
│   ├── handler/            # HTTP handlers (request parsing, response)
│   ├── service/            # Business logic
│   ├── repository/         # Database access (PostgreSQL via GORM)
│   ├── model/              # Domain structs
│   ├── middleware/         # JWT auth, rate limiting, logging, error handling
│   └── testutil/           # Shared test helpers (Testcontainers setup)
├── migrations/             # SQL migration files (auto-run on startup)
├── go.mod
└── go.sum
```

Each service still has its own `go.mod` so it can be built and tested independently. The shared `todo-backend/go.work` lets `go build` and `go test` resolve cross-module references during the unified Docker build.

---

## Services

### auth-service — `:3000`

Handles all identity concerns. All other services validate JWTs independently using a shared secret — no runtime calls to auth-service.

| Method | Route | Auth | Description |
|---|---|---|---|
| POST | `/auth/register` | — | Register with email + password |
| POST | `/auth/login` | — | Login, returns JWT pair |
| POST | `/auth/refresh` | — | Rotate access token using refresh token |
| POST | `/auth/logout` | JWT | Revoke refresh token |
| POST | `/auth/mfa/enroll` | JWT | Generate TOTP secret + QR code |
| POST | `/auth/mfa/verify` | JWT | Verify TOTP code |
| GET | `/health` | — | Health check |

**Key env vars:** `DATABASE_URL`, `JWT_SECRET`, `REDIS_URL`

---

### todo-service — `:3001`

Central domain service. Owns todo lifecycle and delegates to file-service and scheduler-service.

| Method | Route | Description |
|---|---|---|
| POST | `/todos/` | Create todo |
| GET | `/todos/` | List todos (filter by status, priority, tag) |
| GET | `/todos/search?q=` | Full-text search via Elasticsearch |
| GET | `/todos/:id` | Get single todo |
| PUT | `/todos/:id` | Update todo (title, description, status, priority, due date, tags) |
| DELETE | `/todos/:id` | Delete todo (cascades to file-service) |
| POST | `/tags/` | Create tag |
| GET | `/tags/` | List tags |
| DELETE | `/tags/:id` | Delete tag |

All routes require JWT. Status transitions: `pending → in_progress → done` only.

Uses an **outbox worker** (background goroutine) to sync todo changes to Elasticsearch asynchronously. A **circuit breaker** wraps calls to file-service and scheduler-service.

**Key env vars:** `DATABASE_URL`, `JWT_SECRET`, `REDIS_URL`, `ELASTICSEARCH_URL`, `FILE_SERVICE_URL`, `SCHEDULER_SERVICE_URL`

---

### scheduler-service — `:3004`

Owns all time-based concerns. Runs a background cron scheduler that polls PostgreSQL for due reminders.

| Method | Route | Auth | Description |
|---|---|---|---|
| POST | `/reminders` | API key | Schedule a reminder for a todo |
| DELETE | `/reminders/:id` | API key | Cancel a reminder |
| POST | `/todos/:id/recurrence` | API key | Set recurrence pattern (cron expression) |
| POST | `/todos/:id/complete` | API key | Trigger next occurrence on completion |

Routes are internal-only (API key, not JWT). On reminder fire: calls `notification-service`. On recurrence: calls `todo-service` to create the next occurrence.

**Key env vars:** `DATABASE_URL`, `INTERNAL_API_KEY`, `TODO_SERVICE_URL`, `NOTIFICATION_SERVICE_URL`

---

### file-service — `:3002`

Manages file attachment lifecycle. Files stored on a named Docker volume (`file-uploads`). Metadata in PostgreSQL.

| Method | Route | Description |
|---|---|---|
| POST | `/files/` | Upload file (multipart, max 10 MB, allowlisted MIME types) |
| GET | `/files/:id` | Download file (ownership enforced) |
| DELETE | `/files/:id` | Delete file from disk + database |

Allowed MIME types: `image/jpeg`, `image/png`, `image/gif`, `image/webp`, `application/pdf`, `.docx`, `text/plain`. Max 10 attachments per todo.

**Key env vars:** `DATABASE_URL`, `JWT_SECRET`, `FILE_STORAGE_PATH`

---

### notification-service — `:3003`

Real-time delivery layer. Maintains an in-memory WebSocket connection registry per user.

| Method | Route | Auth | Description |
|---|---|---|---|
| GET | `/ws` | JWT (header or `?token=`) | Upgrade to WebSocket |
| POST | `/internal/events` | API key | Ingest event from scheduler-service |

On connect, pending (undelivered) notifications are flushed to the client. If the user is offline, notifications are persisted in PostgreSQL and delivered on next connect.

**Key env vars:** `DATABASE_URL`, `JWT_SECRET`, `INTERNAL_API_KEY`

---

### todo-frontend — `:5173` (dev) / Traefik (prod)

Vue 3 SPA with Pinia state management and Vue Router, packaged into an nginx image at `todo-frontend/`. Source lives under [todo-frontend/vue-app/](todo-frontend/vue-app/). The container serves static assets and reverse-proxies API paths (`/api/auth`, `/api/todos`, `/api/tags`, `/api/files`, `/ws`) to the consolidated backend container — see [todo-frontend/nginx.conf](todo-frontend/nginx.conf). API calls use the `/api` prefix to avoid collisions with SPA client-side routes. Connects to notification-service via WebSocket for real-time updates.

**Views:**

| Route | View | Description |
|---|---|---|
| `/login` | LoginView | Email + password login, MFA redirect |
| `/register` | RegisterView | Account creation with password strength meter |
| `/mfa` | MfaView | TOTP 6-digit verification |
| `/todos` | TodosView | Task list, search, filters, tag manager, stats bar |
| `/todos/:id` | TodoDetailView | Full task detail with inline edit |
| `/settings` | SettingsView | Profile, MFA enrollment with QR code |

Auth guards redirect unauthenticated users to `/login` and authenticated users away from guest-only pages. Token refresh is handled automatically by the Axios interceptor on 401 responses.

---

## Getting Started

### Prerequisites

| Tool | Minimum version |
|---|---|
| Go | 1.22 |
| Docker | 24 |
| Docker Compose | v2 |
| Node.js | 18 |
| npm | 9 |

### 1. Clone

```bash
git clone <repo-url>
cd awslabs-aidlc-demo-2
```

### 2. Configure environment

```bash
cd infra
cp .env.example .env
# Edit .env — fill in all placeholder values (see Environment Variables below)
```

### 3. Start the full stack

```bash
cd infra
docker compose up -d

# Wait for all services to be healthy
docker compose ps
```

Eight containers come up under the `todo-infra-go` compose project: `traefik`, `postgres`, `redis`, `elasticsearch`, `backend`, `frontend`, `prometheus`, `grafana`. A ninth container, `docker-cleanup`, runs daily to prune unused Docker images and build cache automatically. Each of the 5 Go services runs its own migrations on startup against an isolated `<svc>_migrations` table, so first-run boot can take 10–20 s while supervisord brings them up inside the `backend` container.

### 4. Verify

The backend's per-service ports are *not* published to the host — Traefik routes external traffic, and the nginx in `frontend` reverse-proxies API paths internally. To probe each Go service, exec into the container:

```bash
docker exec backend wget -qO- http://localhost:3000/health   # auth
docker exec backend wget -qO- http://localhost:3001/health   # todo
docker exec backend wget -qO- http://localhost:3002/health   # file
docker exec backend wget -qO- http://localhost:3003/health   # notification
docker exec backend wget -qO- http://localhost:3004/health   # scheduler
```

End-to-end through nginx (validates routing + backend wiring):

```bash
docker exec frontend wget -qO- http://localhost/             # SPA index.html
```

---

## Exposed Ports

| Service | Host Port | Purpose |
|---|---|---|
| Frontend (nginx) | `8081` | Vue SPA + API proxy |
| Traefik | `8080` | API gateway |
| PostgreSQL | `5433` | Database (DBeaver, pgAdmin, etc.) |
| Grafana | `3000` | Dashboards & monitoring |

### Connecting DBeaver to PostgreSQL

| Setting | Value |
|---|---|
| Host | `localhost` |
| Port | `5433` |
| Database | `todo_app` |
| Username | `todo_user` |
| Password | `StrongPass123!` (or your `.env` value) |

### Accessing Grafana

Open `http://localhost:3000` and log in with:
- Username: `admin`
- Password: `GrafanaPass123!` (or your `.env` `GRAFANA_ADMIN_PASSWORD` value)

Prometheus is pre-configured as a datasource.

---

## Local Development

### Building a service

The 5 services share `todo-backend/go.work`, so you can build any of them from the workspace root:

```bash
cd todo-backend
go build -o bin/auth ./auth-service/cmd
go build -o bin/todo ./todo-service/cmd
# ...etc
```

Or build a single module on its own:

```bash
cd todo-backend/auth-service
go mod download
go build -o bin/auth-service ./cmd
```

### Running a service locally

Services expect environment variables. Minimal example for auth-service:

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/todoapp?sslmode=disable"
export JWT_SECRET="your-secret-at-least-32-chars"
export REDIS_URL="redis://:password@localhost:6379/0"
export LOG_LEVEL="debug"

./bin/auth
```

For local development you'll typically run the dependent containers (postgres, redis, elasticsearch) via `docker compose up -d postgres redis elasticsearch` from `infra/`, then start the Go service directly on the host so the host can reach `localhost:3000`–`localhost:3004`.

### Running the frontend

```bash
cd todo-frontend/vue-app
npm install
npm run dev        # starts at http://localhost:5173
```

The Vite dev server proxies API calls to the backend services on host ports — for that to work, expose them by adding a temporary `ports:` block to the `backend` service in [infra/docker-compose.yml](infra/docker-compose.yml), or run the Go services directly on the host. The proxy table:

| Path | Proxied to |
|---|---|
| `/api/auth/*` | `localhost:3000` |
| `/api/todos/*`, `/api/tags/*` | `localhost:3001` |
| `/api/files/*` | `localhost:3002` |
| `/ws` | `localhost:3003` (WebSocket) |

---

## Environment Variables

| Variable | Used by | Description |
|---|---|---|
| `DATABASE_URL` | all services | PostgreSQL connection string |
| `JWT_SECRET` | all services | JWT signing key (min 32 chars) |
| `REDIS_URL` | auth, todo | Redis connection URL |
| `ELASTICSEARCH_URL` | todo | Elasticsearch endpoint |
| `FILE_SERVICE_URL` | todo | Internal URL of file-service |
| `SCHEDULER_SERVICE_URL` | todo | Internal URL of scheduler-service |
| `TODO_SERVICE_URL` | scheduler | Internal URL of todo-service |
| `NOTIFICATION_SERVICE_URL` | scheduler | Internal URL of notification-service |
| `INTERNAL_API_KEY` | scheduler, notification | Shared key for internal calls |
| `FILE_STORAGE_PATH` | file | Local disk path for uploads |
| `LOG_LEVEL` | all services | `debug`, `info`, `warn`, `error` |
| `POSTGRES_DB` | postgres container | Database name |
| `POSTGRES_USER` | postgres container | Database user |
| `POSTGRES_PASSWORD` | postgres container | Database password |
| `REDIS_PASSWORD` | redis container | Redis auth password |
| `ELASTIC_PASSWORD` | elasticsearch container | ES password |
| `GRAFANA_ADMIN_PASSWORD` | grafana container | Grafana admin password |
| `GITHUB_ORG` | docker-compose | GitHub org for container image pulls |

---

## Testing

### Unit tests (no Docker required)

```bash
# Run for any service
cd todo-backend/auth-service
go test ./internal/service/... -v

# With race detector
go test -race ./internal/service/... -v

# Property-based tests only
go test ./internal/service/... -run PBT -v
```

**PBT coverage:**
- `auth-service`: JWT round-trip preserves userID; `HashToken` is deterministic
- `todo-service`: All 9 status transition combinations validate correctly
- `scheduler-service`: `NextOccurrence` is always after `from`; deterministic; invalid cron → error

### Integration tests (Testcontainers — requires Docker)

Testcontainers spins up a real PostgreSQL container per test. No external services needed.

```bash
cd todo-backend/auth-service
go test ./internal/repository/... ./internal/handler/... -v -timeout 120s

cd ../todo-service
go test ./internal/repository/... ./internal/handler/... -v -timeout 120s

cd ../file-service
go test ./internal/repository/... -v -timeout 120s

cd ../scheduler-service
go test ./internal/repository/... -v -timeout 120s
```

**Integration test coverage:**

| Service | Tests |
|---|---|
| auth-service | User CRUD, token CRUD, register→login→refresh→logout, duplicate email |
| todo-service | Todo CRUD + status transitions + filters, tag CRUD, JWT enforcement |
| file-service | Insert, FindByID, FindByTodo, CountByTodo, Delete |
| scheduler-service | Reminder Insert/FindDue/MarkFired/Delete/Count, RecurrenceConfig Upsert/Update/Delete |

### Full test suite per service

```bash
cd todo-backend
for svc in auth-service todo-service scheduler-service file-service notification-service; do
  (cd $svc && go test ./... -timeout 180s)
done
```

### End-to-end tests (Cypress)

Requires the full Docker stack running on port 8081.

```bash
cd todo-frontend/vue-app
./node_modules/.bin/cypress run
```

Four spec files:
- `auth_flows.cy.js` — login, wrong password rejection, client-side validation
- `auth_refresh.cy.js` — sessionStorage persistence across reload, logout on session clear
- `tag_management.cy.js` — open tag panel, create/delete tags, create todo with tag
- `todo_filters.cy.js` — filter by status/priority, search, clear filters, empty state

### Vulnerability scan

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest

cd todo-backend
for svc in auth-service todo-service scheduler-service file-service notification-service; do
  echo "=== $svc ===" && (cd $svc && govulncheck ./...)
done
```

---

## Security

Security Baseline extension is **enabled** — all SECURITY-01 through SECURITY-15 rules are enforced as blocking constraints:

- Passwords hashed with **argon2id** (adaptive, memory-hard)
- JWT access tokens (short-lived) + refresh tokens (stored hashed in Redis/PostgreSQL)
- **Rate limiting** on all public endpoints (Fiber middleware)
- **Account lockout** after 5 failed login attempts (Redis-backed)
- **HTTP security headers** enforced via Traefik middleware
- All inter-service traffic stays on the internal Docker network
- `INTERNAL_API_KEY` gates internal endpoints — never exposed publicly
- Input validation on all request bodies via `go-playground/validator`
- Structured logging: no sensitive data (passwords, tokens) in logs
- All data encrypted in transit (TLS 1.2+ via Traefik + Let's Encrypt)
- Prometheus/Grafana accessible only via internal routes (basic auth)

---

## Observability

**Prometheus** scrapes metrics from all services every 15s (`GET /metrics`).

**Grafana** provides pre-provisioned dashboards for:
- Per-service HTTP request rate, latency (p50/p95/p99), error rate
- PostgreSQL connection pool stats
- Redis hit/miss ratio
- Elasticsearch indexing rate

**Alerts configured:**
- Authentication failure spikes
- p95 latency > 500 ms
- Error rate > 1%

**Logs:** zerolog JSON to stdout on all services. Docker log driver `json-file` with `max-size: 100m`, `max-file: 5`. Minimum 90-day retention (host-level archival required).

---

## Data Flow Examples

### Create a todo with a reminder

```
Browser → POST /todos (todo-service)
  → INSERT todo (PostgreSQL)
  → Outbox worker syncs to Elasticsearch
  → POST /reminders (scheduler-service)
      → INSERT reminder (PostgreSQL)
      → schedule goroutine

[At reminder fire time]
scheduler-service goroutine
  → POST /internal/events (notification-service)
      → push to user WebSocket connection (browser)
```

### Complete a recurring todo

```
Browser → PUT /todos/:id {status: "done"} (todo-service)
  → UPDATE todo status = done (PostgreSQL)
  → POST /todos/:id/complete (scheduler-service)
      → compute next occurrence from cron expression
      → POST /todos (todo-service) — create next occurrence
      → DELETE old reminders
  ← 200 OK
```

---

## License

MIT
