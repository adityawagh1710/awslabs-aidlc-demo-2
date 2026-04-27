# Build Instructions — Todo App (Backend)

## Prerequisites

| Requirement | Version |
|---|---|
| Go | 1.22+ |
| Docker | 24+ |
| Docker Compose | v2 |
| Make (optional) | any |

## 1. Clone Repositories

```bash
# Each service is in its own repo (polyrepo)
git clone https://github.com/<org>/todo-auth-service
git clone https://github.com/<org>/todo-todo-service
git clone https://github.com/<org>/todo-scheduler-service
git clone https://github.com/<org>/todo-file-service
git clone https://github.com/<org>/todo-notification-service
git clone https://github.com/<org>/todo-infra
```

## 2. Configure Environment

```bash
cd todo-infra
cp .env.example .env
# Edit .env — fill in all <placeholder> values with real secrets
```

## 3. Build Each Service Locally

```bash
# Run in each service directory
go mod download
go build -ldflags="-s -w" -o bin/<service-name> ./cmd/main.go
```

## 4. Build Docker Images

```bash
# In each service directory
docker build -t ghcr.io/<org>/<service-name>:<git-sha> .

# Or build all via infra docker-compose (uses pre-built images from ghcr.io)
cd todo-infra
docker compose pull   # pull latest images
```

## 5. Start Full Stack

```bash
cd todo-infra
docker compose up -d

# Verify all containers healthy
docker compose ps
```

**Expected**: All 13 containers in `healthy` or `running` state.

## 6. Verify Services

```bash
curl https://<domain>/health          # auth-service
curl https://<domain>/todos/health    # todo-service (via Traefik)
# Or locally:
curl http://localhost:3000/health
curl http://localhost:3001/health
curl http://localhost:3002/health
curl http://localhost:3003/health
curl http://localhost:3004/health
```

## Build Artifacts

| Artifact | Location |
|---|---|
| Go binaries | `bin/` in each service repo |
| Docker images | `ghcr.io/<org>/<service>:<sha>` |
| Migrations | `migrations/` in each service repo (run on startup) |

## Troubleshooting

**`go mod download` fails**: Check Go version (`go version`) — requires 1.22+

**Docker build fails**: Ensure Docker daemon is running; check `go.sum` is committed

**Container unhealthy**: Check logs with `docker compose logs <service>`; verify `.env` values are set correctly

**Elasticsearch fails to start**: Increase Docker memory limit to at least 4GB (ES requires 2GB heap)
