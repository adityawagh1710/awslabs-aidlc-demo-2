# auth-service — Code Summary

## Purpose
Handles all user identity concerns: registration, login, JWT token lifecycle, MFA (TOTP), and brute-force protection.

## Key Files

| File | Purpose |
|---|---|
| `cmd/main.go` | Fiber app bootstrap, middleware chain, graceful shutdown |
| `internal/model/user.go` | User, RefreshToken GORM models (soft delete) |
| `internal/middleware/jwt.go` | JWT validation + Redis blacklist check |
| `internal/middleware/ratelimit.go` | Redis-backed sliding window rate limiter |
| `internal/middleware/errorhandler.go` | Global error handler (fail-closed, no stack traces) |
| `internal/middleware/logger.go` | zerolog request logger with correlation ID |
| `internal/repository/user_repository.go` | User CRUD + soft delete |
| `internal/repository/token_repository.go` | Refresh token store/find/delete |
| `internal/service/auth_service.go` | Register, Login (lockout), Refresh (rotation), Logout (blacklist), MFA enroll/verify |
| `internal/handler/auth_handler.go` | HTTP handlers for all auth endpoints |
| `migrations/` | users, refresh_tokens SQL migrations |
| `Dockerfile` | Multi-stage: golang:1.22-alpine → alpine:3.19 |

## Tests
| File | Type | Coverage |
|---|---|---|
| `repository/user_repository_test.go` | Integration (Testcontainers) | Create, FindByEmail, UpdateMFA, SoftDelete |
| `repository/token_repository_test.go` | Integration (Testcontainers) | Save, Find, Delete, Expired token |
| `service/auth_service_test.go` | Unit (mocks + miniredis) | Duplicate email, wrong password, account lockout |
| `service/auth_service_pbt_test.go` | PBT (gopter) | JWT round-trip, hash determinism |
| `handler/auth_handler_test.go` | E2E (Testcontainers) | register→login→refresh→logout, duplicate register |

## Port: 3000
