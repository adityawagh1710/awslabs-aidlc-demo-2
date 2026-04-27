# NFR Requirements Plan — Unit 1: backend

## Plan Steps
- [x] Answer clarifying questions (user)
- [x] Generate nfr-requirements.md
- [x] Generate tech-stack-decisions.md

---

## Clarifying Questions

Please fill in the `[Answer]:` tag for each question. Let me know when done.

---

## Question 1
Which Go PostgreSQL driver/ORM should be used?

A) `pgx` (raw driver, high performance, no ORM)
B) `sqlx` (lightweight SQL helper over `database/sql`)
C) `GORM` (full ORM, code-first)
D) Other (please describe after [Answer]: tag below)

[Answer]:C

---

## Question 2
How should database migrations be managed?

A) `golang-migrate` (SQL migration files, up/down)
B) `goose` (SQL or Go migration files)
C) `Atlas` (schema-as-code)
D) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 3
Which library should handle JWT creation and validation?

A) `golang-jwt/jwt` (most widely used Go JWT library)
B) `lestrrat-go/jwx` (full JWK/JWE/JWT suite)
C) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 4
Which Elasticsearch/OpenSearch client should be used?

A) Official `elastic/go-elasticsearch` client
B) `opensearch-project/opensearch-go` client
C) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 5
How should structured logging be implemented?

A) `zerolog` (zero-allocation, structured JSON logging)
B) `zap` (Uber's high-performance structured logger)
C) `slog` (Go standard library structured logger, Go 1.21+)
D) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 6
How should inter-service HTTP calls be made?

A) Standard `net/http` client with timeouts configured
B) `resty` (fluent HTTP client with retry support)
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 7
What is the target availability / uptime requirement?

A) Best-effort (no formal SLA — suitable for self-hosted single instance)
B) 99.9% uptime (requires health checks + auto-restart via Docker restart policy)
C) Other (please describe after [Answer]: tag below)

[Answer]:B
