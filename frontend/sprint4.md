# CourseShare Sprint 4 Backend Deliverables

**Base URL:** `http://localhost:8080`

---

## 1. Detail Work Completed in Sprint 4

### Backend Implementation
- Added `GET /my-notes` to return only notes created by the authenticated user
- Added course preloading to My Notes responses so each note includes related course information
- Added note tags to the backend data model
- Added tag support to note creation and note update flows
- Added tag validation so tags are optional but cannot be empty or longer than the allowed limit
- Added tag-based filtering with `GET /courses/{id}/notes?tag=...`
- Added helpful-vote persistence using a dedicated backend model
- Added `POST /notes/{id}/helpful` and `DELETE /notes/{id}/helpful`
- Prevented duplicate helpful votes from the same user on the same note
- Added computed helpful metadata to note responses:
  - `helpfulCount`
  - `isHelpful`
- Added saved-note persistence using a dedicated backend model
- Added `POST /notes/{id}/save`, `DELETE /notes/{id}/save`, and `GET /saved-notes`
- Prevented duplicate saved-note records for the same user and note
- Added computed saved-note metadata to note responses:
  - `isSaved`
- Ensured My Notes, Saved Notes, and note detail/list responses return enriched note metadata
- Expanded backend auth coverage so protected routes are tested through real middleware
- Updated backend deployment readiness by supporting:
  - `DATABASE_URL`
  - `PORT`
  - `ALLOWED_ORIGINS`
- Added backend deployment documentation in [DEPLOYMENT.md](/Users/singh/CourseShare/backend/DEPLOYMENT.md:1)

### Frontend Implementation
- Added `My Notes` dashboard to display notes authored by the current user
- Added `Saved Notes` dashboard to display notes bookmarked by the current user
- Implemented helpful-vote toggling with immediate UI state updates
- Implemented note save/unsave functionality with persistent state
- Added tag support in note creation and edit forms
- Added tag display and keyword search to the course notes list
- Added personal filtering options (e.g., "Written by me") to refine note lists
- Integrated JWT authentication across all new protected frontend routes
- Improved error handling for API interactions using a centralized interceptor

---

## 2. List Backend Unit Tests

### Course Handler Tests
- [course_handler_test.go](/Users/singh/CourseShare/backend/handlers/course_handler_test.go:14)
  - `TestGetCoursesReturnsCourses`
  - `TestCreateCourseSuccess`
  - `TestCreateCourseMissingName`

### User/Auth Handler Tests
- [user_handler_test.go](/Users/singh/CourseShare/backend/handlers/user_handler_test.go:17)
  - `TestRegisterUserSuccess`
  - `TestRegisterUserMissingFields`
  - `TestRegisterUserInvalidEmail`
  - `TestRegisterUserDuplicateEmail`
  - `TestRegisterUserDuplicateEmailCaseInsensitive`
  - `TestRegisterUserPasswordHashed`
  - `TestLoginUserSuccess`
  - `TestLoginUserInvalidEmail`
  - `TestLoginUserWrongPassword`
  - `TestLoginUserMissingFields`

