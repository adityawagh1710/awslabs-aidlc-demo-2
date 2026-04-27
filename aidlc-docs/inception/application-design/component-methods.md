# Component Methods — Todo App

> Detailed business rules are defined later in Functional Design (per-unit, CONSTRUCTION phase).

---

## auth-service

### Handler Layer
| Method | Signature | Purpose |
|---|---|---|
| RegisterHandler | `POST /auth/register (RegisterRequest) → 201 / 4xx` | Validate input, create user, return tokens |
| LoginHandler | `POST /auth/login (LoginRequest) → TokenResponse / 4xx` | Validate credentials, issue JWT pair |
| RefreshHandler | `POST /auth/refresh (RefreshRequest) → TokenResponse / 401` | Validate refresh token, issue new access token |
| LogoutHandler | `POST /auth/logout (AuthHeader) → 204` | Revoke refresh token |
| MFAEnrollHandler | `POST /auth/mfa/enroll (AuthHeader) → MFAEnrollResponse` | Generate TOTP secret for user |
| MFAVerifyHandler | `POST /auth/mfa/verify (MFAVerifyRequest) → 200 / 401` | Validate TOTP code |

### Service Layer
| Method | Signature | Purpose |
|---|---|---|
| RegisterUser | `(ctx, RegisterInput) → (User, error)` | Hash password, persist user, emit verification |
| AuthenticateUser | `(ctx, email, password) → (TokenPair, error)` | Verify credentials, check MFA, issue tokens |
| RefreshTokens | `(ctx, refreshToken) → (TokenPair, error)` | Validate and rotate refresh token |
| RevokeToken | `(ctx, refreshToken) → error` | Invalidate refresh token |
| EnrollMFA | `(ctx, userID) → (secret, qrURL, error)` | Generate and store TOTP secret |
| VerifyMFA | `(ctx, userID, code) → error` | Validate TOTP code against stored secret |

### Repository Layer
| Method | Signature | Purpose |
|---|---|---|
| CreateUser | `(ctx, User) → (User, error)` | Persist new user record |
| FindUserByEmail | `(ctx, email) → (User, error)` | Lookup user for authentication |
| SaveRefreshToken | `(ctx, userID, token, expiry) → error` | Store refresh token |
| DeleteRefreshToken | `(ctx, token) → error` | Remove refresh token on logout |
| UpdateMFASecret | `(ctx, userID, secret) → error` | Persist TOTP secret |

---

## todo-service

### Handler Layer
| Method | Signature | Purpose |
|---|---|---|
| CreateTodoHandler | `POST /todos (CreateTodoRequest) → 201 TodoResponse` | Create new todo for authenticated user |
| ListTodosHandler | `GET /todos (filters) → []TodoResponse` | List todos with optional filters |
| GetTodoHandler | `GET /todos/:id → TodoResponse / 404` | Get single todo (ownership enforced) |
| UpdateTodoHandler | `PUT /todos/:id (UpdateTodoRequest) → TodoResponse` | Update todo fields |
| DeleteTodoHandler | `DELETE /todos/:id → 204` | Delete todo and its attachments |
| SearchTodosHandler | `GET /todos/search?q= → []TodoResponse` | Full-text search |
| CreateTagHandler | `POST /tags (TagRequest) → TagResponse` | Create tag |
| ListTagsHandler | `GET /tags → []TagResponse` | List user's tags |
| DeleteTagHandler | `DELETE /tags/:id → 204` | Delete tag |

### Service Layer
| Method | Signature | Purpose |
|---|---|---|
| CreateTodo | `(ctx, userID, CreateTodoInput) → (Todo, error)` | Validate and persist todo |
| ListTodos | `(ctx, userID, TodoFilter) → ([]Todo, error)` | Fetch todos with filters |
| GetTodo | `(ctx, userID, todoID) → (Todo, error)` | Fetch todo, enforce ownership |
| UpdateTodo | `(ctx, userID, todoID, UpdateTodoInput) → (Todo, error)` | Update, enforce ownership |
| DeleteTodo | `(ctx, userID, todoID) → error` | Delete todo, cascade attachments |
| SearchTodos | `(ctx, userID, query) → ([]Todo, error)` | Full-text search scoped to user |
| CreateTag | `(ctx, userID, name) → (Tag, error)` | Create tag for user |
| DeleteTag | `(ctx, userID, tagID) → error` | Delete tag, unassign from todos |

