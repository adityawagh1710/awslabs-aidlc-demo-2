# Components — Todo App

## Architecture: Separate Microservices (Go + Fiber)

Each service is an independently deployable Go binary containerised with Docker.

---

## Service: auth-service

**Purpose**: Handles all user identity concerns — registration, login, token lifecycle, MFA.

**Responsibilities**:
- User registration and email verification
- Password hashing and validation
- JWT access token + refresh token issuance
- Token refresh and revocation (logout)
- MFA enrollment and verification
- Brute-force protection (rate limiting on login)

**Exposes**: REST endpoints under `/auth`

---

## Service: todo-service

**Purpose**: Core todo domain — CRUD, due dates, priorities, categories/tags, search.

**Responsibilities**:
- Create, read, update, delete todos (scoped to authenticated user)
- Manage due dates and priority levels
- Manage categories and tags (create, rename, delete, assign)
- Full-text search across title, description, tags
- Enforce object-level authorization (user owns todo)

**Exposes**: REST endpoints under `/todos`, `/tags`

---

## Service: scheduler-service

**Purpose**: Handles reminders and recurring task generation using in-process Go goroutine scheduler.

**Responsibilities**:
- Store and manage reminder schedules per todo
- Fire in-app notification events when reminders are due
- Generate next occurrence of recurring todos on completion
- Support recurrence patterns: daily, weekly, monthly, custom

**Exposes**: Internal REST endpoints (called by todo-service); publishes notification events via WebSocket hub

---

## Service: file-service

**Purpose**: Manages file attachments for todos stored on local filesystem (Docker volume).

**Responsibilities**:
- Accept file uploads (validate type and size)
- Store files on local filesystem under a structured path
- Serve files securely (only to owning user)
- Delete files when parent todo is deleted

**Exposes**: REST endpoints under `/files`

---

## Service: notification-service

**Purpose**: Manages real-time in-app notifications delivered to the frontend via WebSockets.

**Responsibilities**:
- Maintain authenticated WebSocket connections per user
- Receive notification events from scheduler-service
- Push notifications to the correct user's WebSocket connection
- Store undelivered notifications for retrieval on reconnect

**Exposes**: WebSocket endpoint `/ws`; internal REST endpoint for event ingestion

---

## Frontend: vue-app

**Purpose**: Vue 3 SPA — the user-facing web interface.

**Responsibilities**:
- Authentication flows (register, login, logout, MFA)
- Todo management UI (CRUD, filters, search, tags, priorities, due dates)
- File attachment upload/download UI
- Reminder and recurrence configuration UI
- Real-time notification display via WebSocket
- Global state management via Pinia

**Communicates with**: All backend services via REST; notification-service via WebSocket
