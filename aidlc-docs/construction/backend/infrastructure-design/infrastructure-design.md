# Infrastructure Design — Unit 1: backend

## Deployment Target
Self-hosted server/VM, Docker Compose v2, all services containerised.

---

## Infrastructure Components

| Component | Image | Purpose |
|---|---|---|
| traefik | `traefik:v3` | API gateway — routing, TLS termination, Let's Encrypt, observability |
| auth-service | `todo-auth-service:latest` | JWT auth, MFA |
| todo-service | `todo-todo-service:latest` | Todo CRUD, tags, search |
| scheduler-service | `todo-scheduler-service:latest` | Reminders, recurrence |
| file-service | `todo-file-service:latest` | File attachments |
| notification-service | `todo-notification-service:latest` | WebSocket notifications |
| postgres | `postgres:16-alpine` | Primary database |
| redis | `redis:7-alpine` | Cache, rate limiting, token blacklist |
| elasticsearch | `docker.elastic.co/elasticsearch/elasticsearch:8.13.0` | Search index |
| prometheus | `prom/prometheus:v2.51.0` | Metrics scraping |
| grafana | `grafana/grafana:10.4.0` | Metrics dashboards + alerting |
| certbot | `certbot/certbot` | Let's Encrypt certificate renewal (sidecar) |

---

## Networking

- All services on a single Docker bridge network: `todo-network`
- Only Traefik exposes ports 80 and 443 to the host
- All other services are internal-only (no host port bindings)
- Inter-service calls use Docker service names (e.g., `http://todo-service:3001`)
- Prometheus and Grafana accessible via Traefik on internal routes (auth-protected)

---

## Routing (Traefik)

| Route | Backend Service | Notes |
|---|---|---|
| `/auth/*` | auth-service:3000 | Public (register, login); rate limited |
| `/todos/*` | todo-service:3001 | JWT required |
| `/tags/*` | todo-service:3001 | JWT required |
| `/files/*` | file-service:3002 | JWT required |
| `/ws` | notification-service:3003 | WebSocket upgrade |
| `/internal/*` | — | Blocked at Traefik (internal only) |
| `/metrics` | prometheus:9090 | Internal, basic auth |
| `/grafana/*` | grafana:3000 | Internal, Grafana auth |

---

## TLS — Let's Encrypt via Traefik

- Traefik handles ACME challenge automatically
- Certificates stored in a Docker named volume: `traefik-certs`
- Auto-renewal handled by Traefik (no Certbot sidecar needed with Traefik v3)
- HTTP → HTTPS redirect enforced at Traefik level
- HSTS header set via Traefik middleware

---

## Storage Volumes

| Volume | Used By | Contents |
|---|---|---|
| `postgres-data` | postgres | PostgreSQL data directory |
| `redis-data` | redis | Redis AOF persistence |
| `es-data` | elasticsearch | Elasticsearch indices |
| `file-uploads` | file-service | User file attachments |
| `traefik-certs` | traefik | Let's Encrypt certificates |
| `prometheus-data` | prometheus | Metrics time-series data |
| `grafana-data` | grafana | Dashboards, datasources |

---

## Observability

### Prometheus
- Each Go service exposes `GET /metrics` (Prometheus format via `fiberprom` middleware)
- Prometheus scrapes all services every 15s
- Metrics: HTTP request rate, latency histograms, error rates, DB pool stats, Redis hit/miss

### Grafana
- Datasource: Prometheus
- Pre-provisioned dashboards: per-service HTTP metrics, DB pool, Redis, Elasticsearch
- Alerting: authentication failure spikes, p95 latency > 500ms, error rate > 1% (SECURITY-14)

### Logging
- zerolog JSON to stdout on all services
- Docker log driver: `json-file` with `max-size: 100m`, `max-file: 5`
- Log retention: 90 days minimum (SECURITY-14) — achieved via log rotation + host-level archival

---

## CI/CD — GitHub Actions

### Pipeline per service repository

```
on: push to main / pull_request

jobs:
  test:
    - go test ./... (unit + integration via Testcontainers)
    - govulncheck ./...

  build:
    - docker build (multi-stage, pinned base images)
    - docker push to GitHub Container Registry (ghcr.io)

  deploy:
    - SSH to server
    - docker compose pull
    - docker compose up -d --no-deps <service>
    - Health check: curl /health until 200
```

### Security controls (SECURITY-10, SECURITY-13)
- Pinned action versions (no `@latest`)
- Secrets via GitHub Actions secrets (never in code)
- Separate deploy approval for production (environment protection rules)

---

## Health Checks

Each service implements `GET /health` returning:
```json
{ "status": "ok", "service": "<name>", "version": "<git-sha>" }
```

Docker Compose `healthcheck` configured on all services:
```yaml
healthcheck:
  test: ["CMD", "wget", "-qO-", "http://localhost:<port>/health"]
  interval: 30s
  timeout: 5s
  retries: 3
  start_period: 10s
```

`restart: unless-stopped` on all containers for 99.9% uptime target.
