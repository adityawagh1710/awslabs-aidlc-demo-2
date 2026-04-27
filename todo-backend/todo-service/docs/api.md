# todo-service API

All endpoints require `Authorization: Bearer <access_token>`.

## Todos

### POST /todos
Create a todo.
```json
{ "title": "Buy milk", "description": "...", "priority": "high", "due_date": "2026-05-01T09:00:00Z", "tag_ids": [] }
```
**Response** `201` Todo object. **Errors**: `422` validation

### GET /todos
List todos. Query params: `status`, `priority`, `tag_id`.
**Response** `200` `[]Todo`

### GET /todos/search?q=
Search todos via Elasticsearch.
**Response** `200` `[]Todo`. **Errors**: `422` missing q

### GET /todos/:id
**Response** `200` Todo. **Errors**: `404`, `403`

### PUT /todos/:id
Update todo fields. Status transitions enforced: `pending→in_progress→done`.
```json
{ "title": "...", "status": "in_progress", "priority": "low", "tag_ids": [] }
```
**Response** `200` Todo. **Errors**: `422` invalid transition, `403`, `404`

### DELETE /todos/:id
Soft delete. **Response** `204`. **Errors**: `403`, `404`

---

## Tags

### POST /tags
```json
{ "name": "work" }
```
**Response** `201` Tag. **Errors**: `422`

### GET /tags
**Response** `200` `[]Tag`

### DELETE /tags/:id
**Response** `204`. **Errors**: `403`, `404`

---

## GET /health
```json
{ "status": "ok", "service": "todo-service" }
```
