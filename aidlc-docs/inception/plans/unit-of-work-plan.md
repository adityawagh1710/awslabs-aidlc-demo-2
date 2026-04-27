# Unit of Work Plan — Todo App

## Plan Steps
- [x] Answer clarifying questions (user)
- [x] Generate unit-of-work.md
- [x] Generate unit-of-work-dependency.md
- [x] Generate unit-of-work-story-map.md
- [x] Validate unit boundaries

---

## Clarifying Questions

Please fill in the `[Answer]:` tag for each question. Let me know when done.

---

## Question 1
The application design identified 5 microservices. How should these map to units of work for development?

A) One unit per service — each service developed and delivered independently (auth, todo, scheduler, file, notification + frontend as separate units)
B) Grouped units — related services bundled together (e.g., Unit 1: auth; Unit 2: todo + scheduler + file; Unit 3: notification + frontend)
C) Two units — backend (all 5 services) and frontend (vue-app)
D) Other (please describe after [Answer]: tag below)

[Answer]:C

---

## Question 2
What should the development sequence / priority order be?

A) Foundation first — auth-service → todo-service → scheduler-service → file-service → notification-service → vue-app
B) Core value first — todo-service (with stub auth) → auth-service → remaining services → vue-app
C) Parallel — all backend services developed simultaneously, then frontend
D) Other (please describe after [Answer]: tag below)

[Answer]:C

---

## Question 3
How should the repository be structured?

A) Monorepo — all services in one repository with a top-level directory per service
B) Polyrepo — each service in its own repository
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 4
Should each unit have its own PostgreSQL schema or database?

A) Shared database, separate schemas per service (simpler ops, still logical isolation)
B) Separate database per service (stronger isolation, more complex ops)
C) Shared database, shared schema (simplest, no isolation)
D) Other (please describe after [Answer]: tag below)

[Answer]:C
