# Integration Test Instructions

## Purpose
Test interactions between all 5 backend services running together via Docker Compose.

## Setup

```bash
cd todo-infra
cp .env.example .env.test
# Set test-specific values (separate DB, test secrets)
docker compose --env-file .env.test up -d
docker compose ps   # wait until all healthy
```

## Integration Test Scenarios

### Scenario 1: Auth → Todo (JWT validation across services)
```bash
# 1. Register user
TOKEN=$(curl -s -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.access_token')

# 2. Create todo using token (validates JWT in todo-service)
curl -s -X POST http://localhost:3001/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Integration test todo","priority":"high"}'
```
**Expected**: 201 with todo object

### Scenario 2: Todo → Scheduler (reminder creation)
```bash
TODO_ID=<id from above>

curl -s -X POST http://localhost:3004/reminders \
  -H "X-Internal-API-Key: $INTERNAL_API_KEY" \
  -H "Content-Type: application/json" \
  -d "{\"todo_id\":\"$TODO_ID\",\"user_id\":\"<user_id>\",\"fire_at\":\"$(date -u -d '+1 hour' +%Y-%m-%dT%H:%M:%SZ)\"}"
```
**Expected**: 201 with reminder object

### Scenario 3: Todo → File (attachment upload)
```bash
curl -s -X POST http://localhost:3002/files \
  -H "Authorization: Bearer $TOKEN" \
  -F "todo_id=$TODO_ID" \
  -F "file=@/tmp/test.pdf;type=application/pdf"
```
**Expected**: 201 with file metadata

### Scenario 4: Scheduler → Notification (event delivery)
```bash
curl -s -X POST http://localhost:3003/internal/events \
  -H "X-Internal-API-Key: $INTERNAL_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"<user_id>","todo_id":"<todo_id>","message":"Test reminder"}'
```
**Expected**: 202 Accepted; notification stored in DB

### Scenario 5: Todo soft delete cascade
```bash
curl -s -X DELETE http://localhost:3001/todos/$TODO_ID \
  -H "Authorization: Bearer $TOKEN"
```
**Expected**: 204; verify file deleted from disk; verify reminders deleted

## Cleanup

```bash
docker compose --env-file .env.test down -v
```