### Note Handler Tests
- [note_handler_test.go](/Users/singh/CourseShare/backend/handlers/note_handler_test.go:20)
  - Basic note retrieval:
    - `TestGetNotesByCourseFiltersByCourseID`
    - `TestGetNoteReturnsNote`
    - `TestGetNoteNotFound`
  - Create/update/delete note flows:
    - `TestCreateNoteSuccess`
    - `TestCreateNoteValidation`
    - `TestDeleteNoteSuccess`
    - `TestDeleteNoteNotFound`
    - `TestDeleteNoteUnauthorized`
    - `TestUpdateNoteSuccess`
    - `TestUpdateNoteNotFound`
    - `TestUpdateNoteUnauthorized`
    - `TestUpdateNoteMissingFields`
  - Author preloading:
    - `TestGetNoteIncludesAuthor`
    - `TestGetNotesByCourseIncludesAuthors`
    - `TestCreateNoteIncludesAuthor`
  - My Notes:
    - `TestGetMyNotesSuccess`
    - `TestGetMyNotesEmpty`
    - `TestGetMyNotesUnauthenticated`
    - `TestGetMyNotesIncludesCourseInfo`
  - Tags:
    - `TestCreateNoteWithTags`
    - `TestCreateNoteWithoutTags`
    - `TestUpdateNoteWithTags`
    - `TestCreateNoteRejectsEmptyTag`
    - `TestCreateNoteRejectsLongTag`
    - `TestUpdateNoteRejectsInvalidTags`
    - `TestGetNotesByCourseFiltersByTag`
  - Helpful votes:
    - `TestMarkNoteHelpfulSuccess`
    - `TestMarkNoteHelpfulPreventsDuplicateVotes`
    - `TestUnmarkNoteHelpfulSuccess`
    - `TestGetNoteIncludesHelpfulMetadata`
    - `TestGetNotesByCourseIncludesHelpfulMetadata`
  - Saved notes:
    - `TestSaveNoteSuccess`
    - `TestSaveNotePreventsDuplicateSaves`
    - `TestUnsaveNoteSuccess`
    - `TestGetSavedNotesSuccess`
    - `TestGetSavedNotesEmpty`
    - `TestGetNoteIncludesSavedMetadata`
  - Protected-route auth coverage:
    - `TestProtectedNoteEndpointsRequireAuthentication`
    - `TestProtectedNoteEndpointsRejectInvalidToken`

---

## 3. Updated Documentation for Backend API

### Endpoint Summary

| Method | Path | Auth Required | Description |
|--------|------|:---:|-------------|
| GET | `/` | No | Health check |
| POST | `/auth/register` | No | Register a new user |
| POST | `/auth/login` | No | Login and receive JWT token |
| GET | `/courses` | No | List all courses |
| POST | `/courses` | No | Create a new course |
| GET | `/courses/{id}/notes` | No | List notes for a course |
| GET | `/notes/{id}` | No | Get a single note |
| GET | `/my-notes` | **Yes** | Get notes created by the authenticated user |
| GET | `/saved-notes` | **Yes** | Get notes saved by the authenticated user |
| POST | `/notes` | **Yes** | Create a note |
| PUT | `/notes/{id}` | **Yes** | Update a note (owner only) |
| DELETE | `/notes/{id}` | **Yes** | Delete a note (owner only) |
| POST | `/notes/{id}/helpful` | **Yes** | Mark a note as helpful |
| DELETE | `/notes/{id}/helpful` | **Yes** | Remove a helpful vote |
| POST | `/notes/{id}/save` | **Yes** | Save a note |
| DELETE | `/notes/{id}/save` | **Yes** | Unsave a note |

### Authentication

Protected endpoints require a JWT Bearer token:

```text
Authorization: Bearer <your_jwt_token>
```

Tokens are issued by `POST /auth/login` and expire after 24 hours.

---

### Courses

#### GET /courses
Returns all courses.

#### POST /courses
Creates a new course.

**Request Body**
```json
{ "name": "Database Systems" }
```

**Responses**
- `201 Created` - returns the created course object
- `400 Bad Request` - course name is empty

---

### Notes

#### GET /courses/{id}/notes
Returns all notes for a course. Notes include:
- author information
- tags
- helpful count
- `isHelpful` when an authenticated user's token is present
- `isSaved` when an authenticated user's token is present

**Optional Query Parameter**
- `tag` - filters notes by tag value

**Example**
```text
GET /courses/1/notes?tag=exam
```

#### GET /notes/{id}
Returns a single note by ID with:
- author information
- tags
- helpful count
- `isHelpful`
- `isSaved`

**Responses**
- `200 OK`
- `400 Bad Request` - invalid ID format
- `404 Not Found` - note does not exist

#### POST /notes
Creates a new note for the authenticated user.

**Request Body**
```json
{
  "title": "Chapter 1: Algebra Basics",
  "content": "Introduction to algebraic expressions.",
  "courseId": 1,
  "tags": ["algebra", "exam"]
}
```

