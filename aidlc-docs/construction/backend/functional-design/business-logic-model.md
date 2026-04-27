# Business Logic Model — Unit 1: backend

## Core Flows

---

### Flow 1: User Registration

```
Input: email, password
1. Validate email format and max length
2. Check email uniqueness (AUTH-01)
3. Validate password length >= 8 (AUTH-02)
4. Check password against breached list (AUTH-02)
5. Hash password with argon2id (AUTH-03)
6. Persist User record (is_active=true, mfa_enabled=false)
7. Issue access token (1h) + refresh token (30d) (AUTH-04)
8. Persist refresh token hash
Output: TokenPair
```

---

### Flow 2: User Login

```
Input: email, password, [mfa_code]
1. Lookup user by email; reject if not found or deleted_at IS NOT NULL (AUTH-08)
2. Check login attempt count — reject with 429 if >= 5 in last 15 min (AUTH-06)
3. Verify password hash
4. On failure: increment attempt counter → return 401
5. On success: reset attempt counter
6. If mfa_enabled=true: validate mfa_code via TOTP (AUTH-07) → return 401 if invalid
7. Issue access token (1h) + refresh token (30d)
8. Persist refresh token hash
Output: TokenPair
```

---

### Flow 3: Token Refresh

```
Input: refresh_token
1. Hash incoming token; lookup in DB
2. Reject if not found or expired
3. Delete old refresh token (single-use rotation, AUTH-05)
4. Issue new access token + new refresh token
5. Persist new refresh token hash
Output: TokenPair
```

---

### Flow 4: Create Todo

```
Input: userID, title, description?, priority?, due_date?, tag_ids?, recurrence_cron?
1. Validate title (required, max 255) (TODO-02)
2. Validate description max 5000 (TODO-03)
3. Validate due_date is in future if provided (TODO-08)
4. Validate tag_ids belong to userID (TODO-01)
5. Validate tag count <= 20 (TAG-04)
6. Persist Todo (status=pending)
7. Persist TodoTag rows
8. If recurrence_cron provided: validate cron expression; persist RecurrenceConfig with next_occurrence computed
9. Sync todo document to Elasticsearch (SRCH-03)
Output: Todo
```

---

### Flow 5: Update Todo Status

```
Input: userID, todoID, newStatus
1. Fetch todo; verify user_id == userID (TODO-01)
2. Verify todo not soft-deleted
3. Enforce transition rules (TODO-04):
   - pending → in_progress: allowed
   - in_progress → done: allowed
   - Any other transition: reject with 422
4. If newStatus == done AND RecurrenceConfig exists:
   a. Compute next_occurrence from cron expression
   b. Create new Todo (title, description, priority, tags copied; status=pending)
   c. Persist new RecurrenceConfig for new todo
   d. Sync new todo to Elasticsearch
5. Update todo status + updated_at
6. Sync updated document to Elasticsearch
Output: Todo
```

---

### Flow 6: Soft Delete Todo

```
Input: userID, todoID
1. Fetch todo; verify user_id == userID (TODO-01)
2. Set deleted_at = now() (TODO-06)
3. Cancel all reminders for todo (delete Reminder rows) (TODO-07)
4. Delete RecurrenceConfig if exists (TODO-07)
5. Delete all FileAttachment records + files from disk (TODO-07, FILE-05)
6. Update Elasticsearch document with deleted_at (SRCH-04)
Output: void
```

---

### Flow 7: Schedule Reminder

```
Input: userID, todoID, fire_at
1. Fetch todo; verify ownership (TODO-01)
2. Validate fire_at is in future (REM-01)
3. Count existing reminders for todo; reject if >= 10 (REM-02)
4. Persist Reminder (fired=false)
5. Goroutine scheduler picks it up on next poll cycle (REM-05)
Output: Reminder
```

---

### Flow 8: Fire Reminder (Scheduler Goroutine)

```
Runs every 30 seconds (REM-05):
1. Query Reminder WHERE fire_at <= now() AND fired = false
2. For each due reminder:
   a. Mark fired = true
   b. Create Notification record (delivered=false)
   c. POST /internal/events to notification-service with {userID, todoID, message}
3. notification-service:
   a. If user has active WebSocket: push notification, mark delivered=true
   b. Else: store as undelivered (NOTIF-02)
```

---

### Flow 9: Upload File Attachment

```
Input: userID, todoID, file (multipart)
1. Fetch todo; verify ownership (TODO-01)
2. Validate MIME type against allowlist (FILE-02)
3. Validate file size <= 10MB (FILE-01)
4. Count existing attachments; reject if >= 10 (FILE-03)
5. Generate storage path: uploads/{userID}/{todoID}/{uuid}_{filename} (FILE-06)
6. Write file to disk
7. Persist FileAttachment metadata
Output: FileAttachment
```

---

### Flow 10: Search Todos

```
Input: userID, query string, filters (status?, priority?, tag_ids?, due_date_range?)
1. Build Elasticsearch query:
   - Filter: user_id == userID AND deleted_at IS NULL
   - Full-text: match on title, description (stemmed, scored)
   - Optional filters: status, priority, tags, due_date range
2. Execute query against todos index
3. Return results ordered by relevance score
Output: []Todo
```

---

## Error Handling Model

| Scenario | HTTP Status | Response |
|---|---|---|
| Validation failure | 422 | Generic field error message |
| Unauthorized (no/invalid token) | 401 | "Unauthorized" |
| Forbidden (wrong owner) | 403 | "Forbidden" |
| Resource not found | 404 | "Not found" |
| Rate limit exceeded | 429 | "Too many requests" |
| Invalid status transition | 422 | "Invalid status transition" |
| Internal error | 500 | "Internal server error" (no stack trace) |

All error responses are generic — no internal details, stack traces, or DB errors exposed (SECURITY-09, SECURITY-15).
