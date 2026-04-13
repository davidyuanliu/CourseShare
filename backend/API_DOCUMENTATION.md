# CourseShare Backend API Documentation

**Base URL:** `http://localhost:8080`

**Last Updated:** April 12, 2026

---

## Table of Contents
1. [Authentication Endpoints](#authentication-endpoints)
2. [Course Endpoints](#course-endpoints)
3. [Note Endpoints](#note-endpoints)
4. [Error Responses](#error-responses)
5. [Testing Guide](#testing-guide)

---

## Authentication Endpoints

### 1. Register User
**Endpoint:** `POST /auth/register`

**Authentication:** None (Public)

**Request Body:**
```json
{
  "email": "student@example.com",
  "password": "securepassword123"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "email": "student@example.com"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Missing fields or invalid email format
- `400 Bad Request` - Email already registered

---

### 2. Login User
**Endpoint:** `POST /auth/login`

**Authentication:** None (Public)

**Request Body:**
```json
{
  "email": "student@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "student@example.com"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Missing email or password
- `401 Unauthorized` - Invalid credentials

**Note:** Save the `token` value for use in subsequent authenticated requests

---

## Course Endpoints

### 1. Get All Courses
**Endpoint:** `GET /courses`

**Authentication:** None (Public)

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Mathematics",
    "notes": null
  },
  {
    "id": 2,
    "name": "Physics",
    "notes": null
  }
]
```

---

### 2. Create Course
**Endpoint:** `POST /courses`

**Authentication:** None (Public for MVP)

**Request Body:**
```json
{
  "name": "Chemistry"
}
```

**Response (201 Created):**
```json
{
  "message": "Course created successfully",
  "course": {
    "id": 3,
    "name": "Chemistry",
    "notes": null
  }
}
```

**Error Responses:**
- `400 Bad Request` - Course name is required

---

## Note Endpoints

### 1. Get All Notes for a Course
**Endpoint:** `GET /courses/{courseId}/notes`

**Authentication:** None (Public - Read Only)

**Path Parameters:**
- `courseId` (integer) - The ID of the course

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "title": "Chapter 1: Introduction",
    "content": "This is the introduction to algebra...",
    "courseId": 1,
    "userId": 1,
    "author": {
      "id": 1,
      "email": "student@example.com"
    }
  },
  {
    "id": 2,
    "title": "Chapter 2: Equations",
    "content": "Solving linear equations...",
    "courseId": 1,
    "userId": 2,
    "author": {
      "id": 2,
      "email": "student2@example.com"
    }
  }
]
```

**Error Responses:**
- `400 Bad Request` - Invalid course ID

---

### 2. Get Single Note
**Endpoint:** `GET /notes/{noteId}`

**Authentication:** None (Public - Read Only)

**Path Parameters:**
- `noteId` (integer) - The ID of the note

**Response (200 OK):**
```json
{
  "id": 1,
  "title": "Chapter 1: Introduction",
  "content": "This is the introduction to algebra...",
  "courseId": 1,
  "userId": 1,
  "author": {
    "id": 1,
    "email": "student@example.com"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid note ID
- `404 Not Found` - Note does not exist

---

### 3. Create Note
**Endpoint:** `POST /notes`

**Authentication:** Required (Bearer Token)

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Chapter 1: Introduction",
  "content": "This is the introduction to algebra...",
  "courseId": 1
}
```

**Response (201 Created):**
```json
{
  "message": "Note created successfully",
  "note": {
    "id": 1,
    "title": "Chapter 1: Introduction",
    "content": "This is the introduction to algebra...",
    "courseId": 1,
    "userId": 1,
    "author": {
      "id": 1,
      "email": "student@example.com"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` - Missing required fields (title, content, courseId)
- `401 Unauthorized` - Missing or invalid authorization header
- `500 Internal Server Error` - Database error

**Note:** The `userId` and `author` are automatically set to the authenticated user

---

### 4. Update Note
**Endpoint:** `PUT /notes/{noteId}`

**Authentication:** Required (Bearer Token)

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

**Path Parameters:**
- `noteId` (integer) - The ID of the note to update

**Request Body:**
```json
{
  "title": "Updated Chapter 1: Introduction",
  "content": "Updated content for the introduction to algebra..."
}
```

**Response (200 OK):**
```json
{
  "message": "Note updated successfully",
  "note": {
    "id": 1,
    "title": "Updated Chapter 1: Introduction",
    "content": "Updated content for the introduction to algebra...",
    "courseId": 1,
    "userId": 1,
    "author": {
      "id": 1,
      "email": "student@example.com"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid note ID or missing required fields
- `401 Unauthorized` - Missing or invalid authorization header
- `403 Forbidden` - You are not the note owner
- `404 Not Found` - Note does not exist

**Note:** Only the note owner can update their notes. Title and content must not be empty.

---

### 5. Delete Note
**Endpoint:** `DELETE /notes/{noteId}`

**Authentication:** Required (Bearer Token)

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters:**
- `noteId` (integer) - The ID of the note to delete

**Response (200 OK):**
```json
{
  "message": "Note deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid note ID
- `401 Unauthorized` - Missing or invalid authorization header
- `403 Forbidden` - You are not the note owner
- `404 Not Found` - Note does not exist

**Note:** Only the note owner can delete their notes.

---

## Error Responses

### Standard Error Format

**400 Bad Request:**
```json
{
  "error": "Missing required fields",
  "missingFields": ["title", "content"]
}
```

**401 Unauthorized:**
```
Missing authorization header
```

**403 Forbidden:**
```
Unauthorized: you can only delete your own notes
```

**404 Not Found:**
```
Note not found
```

**500 Internal Server Error:**
```
Internal server error
```

---

## Testing Guide

### Postman Collection Setup

#### Step 1: Import Environment Variables

Create a Postman Environment with the following variables:
```
base_url = http://localhost:8080
token = (will be populated after login)
student1_email = student1@example.com
student1_password = password123
student2_email = student2@example.com
student2_password = password456
```

---

### Complete Testing Workflow

#### Test Sequence 1: User Registration and Login

**1. Register First User**
```
POST http://localhost:8080/auth/register
Headers: Content-Type: application/json
Body:
{
  "email": "student1@example.com",
  "password": "password123"
}
```

Expected Response: `201 Created` with user ID

**2. Register Second User**
```
POST http://localhost:8080/auth/register
Headers: Content-Type: application/json
Body:
{
  "email": "student2@example.com",
  "password": "password456"
}
```

Expected Response: `201 Created` with user ID

**3. Login First User**
```
POST http://localhost:8080/auth/login
Headers: Content-Type: application/json
Body:
{
  "email": "student1@example.com",
  "password": "password123"
}
```

Expected Response: `200 OK` with JWT token
- **Save the token** in Postman environment as `token`

---

#### Test Sequence 2: Course and Note CRUD

**1. Create Course**
```
POST http://localhost:8080/courses
Headers: Content-Type: application/json
Body:
{
  "name": "Mathematics 101"
}
```

Expected Response: `201 Created` with course ID
- **Save the course ID** for later use

**2. Get All Courses**
```
GET http://localhost:8080/courses
```

Expected Response: `200 OK` with array of courses

**3. Create Note (as Student 1)**
```
PUT http://localhost:8080/notes
Headers: 
  Content-Type: application/json
  Authorization: Bearer {{token}}
Body:
{
  "title": "Chapter 1: Algebra Basics",
  "content": "Introduction to algebraic expressions and equations",
  "courseId": 1
}
```

Expected Response: `201 Created` with note including author info
- **Save the note ID** for later use

**4. Get Note by ID**
```
GET http://localhost:8080/notes/1
```

Expected Response: `200 OK` with note details including author

**5. Get All Notes for Course**
```
GET http://localhost:8080/courses/1/notes
```

Expected Response: `200 OK` with array of notes for the course

**6. Create Another Note (Student 1)**
```
POST http://localhost:8080/notes
Headers:
  Content-Type: application/json
  Authorization: Bearer {{token}}
Body:
{
  "title": "Chapter 2: Solving Equations",
  "content": "Methods for solving linear and quadratic equations",
  "courseId": 1
}
```

Expected Response: `201 Created` with new note
- **Save the note ID**

---

#### Test Sequence 3: Note Editing

**1. Update Own Note (as Student 1)**
```
PUT http://localhost:8080/notes/1
Headers:
  Content-Type: application/json
  Authorization: Bearer {{token}}
Body:
{
  "title": "Chapter 1: Algebra Fundamentals",
  "content": "Updated content about algebra fundamentals"
}
```

Expected Response: `200 OK` with updated note

**2. Attempt to Update Another User's Note (as Student 1)**

First, login as Student 2 and create a note:
```
POST http://localhost:8080/auth/login
Body: {"email": "student2@example.com", "password": "password456"}
```

Create a note as Student 2:
```
POST http://localhost:8080/notes
Headers:
  Content-Type: application/json
  Authorization: Bearer <STUDENT2_TOKEN>
Body:
{
  "title": "Student 2's Note",
  "content": "This is Student 2's note",
  "courseId": 1
}
```

Now try to update Student 2's note as Student 1:
```
PUT http://localhost:8080/notes/3
Headers:
  Content-Type: application/json
  Authorization: Bearer {{token}}
Body:
{
  "title": "Hacked Note",
  "content": "This should fail"
}
```

Expected Response: `403 Forbidden` (Unauthorized: you can only update your own notes)

---

#### Test Sequence 4: Note Deletion

**1. Delete Own Note (as Student 1)**
```
DELETE http://localhost:8080/notes/2
Headers: Authorization: Bearer {{token}}
```

Expected Response: `200 OK` with message "Note deleted successfully"

**2. Verify Note is Deleted**
```
GET http://localhost:8080/notes/2
```

Expected Response: `404 Not Found`

**3. Attempt to Delete Another User's Note (as Student 1)**
```
DELETE http://localhost:8080/notes/3
Headers: Authorization: Bearer {{token}}
```

Expected Response: `403 Forbidden`

---

#### Test Sequence 5: Validation Testing

**1. Create Note with Missing Fields**
```
POST http://localhost:8080/notes
Headers:
  Content-Type: application/json
  Authorization: Bearer {{token}}
Body:
{
  "title": "",
  "content": "",
  "courseId": 0
}
```

Expected Response: `400 Bad Request` with missingFields array

**2. Update Note with Missing Content**
```
PUT http://localhost:8080/notes/1
Headers:
  Content-Type: application/json
  Authorization: Bearer {{token}}
Body:
{
  "title": "Chapter 1",
  "content": ""
}
```

Expected Response: `400 Bad Request` with missingFields array

---

#### Test Sequence 6: Authentication Testing

**1. Access Protected Endpoint Without Token**
```
POST http://localhost:8080/notes
Headers: Content-Type: application/json
Body:
{
  "title": "No Auth Note",
  "content": "This should fail",
  "courseId": 1
}
```

Expected Response: `401 Unauthorized`

**2. Access Protected Endpoint with Invalid Token**
```
POST http://localhost:8080/notes
Headers:
  Content-Type: application/json
  Authorization: Bearer invalid_token_123
Body:
{
  "title": "Invalid Token",
  "content": "This should fail",
  "courseId": 1
}
```

Expected Response: `401 Unauthorized`

---

## Quick Postman Request Templates

Copy these directly into Postman:

### Register User Template
```
POST {{base_url}}/auth/register
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

### Create Note Template
```
POST {{base_url}}/notes
Content-Type: application/json
Authorization: Bearer {{token}}

{
  "title": "Note Title",
  "content": "Note content",
  "courseId": 1
}
```

### Update Note Template
```
PUT {{base_url}}/notes/1
Content-Type: application/json
Authorization: Bearer {{token}}

{
  "title": "Updated Title",
  "content": "Updated content"
}
```

### Delete Note Template
```
DELETE {{base_url}}/notes/1
Authorization: Bearer {{token}}
```

---

## HTTP Status Codes Reference

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful query or update |
| 201 | Created | Successful resource creation |
| 400 | Bad Request | Invalid input or missing fields |
| 401 | Unauthorized | Missing or invalid auth token |
| 403 | Forbidden | Authenticated but not authorized (e.g., not note owner) |
| 404 | Not Found | Resource doesn't exist |
| 500 | Internal Server Error | Server-side error |

---

## Notes for Testing

1. **Always save tokens** after login for use in subsequent requests
2. **Course 1 must exist** - Create it before testing notes
3. **Test with multiple users** to verify ownership restrictions
4. **Verify 403 responses** when trying to edit/delete others' notes
5. **Test all validation scenarios** - empty fields, missing fields, invalid IDs
6. **Check response format** - All responses include proper JSON structure
7. **Verify timestamps** - Author info should appear in all note responses

---

## Environment Variables Setup in Postman

1. Click "Environments" in Postman
2. Create new environment "CourseShare"
3. Add variables:
   - `base_url`: `http://localhost:8080`
   - `token`: `(empty - will be set after login)`
   - `course_id`: `(empty - will be set after course creation)`
   - `note_id`: `(empty - will be set after note creation)`
4. Use `{{variable_name}}` in requests

---

## Example cURL Commands

```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"student@example.com","password":"pass123"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"student@example.com","password":"pass123"}'

# Create Note
curl -X POST http://localhost:8080/notes \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"My Note","content":"Content","courseId":1}'

# Get Note
curl -X GET http://localhost:8080/notes/1

# Update Note
curl -X PUT http://localhost:8080/notes/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","content":"Updated content"}'

# Delete Note
curl -X DELETE http://localhost:8080/notes/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

**Happy Testing! 🚀**
