# Tech Stack Decisions — Unit 1: backend

## Per-Service Stack (all 5 Go services)

| Layer | Technology | Version | Rationale |
|---|---|---|---|
| Language | Go | 1.22+ | Performance, concurrency, small Docker images |
| HTTP Framework | Fiber | v2 | Fast, Express-like API, built-in middleware ecosystem |
| ORM / DB | GORM | v2 | Code-first models, auto-migrations support, familiar API |
| DB Driver (GORM) | `gorm.io/driver/postgres` | latest | PostgreSQL dialect for GORM |
| Migrations | `golang-migrate` | v4 | SQL migration files, up/down, Docker-friendly |
| JWT | `golang-jwt/jwt` | v5 | Most widely used, well-maintained, HS256/RS256 support |
| Password hashing | `alexedwards/argon2id` | latest | argon2id wrapper, OWASP-recommended |
| Structured logging | `zerolog` | v1 | Zero-allocation JSON logging, fast, simple API |
| HTTP client | `resty` | v2 | Fluent API, built-in retry with backoff |
| Search client | `elastic/go-elasticsearch` | v8 | Official Elasticsearch v8 client |
| TOTP (MFA) | `pquerna/otp` | latest | TOTP/HOTP generation and validation |
| Cron parsing | `robfig/cron` | v3 | Standard cron expression parsing and scheduling |
| Config | `spf13/viper` | v1 | Environment variable + config file support |
| Validation | `go-playground/validator` | v10 | Struct tag-based input validation |
| UUID | `google/uuid` | v1 | UUID v4 generation |
| Testing | `testify` | v1 | Assertions, mocks, test suites |
| Testcontainers | `testcontainers-go` | latest | Real PostgreSQL + Elasticsearch in tests |
| PBT | `leanovate/gopter` | latest | Property-based testing for pure functions |

---

## Infrastructure Stack

| Component | Technology | Version | Rationale |
|---|---|---|---|
| Database | PostgreSQL | 16 | Reliable, ACID, full-text search fallback |
| Search | Elasticsearch | 8.x | Advanced search, stemming, relevance scoring |
| Container runtime | Docker + Docker Compose | v2 | Self-hosted deployment, all services containerised |
| Reverse proxy | Nginx | alpine | TLS termination, HTTP security headers, static file serving |
| File storage | Docker named volume | — | Local filesystem, simple ops |

---

## Docker Image Strategy

| Service | Base Image | Notes |
|---|---|---|
| All Go services | `golang:1.22-alpine` (build) → `alpine:3.19` (runtime) | Multi-stage build, minimal runtime image |
| vue-app | `node:20-alpine` (build) → `nginx:alpine` (runtime) | Static files served by Nginx |
| PostgreSQL | `postgres:16-alpine` | Pinned version, no `latest` |
| Elasticsearch | `docker.elastic.co/elasticsearch/elasticsearch:8.13.0` | Pinned version |

---

## Environment Configuration

All services configured via environment variables (no hardcoded secrets):

| Variable | Used By | Description |
|---|---|---|
| `DATABASE_URL` | all services | PostgreSQL connection string |
| `JWT_SECRET` | auth-service, all services | Shared JWT signing secret |
| `ELASTICSEARCH_URL` | todo-service | Elasticsearch endpoint |
| `FILE_STORAGE_PATH` | file-service | Mount path for Docker volume |
| `INTERNAL_API_KEY` | scheduler, notification | Internal service auth header |
| `LOG_LEVEL` | all services | zerolog level (debug/info/warn/error) |
