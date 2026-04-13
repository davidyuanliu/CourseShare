# CourseShare Backend API Documentation

**Sprint 3 | April 2026**
**Base URL:** `http://localhost:8080`

---

## What's New in Sprint 3

- User registration and login via JWT authentication
- Notes are now author-tracked — each note is linked to its creator
- Update and Delete note endpoints added (owner-only)
- All protected endpoints require a Bearer token in the `Authorization` header
- Author information is preloaded and returned in all note responses
- Improved validation: whitespace-only values are rejected; errors include a `missingFields` array

---

## Endpoint Summary

| Method | Path | Auth Required | Description |
|--------|------|:---:|-------------|
| GET | `/` | No | Health check |
| POST | `/auth/register` | No | Register a new user |
| POST | `/auth/login` | No | Login and receive JWT token |
| GET | `/courses` | No | List all courses |
| POST | `/courses` | No | Create a new course |
| GET | `/courses/{id}/notes` | No | Get all notes for a course |
| GET | `/notes/{id}` | No | Get a single note by ID |
| POST | `/notes` | **Yes** | Create a note (authenticated) |
| PUT | `/notes/{id}` | **Yes** | Update a note (owner only) |
| DELETE | `/notes/{id}` | **Yes** | Delete a note (owner only) |

---

## Authentication

Protected endpoints require a JWT Bearer token in the `Authorization` header. Tokens are issued on login and expire after 24 hours.

```
Authorization: Bearer <your_jwt_token>
```

---

### POST /auth/register

Registers a new user account.

**Request Body**
```json
{
  "email": "student@example.com",
  "password": "securepassword123"
}
```

**Responses**

| Status | Condition | Body |
|--------|-----------|------|
| `201 Created` | Valid email & password | `{ "message": "User registered successfully", "user": { "id": 1, "email": "..." } }` |
| `400 Bad Request` | Missing fields | `{ "error": "Missing required fields", "missingFields": ["email"] }` |
| `400 Bad Request` | Invalid email format | `{ "error": "Invalid email format" }` |
| `400 Bad Request` | Duplicate email | `{ "error": "Email already registered" }` |

---

### POST /auth/login

Authenticates a user and returns a JWT token valid for 24 hours.

**Request Body**
```json
{
  "email": "student@example.com",
  "password": "securepassword123"
}
```

**Success Response — 200 OK**
```json
{
  "message": "Login successful",
  "token": "<jwt_token>",
  "user": {
    "id": 1,
    "email": "student@example.com"
  }
}
```

**Error Responses**
- `400 Bad Request` — missing email or password
- `401 Unauthorized` — invalid credentials

---

## Courses

### GET /courses

Returns all courses. No authentication required.

**Response — 200 OK**
```json
[
  {
    "ID": 1,
    "CreatedAt": "2026-04-01T12:00:00Z",
    "UpdatedAt": "2026-04-01T12:00:00Z",
    "DeletedAt": null,
    "name": "Database Systems",
    "notes": []
  }
]
```

---

### POST /courses

Creates a new course. No authentication required.

**Request Body**
```json
{ "name": "Database Systems" }
```

**Responses**
- `201 Created` — returns the created course object
- `400 Bad Request` — course name is empty

---

## Notes

### GET /courses/{id}/notes

Returns all notes for a given course. Author information is included in each note.

**Example**
```
GET /courses/1/notes
```

**Response — 200 OK**
```json
[
  {
    "ID": 1,
    "title": "Chapter 1: Intro",
    "content": "Notes content here...",
    "courseId": 1,
    "userId": 2,
    "author": {
      "ID": 2,
      "email": "student@example.com"
    }
  }
]
```

**Error Responses**
- `400 Bad Request` — invalid course ID format

---

### GET /notes/{id}

Returns a single note by ID, including author information.

**Responses**
- `200 OK` — note object with `author` field
- `400 Bad Request` — invalid ID format
- `404 Not Found` — note does not exist

---

### POST /notes 🔒 Authenticated

Creates a new note. The authenticated user is automatically set as the author.

**Request Body**
```json
{
  "title": "Chapter 1: Algebra Basics",
  "content": "Introduction to algebraic expressions.",
  "courseId": 1
}
```

