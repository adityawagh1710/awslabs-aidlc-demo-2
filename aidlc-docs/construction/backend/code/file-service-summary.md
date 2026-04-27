# file-service — Code Summary

## Purpose
Manages file attachments: upload (MIME/size/count validation), download (ownership enforced), delete (disk + DB cleanup). Files stored on Docker named volume.

## Key Files

| File | Purpose |
|---|---|
| `cmd/main.go` | Fiber app (11MB body limit), graceful shutdown |
| `internal/model/file.go` | FileAttachment model; AllowedMimeTypes allowlist; MaxFileSizeBytes (10MB); MaxAttachmentsPerTodo (10) |
| `internal/repository/file_repository.go` | File metadata CRUD + CountByTodo |
| `internal/service/file_service.go` | Upload (validate → write disk → persist), GetPath (ownership), Delete (disk + DB) |
| `internal/handler/file_handler.go` | Upload (multipart), Download (SendFile), Delete |
| `internal/middleware/middleware.go` | JWT auth, ErrorHandler, Recover |
| `migrations/` | file_attachments SQL migration |
| `Dockerfile` | Multi-stage: golang:1.22-alpine → alpine:3.19 |

## Tests
| File | Type | Coverage |
|---|---|---|
| `service/file_service_test.go` | Unit (mocks) | Invalid MIME, too large, success (disk write verified), forbidden delete |

## Port: 3002
