# Infrastructure Design Plan — Unit 1: backend

## Plan Steps
- [x] Answer clarifying questions (user)
- [x] Generate infrastructure-design.md
- [x] Generate deployment-architecture.md

---

## Clarifying Questions

Please fill in the `[Answer]:` tag for each question. Let me know when done.

---

## Question 1
How should TLS be handled for the Nginx reverse proxy?

A) Self-signed certificate (development/internal use only)
B) Let's Encrypt via Certbot (free, auto-renewing, requires a public domain)
C) Bring your own certificate (manually mounted into Nginx container)
D) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 2
How should PostgreSQL data be persisted?

A) Docker named volume (managed by Docker, survives container restarts)
B) Bind mount to a specific host directory (e.g., `/data/postgres`)
C) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 3
How should observability (metrics + monitoring) be handled?

A) None for now — Docker logs + zerolog JSON output is sufficient
B) Prometheus + Grafana (self-hosted, Docker containers, scrape `/metrics` from each service)
C) Other (please describe after [Answer]: tag below)

[Answer]:B

---

## Question 4
How should the CI/CD pipeline be set up?

A) GitHub Actions — build, test, push Docker images, deploy via SSH to server
B) GitLab CI — same flow
C) No CI/CD for now — manual `docker compose up` on server
D) Other (please describe after [Answer]: tag below)

[Answer]:A

---

## Question 5
Should an API gateway sit in front of the microservices, or should Nginx route directly to each service?

A) Nginx routes directly to each service (simple, no extra layer)
B) Add a lightweight API gateway (e.g., Traefik) for routing, auth middleware, and observability
C) Other (please describe after [Answer]: tag below)

[Answer]:B
