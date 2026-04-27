# Build stage: build all 5 service binaries via the Go workspace.
FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.work ./
COPY auth-service/go.mod auth-service/go.sum ./auth-service/
COPY todo-service/go.mod todo-service/go.sum ./todo-service/
COPY file-service/go.mod file-service/go.sum ./file-service/
COPY notification-service/go.mod notification-service/go.sum ./notification-service/
COPY scheduler-service/go.mod scheduler-service/go.sum ./scheduler-service/
RUN go work sync

COPY auth-service ./auth-service
COPY todo-service ./todo-service
COPY file-service ./file-service
COPY notification-service ./notification-service
COPY scheduler-service ./scheduler-service

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o /out/auth-service          ./auth-service/cmd \
 && go build -ldflags="-s -w" -o /out/todo-service          ./todo-service/cmd \
 && go build -ldflags="-s -w" -o /out/file-service          ./file-service/cmd \
 && go build -ldflags="-s -w" -o /out/notification-service  ./notification-service/cmd \
 && go build -ldflags="-s -w" -o /out/scheduler-service     ./scheduler-service/cmd

# Runtime stage: supervisord runs all 5 binaries; each has its own workdir
# so the `file://migrations` relative path used by golang-migrate resolves.
FROM alpine:3.19
RUN apk --no-cache add ca-certificates wget supervisor

WORKDIR /app

COPY --from=builder /out/auth-service          /app/auth/auth-service
COPY --from=builder /out/todo-service          /app/todo/todo-service
COPY --from=builder /out/file-service          /app/file/file-service
COPY --from=builder /out/notification-service  /app/notification/notification-service
COPY --from=builder /out/scheduler-service     /app/scheduler/scheduler-service

COPY --from=builder /src/auth-service/migrations          /app/auth/migrations
COPY --from=builder /src/todo-service/migrations          /app/todo/migrations
COPY --from=builder /src/file-service/migrations          /app/file/migrations
COPY --from=builder /src/notification-service/migrations  /app/notification/migrations
COPY --from=builder /src/scheduler-service/migrations     /app/scheduler/migrations

COPY supervisord.conf /etc/supervisor/conf.d/todo-backend.conf

EXPOSE 3000 3001 3002 3003 3004

ENTRYPOINT ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/conf.d/todo-backend.conf"]
