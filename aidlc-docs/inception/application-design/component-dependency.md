# Component Dependencies — Todo App

## Dependency Matrix

| Component | Depends On | Communication |
|---|---|---|
| vue-app | auth-service | REST |
| vue-app | todo-service | REST |
| vue-app | file-service | REST |
| vue-app | notification-service | WebSocket |
| todo-service | file-service | Internal REST |
| todo-service | scheduler-service | Internal REST |
| scheduler-service | notification-service | Internal REST |
| scheduler-service | todo-service | Internal REST |
| auth-service | — | Standalone |
| file-service | — | Standalone |

## Dependency Diagram

```
                    +------------+
                    |  vue-app   |
                    +-----+------+
                          |
          +---------------+------------------+----------+
          |               |                  |          |
          v               v                  v          v
   +-----------+   +-----------+   +----------+  +------------------+
   |auth-svc   |   |todo-svc   |   |file-svc  |  |notification-svc  |
   +-----------+   +-----+-----+   +----------+  +--------+---------+
                         |                                 ^
                    +----+----+                            |
                    |         |                            |
                    v         v                            |
             +----------+  +----------+                   |
             |file-svc  |  |scheduler |-------------------+
             +----------+  |  -svc    |
                           +----+-----+
                                |
                                v
                          +----------+
                          |todo-svc  |
                          |(recur.)  |
                          +----------+
```

## Data Flow: Create Todo with Reminder

```
vue-app
  --> POST /todos (todo-service)
        --> INSERT todo (PostgreSQL)
        --> POST /reminders (scheduler-service)
              --> INSERT reminder (PostgreSQL)
              --> schedule goroutine
  <-- 201 TodoResponse

[at reminder fire time]
scheduler-service goroutine
  --> POST /internal/events (notification-service)
        --> push to user WebSocket (vue-app)
```

## Data Flow: Complete Recurring Todo

```
vue-app
  --> POST /todos/:id/complete (todo-service)
        --> UPDATE todo status = done
        --> POST /todos/:id/complete (scheduler-service)
              --> compute next occurrence date
              --> POST /todos (todo-service)
                    --> INSERT new todo
              --> DELETE old reminders
  <-- 200 OK
```

## Coupling Notes

- **auth-service** is fully decoupled at runtime — JWT validation is stateless (shared secret)
- **scheduler-service ↔ todo-service** have a circular dependency (scheduler creates todos for recurrence); this is intentional and managed via internal REST to avoid tight coupling
- **file-service** is a pure leaf service with no outbound dependencies
- All inter-service calls use internal Docker network (not exposed publicly)
