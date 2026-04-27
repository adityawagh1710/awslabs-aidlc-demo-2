# Application Design Plan — Todo App

## Plan Steps
- [x] Answer clarifying questions (user)
- [x] Generate components.md
- [x] Generate component-methods.md
- [x] Generate services.md
- [x] Generate component-dependency.md
- [x] Generate application-design.md (consolidated)

---

## Clarifying Questions

Please fill in the `[Answer]:` tag for each question. Let me know when done.

---

## Question 1
How should the backend be structured architecturally?

A) Layered monolith — single Go binary with clear internal layers (handler → service → repository)
B) Modular monolith — single binary but split into self-contained feature modules (auth, todos, files, notifications)
C) Separate microservices — independent deployable services per domain
D) Other (please describe after [Answer]: tag below)

[Answer]:C

---

## Question 2
How should authentication be handled?

A) Built-in — implement JWT auth directly in the Go/Fiber backend (register, login, refresh, logout endpoints)
B) External identity provider — delegate to a third-party (Keycloak, Auth0, etc.)
C) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 3
How should file attachments be stored?

A) Local filesystem on the server (simple, works with Docker volumes)
B) MinIO (S3-compatible object storage, self-hosted, runs as a Docker container)
C) Other (please describe after [Answer]: tag below)

[Answer]:A
---

## Question 4
How should reminders and recurring task scheduling be handled?

A) In-process scheduler — Go goroutine-based scheduler within the backend (simple, no extra infra)
B) Database-driven polling — a background worker polls the DB for due reminders/recurrences
C) External job queue — use a message broker (e.g., Redis + worker) for scheduling
D) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 5
How should the Vue 3 frontend communicate with the backend?

A) REST API only — standard HTTP/JSON endpoints
B) REST API + WebSockets — REST for CRUD, WebSockets for real-time in-app notifications
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 6
Should the frontend use a state management library?

A) Pinia (recommended for Vue 3)
B) No global state management — component-local state + composables only
C) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 7
How should the API be versioned?

A) URL path versioning — `/api/v1/...`
B) No versioning — keep it simple for now
C) Other (please describe after [Answer]: tag below)

[Answer]:B
