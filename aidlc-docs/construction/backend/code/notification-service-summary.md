# notification-service — Code Summary

## Purpose
Real-time in-app notifications via WebSocket. Delivers to live connections; stores undelivered notifications for retrieval on reconnect.

## Key Files

| File | Purpose |
|---|---|
| `cmd/main.go` | Fiber app, graceful shutdown |
| `internal/model/notification.go` | Notification GORM model |
| `internal/repository/notification_repository.go` | Insert, FindPendingByUser, MarkDelivered |
| `internal/service/hub.go` | Thread-safe WebSocket connection registry (sync.RWMutex) |
| `internal/service/notification_service.go` | Deliver (live WS → persist); GetPending |
| `internal/handler/ws_handler.go` | JWT auth (header or ?token=), register connection, push pending on connect, read loop |
| `internal/handler/event_handler.go` | Internal `POST /internal/events` (API key protected) |
| `internal/middleware/middleware.go` | ErrorHandler, Recover, InternalAPIKey |
| `migrations/` | notifications SQL migration |
| `Dockerfile` | Multi-stage: golang:1.22-alpine → alpine:3.19 |

## Tests
| File | Type | Coverage |
|---|---|---|
| `service/notification_service_test.go` | Unit (mocks) | Offline user stores undelivered, GetPending returns results |

## Port: 3003 (WebSocket: `GET /ws`, Internal: `POST /internal/events`)
