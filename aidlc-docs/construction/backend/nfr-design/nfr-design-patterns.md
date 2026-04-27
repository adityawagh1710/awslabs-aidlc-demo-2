# NFR Design Patterns — Unit 1: backend

## Pattern 1: Redis Caching (Token Blacklist + Hot Data)

**Applies to**: auth-service, todo-service

**Token Blacklist** (auth-service):
- On logout/token revocation: write `token_hash → expiry` to Redis with TTL = remaining token lifetime
- On every authenticated request: check Redis blacklist before accepting token
- Prevents use of revoked tokens without DB lookup on every request

**Hot Data Cache** (todo-service):
- Cache todo list responses per user: key = `todos:{userID}:{filterHash}`, TTL = 60s
- Cache user profile: key = `user:{userID}`, TTL = 300s
- Cache invalidation: on todo create/update/delete, delete all `todos:{userID}:*` keys
- Read-through: cache miss → DB query → populate cache → return

---

## Pattern 2: Redis-Backed Rate Limiting

**Applies to**: all services (public-facing endpoints)

**Implementation**: Fiber rate limiter middleware using Redis as shared store
- Algorithm: sliding window counter
- Key: `ratelimit:{IP}:{endpoint}`
- Limits:
  - Login endpoint: 5 requests / 15 min per IP (AUTH-06)
  - Register endpoint: 10 requests / hour per IP
  - General API: 100 requests / min per authenticated user
- On limit exceeded: return 429 with `Retry-After` header
- Redis-backed ensures limits are shared across horizontally scaled instances (SECURITY-11)

---

## Pattern 3: Circuit Breaker on Inter-Service Calls

**Applies to**: todo-service, scheduler-service (outbound internal REST calls)

**Implementation**: `sony/gobreaker` wrapping `resty` HTTP calls

**Configuration per circuit**:
| Circuit | Threshold | Timeout | Half-Open Requests |
|---|---|---|---|
| todo → file-service | 5 failures | 30s | 1 |
| todo → scheduler-service | 5 failures | 30s | 1 |
| scheduler → notification-service | 5 failures | 30s | 1 |
| scheduler → todo-service (recurrence) | 5 failures | 30s | 1 |

**States**: Closed (normal) → Open (failing, reject fast) → Half-Open (probe) → Closed
**On Open**: return degraded response (e.g., skip reminder creation, log warning) — never fail the primary operation

---

## Pattern 4: Outbox Pattern for Elasticsearch Sync

**Applies to**: todo-service

**Design**:
- `search_outbox` table in PostgreSQL:
  ```
  id          UUID PK
  todo_id     UUID
  operation   enum (upsert | delete)
  payload     JSONB
  processed   bool default false
  created_at  timestamp
  ```
- On todo create/update/soft-delete: write to `search_outbox` in the **same DB transaction** as the todo mutation
- Background worker goroutine polls `search_outbox WHERE processed = false ORDER BY created_at` every 5s
- Worker sends to Elasticsearch; on success marks `processed = true`
- On Elasticsearch failure: log error, leave `processed = false` for next poll cycle (automatic retry)
- Guarantees: no sync event lost even if Elasticsearch is temporarily down; eventual consistency

---

## Pattern 5: Graceful Shutdown

**Applies to**: all services

```
1. Receive SIGTERM
2. Stop accepting new connections (Fiber ShutdownWithTimeout)
3. Wait for in-flight requests to complete (max 30s)
4. Close DB connection pool
5. Close Redis connection
6. Exit 0
```

---

## Pattern 6: Structured Logging with Correlation IDs

**Applies to**: all services

- Each request assigned a `X-Request-ID` (UUID) at entry (generated if not provided by client)
- zerolog logger enriched with: `request_id`, `service`, `level`, `timestamp`, `user_id` (when authenticated)
- All log lines for a request share the same `request_id` for traceability (SECURITY-03)
- No PII, passwords, or tokens logged at any level

---

## Pattern 7: Global Error Handler (Fail-Closed)

**Applies to**: all services

- Fiber `app.Use(recover.New())` catches panics → logs with zerolog → returns 500 generic response
- Custom Fiber error handler maps domain errors to HTTP status codes
- Auth/authz errors always return 401/403 — never leak internal state (SECURITY-15)
- DB/external errors return 500 with generic message — no stack traces (SECURITY-09)
