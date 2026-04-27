# Business Rules — Unit 1: backend

## Authentication Rules

| ID | Rule |
|---|---|
| AUTH-01 | Email must be unique across all active users |
| AUTH-02 | Password minimum 8 characters; checked against breached password list (HIBP API or local list) |
| AUTH-03 | Password stored as bcrypt/argon2 hash — never plaintext |
| AUTH-04 | Access token expires in 1 hour; refresh token expires in 30 days |
| AUTH-05 | Refresh token is single-use — rotated on every refresh |
| AUTH-06 | Login endpoint locked for 15 minutes after 5 consecutive failed attempts per email |
| AUTH-07 | MFA verification required on login if mfa_enabled = true |
| AUTH-08 | Soft-deleted users (deleted_at IS NOT NULL) cannot log in |
| AUTH-09 | Session invalidated on logout — refresh token deleted from DB |

---

## Todo Rules

| ID | Rule |
|---|---|
| TODO-01 | A user can only read, update, or delete their own todos (user_id match enforced) |
| TODO-02 | Title is required; max 255 characters |
| TODO-03 | Description max 5000 characters |
| TODO-04 | Status transitions enforced: pending → in_progress → done only (no skipping, no reversal) |
| TODO-05 | Priority can be changed freely at any time |
| TODO-06 | Soft delete: sets deleted_at = now(); excluded from all list/search queries |
| TODO-07 | Cascade on soft delete: associated reminders cancelled, recurrence config deleted, files deleted from disk |
| TODO-08 | Due date must be in the future when set |

---

## Tag Rules

| ID | Rule |
|---|---|
| TAG-01 | Tag names are unique per user (case-insensitive) |
| TAG-02 | Tag name max 50 characters |
| TAG-03 | Deleting a tag removes it from all associated todos (TodoTag rows deleted) |
| TAG-04 | A todo can have at most 20 tags |

---

## Recurrence Rules

| ID | Rule |
|---|---|
| REC-01 | Recurrence uses standard 5-field cron expression (minute hour day month weekday) |
| REC-02 | On todo completion, if RecurrenceConfig exists: compute next_occurrence from cron, create new todo with same title/description/priority/tags/recurrence |
| REC-03 | New recurring todo starts with status = pending |
| REC-04 | Reminders are NOT copied to the new recurrence todo automatically |
| REC-05 | Deleting a todo deletes its RecurrenceConfig |

---

## Reminder Rules

| ID | Rule |
|---|---|
| REM-01 | fire_at must be in the future at creation time |
| REM-02 | A todo can have at most 10 reminders |
| REM-03 | Fired reminders (fired = true) are not re-fired |
| REM-04 | Deleting a todo cancels (deletes) all its reminders |
| REM-05 | Scheduler polls for due reminders every 30 seconds |

---

## File Attachment Rules

| ID | Rule |
|---|---|
| FILE-01 | Max file size: 10MB (10,485,760 bytes) |
| FILE-02 | Allowed MIME types: image/jpeg, image/png, image/gif, image/webp, application/pdf, application/vnd.openxmlformats-officedocument.wordprocessingml.document, text/plain |
| FILE-03 | A todo can have at most 10 attachments |
| FILE-04 | Only the owning user can download or delete a file |
| FILE-05 | Deleting a todo deletes all its file attachments from disk and DB |
| FILE-06 | Storage path format: `uploads/{user_id}/{todo_id}/{uuid}_{original_filename}` |

---

## Search Rules

| ID | Rule |
|---|---|
| SRCH-01 | Search is scoped to the authenticated user's non-deleted todos only |
| SRCH-02 | Search uses Elasticsearch/OpenSearch with stemming and relevance scoring |
| SRCH-03 | Search index synced on todo create, update, and soft delete |
| SRCH-04 | Soft-deleted todos excluded from search index via deleted_at filter |

---

## Account Deletion Rules

| ID | Rule |
|---|---|
| ACC-01 | Account deletion sets is_active = false, deleted_at = now() |
| ACC-02 | Soft-deleted user data retained for 30-day grace period |
| ACC-03 | After grace period, scheduler purges: todos (soft-deleted), files (disk + DB), reminders, notifications, refresh tokens |
| ACC-04 | User cannot log in after soft delete |

---

## Notification Rules

| ID | Rule |
|---|---|
| NOTIF-01 | Notifications delivered via WebSocket if user has active connection |
| NOTIF-02 | If user offline, notification stored (delivered = false) and sent on next WebSocket connect |
| NOTIF-03 | Pending notifications fetched and pushed on WebSocket connection establishment |
