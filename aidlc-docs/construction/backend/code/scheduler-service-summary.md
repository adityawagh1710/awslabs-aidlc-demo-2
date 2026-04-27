# scheduler-service — Code Summary

## Purpose
Manages reminders (schedule, cancel, fire) and recurring task generation using an in-process goroutine scheduler and cron expressions.

## Key Files

| File | Purpose |
|---|---|
| `cmd/main.go` | Fiber app, RunScheduler goroutine, graceful shutdown |
| `internal/model/scheduler.go` | Reminder, RecurrenceConfig GORM models |
| `internal/repository/scheduler_repository.go` | ReminderRepository + RecurrenceRepository |
| `internal/service/recurrence_engine.go` | Pure `NextOccurrence(cron, from)` — no side effects |
| `internal/service/scheduler_service.go` | ScheduleReminder (10-limit), CancelReminder, SetRecurrence, HandleTodoCompletion, RunScheduler (30s poll) |
| `internal/handler/scheduler_handler.go` | Internal REST handlers (API key protected) |
| `internal/middleware/middleware.go` | ErrorHandler, Recover, InternalAPIKey |
| `migrations/` | reminders, recurrence_configs SQL migrations |
| `Dockerfile` | Multi-stage: golang:1.22-alpine → alpine:3.19 |

## Tests
| File | Type | Coverage |
|---|---|---|
| `service/recurrence_engine_pbt_test.go` | PBT (gopter) | next > from, deterministic, invalid cron → ErrInvalidCron |

## Port: 3004 (internal only, not exposed via Traefik)
