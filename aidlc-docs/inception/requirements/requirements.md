# Requirements — Todo App

## Intent Analysis

- **User Request**: Build a todo application
- **Request Type**: New Project (Greenfield)
- **Scope Estimate**: Multiple Components (frontend + backend + database)
- **Complexity Estimate**: Complex — full-featured, multi-user, production-grade

---

## Functional Requirements

### FR-01: User Authentication & Account Management
- Users must be able to register with email and password
- Users must be able to log in and log out
- Sessions must expire and be invalidated on logout
- Password reset via email must be supported
- MFA must be supported for all users (mandatory for admin accounts)

### FR-02: Todo Management (CRUD)
- Authenticated users can create, read, update, and delete their own todos
- Each todo has: title, description, status (pending/in-progress/done), due date, priority (low/medium/high), category/tags
- Users can mark todos as complete/incomplete
- Todos are private to the owning user (no cross-user access)

### FR-03: Due Dates & Priorities
- Users can set and update due dates on todos
- Users can assign priority levels: low, medium, high
- Overdue todos must be visually distinguishable in the UI

### FR-04: Categories & Tags
- Users can create, rename, and delete categories/tags
- Todos can be assigned one or more tags
- Users can filter todos by tag/category

### FR-05: Search
- Users can search todos by title, description, or tag
- Search results are scoped to the authenticated user's todos only

### FR-06: Reminders
- Users can set one or more reminder times per todo
- Reminders are delivered as in-app notifications (email notifications are a stretch goal)

### FR-07: Recurring Tasks
- Users can configure a todo to recur on a schedule (daily, weekly, monthly, custom)
- On completion of a recurring todo, the next occurrence is automatically created

### FR-08: File Attachments
- Users can attach files to todos (images, documents)
- Maximum file size and allowed types must be enforced
- Files are stored securely and accessible only to the owning user

---

## Non-Functional Requirements

### NFR-01: Performance
- API response time < 200ms at p95 under normal load
- Support hundreds of concurrent users
- File upload/download must not block the main API thread

### NFR-02: Security
- Security Baseline extension: **ENABLED** (all SECURITY-01 through SECURITY-15 rules enforced as blocking constraints)
- Authentication: JWT-based, validated server-side on every request
- Passwords hashed with adaptive algorithm (bcrypt or argon2)
- All data encrypted at rest and in transit (TLS 1.2+)
- HTTP security headers enforced on all HTML-serving endpoints
- Input validation on all API parameters
- Rate limiting on all public-facing endpoints
- CORS restricted to allowed origins only

### NFR-03: Scalability
- Stateless backend to allow horizontal scaling
- Database connection pooling

### NFR-04: Reliability
- Graceful error handling — no stack traces exposed to users
- Global error handler at application entry point
- Database transactions for multi-step operations

### NFR-05: Maintainability
- Structured logging with correlation IDs, log levels, timestamps
- No sensitive data in logs
- Log retention minimum 90 days

### NFR-06: End-to-End Testing
- Backend E2E tests must use Testcontainers to spin up real dependencies (PostgreSQL, etc.) in Docker containers
- Backend E2E tests must cover critical API flows: registration, login, todo CRUD, file attachment, recurring tasks
- Frontend E2E tests must use Cypress to cover critical user journeys in the browser
- Cypress tests must cover: registration, login, todo creation/completion/deletion, search, tag filtering
- All tests must be runnable locally and in CI without any pre-existing external services

### NFR-07: Accessibility
- Frontend must meet WCAG 2.1 AA standards

---

## Technology Stack

| Layer | Technology |
|---|---|
| Backend | Go with Fiber framework |
| Frontend | Vue 3 (latest) |
| Database | PostgreSQL (via REST API) |
| Deployment | Docker-based (self-hosted server/VM) |

---

## Extension Configuration

| Extension | Enabled | Decided At |
|---|---|---|
| Security Baseline | Yes | Requirements Analysis |
| Property-Based Testing | Partial (pure functions + serialization) | Requirements Analysis |

---

## Constraints & Assumptions

- Deployment target: Docker-based on a self-hosted server or VM (Dockerfiles + docker-compose for all services)
- No mobile app in scope (web only)
- Email reminders are a stretch goal; in-app notifications are in scope
- File storage location (local disk vs object storage) to be decided in Infrastructure Design
