# API Reference

## Base URL

```
http://localhost:4041
```

## Response Format

All API responses follow a consistent JSON format.

### Success Response

```json
{
  "data": {
    // Response payload
  }
}
```

### Error Response

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": [
      {
        "field": "field_name",
        "message": "Field-specific error message"
      }
    ]
  }
}
```

## Status Codes

| Code | Description |
|------|-------------|
| 200  | OK - Request successful |
| 201  | Created - Resource created successfully |
| 400  | Bad Request - Invalid request format or validation error |
| 404  | Not Found - Resource not found |
| 500  | Internal Server Error - Server error |

## Todo Endpoints

### Create Todo

Create a new todo item.

**Endpoint:** `POST /todos`

**Request Body:**
```json
{
  "title": "string (required, min: 3, max: 255)",
  "description": "string (optional, max: 5000)"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:4041/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread"
  }'
```

**Example Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "status": "pending",
  "created_at": "2025-12-18T10:00:00Z",
  "updated_at": null,
  "deleted_at": null
}
```

### List Todos

Retrieve a paginated list of todos.

**Endpoint:** `GET /todos`

**Query Parameters:**
| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| page | integer | Page number | 1 |
| limit | integer | Items per page | 10 |
| status | string | Filter by status (pending, completed, cancelled, archived) | - |
| sort | string | Sort field | created_at |
| order | string | Sort order (asc, desc) | desc |

**Example Request:**
```bash
curl "http://localhost:4041/todos?page=1&limit=10&status=pending&sort=created_at&order=desc"
```

**Example Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "status": "pending",
    "created_at": "2025-12-18T10:00:00Z",
    "updated_at": null,
    "deleted_at": null
  },
  {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "title": "Clean house",
    "description": "",
    "status": "pending",
    "created_at": "2025-12-18T09:00:00Z",
    "updated_at": null,
    "deleted_at": null
  }
]
```

### Get Todo

Retrieve a specific todo by ID.

**Endpoint:** `GET /todos/:id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | Todo ID |

**Example Request:**
```bash
curl http://localhost:4041/todos/550e8400-e29b-41d4-a716-446655440000
```

**Example Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "status": "pending",
  "created_at": "2025-12-18T10:00:00Z",
  "updated_at": null,
  "deleted_at": null
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Todo not found"
  }
}
```

### Update Todo

Update an existing todo.

**Endpoint:** `PATCH /todos/:id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | Todo ID |

**Request Body:**
All fields are optional. Only provided fields will be updated.

```json
{
  "title": "string (optional, min: 3, max: 255)",
  "description": "string (optional, max: 5000)",
  "status": "string (optional, enum: pending, completed, cancelled, archived)"
}
```

**Example Request:**
```bash
curl -X PATCH http://localhost:4041/todos/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

**Example Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "status": "completed",
  "created_at": "2025-12-18T10:00:00Z",
  "updated_at": "2025-12-18T11:00:00Z",
  "deleted_at": null
}
```

### Delete Todo

Soft delete a todo (marks as archived).

**Endpoint:** `DELETE /todos/:id`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | UUID | Todo ID |

**Example Request:**
```bash
curl -X DELETE http://localhost:4041/todos/550e8400-e29b-41d4-a716-446655440000
```

**Example Response:** `204 No Content`

## Validation Errors

When validation fails, the API returns a `400 Bad Request` with details:

**Example Validation Error:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "title",
        "message": "title is required"
      },
      {
        "field": "description",
        "message": "description must be at most 5000 characters"
      }
    ]
  }
}
```

## Common Validation Rules

### Todo

**title:**
- Required
- Minimum length: 3 characters
- Maximum length: 255 characters

**description:**
- Optional
- Maximum length: 5000 characters

**status:**
- Optional (defaults to "pending")
- Allowed values: `pending`, `completed`, `cancelled`, `archived`

## Headers

### Request Headers

| Header | Required | Description |
|--------|----------|-------------|
| Content-Type | Yes (for POST/PATCH) | Must be `application/json` |
| Accept-Language | No | Preferred language (en, tr) |

### Response Headers

| Header | Description |
|--------|-------------|
| Content-Type | Always `application/json` |

## Internationalization

The API supports multiple languages based on the `Accept-Language` header.

**Supported Languages:**
- English (`en`)
- Turkish (`tr`)

**Example:**
```bash
curl http://localhost:4041/todos/invalid-id \
  -H "Accept-Language: tr"
```

Error messages will be returned in Turkish.

## Pagination

List endpoints support pagination via query parameters:

**Parameters:**
- `page`: Page number (starts at 1)
- `limit`: Items per page (default: 10, max: 100)

**Example:**
```bash
curl "http://localhost:4041/todos?page=2&limit=20"
```

## Filtering

List endpoints support filtering via query parameters:

**Todo Filters:**
- `status`: Filter by status (pending, completed, cancelled, archived)

**Example:**
```bash
curl "http://localhost:4041/todos?status=completed"
```

## Sorting

List endpoints support sorting via query parameters:

**Parameters:**
- `sort`: Field name to sort by
- `order`: Sort order (`asc` or `desc`)

**Example:**
```bash
curl "http://localhost:4041/todos?sort=created_at&order=asc"
```

## Rate Limiting

Currently, there is no rate limiting implemented. This should be added for production use.

## Authentication

Currently, the API does not require authentication. For production use, implement:
- JWT authentication
- API key authentication
- OAuth 2.0

## CORS

Configure CORS settings in production based on your requirements.

## Examples

### Complete Workflow

```bash
# 1. Create a todo
TODO_ID=$(curl -X POST http://localhost:4041/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"My Task","description":"Do something"}' \
  | jq -r '.id')

# 2. Get the todo
curl http://localhost:4041/todos/$TODO_ID

# 3. Update the todo
curl -X PATCH http://localhost:4041/todos/$TODO_ID \
  -H "Content-Type: application/json" \
  -d '{"status":"completed"}'

# 4. List all todos
curl http://localhost:4041/todos

# 5. Delete the todo
curl -X DELETE http://localhost:4041/todos/$TODO_ID
```

### Using with JavaScript

```javascript
// Create todo
const response = await fetch('http://localhost:4041/todos', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: 'My Task',
    description: 'Do something'
  })
});

const todo = await response.json();
console.log('Created:', todo);

// Update todo
const updateResponse = await fetch(`http://localhost:4041/todos/${todo.id}`, {
  method: 'PATCH',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    status: 'completed'
  })
});

const updatedTodo = await updateResponse.json();
console.log('Updated:', updatedTodo);
```

### Using with Python

```python
import requests

# Create todo
response = requests.post(
    'http://localhost:4041/todos',
    json={
        'title': 'My Task',
        'description': 'Do something'
    }
)
todo = response.json()
print('Created:', todo)

# Update todo
response = requests.patch(
    f'http://localhost:4041/todos/{todo["id"]}',
    json={'status': 'completed'}
)
updated_todo = response.json()
print('Updated:', updated_todo)
```

---

For more examples, see the [Quick Start Guide](QUICKSTART.md).