**Success Response**
```json
{
  "message": "Note created successfully",
  "note": {
    "ID": 5,
    "title": "Chapter 1: Algebra Basics",
    "content": "Introduction to algebraic expressions.",
    "courseId": 1,
    "userId": 2,
    "tags": ["algebra", "exam"],
    "author": { "ID": 2, "email": "student@example.com" },
    "course": { "ID": 1, "name": "Database Systems" },
    "helpfulCount": 0,
    "isHelpful": false,
    "isSaved": false
  }
}
```

**Error Responses**
- `400 Bad Request` - missing required fields or invalid tags
- `401 Unauthorized` - missing or invalid token

#### PUT /notes/{id}
Updates an existing note. Only the note owner can update it.

**Request Body**
```json
{
  "title": "Updated Chapter 1",
  "content": "Updated content with more examples.",
  "tags": ["review", "chapter-1"]
}
```

**Error Responses**
- `400 Bad Request` - missing title/content or invalid tags
- `401 Unauthorized` - missing or invalid token
- `403 Forbidden` - authenticated user is not the owner
- `404 Not Found` - note does not exist

#### DELETE /notes/{id}
Deletes a note. Only the note owner can delete it.

**Responses**
- `200 OK` - `{ "message": "Note deleted successfully" }`
- `401 Unauthorized`
- `403 Forbidden`
- `404 Not Found`

---

### My Notes

#### GET /my-notes
Returns only notes owned by the authenticated user.

Each note includes:
- author information
- course information
- tags
- helpful count
- `isHelpful`
- `isSaved`

**Responses**
- `200 OK`
- `401 Unauthorized`

---

### Saved Notes

#### POST /notes/{id}/save
Saves a note for the authenticated user.

**Responses**
- `200 OK` - `{ "message": "Note saved successfully", "isSaved": true }`
- `401 Unauthorized`
- `404 Not Found`
- `409 Conflict` - note already saved

#### DELETE /notes/{id}/save
Removes a saved note for the authenticated user.

**Responses**
- `200 OK` - `{ "message": "Note unsaved successfully", "isSaved": false }`
- `401 Unauthorized`
- `404 Not Found`

#### GET /saved-notes
Returns notes saved by the authenticated user.

Each saved note includes:
- title
- author information
- course information
- tags
- helpful count
- `isHelpful`
- `isSaved`

**Responses**
- `200 OK`
- `401 Unauthorized`

---

### Helpful Votes

#### POST /notes/{id}/helpful
Marks a note as helpful for the authenticated user.

**Responses**
- `200 OK`
```json
{
  "message": "Note marked as helpful",
  "helpfulCount": 1,
  "isHelpful": true
}
```
- `401 Unauthorized`
- `404 Not Found`
- `409 Conflict` - already marked helpful

#### DELETE /notes/{id}/helpful
Removes a helpful vote.

**Responses**
- `200 OK`
```json
{
  "message": "Helpful vote removed",
  "helpfulCount": 0,
  "isHelpful": false
}
```
- `401 Unauthorized`
- `404 Not Found`

---

### Data Models

#### User
| Field | Type | Notes |
|-------|------|-------|
| `ID` | uint | Auto-incremented primary key |
| `email` | string | Unique, case-insensitive |
| `password` | string | bcrypt hashed; omitted from API responses |
| `notes` | []Note | Notes created by the user |

#### Course
| Field | Type | Notes |
|-------|------|-------|
| `ID` | uint | Auto-incremented primary key |
| `name` | string | Required |
| `notes` | []Note | Has-many relationship |

#### Note
| Field | Type | Notes |
|-------|------|-------|
| `ID` | uint | Auto-incremented primary key |
| `title` | string | Required |
| `content` | string | Required |
| `courseId` | uint | Required |
| `userId` | uint | Owner ID |
| `tags` | JSON array of strings | Optional |
| `author` | *User | Preloaded on read responses |
| `course` | *Course | Preloaded where needed |
| `helpfulCount` | int64 | Computed response field |
| `isHelpful` | bool | Computed response field |
| `isSaved` | bool | Computed response field |

#### HelpfulVote
| Field | Type | Notes |
|-------|------|-------|
| `noteId` | uint | Part of unique note/user pair |
| `userId` | uint | Part of unique note/user pair |

