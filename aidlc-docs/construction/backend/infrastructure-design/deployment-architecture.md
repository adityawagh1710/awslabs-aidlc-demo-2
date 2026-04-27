# Deployment Architecture — Unit 1: backend

## Architecture Diagram

```
Internet
    |
    | :443 / :80
    v
+------------------+
|     Traefik v3   |  API Gateway + TLS (Let's Encrypt) + HTTP->HTTPS redirect
|  (todo-network)  |  Routes: /auth /todos /tags /files /ws
+--------+---------+
         |  (internal Docker network: todo-network)
         |
    +----+----------------------------------------------+
    |    |              |              |                 |
    v    v              v              v                 v
+-------+--+  +--------+--+  +-------+---+  +----------+--------+
|auth-service|  |todo-service|  |file-service|  |notification-svc  |
|  :3000     |  |  :3001     |  |  :3002     |  |  :3003 (WS)      |
+-----+------+  +-----+------+  +-----+------+  +------------------+
      |                |               |
      |         +------+------+        |
      |         |             |        |
      v         v             v        v
  +-------+  +------+  +----------+  +--------+
  | Redis |  | PG   |  |   ES     |  | PG     |
  | :6379 |  | :5432|  | :9200    |  | :5432  |
  +-------+  +------+  +----------+  +--------+
      ^           ^
      |           |
+-----+-----+  +--+--+
|scheduler  |  | PG  |
|svc :3004  |  |:5432|
+-----------+  +-----+

Observability (internal, Traefik-protected):
+------------+     +----------+
| Prometheus |---->| Grafana  |
|  :9090     |     |  :3000   |
+------------+     +----------+
  (scrapes all services /metrics)
```

## Docker Compose Service Definitions (Summary)

```yaml
services:

  traefik:
    image: traefik:v3
    ports: ["80:80", "443:443"]
    volumes:
      - traefik-certs:/certs
      - /var/run/docker.sock:/var/run/docker.sock:ro
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    volumes: [postgres-data:/var/lib/postgresql/data]
    environment: [POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD]
    healthcheck: pg_isready
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes: [redis-data:/data]
    command: redis-server --appendonly yes
    restart: unless-stopped

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.13.0
    volumes: [es-data:/usr/share/elasticsearch/data]
    environment: [discovery.type=single-node, xpack.security.enabled=true]
    restart: unless-stopped

  auth-service:
    image: ghcr.io/<org>/todo-auth-service:<sha>
    environment: [DATABASE_URL, JWT_SECRET, REDIS_URL, LOG_LEVEL]
    depends_on: [postgres, redis]
    healthcheck: GET /health
    restart: unless-stopped
    labels: [traefik routing for /auth/*]

  todo-service:
    image: ghcr.io/<org>/todo-todo-service:<sha>
    environment: [DATABASE_URL, JWT_SECRET, REDIS_URL, ELASTICSEARCH_URL, LOG_LEVEL]
    depends_on: [postgres, redis, elasticsearch]
    healthcheck: GET /health
    restart: unless-stopped
    labels: [traefik routing for /todos/*, /tags/*]

  scheduler-service:
    image: ghcr.io/<org>/todo-scheduler-service:<sha>
    environment: [DATABASE_URL, JWT_SECRET, INTERNAL_API_KEY, LOG_LEVEL]
    depends_on: [postgres]
    healthcheck: GET /health
    restart: unless-stopped

  file-service:
    image: ghcr.io/<org>/todo-file-service:<sha>
    volumes: [file-uploads:/uploads]
    environment: [DATABASE_URL, JWT_SECRET, FILE_STORAGE_PATH=/uploads, LOG_LEVEL]
    depends_on: [postgres]
    healthcheck: GET /health
    restart: unless-stopped
    labels: [traefik routing for /files/*]

  notification-service:
    image: ghcr.io/<org>/todo-notification-service:<sha>
    environment: [DATABASE_URL, JWT_SECRET, INTERNAL_API_KEY, LOG_LEVEL]
    depends_on: [postgres]
    healthcheck: GET /health
    restart: unless-stopped
    labels: [traefik routing for /ws]

  prometheus:
    image: prom/prometheus:v2.51.0
    volumes: [prometheus-data:/prometheus, ./prometheus.yml:/etc/prometheus/prometheus.yml:ro]
    restart: unless-stopped

  grafana:
    image: grafana/grafana:10.4.0
    volumes: [grafana-data:/var/lib/grafana]
    environment: [GF_SECURITY_ADMIN_PASSWORD]
    restart: unless-stopped

volumes:
  postgres-data:
  redis-data:
  es-data:
  file-uploads:
  traefik-certs:
  prometheus-data:
  grafana-data:

networks:
  default:
    name: todo-network
```

## Environment Variables (per service, injected via Docker Compose `.env`)

| Variable | Services | Source |
|---|---|---|
| `DATABASE_URL` | all backend | `.env` (never in code) |
| `JWT_SECRET` | all backend | `.env` |
| `REDIS_URL` | auth, todo | `.env` |
| `ELASTICSEARCH_URL` | todo | `.env` |
| `INTERNAL_API_KEY` | scheduler, notification | `.env` |
| `FILE_STORAGE_PATH` | file | `.env` |
| `LOG_LEVEL` | all | `.env` |
| `POSTGRES_*` | postgres | `.env` |
| `GF_SECURITY_ADMIN_PASSWORD` | grafana | `.env` |

`.env` file is gitignored. A `.env.example` with placeholder values is committed.

## Migration Strategy

- `golang-migrate` runs automatically on service startup via `migrate up`
- Migration files in `migrations/` directory per service repository
- Idempotent — safe to run on every deploy
