# Services — Todo App

## Service Orchestration Overview

Each microservice owns its domain. Inter-service communication is synchronous REST (internal calls) except for notification delivery which is event-driven via the notification-service ingestion endpoint.

---

## auth-service

**Role**: Identity provider for the entire system.

**Orchestration**:
- All other services validate JWT tokens independently (shared secret / public key)
- No service calls auth-service at runtime for token validation — validation is stateless

**Internal interactions**: None (standalone)

---

## todo-service

**Role**: Central domain service. Orchestrates todo lifecycle including delegation to scheduler and file services.

**Orchestration**:
- On todo deletion → calls `file-service DELETE /files` for each attachment
- On todo completion → calls `scheduler-service POST /todos/:id/complete` to trigger recurrence
- On reminder creation → calls `scheduler-service POST /reminders`

**Consumed by**: vue-app (REST)

---

## scheduler-service

**Role**: Owns all time-based concerns — reminders and recurrence.

**Orchestration**:
- Runs an internal goroutine scheduler that polls for due reminders
- On reminder fire → calls `notification-service POST /internal/events`
- On todo completion with recurrence → creates new todo via `todo-service POST /todos`

**Consumed by**: todo-service (internal REST)

---

## file-service

**Role**: Manages file storage lifecycle.

**Orchestration**:
- Standalone — no outbound service calls
- Files stored on Docker volume; metadata in PostgreSQL

**Consumed by**: todo-service (internal REST for cascade delete), vue-app (REST for upload/download)

---

## notification-service

**Role**: Real-time delivery layer.

**Orchestration**:
- Maintains in-memory WebSocket connection registry per user
- Receives events from scheduler-service via internal REST
- Delivers to live WebSocket connection or stores for later retrieval

**Consumed by**: vue-app (WebSocket), scheduler-service (internal REST)

---

## Service Communication Summary

| From | To | Protocol | Purpose |
|---|---|---|---|
| vue-app | auth-service | REST | Auth flows |
| vue-app | todo-service | REST | Todo CRUD, tags, search |
| vue-app | file-service | REST | Upload / download attachments |
| vue-app | notification-service | WebSocket | Real-time notifications |
| todo-service | file-service | Internal REST | Cascade delete attachments |
| todo-service | scheduler-service | Internal REST | Create/cancel reminders, complete todo |
| scheduler-service | notification-service | Internal REST | Fire reminder notification event |
| scheduler-service | todo-service | Internal REST | Create next recurrence todo |
