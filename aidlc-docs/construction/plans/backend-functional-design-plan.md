# Functional Design Plan — Unit 1: backend

## Plan Steps
- [x] Answer clarifying questions (user)
- [x] Generate domain-entities.md
- [x] Generate business-rules.md
- [x] Generate business-logic-model.md

---

## Clarifying Questions

Please fill in the `[Answer]:` tag for each question. Let me know when done.

---

## Question 1
How should todo ownership and soft delete work?

A) Hard delete — todos are permanently removed from the database on delete
B) Soft delete — todos are marked as deleted (deleted_at timestamp) and excluded from queries, but retained in DB
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 2
What recurrence patterns need to be supported?

A) Fixed set — daily, weekly, monthly, yearly only
B) Fixed set + custom interval — e.g., every N days/weeks/months
C) Full cron-style — arbitrary schedule expressions
D) Other (please describe after [Answer]: tag below)

[Answer]:C

---

## Question 3
How should todo priority and status transitions work?

A) Free-form — any status/priority can be set at any time with no enforced transitions
B) Enforced transitions — status must follow: pending → in-progress → done (no skipping)
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 4
What are the file attachment constraints?

A) Max 10MB per file, allowed types: images (jpg, png, gif, webp) and documents (pdf, docx, txt)
B) Max 25MB per file, any file type allowed
C) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 5
How should JWT token expiry be configured?

A) Short-lived access token (15 min) + long-lived refresh token (7 days)
B) Short-lived access token (1 hour) + long-lived refresh token (30 days)
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 6
What should happen to a user's todos, files, and reminders when their account is deleted?

A) Cascade delete — all user data (todos, files, reminders, notifications) permanently deleted
B) Soft delete user — user marked inactive, data retained for a grace period then purged
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 7
How should search work across todos?

A) Simple LIKE/ILIKE search on title and description (no full-text index)
B) PostgreSQL full-text search (tsvector/tsquery) on title, description, and tags
C) Other (please describe after [Answer]: tag below)

[Answer]: Elasticsearch or OpenSearch for advanced search capabilities (stemming, relevance scoring, etc.)
