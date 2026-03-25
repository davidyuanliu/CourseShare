# Sprint 2 Report – CourseShare Backend

## Overview

This document describes the work completed during Sprint 2 and provides detailed documentation of the CourseShare backend API. The CourseShare application allows users to create courses and upload notes associated with each course.

## Sprint 2 Work Completed

During Sprint 2, the following major tasks were completed:

### 1. Integration of Frontend and Backend

The frontend application was successfully connected to the backend API. The frontend can now:

* Fetch courses from the backend
* Create new notes
* Display notes for a specific course
* Send HTTP requests to backend endpoints
* Receive JSON responses from the backend

This integration ensures that the application is fully functional end-to-end, with the frontend communicating with the backend server through REST API calls.

### 2. Unit Tests for All Features

Unit tests were written for the backend features to ensure that:

* API endpoints return correct responses
* Database operations work correctly
* Course creation works properly
* Note creation works properly
* Error handling works for invalid inputs

These unit tests help ensure the reliability and correctness of the backend system.

### 3. Validation for Required Fields on Note Creation Form

Validation was implemented to ensure that users cannot create a note without required fields.

The required fields for note creation are:

* Title
* Content
* Course ID

If any of these fields are missing, the backend returns a `400 Bad Request` error with an appropriate error message. This prevents invalid or incomplete data from being stored in the database.

### 4. Success Message After Submitting a Note

After successfully creating a note, the backend now returns a success response.

The response includes:

* Note ID
* Title
* Content
* Course ID
* Created timestamp

This allows the frontend to display a success message to the user after submitting a note.

### 5. Frontend Application Features

In addition to backend integration, the following front-end components and flows were established in Angular during Sprint 2:

* **Reactive Form for Note Creation**: Implemented a robust Angular Reactive Form handling user inputs for note title and content.
* **Loading States & User Feedback**: Implemented spinners indicating loading states while courses are being fetched, and integrated Material Design snack bars to provide immediate, non-intrusive feedback upon successful note creation.
* **Course Navigation Flow**: Finalized the routing and structural layout using Angular Router, allowing users to select a course from the list and immediately view its corresponding notes on a dedicated page.

---

## Detailed Documentation of Backend API

**Base URL**: `http://localhost:8080`

### API Endpoints

#### 1. Health Check
* **Endpoint**: `GET /`
* **Description**: Checks if the server is running.
* **Response**:
  ```text
  Server is running
  ```

#### 2. Get All Courses
* **Endpoint**: `GET /courses`
* **Description**: Returns a list of all courses stored in the database.
* **Response Example**:
  ```json
  [
    {
      "ID": 1,
      "CreatedAt": "2026-02-16T19:30:00Z",
      "UpdatedAt": "2026-02-16T19:30:00Z",
      "DeletedAt": null,
      "name": "Go Programming"
    }
  ]
  ```

#### 3. Create Course
* **Endpoint**: `POST /courses`
* **Request Body**:
  ```json
  {
    "name": "Database Systems"
  }
  ```
* **Description**: Creates a new course in the database.
* **Success Response**:
  ```json
  {
    "ID": 2,
    "name": "Database Systems"
  }
  ```

#### 4. Get Notes By Course
* **Endpoint**: `GET /courses/{id}/notes`
* **Example**: `GET /courses/1/notes`
* **Description**: Returns all notes associated with a specific course.
* **Response Example**:
  ```json
  [
    {
      "ID": 1,
      "title": "Introduction",
      "content": "This is the first note",
      "courseId": 1
    }
  ]
  ```

#### 5. Get Single Note
* **Endpoint**: `GET /notes/{id}`
* **Example**: `GET /notes/1`
* **Description**: Returns a specific note by its ID.
* **Response Example**:
  ```json
  {
    "ID": 1,
    "title": "Introduction",
    "content": "This is the first note",
    "courseId": 1
  }
  ```

#### 6. Create Note
* **Endpoint**: `POST /notes`
* **Request Body**:
  ```json
  {
    "title": "Golang Basics",
    "content": "Go is a statically typed language.",
    "courseId": 1
  }
  ```
* **Description**: Creates a new note associated with a course.

---

## Unit Testing Implementation

### Testing Overview

During Sprint 2, unit testing was added to the backend to ensure that all API endpoints work correctly and to prevent future changes from breaking existing functionality. The tests were written for the course and note handlers.

To make testing independent from the production database, an SQLite in-memory test database was introduced. This allows tests to run quickly without connecting to the PostgreSQL production database.

### Test Database Setup

A reusable SQLite in-memory database helper was implemented for testing. This helper initializes a temporary database and seeds the global database variable (`config.DB`) so that all handler tests run in an isolated environment.

This ensures:
* Tests do not affect the production database
* Tests run faster
* Tests are repeatable and isolated
* Each test starts with a clean database

### Course Handler Tests

Unit tests were added for the course handlers to verify the following functionality:

#### Test Cases Covered
* **Get Courses**: Ensures that the API returns a list of courses stored in the database.
* **Create Course – Success**: Ensures that a course is successfully created when valid input is provided. Verifies that the API returns HTTP `201 Created`.
* **Create Course – Validation Failure**: Ensures that the API returns HTTP `400 Bad Request` when the course name is missing.

These tests confirm that the course endpoints work correctly and handle invalid input properly.

### Note Handler Tests

Comprehensive tests were written for the note handlers covering multiple scenarios.