**Success Response — 201 Created**
```json
{
  "message": "Note created successfully",
  "note": {
    "ID": 5,
    "title": "Chapter 1: Algebra Basics",
    "content": "Introduction to algebraic expressions.",
    "courseId": 1,
    "userId": 2,
    "author": { "ID": 2, "email": "student@example.com" }
  }
}
```

**Error Responses**
- `400 Bad Request` — missing or whitespace-only `title`, `content`, or `courseId`
- `401 Unauthorized` — missing or invalid token

**Validation Error Format**
```json
{
  "error": "Missing required fields",
  "missingFields": ["title", "content"]
}
```

> All three fields are required. Whitespace-only values are treated as empty.

---

### PUT /notes/{id} 🔒 Owner Only

Updates an existing note. Only the original author may update their note.

**Request Body**
```json
{
  "title": "Updated Chapter 1",
  "content": "Updated content with more examples."
}
```

**Responses**
- `200 OK` — returns updated note with author information
- `400 Bad Request` — missing or empty `title` or `content`
- `401 Unauthorized` — missing or invalid token
- `403 Forbidden` — authenticated user is not the note owner
- `404 Not Found` — note does not exist

> `courseId` cannot be changed via update. Only `title` and `content` are mutable.

---

### DELETE /notes/{id} 🔒 Owner Only

Deletes a note. Only the original author may delete their note.

**Responses**
- `200 OK` — `{ "message": "Note deleted successfully" }`
- `401 Unauthorized` — missing or invalid token
- `403 Forbidden` — authenticated user is not the note owner
- `404 Not Found` — note does not exist

---

## Data Models

### User

| Field | Type | Notes |
|-------|------|-------|
| `ID` | uint | Auto-incremented primary key |
| `email` | string | Unique, case-insensitive |
| `password` | string | bcrypt hashed; omitted from all API responses |
| `notes` | []Note | Has-many relationship |
| `CreatedAt` / `UpdatedAt` | timestamp | Managed by GORM |

### Course

| Field | Type | Notes |
|-------|------|-------|
| `ID` | uint | Auto-incremented primary key |
| `name` | string | Required, cannot be empty |
| `notes` | []Note | Has-many relationship |
| `CreatedAt` / `UpdatedAt` / `DeletedAt` | timestamp | Managed by GORM (soft delete) |

### Note

| Field | Type | Notes |
|-------|------|-------|
| `ID` | uint | Auto-incremented primary key |
| `title` | string | Required, cannot be blank |
| `content` | string | Required, cannot be blank |
| `courseId` | uint | Required, must be non-zero |
| `userId` | uint | Set automatically from JWT on create |
| `author` | *User | Preloaded on all read responses |
| `CreatedAt` / `UpdatedAt` / `DeletedAt` | timestamp | Managed by GORM |

---

## HTTP Status Code Reference

| Code | Meaning |
|------|---------|
| `200 OK` | Successful read or update |
| `201 Created` | Resource successfully created |
| `400 Bad Request` | Validation failure or malformed input |
| `401 Unauthorized` | Missing or invalid JWT token |
| `403 Forbidden` | Authenticated, but user does not own the resource |
| `404 Not Found` | Resource does not exist |

---

## Manual Testing with curl

**Register & Login**
```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"student@example.com","password":"pass123"}'

# Login (save token from response)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"student@example.com","password":"pass123"}'
```

**Create, Update & Delete Notes**
```bash
# Create note
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"Chapter 1","content":"Intro notes","courseId":1}'

# Update note
curl -X PUT http://localhost:8080/notes/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"Updated Chapter 1","content":"Updated content"}'

# Delete note
curl -X DELETE http://localhost:8080/notes/1 \
  -H "Authorization: Bearer <token>"
```

**Test Validation & Auth Errors**
```bash
# Missing fields (400)
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"","content":"","courseId":0}'

# No token (401)
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","content":"Body","courseId":1}'

# Non-owner update (403)
curl -X PUT http://localhost:8080/notes/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <other_users_token>" \
  -d '{"title":"Hacked","content":"..."}'
```