#### SavedNote
| Field | Type | Notes |
|-------|------|-------|
| `noteId` | uint | Part of unique note/user pair |
| `userId` | uint | Part of unique note/user pair |

---

### HTTP Status Code Reference

| Code | Meaning |
|------|---------|
| `200 OK` | Successful read, update, vote, or save action |
| `201 Created` | Resource created successfully |
| `400 Bad Request` | Validation failure or malformed input |
| `401 Unauthorized` | Missing or invalid JWT token |
| `403 Forbidden` | Authenticated user does not own the note |
| `404 Not Found` | Resource does not exist |
| `409 Conflict` | Duplicate helpful vote or duplicate saved note |

---

## 4. List Frontend Unit Tests

### Component Tests
- [courses.component.spec.ts](/Users/singh/CourseShare/frontend/src/app/pages/courses/courses.component.spec.ts:1)
- [create-note.component.spec.ts](/Users/singh/CourseShare/frontend/src/app/pages/create-note/create-note.component.spec.ts:1)
- [edit-note.component.spec.ts](/Users/singh/CourseShare/frontend/src/app/pages/edit-note/edit-note.component.spec.ts:1)
- [login.component.spec.ts](/Users/singh/CourseShare/frontend/src/app/pages/login/login.component.spec.ts:1)
- [my-notes.component.spec.ts](/Users/singh/CourseShare/frontend/src/app/pages/my-notes/my-notes.component.spec.ts:1)
- [note-detail.component.spec.ts](/Users/singh/CourseShare/frontend/src/app/pages/note-detail/note-detail.component.spec.ts:1)
- [notes-list.component.spec.ts](/Users/singh/CourseShare/frontend/src/app/pages/notes-list/notes-list.component.spec.ts:1)
- [register.component.spec.ts](/Users/singh/CourseShare/frontend/src/app/pages/register/register.component.spec.ts:1)
- [saved-notes.component.spec.ts](/Users/singh/CourseShare/frontend/src/app/pages/saved-notes/saved-notes.component.spec.ts:1)

### Service & Guard Tests
- [auth.guard.spec.ts](/Users/singh/CourseShare/frontend/src/app/services/auth.guard.spec.ts:1)
- [auth.interceptor.spec.ts](/Users/singh/CourseShare/frontend/src/app/services/auth.interceptor.spec.ts:1)
- [auth.service.spec.ts](/Users/singh/CourseShare/frontend/src/app/services/auth.service.spec.ts:1)
- [course-share-api.service.spec.ts](/Users/singh/CourseShare/frontend/src/app/services/course-share-api.service.spec.ts:1)
- [demo-control.service.spec.ts](/Users/singh/CourseShare/frontend/src/app/services/demo-control.service.spec.ts:1)

---

## 5. List Frontend Cypress E2E Tests

### Feature Tests
- [auth.cy.ts](/Users/singh/CourseShare/frontend/cypress/e2e/auth.cy.ts:1)
  - User registration flow
  - User login/logout flow
- [courses.cy.ts](/Users/singh/CourseShare/frontend/cypress/e2e/courses.cy.ts:1)
  - Course listing and selection
  - Course creation
- [notes.cy.ts](/Users/singh/CourseShare/frontend/cypress/e2e/notes.cy.ts:1)
  - Note creation and editing
  - Note deletion (owner only)
  - Note viewing details
- [sprint4.cy.ts](/Users/singh/CourseShare/frontend/cypress/e2e/sprint4.cy.ts:1)
  - Helpful vote toggling
  - Note save/unsave toggling
  - Dashboard navigation (My Notes, Saved Notes)
  - Search and personal filtering

---

### Deployment Notes

Sprint 4 also introduced backend deployment readiness improvements:
- `DATABASE_URL` support for hosted PostgreSQL
- `PORT` support for cloud deployment
- `ALLOWED_ORIGINS` support for configurable CORS
- automatic migration for `helpful_votes` and `saved_notes`

See [DEPLOYMENT.md](/Users/singh/CourseShare/backend/DEPLOYMENT.md:1) for deployment instructions.
