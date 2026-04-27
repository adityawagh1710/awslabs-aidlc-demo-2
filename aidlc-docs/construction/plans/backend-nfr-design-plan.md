# NFR Design Plan — Unit 1: backend

## Plan Steps
- [x] Answer clarifying questions (user)
- [x] Generate nfr-design-patterns.md
- [x] Generate logical-components.md

---

## Clarifying Questions

Please fill in the `[Answer]:` tag for each question. Let me know when done.

---

## Question 1
Should a caching layer (e.g., Redis) be added for frequently accessed data (user sessions, todo lists)?

A) No cache — PostgreSQL + GORM connection pooling is sufficient for the target scale
B) Redis for session/token blacklist only (revoked JWT tracking)
C) Redis for both token blacklist and hot data caching (todo lists, user profiles)
D) Other (please describe after [Answer]: tag below)

[Answer]:C

---

## Question 2
How should rate limiting be implemented?

A) In-process Fiber middleware (memory-based, per-instance — sufficient for single Docker instance)
B) Redis-backed rate limiter (shared across instances if scaled horizontally)
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 3
Should a circuit breaker pattern be applied to inter-service calls?

A) No circuit breaker — `resty` retries with backoff are sufficient; services degrade gracefully
B) Yes — use `sony/gobreaker` circuit breaker on inter-service HTTP calls
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 4
How should the Elasticsearch sync (write-through) be handled when Elasticsearch is unavailable?

A) Best-effort sync — log the failure, continue; search falls back to PostgreSQL ILIKE
B) Retry queue — failed sync events queued in-memory and retried by a background goroutine
C) Other (please describe after [Answer]: tag below)

[Answer]: Outbox + async worker (most reliable) — write sync events to an outbox table in PostgreSQL, processed by a separate worker goroutine with retry logic