#### Test Cases Covered
* **Get Notes By Course**: Ensures that notes are filtered correctly by course ID.
* **Get Single Note – Success**: Ensures that fetching a valid note ID returns HTTP `200 OK`.
* **Get Single Note – Not Found**: Ensures that requesting a non-existing note returns HTTP `404 Not Found`.
* **Create Note – Success**: Ensures that a note is created successfully with valid input. Verifies the response structure:
  ```json
  {
    "message": "Note created successfully",
    "note": { ... }
  }
  ```
* **Create Note – Validation Failure**: Ensures that missing required fields return HTTP `400 Bad Request`. Validation checks for:
  * Missing title
  * Missing content
  * Missing courseId

These tests ensure that the note endpoints behave correctly under both valid and invalid conditions.

### Dependencies for Testing

The SQLite driver dependency was added to the project so that tests can run without using PostgreSQL.

The following files were updated:
* `go.mod`
* `go.sum`

This ensures that the project can install the SQLite driver and run tests successfully on any system.

### Running Unit Tests

To run all unit tests, use the following steps:

**Run All Tests**
```bash
cd backend
go test ./...
```

**Run Only Handler Tests**
```bash
go test ./handlers -v
```

These tests compile the SQLite driver and then execute all unit tests for the handlers.

### Automated Test Coverage

The unit tests automatically verify the following:
* Course listing returns stored courses
* Course creation returns `201 Created`
* Course creation without name returns `400 Bad Request`
* Notes are filtered correctly by course ID
* Single note retrieval returns `200 OK` or `404 Not Found`
* Successful note creation returns success message and note object
* Missing required fields return validation errors

### Manual API Testing (Optional)

After starting the backend server, the API can also be tested manually using curl or Postman.

**Get Courses**
```bash
curl http://localhost:8080/courses
```

**Create Note**
```bash
curl -X POST http://localhost:8080/notes \
-H "Content-Type: application/json" \
-d '{"title":"Test","content":"Body","courseId":1}'
```

**Validation Test (Missing Fields)**
```bash
curl -X POST http://localhost:8080/notes \
-H "Content-Type: application/json" \
-d '{"title":"Test"}'
```
This should return a `400 Bad Request` response.

---

## Frontend Testing Implementation

### Testing Overview

During Sprint 2, comprehensive testing was also set up for the frontend Angular application. This ensures that the user interface correctly interacts with the backend and handles user inputs as expected. The testing approach utilizes Cypress for End-to-End (E2E) UI testing and Jasmine/Karma for Angular component unit testing.

### Cypress End-to-End (E2E) Test

A Cypress test was implemented to simulate a real user interacting with the application in a browser environment.

#### Test Cases Covered
* **Courses Display and Navigation**: 
  * Ensures that the initial project page loads successfully and displays the list of loaded courses.
  * Verifies that the course cards are correctly populated.
  * Ensures that clicking the 'VIEW NOTES' button successfully navigates the user to the correct notes route.

### Angular Unit Tests

Unit tests were added for the frontend features to ensure that components and services function accurately in isolation. The test coverage maintains a robust unit test to function ratio.

* **CourseShareApiService Tests**: 
  * **Get Courses**: Ensures that the API service successfully executes an HTTP `GET` request to the backend and retrieves the course list.
* **CoursesComponent Tests**:
  * **Initialization**: Ensures that the component calls the API service on initialization and correctly populates the courses array upon success.
* **CreateNoteComponent Tests**:
  * **Form Validation (Missing Fields)**: Ensures that the component marks the form as invalid when required fields (title, content) are empty.
  * **Form Validation (Success)**: Ensures the form is evaluated as valid when both the title and content are properly populated.
  * **Create Note – Success**: Ensures that submitting a valid form correctly triggers the API service's `createNote` method, displays a success snackbar message, and navigates the user back to the course notes page.
  * **Create Note – Validation Failure**: Ensures that if the backend returns an error, the component correctly extracts the error message, displays it in the UI, and prevents navigation.

### Running Frontend Tests

To run the frontend tests, use the following commands from the frontend directory:

**Run Angular Unit Tests**
```bash
cd frontend
npm run test
```

**Run Cypress E2E Tests**
```bash
npm run cypress:open
```

---

## Summary

Unit testing was successfully implemented in Sprint 2 to ensure backend reliability. The tests cover all major course and note API functionalities, including success cases, validation errors, and data retrieval. The use of an SQLite in-memory database ensures that tests run independently from the production database and can be executed quickly during development.

### Validation Rules

The following fields are required:
* `title`
* `content`
* `courseId`

If any field is missing, the API returns:
* `400 Bad Request`

**Success Response**: `201 Created`

### Database Schema

#### Course Table

| Field | Type |
|-------|------|
| ID | Integer |
| Name | String |
| CreatedAt | Timestamp |
| UpdatedAt | Timestamp |

#### Note Table

| Field | Type |
|-------|------|
| ID | Integer |
| Title | String |
| Content | Text |
| CourseID | Integer |
| CreatedAt | Timestamp |

## Conclusion

In Sprint 2, the frontend and backend were successfully integrated, unit tests were implemented, validation was added for note creation, and success responses were implemented. Additionally, the backend API was fully documented, including all endpoints, request formats, and responses.

The CourseShare backend is now capable of handling course and note management through RESTful API endpoints and is ready for further feature development.