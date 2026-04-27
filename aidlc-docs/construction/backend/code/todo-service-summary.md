# todo-service — Code Summary

## Purpose
Core todo domain: CRUD, tags, search (Elasticsearch), status transition enforcement, outbox-based ES sync, circuit breakers for inter-service calls.

## Key Files

| File | Purpose |
|---|---|
| `cmd/main.go` | Fiber app, outbox worker goroutine, graceful shutdown |
| `internal/model/todo.go` | Todo, Tag, SearchOutbox GORM models; ValidTransitions map |
| `internal/repository/todo_repository.go` | Todo CRUD + soft delete + filter queries |
| `internal/repository/tag_repository.go` | Tag CRUD scoped to user |
| `internal/repository/outbox_repository.go` | SearchOutbox insert/find-unprocessed/mark-processed |
| `internal/service/todo_service.go` | Create, List, Get, Update (transition enforcement), Delete, Search |
| `internal/service/tag_service.go` | Tag create/list/delete with ownership check |
| `internal/service/outbox_worker.go` | 5s poll goroutine: sync outbox events to Elasticsearch |
| `internal/service/circuit_breaker.go` | gobreaker + resty wrappers for file/scheduler calls |
| `internal/handler/todo_handler.go` | Todo CRUD + search HTTP handlers |
| `internal/handler/tag_handler.go` | Tag HTTP handlers |
| `migrations/` | todos, tags, todo_tags, search_outbox SQL migrations |
| `Dockerfile` | Multi-stage: golang:1.22-alpine → alpine:3.19 |

## Tests
| File | Type | Coverage |
|---|---|---|
| `repository/todo_repository_test.go` | Integration (Testcontainers) | Insert, FindByUser, Update, SoftDelete, FilterByStatus |
| `repository/tag_repository_test.go` | Integration (Testcontainers) | Insert, FindByUser, Delete |
| `service/todo_service_test.go` | Unit (mocks) | Create, invalid transition, forbidden access |
| `service/todo_service_pbt_test.go` | PBT (gopter) | All 9 status transition combinations |
| `handler/todo_handler_test.go` | E2E (Testcontainers) | create→list→get→update→invalid transition→delete→404, 401 |

## Port: 3001
