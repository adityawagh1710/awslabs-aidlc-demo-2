# Unit of Work Story Map — Todo App

> User Stories stage was skipped. This map is derived directly from functional requirements.

## Unit 1: backend

| Requirement | Service | Priority |
|---|---|---|
| FR-01: User registration, login, logout | auth-service | High |
| FR-01: Password reset via email | auth-service | Medium |
| FR-01: MFA (TOTP) enrollment and verification | auth-service | High |
| FR-02: Todo CRUD (create, read, update, delete) | todo-service | High |
| FR-02: Todo status management (pending/in-progress/done) | todo-service | High |
| FR-03: Due dates and priority levels | todo-service | High |
| FR-04: Categories and tags (create, assign, filter) | todo-service | Medium |
| FR-05: Full-text search across todos | todo-service | Medium |
| FR-06: Reminder scheduling and firing | scheduler-service | Medium |
| FR-07: Recurring task generation on completion | scheduler-service | Medium |
| FR-08: File attachment upload and storage | file-service | Medium |
| FR-08: File download (ownership enforced) | file-service | Medium |
| FR-08: File deletion on todo delete | file-service | Medium |
| Real-time in-app notifications (WebSocket) | notification-service | Medium |
| Undelivered notification storage and retrieval | notification-service | Low |

## Unit 2: frontend

| Requirement | Component | Priority |
|---|---|---|
| FR-01: Registration, login, logout UI | Auth views | High |
| FR-01: MFA setup and verification UI | Auth views | High |
| FR-02: Todo list, create, edit, delete UI | Todo views | High |
| FR-03: Due date picker, priority selector | Todo views | High |
| FR-04: Tag management and filter UI | Tag components | Medium |
| FR-05: Search bar and results UI | Search component | Medium |
| FR-06: Reminder configuration UI | Reminder component | Medium |
| FR-07: Recurrence configuration UI | Recurrence component | Medium |
| FR-08: File attachment upload/download UI | File component | Medium |
| Real-time notification display | Notification component | Medium |
| Cypress E2E tests (NFR-06) | cypress/ | High |

## Delivery Milestones

| Milestone | Deliverables |
|---|---|
| M1 — Backend foundation | auth-service + todo-service (core CRUD) fully tested |
| M2 — Backend complete | All 5 services integrated, Testcontainers E2E passing |
| M3 — Frontend complete | Vue app feature-complete, Cypress E2E passing |