### Repository Layer
| Method | Signature | Purpose |
|---|---|---|
| InsertTodo | `(ctx, Todo) → (Todo, error)` | Persist todo |
| FindTodosByUser | `(ctx, userID, TodoFilter) → ([]Todo, error)` | Query todos with filters |
| FindTodoByID | `(ctx, todoID) → (Todo, error)` | Fetch single todo |
| UpdateTodo | `(ctx, Todo) → (Todo, error)` | Update todo record |
| DeleteTodo | `(ctx, todoID) → error` | Delete todo record |
| SearchTodos | `(ctx, userID, query) → ([]Todo, error)` | Parameterised full-text search |
| InsertTag | `(ctx, Tag) → (Tag, error)` | Persist tag |
| FindTagsByUser | `(ctx, userID) → ([]Tag, error)` | List user tags |
| DeleteTag | `(ctx, tagID) → error` | Delete tag |

---

## scheduler-service

### Handler Layer
| Method | Signature | Purpose |
|---|---|---|
| CreateReminderHandler | `POST /reminders (ReminderRequest) → 201` | Schedule a reminder for a todo |
| DeleteReminderHandler | `DELETE /reminders/:id → 204` | Cancel a reminder |
| SetRecurrenceHandler | `POST /todos/:id/recurrence (RecurrenceRequest) → 200` | Configure recurrence for a todo |
| CompleteTodoHandler | `POST /todos/:id/complete → 200` | Mark done, trigger next recurrence if applicable |

### Service Layer
| Method | Signature | Purpose |
|---|---|---|
| ScheduleReminder | `(ctx, todoID, userID, fireAt) → (Reminder, error)` | Persist and schedule reminder |
| CancelReminder | `(ctx, reminderID) → error` | Remove reminder from scheduler |
| SetRecurrence | `(ctx, todoID, RecurrenceConfig) → error` | Persist recurrence config |
| HandleTodoCompletion | `(ctx, todoID, userID) → error` | Generate next occurrence if recurring |
| RunScheduler | `(ctx) → error` | Background goroutine: poll and fire due reminders |

### Repository Layer
| Method | Signature | Purpose |
|---|---|---|
| InsertReminder | `(ctx, Reminder) → (Reminder, error)` | Persist reminder |
| FindDueReminders | `(ctx, now) → ([]Reminder, error)` | Fetch reminders due for firing |
| DeleteReminder | `(ctx, reminderID) → error` | Remove reminder |
| UpsertRecurrence | `(ctx, RecurrenceConfig) → error` | Save/update recurrence config |
| FindRecurrence | `(ctx, todoID) → (RecurrenceConfig, error)` | Load recurrence config |

---

## file-service

### Handler Layer
| Method | Signature | Purpose |
|---|---|---|
| UploadFileHandler | `POST /files (multipart, todoID) → FileResponse` | Validate and store file |
| DownloadFileHandler | `GET /files/:id → file stream` | Serve file (ownership enforced) |
| DeleteFileHandler | `DELETE /files/:id → 204` | Delete file from disk and DB |

### Service Layer
| Method | Signature | Purpose |
|---|---|---|
| UploadFile | `(ctx, userID, todoID, FileInput) → (File, error)` | Validate type/size, write to disk, persist metadata |
| GetFile | `(ctx, userID, fileID) → (FilePath, error)` | Verify ownership, return path |
| DeleteFile | `(ctx, userID, fileID) → error` | Delete from disk and DB |

### Repository Layer
| Method | Signature | Purpose |
|---|---|---|
| InsertFile | `(ctx, File) → (File, error)` | Persist file metadata |
| FindFileByID | `(ctx, fileID) → (File, error)` | Fetch file metadata |
| DeleteFile | `(ctx, fileID) → error` | Remove file metadata |
| FindFilesByTodo | `(ctx, todoID) → ([]File, error)` | List attachments for a todo |

---

## notification-service

### Handler Layer
| Method | Signature | Purpose |
|---|---|---|
| WebSocketHandler | `GET /ws (AuthHeader) → WebSocket upgrade` | Establish authenticated WS connection |
| IngestEventHandler | `POST /internal/events (NotificationEvent) → 202` | Receive event from scheduler-service |

### Service Layer
| Method | Signature | Purpose |
|---|---|---|
| RegisterConnection | `(userID, conn) → error` | Track active WS connection for user |
| RemoveConnection | `(userID) → error` | Clean up on disconnect |
| DeliverNotification | `(ctx, userID, Notification) → error` | Push to live connection or store for later |
| StoreUndelivered | `(ctx, userID, Notification) → error` | Persist if user offline |
| GetPendingNotifications | `(ctx, userID) → ([]Notification, error)` | Return undelivered on reconnect |

### Repository Layer
| Method | Signature | Purpose |
|---|---|---|
| InsertNotification | `(ctx, Notification) → error` | Persist undelivered notification |
| FindPendingByUser | `(ctx, userID) → ([]Notification, error)` | Fetch undelivered notifications |
| MarkDelivered | `(ctx, notificationID) → error` | Mark notification as delivered |
