# Domain Entities — Unit 1: backend

## Entity: User

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK, not null |
| email | string | unique, not null, max 255 |
| password_hash | string | not null |
| mfa_secret | string | nullable (set on MFA enroll) |
| mfa_enabled | bool | default false |
| is_active | bool | default true |
| deleted_at | timestamp | nullable (soft delete) |
| created_at | timestamp | not null |
| updated_at | timestamp | not null |

**Soft delete**: `deleted_at IS NULL` filter applied on all queries.
**Account deletion**: sets `is_active = false`, `deleted_at = now()`. Grace period purge handled by scheduler.

---

## Entity: RefreshToken

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK → User.id |
| token_hash | string | unique, not null |
| expires_at | timestamp | not null |
| created_at | timestamp | not null |

---

## Entity: Todo

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK → User.id, not null |
| title | string | not null, max 255 |
| description | string | nullable, max 5000 |
| status | enum | pending \| in_progress \| done, default pending |
| priority | enum | low \| medium \| high, default medium |
| due_date | timestamp | nullable |
| deleted_at | timestamp | nullable (soft delete) |
| created_at | timestamp | not null |
| updated_at | timestamp | not null |

**Soft delete**: `deleted_at IS NULL` filter on all queries.
**Status transitions**: pending → in_progress → done (enforced, no skipping).

---

## Entity: Tag

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK → User.id, not null |
| name | string | not null, max 50, unique per user |
| created_at | timestamp | not null |

---

## Entity: TodoTag (join table)

| Field | Type | Constraints |
|---|---|---|
| todo_id | UUID | FK → Todo.id |
| tag_id | UUID | FK → Tag.id |
| PK | (todo_id, tag_id) | composite |

---

## Entity: Reminder

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| todo_id | UUID | FK → Todo.id |
| user_id | UUID | FK → User.id |
| fire_at | timestamp | not null |
| fired | bool | default false |
| created_at | timestamp | not null |

---

## Entity: RecurrenceConfig

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| todo_id | UUID | FK → Todo.id, unique |
| cron_expression | string | not null (standard cron format) |
| next_occurrence | timestamp | not null |
| created_at | timestamp | not null |
| updated_at | timestamp | not null |

**Cron-style recurrence**: full cron expression support (e.g., `0 9 * * 1` = every Monday at 9am).

---

## Entity: FileAttachment

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| todo_id | UUID | FK → Todo.id |
| user_id | UUID | FK → User.id |
| filename | string | not null, max 255 |
| storage_path | string | not null |
| mime_type | string | not null |
| size_bytes | int | not null, max 10485760 (10MB) |
| created_at | timestamp | not null |

**Allowed MIME types**: image/jpeg, image/png, image/gif, image/webp, application/pdf, application/vnd.openxmlformats-officedocument.wordprocessingml.document, text/plain

---

## Entity: Notification

| Field | Type | Constraints |
|---|---|---|
| id | UUID | PK |
| user_id | UUID | FK → User.id |
| todo_id | UUID | FK → Todo.id, nullable |
| message | string | not null, max 500 |
| delivered | bool | default false |
| created_at | timestamp | not null |

---

## Entity Relationships

```
User ──< Todo ──< TodoTag >── Tag
User ──< Tag
Todo ──< Reminder
Todo ──1 RecurrenceConfig
Todo ──< FileAttachment
Todo ──< Notification (via user_id)
User ──< RefreshToken
User ──< Notification
```

---

## Search Index (Elasticsearch/OpenSearch)

**Index**: `todos`

| Field | Type | Notes |
|---|---|---|
| id | keyword | |
| user_id | keyword | filter |
| title | text | analyzed, stemming |
| description | text | analyzed, stemming |
| tags | keyword[] | filter + search |
| status | keyword | filter |
| priority | keyword | filter |
| due_date | date | range filter |
| deleted_at | date | exclude deleted |

Documents synced from PostgreSQL on create/update/soft-delete.
