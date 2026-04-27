# NFR Requirements — Unit 1: backend

## Performance

| Requirement | Target |
|---|---|
| API p95 response time | < 200ms under normal load |
| Search query response time | < 500ms (Elasticsearch) |
| File upload throughput | Non-blocking (async write, immediate metadata response) |
| Concurrent users | Hundreds (production-grade) |
| DB connection pool | Configured per service (GORM connection pool settings) |

---

## Availability

| Requirement | Target |
|---|---|
| Uptime target | 99.9% |
| Health check endpoint | `GET /health` on each service (liveness + readiness) |
| Auto-restart | Docker `restart: unless-stopped` policy on all containers |
| Graceful shutdown | Each service handles SIGTERM — drains in-flight requests before exit |

---

## Security

All 15 Security Baseline rules enforced (SECURITY-01 through SECURITY-15). Key requirements:

| Rule | Requirement |
|---|---|
| SECURITY-01 | PostgreSQL TLS enforced; Elasticsearch TLS enforced |
| SECURITY-03 | Structured JSON logging via zerolog; no PII/secrets in logs |
| SECURITY-04 | HTTP security headers on all HTML-serving endpoints (handled by vue-app reverse proxy) |
| SECURITY-05 | Input validation on all API parameters; GORM parameterised queries (no raw SQL concat) |
| SECURITY-06 | JWT scoped to user; no wildcard permissions |
| SECURITY-08 | JWT validated server-side on every request; IDOR prevention via user_id ownership checks |
| SECURITY-09 | Generic error responses; no stack traces exposed |
| SECURITY-10 | `go.sum` lock file committed; `govulncheck` in CI |
| SECURITY-11 | Rate limiting on all public endpoints (Fiber middleware) |
| SECURITY-12 | argon2id password hashing; session invalidated on logout; brute-force lockout |
| SECURITY-15 | Global error handler at Fiber app level; fail-closed on auth errors |

---

## Reliability

| Requirement | Detail |
|---|---|
| Global error handler | Fiber `app.Use(recover middleware)` + custom error handler on each service |
| DB transaction safety | Multi-step operations (e.g., create todo + tags + recurrence) wrapped in DB transactions |
| Scheduler fault tolerance | Reminder goroutine recovers from panics; logs errors and continues |
| Inter-service retries | `resty` retry on transient 5xx (max 3 retries, exponential backoff) |
| Graceful degradation | Search unavailable → fallback to PostgreSQL ILIKE search |

---

## Scalability

| Requirement | Detail |
|---|---|
| Stateless services | No in-process session state; JWT-based auth |
| Horizontal scaling | Each service independently scalable via Docker Compose `scale` or future orchestrator |
| Elasticsearch scaling | Separate container, independently scalable |

---

## Maintainability

| Requirement | Detail |
|---|---|
| Structured logging | zerolog JSON output with timestamp, level, correlation ID, service name |
| Log retention | Minimum 90 days (SECURITY-14) |
| Migrations | golang-migrate SQL files; run on service startup |
| Dependency pinning | `go.sum` committed; no `latest` tags in Dockerfiles |
| Vulnerability scanning | `govulncheck` in CI pipeline |

---

## Testing

| Requirement | Detail |
|---|---|
| Unit tests | Per service, table-driven tests for business logic |
| Integration tests | Testcontainers (real PostgreSQL + Elasticsearch containers) |
| Property-based tests | Partial — pure functions (cron computation, token validation, file MIME detection) + serialisation round-trips |
| E2E (frontend) | Cypress (Unit 2: frontend) |
