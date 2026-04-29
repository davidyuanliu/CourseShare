# CourseShare

CourseShare is a collaborative web platform where students can share, rate, and improve class notes over time.

## Quick Access
You can access the live version of the application here:
[https://davidyuanliu.github.io/CourseShare/courses](https://davidyuanliu.github.io/CourseShare/courses)

---

## Local Development Setup

To run the application locally, you will need to set up both the frontend and the backend.

### Prerequisites
- **Node.js** (v18 or later recommended)
- **Go** (v1.25 or later as specified in `go.mod`)
- **Angular CLI** (`npm install -g @angular/cli`)

### 1. Backend Setup (Go)
The backend uses Go with GORM and can run with either SQLite (local) or PostgreSQL.

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the backend:
   ```bash
   go run .
   ```
   *By default, it will create a local `courseshare.db` SQLite database.*

### 2. Frontend Setup (Angular)
The frontend is built with Angular and communicates with the backend API.

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm start
   ```
4. Access the application at `http://localhost:4200`.

---

## Project Structure
- `/frontend`: Angular source code, components, and styles.
- `/backend`: Go source code, models, handlers, and database configuration.
- `/backend/docs`: API documentation and collections.

## Features
- **User Authentication**: Secure registration and login using JWT.
- **Course Management**: Create and list courses.
- **Note Sharing**: Upload, edit, and delete class notes.
- **Collaboration**: Rate notes and save them for future reference.
- **Tag Filtering**: Easily find notes using tags.

## Members
- Bakshish Singh – Backend (Go)
- Rajesh Kumar Bathula – Backend (Go)
- David Liu – Frontend (Angular)
