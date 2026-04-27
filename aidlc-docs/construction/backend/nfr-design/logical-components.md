# Logical Components — Unit 1: backend

## Infrastructure Components

| Component | Technology | Purpose |
|---|---|---|
| PostgreSQL | postgres:16-alpine | Primary data store for all services |
| Redis | redis:7-alpine | Token blacklist, hot data cache, rate limit counters |
| Elasticsearch | elasticsearch:8.13.0 | Full-text search with stemming and relevance scoring |
| Nginx | nginx:alpine | Reverse proxy, TLS termination, HTTP security headers |
| Docker Compose | v2 | Container orchestration for all services |

---

## Updated Service Component Map

### auth-service
| Component | Type | Responsibility |
|---|---|---|
| Fiber app | HTTP server | Route handling, middleware chain |
| JWT middleware | Middleware | Validate access token + Redis blacklist check |
| Rate limiter | Middleware | Redis-backed sliding window (login: 5/15min, register: 10/hr) |
| Security headers | Middleware | HSTS, CSP, X-Frame-Options, etc. (via Nginx for HTML; Fiber for API) |
| Error handler | Middleware | Global panic recovery + domain error mapping |
| Auth service | Business logic | Registration, login, token lifecycle, MFA |
| Redis client | Infrastructure | Token blacklist reads/writes |
| GORM + postgres | Infrastructure | User and refresh token persistence |

### todo-service
| Component | Type | Responsibility |
|---|---|---|
| Fiber app | HTTP server | Route handling, middleware chain |
| JWT middleware | Middleware | Validate token + Redis blacklist check |
| Rate limiter | Middleware | Redis-backed (100 req/min per user) |
| Error handler | Middleware | Global panic recovery |
| Todo service | Business logic | CRUD, tags, ownership enforcement |
| Cache layer | Infrastructure | Redis read-through/write-invalidate for todo lists + user profiles |
| Circuit breaker | Resilience | `gobreaker` wrapping calls to file-service and scheduler-service |
| Outbox writer | Infrastructure | Write sync events to `search_outbox` in same DB transaction |
| Outbox worker | Background goroutine | Poll `search_outbox`, sync to Elasticsearch, mark processed |
| Elasticsearch client | Infrastructure | `go-elasticsearch` v8 for search queries and index sync |
| GORM + postgres | Infrastructure | Todo, tag, attachment metadata persistence |

### scheduler-service
| Component | Type | Responsibility |
|---|---|---|
| Fiber app | HTTP server | Internal REST endpoints |
| JWT middleware | Middleware | Validate token on internal calls |
| Reminder scheduler | Background goroutine | Poll DB every 30s for due reminders, fire events |
| Recurrence engine | Business logic | Parse cron expression, compute next occurrence |
| Circuit breaker | Resilience | `gobreaker` on calls to notification-service and todo-service |
| Panic recovery | Resilience | Scheduler goroutine recovers from panics, logs, continues |
| GORM + postgres | Infrastructure | Reminder and recurrence config persistence |

### file-service
| Component | Type | Responsibility |
|---|---|---|
| Fiber app | HTTP server | Upload, download, delete endpoints |
| JWT middleware | Middleware | Validate token, enforce ownership |
| File validator | Business logic | MIME type allowlist, size limit enforcement |
| Disk writer | Infrastructure | Write/read/delete files on Docker volume |
| GORM + postgres | Infrastructure | File metadata persistence |

### notification-service
| Component | Type | Responsibility |
|---|---|---|
| Fiber app | HTTP server | Internal event ingestion endpoint |
| WebSocket hub | Infrastructure | In-memory map of `userID → ws.Conn` |
| WebSocket handler | Handler | Upgrade HTTP → WS, register connection, push pending notifications |
| Event processor | Business logic | Route incoming events to live connection or store as undelivered |
| GORM + postgres | Infrastructure | Undelivered notification persistence |

---

## Component Interaction Diagram

```
Internet
    |
  [Nginx] ── TLS termination, security headers, reverse proxy
    |
    +──> auth-service     [Fiber + JWT + Redis blacklist + Rate limiter]
    |         |
    |       [Redis] <──── token blacklist, rate limit counters
    |       [PostgreSQL] ─ users, refresh_tokens
    |
    +──> todo-service     [Fiber + JWT + Redis cache + Circuit breaker + Outbox]
    |         |
    |       [Redis] <──── todo list cache, user profile cache, rate limit
    |       [PostgreSQL] ─ todos, tags, todo_tags, search_outbox
    |       [Elasticsearch] ── search index (synced via outbox worker)
    |       [file-service] ── cascade delete (circuit breaker)
    |       [scheduler-service] ── reminders, completion (circuit breaker)
    |
    +──> file-service     [Fiber + JWT + disk I/O]
    |         |
    |       [Docker volume] ── file storage
    |       [PostgreSQL] ─ file_attachments
    |
    +──> notification-service  [Fiber + WebSocket hub]
    |         |
    |       [PostgreSQL] ─ notifications
    |
    scheduler-service     [Fiber + goroutine scheduler + circuit breaker]
              |
            [PostgreSQL] ─ reminders, recurrence_configs
            [notification-service] ── fire events (circuit breaker)
            [todo-service] ── create recurrence (circuit breaker)
```

---

## New Infrastructure Components Added

| Component | Added By | Reason |
|---|---|---|
| Redis | NFR Design | Token blacklist, hot data cache, Redis-backed rate limiting |
| `search_outbox` table | NFR Design | Reliable Elasticsearch sync via outbox pattern |
| `sony/gobreaker` | NFR Design | Circuit breaker on inter-service calls |
| Outbox worker goroutine | NFR Design | Background Elasticsearch sync processor |
