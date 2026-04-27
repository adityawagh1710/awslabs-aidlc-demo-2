# Unit Test Execution Instructions

## Prerequisites
- Go 1.22+
- Docker (required for Testcontainers)
- Docker daemon running

## Run All Tests Per Service

```bash
# In each service directory:
go test ./... -v -timeout 120s
```

## Run Specific Test Types

### Unit tests only (no Testcontainers — fast)
```bash
go test ./internal/service/... -v
```

### Repository tests (Testcontainers — requires Docker)
```bash
go test ./internal/repository/... -v -timeout 120s
```

### Handler E2E tests (Testcontainers — requires Docker)
```bash
go test ./internal/handler/... -v -timeout 120s
```

### PBT tests only
```bash
go test ./internal/service/... -run PBT -v
```

## Per-Service Test Summary

| Service | Unit Tests | Integration (Testcontainers) | PBT |
|---|---|---|---|
| auth-service | service/auth_service_test.go | repository/*_test.go, handler/auth_handler_test.go | service/auth_service_pbt_test.go |
| todo-service | service/todo_service_test.go | repository/*_test.go, handler/todo_handler_test.go | service/todo_service_pbt_test.go |
| scheduler-service | — | — | service/recurrence_engine_pbt_test.go |
| file-service | service/file_service_test.go | — | — |
| notification-service | service/notification_service_test.go | — | — |

## Vulnerability Scan

```bash
# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# Run in each service directory
govulncheck ./...
```

**Expected**: No known vulnerabilities. Address any HIGH/CRITICAL findings before deployment.

## Expected Results
- All tests pass (`PASS`)
- No data races (`go test -race ./...`)
- No vulnerabilities reported by `govulncheck`
