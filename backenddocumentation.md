# Backend Documentation

## Base URL
http://localhost:8080

## Tech Stack
- Spring Boot
- Java
- MongoDB

## Endpoints

### Authentication
- POST /register → Register user
- POST /login → Login user

### Courses
- GET /courses → Fetch all courses
- POST /courses → Add a new course

### Notes
- POST /notes → Upload notes
- GET /notes/{courseId} → Get notes for a course

## Error Handling
- 400: Bad Request
- 401: Unauthorized
- 500: Internal Server Error

## Notes
All endpoints return JSON responses.