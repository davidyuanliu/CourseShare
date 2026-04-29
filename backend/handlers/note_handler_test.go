package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"gorm.io/datatypes"

	"courseshare/config"
	authMiddleware "courseshare/middleware"
	"courseshare/models"
)

func TestGetNotesByCourseFiltersByCourseID(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "History"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	otherCourse := models.Course{Name: "Math"}
	if err := config.DB.Create(&otherCourse).Error; err != nil {
		t.Fatalf("failed to seed other course: %v", err)
	}

	notes := []models.Note{
		{Title: "H1", Content: "C1", CourseID: course.ID},
		{Title: "H2", Content: "C2", CourseID: course.ID},
		{Title: "M1", Content: "C3", CourseID: otherCourse.ID},
	}

	for _, note := range notes {
		if err := config.DB.Create(&note).Error; err != nil {
			t.Fatalf("failed to seed notes: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/courses/%d/notes", course.ID), nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/courses/{id}/notes", GetNotesByCourse)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []models.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(resp))
	}
}

func TestGetNoteReturnsNote(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Chem"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	note := models.Note{Title: "Read", Content: "Ch1", CourseID: course.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notes/%d", note.ID), nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/notes/{id}", GetNote)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp models.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != note.ID {
		t.Fatalf("expected note ID %d, got %d", note.ID, resp.ID)
	}
}

func TestGetNoteNotFound(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/notes/999", nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/notes/{id}", GetNote)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestCreateNoteSuccess(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Physics"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "test@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	body := strings.NewReader(fmt.Sprintf(`{"title":"New","content":"Body","courseId":%d}`, course.ID))
	req := httptest.NewRequest(http.MethodPost, "/notes", body)
	req.Header.Set("Content-Type", "application/json")
	// Add user ID to context
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	CreateNote(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp struct {
		Message string      `json:"message"`
		Note    models.Note `json:"note"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Note created successfully" {
		t.Fatalf("unexpected success message: %s", resp.Message)
	}

	if resp.Note.Title != "New" {
		t.Fatalf("expected note title New, got %s", resp.Note.Title)
	}

	if resp.Note.UserID != user.ID {
		t.Fatalf("expected note owner %d, got %d", user.ID, resp.Note.UserID)
	}
}

func TestCreateNoteValidation(t *testing.T) {
	setupTestDB(t)

	user := models.User{Email: "test@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	body := strings.NewReader(`{"title":" ","content":"","courseId":0}`)
	req := httptest.NewRequest(http.MethodPost, "/notes", body)
	req.Header.Set("Content-Type", "application/json")
	// Add user ID to context
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	CreateNote(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	missing, ok := resp["missingFields"].([]interface{})
	if !ok {
		t.Fatalf("missingFields not returned")
	}

	if len(missing) != 3 {
		t.Fatalf("expected 3 missing fields, got %d", len(missing))
	}
}

func TestDeleteNoteSuccess(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Math"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "test@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	note := models.Note{Title: "Algebra", Content: "Ch1", CourseID: course.ID, UserID: user.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/notes/%d", note.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %d", user.ID))
	// Manually add userID to context
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	// Add URL parameters
	req = req.Clone(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))

	rr := httptest.NewRecorder()
	DeleteNote(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp models.Note
	result := config.DB.First(&resp, note.ID)
	if result.Error == nil {
		t.Fatalf("expected note to be deleted, but it still exists")
	}
}

func TestDeleteNoteNotFound(t *testing.T) {
	setupTestDB(t)

	user := models.User{Email: "test@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/notes/999", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %d", user.ID))
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{"999"}}}))

	rr := httptest.NewRecorder()
	DeleteNote(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestDeleteNoteUnauthorized(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Math"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	owner := models.User{Email: "owner@example.com"}
	if err := config.DB.Create(&owner).Error; err != nil {
		t.Fatalf("failed to seed owner user: %v", err)
	}

	attacker := models.User{Email: "attacker@example.com"}
	if err := config.DB.Create(&attacker).Error; err != nil {
		t.Fatalf("failed to seed attacker user: %v", err)
	}

	note := models.Note{Title: "Algebra", Content: "Ch1", CourseID: course.ID, UserID: owner.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/notes/%d", note.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %d", attacker.ID))
	ctx := context.WithValue(req.Context(), "userID", attacker.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))

	rr := httptest.NewRecorder()
	DeleteNote(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rr.Code)
	}
}

func TestUpdateNoteSuccess(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Math"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "owner@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	note := models.Note{Title: "Old Title", Content: "Old Content", CourseID: course.ID, UserID: user.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	body := strings.NewReader(`{"title":"New Title","content":"New Content"}`)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/notes/%d", note.ID), body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))

	rr := httptest.NewRecorder()
	UpdateNote(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		Message string      `json:"message"`
		Note    models.Note `json:"note"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Note.Title != "New Title" {
		t.Fatalf("expected title 'New Title', got '%s'", resp.Note.Title)
	}

	if resp.Note.Content != "New Content" {
		t.Fatalf("expected content 'New Content', got '%s'", resp.Note.Content)
	}
}

func TestUpdateNoteNotFound(t *testing.T) {
	setupTestDB(t)

	user := models.User{Email: "user@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	body := strings.NewReader(`{"title":"New","content":"Content"}`)
	req := httptest.NewRequest(http.MethodPut, "/notes/999", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{"999"}}}))

	rr := httptest.NewRecorder()
	UpdateNote(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestUpdateNoteUnauthorized(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Math"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	owner := models.User{Email: "owner@example.com"}
	if err := config.DB.Create(&owner).Error; err != nil {
		t.Fatalf("failed to seed owner: %v", err)
	}

	attacker := models.User{Email: "attacker@example.com"}
	if err := config.DB.Create(&attacker).Error; err != nil {
		t.Fatalf("failed to seed attacker: %v", err)
	}

	note := models.Note{Title: "Title", Content: "Content", CourseID: course.ID, UserID: owner.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	body := strings.NewReader(`{"title":"Updated","content":"Updated Content"}`)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/notes/%d", note.ID), body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", attacker.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))

	rr := httptest.NewRecorder()
	UpdateNote(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rr.Code)
	}
}

func TestUpdateNoteMissingFields(t *testing.T) {
	setupTestDB(t)

	user := models.User{Email: "user@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	body := strings.NewReader(`{"title":"","content":""}`)
	req := httptest.NewRequest(http.MethodPut, "/notes/1", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{"1"}}}))

	rr := httptest.NewRecorder()
	UpdateNote(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	missing, ok := resp["missingFields"].([]interface{})
	if !ok {
		t.Fatalf("missingFields not returned")
	}

	if len(missing) != 2 {
		t.Fatalf("expected 2 missing fields, got %d", len(missing))
	}
}

func TestGetNoteIncludesAuthor(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "History"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "author@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	note := models.Note{Title: "WWII", Content: "Timeline", CourseID: course.ID, UserID: user.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notes/%d", note.ID), nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/notes/{id}", GetNote)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		ID     uint   `json:"id"`
		Title  string `json:"title"`
		Author *struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"author"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Author == nil {
		t.Fatalf("author is nil in response")
	}

	if resp.Author.Email != "author@example.com" {
		t.Fatalf("expected author email 'author@example.com', got '%s'", resp.Author.Email)
	}

	if resp.Author.ID != user.ID {
		t.Fatalf("expected author ID %d, got %d", user.ID, resp.Author.ID)
	}
}

func TestGetNotesByCourseIncludesAuthors(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Science"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user1 := models.User{Email: "student1@example.com"}
	if err := config.DB.Create(&user1).Error; err != nil {
		t.Fatalf("failed to seed user1: %v", err)
	}

	user2 := models.User{Email: "student2@example.com"}
	if err := config.DB.Create(&user2).Error; err != nil {
		t.Fatalf("failed to seed user2: %v", err)
	}

	notes := []models.Note{
		{Title: "Physics", Content: "Mechanics", CourseID: course.ID, UserID: user1.ID},
		{Title: "Chemistry", Content: "Reactions", CourseID: course.ID, UserID: user2.ID},
	}

	for _, note := range notes {
		if err := config.DB.Create(&note).Error; err != nil {
			t.Fatalf("failed to seed note: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/courses/%d/notes", course.ID), nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/courses/{id}/notes", GetNotesByCourse)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []struct {
		ID     uint   `json:"id"`
		Title  string `json:"title"`
		Author *struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
		} `json:"author"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(resp))
	}

	// Check first note author
	if resp[0].Author == nil {
		t.Fatalf("first note author is nil")
	}
	if resp[0].Author.Email != "student1@example.com" {
		t.Fatalf("expected first author email 'student1@example.com', got '%s'", resp[0].Author.Email)
	}

	// Check second note author
	if resp[1].Author == nil {
		t.Fatalf("second note author is nil")
	}
	if resp[1].Author.Email != "student2@example.com" {
		t.Fatalf("expected second author email 'student2@example.com', got '%s'", resp[1].Author.Email)
	}
}

func TestCreateNoteIncludesAuthor(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Art"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "artist@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	body := strings.NewReader(fmt.Sprintf(`{"title":"Painting","content":"Oil on canvas","courseId":%d}`, course.ID))
	req := httptest.NewRequest(http.MethodPost, "/notes", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	CreateNote(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp struct {
		Note struct {
			Title  string `json:"title"`
			Author *struct {
				ID    uint   `json:"id"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"note"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Note.Author == nil {
		t.Fatalf("author is nil in created note response")
	}

	if resp.Note.Author.Email != "artist@example.com" {
		t.Fatalf("expected author email 'artist@example.com', got '%s'", resp.Note.Author.Email)
	}
}

func TestGetMyNotesSuccess(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Science"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "student@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	otherUser := models.User{Email: "other@example.com"}
	if err := config.DB.Create(&otherUser).Error; err != nil {
		t.Fatalf("failed to seed other user: %v", err)
	}

	// Create notes for different users
	notes := []models.Note{
		{Title: "User's Note 1", Content: "Content 1", CourseID: course.ID, UserID: user.ID},
		{Title: "User's Note 2", Content: "Content 2", CourseID: course.ID, UserID: user.ID},
		{Title: "Other User's Note", Content: "Other content", CourseID: course.ID, UserID: otherUser.ID},
	}

	for _, note := range notes {
		if err := config.DB.Create(&note).Error; err != nil {
			t.Fatalf("failed to seed note: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/my-notes", nil)
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	GetMyNotes(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []models.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Should only return user's notes (2), not other user's note
	if len(resp) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(resp))
	}

	// Verify titles
	if resp[0].Title != "User's Note 1" {
		t.Fatalf("expected first note title 'User's Note 1', got '%s'", resp[0].Title)
	}

	if resp[1].Title != "User's Note 2" {
		t.Fatalf("expected second note title 'User's Note 2', got '%s'", resp[1].Title)
	}

	// Verify author information is included
	if resp[0].User == nil {
		t.Fatalf("first note author is nil")
	}

	if resp[0].User.Email != "student@example.com" {
		t.Fatalf("expected author email 'student@example.com', got '%s'", resp[0].User.Email)
	}
}

func TestGetMyNotesEmpty(t *testing.T) {
	setupTestDB(t)

	user := models.User{Email: "notnotes@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/my-notes", nil)
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	GetMyNotes(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []models.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 0 {
		t.Fatalf("expected 0 notes, got %d", len(resp))
	}
}

func TestGetMyNotesUnauthenticated(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/my-notes", nil)
	// No userID in context - simulating unauthenticated request
	rr := httptest.NewRecorder()

	GetMyNotes(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}
}

func TestGetMyNotesIncludesCourseInfo(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Computer Science"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "student@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	note := models.Note{Title: "Algorithms", Content: "Sorting algorithms", CourseID: course.ID, UserID: user.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/my-notes", nil)
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	GetMyNotes(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []struct {
		ID     uint   `json:"id"`
		Title  string `json:"title"`
		Course *struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"course"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 1 {
		t.Fatalf("expected 1 note, got %d", len(resp))
	}

	if resp[0].Course == nil {
		t.Fatalf("course is nil in response")
	}

	if resp[0].Course.Name != "Computer Science" {
		t.Fatalf("expected course name 'Computer Science', got '%s'", resp[0].Course.Name)
	}
}

func TestCreateNoteWithTags(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Physics"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "test@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	body := strings.NewReader(fmt.Sprintf(`{"title":"Mechanics","content":"Velocity and acceleration","courseId":%d,"tags":["physics","motion","exam"]}`, course.ID))
	req := httptest.NewRequest(http.MethodPost, "/notes", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	CreateNote(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp struct {
		Note struct {
			Title string          `json:"title"`
			Tags  json.RawMessage `json:"tags"`
		} `json:"note"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Note.Title != "Mechanics" {
		t.Fatalf("expected note title Mechanics, got %s", resp.Note.Title)
	}

	var tags []string
	if err := json.Unmarshal(resp.Note.Tags, &tags); err != nil {
		t.Fatalf("failed to decode tags: %v", err)
	}

	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}

	if tags[0] != "physics" {
		t.Fatalf("expected first tag 'physics', got '%s'", tags[0])
	}
}

func TestCreateNoteWithoutTags(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Chemistry"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "test@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	body := strings.NewReader(fmt.Sprintf(`{"title":"Reactions","content":"Chemical bonds","courseId":%d}`, course.ID))
	req := httptest.NewRequest(http.MethodPost, "/notes", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	CreateNote(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp struct {
		Note struct {
			Tags json.RawMessage `json:"tags"`
		} `json:"note"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	var tags []string
	if err := json.Unmarshal(resp.Note.Tags, &tags); err != nil {
		t.Fatalf("failed to decode tags: %v", err)
	}

	if len(tags) != 0 {
		t.Fatalf("expected empty tags, got %d", len(tags))
	}
}

func TestUpdateNoteWithTags(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Math"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "owner@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	note := models.Note{Title: "Algebra", Content: "Equations", CourseID: course.ID, UserID: user.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	body := strings.NewReader(`{"title":"Algebra Advanced","content":"Complex equations","tags":["algebra","advanced","exam"]}`)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/notes/%d", note.ID), body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))

	rr := httptest.NewRecorder()
	UpdateNote(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		Note struct {
			Title string          `json:"title"`
			Tags  json.RawMessage `json:"tags"`
		} `json:"note"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Note.Title != "Algebra Advanced" {
		t.Fatalf("expected updated title 'Algebra Advanced', got '%s'", resp.Note.Title)
	}

	var tags []string
	if err := json.Unmarshal(resp.Note.Tags, &tags); err != nil {
		t.Fatalf("failed to decode tags: %v", err)
	}

	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}
}

func TestCreateNoteRejectsEmptyTag(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Biology"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "student@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	body := strings.NewReader(fmt.Sprintf(`{"title":"Cells","content":"Cell structure","courseId":%d,"tags":["exam","   "]}`, course.ID))
	req := httptest.NewRequest(http.MethodPost, "/notes", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	CreateNote(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestCreateNoteRejectsLongTag(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Biology"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "student@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	longTag := strings.Repeat("a", 51)
	body := strings.NewReader(fmt.Sprintf(`{"title":"Cells","content":"Cell structure","courseId":%d,"tags":["%s"]}`, course.ID, longTag))
	req := httptest.NewRequest(http.MethodPost, "/notes", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	CreateNote(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestUpdateNoteRejectsInvalidTags(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Physics"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "owner@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	note := models.Note{Title: "Motion", Content: "Kinematics", CourseID: course.ID, UserID: user.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	body := strings.NewReader(`{"title":"Motion","content":"Kinematics","tags":["valid",""]}`)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/notes/%d", note.ID), body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))

	rr := httptest.NewRecorder()
	UpdateNote(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestGetNotesByCourseFiltersByTag(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Computer Science"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "student@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	taggedNotes := []models.Note{
		{Title: "Sorting", Content: "Merge sort", CourseID: course.ID, UserID: user.ID, Tags: datatypes.JSON([]byte(`["exam","algorithms"]`))},
		{Title: "Graphs", Content: "Dijkstra", CourseID: course.ID, UserID: user.ID, Tags: datatypes.JSON([]byte(`["concept"]`))},
		{Title: "Dynamic Programming", Content: "Knapsack", CourseID: course.ID, UserID: user.ID, Tags: datatypes.JSON([]byte(`["Algorithms","practice"]`))},
	}

	for _, note := range taggedNotes {
		if err := config.DB.Create(&note).Error; err != nil {
			t.Fatalf("failed to seed note: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/courses/%d/notes?tag=algorithms", course.ID), nil)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/courses/{id}/notes", GetNotesByCourse)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []models.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(resp))
	}
}

func TestMarkNoteHelpfulSuccess(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "History"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	author := models.User{Email: "author@example.com"}
	if err := config.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to seed author: %v", err)
	}

	voter := models.User{Email: "voter@example.com"}
	if err := config.DB.Create(&voter).Error; err != nil {
		t.Fatalf("failed to seed voter: %v", err)
	}

	note := models.Note{Title: "Lecture", Content: "Summary", CourseID: course.ID, UserID: author.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/notes/%d/helpful", note.ID), nil)
	ctx := context.WithValue(req.Context(), "userID", voter.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))
	rr := httptest.NewRecorder()

	MarkNoteHelpful(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		Message      string `json:"message"`
		HelpfulCount int64  `json:"helpfulCount"`
		IsHelpful    bool   `json:"isHelpful"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.HelpfulCount != 1 {
		t.Fatalf("expected helpful count 1, got %d", resp.HelpfulCount)
	}

	if !resp.IsHelpful {
		t.Fatalf("expected isHelpful true")
	}
}

func TestMarkNoteHelpfulPreventsDuplicateVotes(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "History"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "student@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	note := models.Note{Title: "Lecture", Content: "Summary", CourseID: course.ID, UserID: user.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	vote := models.HelpfulVote{NoteID: note.ID, UserID: user.ID}
	if err := config.DB.Create(&vote).Error; err != nil {
		t.Fatalf("failed to seed helpful vote: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/notes/%d/helpful", note.ID), nil)
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))
	rr := httptest.NewRecorder()

	MarkNoteHelpful(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", rr.Code)
	}
}

func TestUnmarkNoteHelpfulSuccess(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "History"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	author := models.User{Email: "author@example.com"}
	if err := config.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to seed author: %v", err)
	}

	voter := models.User{Email: "voter@example.com"}
	if err := config.DB.Create(&voter).Error; err != nil {
		t.Fatalf("failed to seed voter: %v", err)
	}

	note := models.Note{Title: "Lecture", Content: "Summary", CourseID: course.ID, UserID: author.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	vote := models.HelpfulVote{NoteID: note.ID, UserID: voter.ID}
	if err := config.DB.Create(&vote).Error; err != nil {
		t.Fatalf("failed to seed helpful vote: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/notes/%d/helpful", note.ID), nil)
	ctx := context.WithValue(req.Context(), "userID", voter.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))
	rr := httptest.NewRecorder()

	UnmarkNoteHelpful(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		HelpfulCount int64 `json:"helpfulCount"`
		IsHelpful    bool  `json:"isHelpful"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.HelpfulCount != 0 {
		t.Fatalf("expected helpful count 0, got %d", resp.HelpfulCount)
	}

	if resp.IsHelpful {
		t.Fatalf("expected isHelpful false")
	}
}

func TestGetNoteIncludesHelpfulMetadata(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Economics"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	author := models.User{Email: "author@example.com"}
	if err := config.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to seed author: %v", err)
	}

	voter1 := models.User{Email: "voter1@example.com"}
	if err := config.DB.Create(&voter1).Error; err != nil {
		t.Fatalf("failed to seed voter1: %v", err)
	}

	voter2 := models.User{Email: "voter2@example.com"}
	if err := config.DB.Create(&voter2).Error; err != nil {
		t.Fatalf("failed to seed voter2: %v", err)
	}

	note := models.Note{Title: "Supply", Content: "Demand curve", CourseID: course.ID, UserID: author.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	votes := []models.HelpfulVote{
		{NoteID: note.ID, UserID: voter1.ID},
		{NoteID: note.ID, UserID: voter2.ID},
	}
	for _, vote := range votes {
		if err := config.DB.Create(&vote).Error; err != nil {
			t.Fatalf("failed to seed helpful vote: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notes/%d", note.ID), nil)
	ctx := context.WithValue(req.Context(), "userID", voter1.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/notes/{id}", GetNote)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		HelpfulCount int64 `json:"helpfulCount"`
		IsHelpful    bool  `json:"isHelpful"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.HelpfulCount != 2 {
		t.Fatalf("expected helpful count 2, got %d", resp.HelpfulCount)
	}

	if !resp.IsHelpful {
		t.Fatalf("expected isHelpful true")
	}
}

func TestGetNotesByCourseIncludesHelpfulMetadata(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Economics"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	author := models.User{Email: "author@example.com"}
	if err := config.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to seed author: %v", err)
	}

	voter := models.User{Email: "voter@example.com"}
	if err := config.DB.Create(&voter).Error; err != nil {
		t.Fatalf("failed to seed voter: %v", err)
	}

	note1 := models.Note{Title: "Supply", Content: "Demand curve", CourseID: course.ID, UserID: author.ID}
	note2 := models.Note{Title: "Inflation", Content: "CPI", CourseID: course.ID, UserID: author.ID}
	for _, note := range []models.Note{note1, note2} {
		if err := config.DB.Create(&note).Error; err != nil {
			t.Fatalf("failed to seed note: %v", err)
		}
	}

	var savedNotes []models.Note
	if err := config.DB.Order("id asc").Find(&savedNotes).Error; err != nil {
		t.Fatalf("failed to fetch seeded notes: %v", err)
	}

	vote := models.HelpfulVote{NoteID: savedNotes[0].ID, UserID: voter.ID}
	if err := config.DB.Create(&vote).Error; err != nil {
		t.Fatalf("failed to seed helpful vote: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/courses/%d/notes", course.ID), nil)
	ctx := context.WithValue(req.Context(), "userID", voter.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/courses/{id}/notes", GetNotesByCourse)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []struct {
		Title        string `json:"title"`
		HelpfulCount int64  `json:"helpfulCount"`
		IsHelpful    bool   `json:"isHelpful"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(resp))
	}

	if resp[0].HelpfulCount != 1 || !resp[0].IsHelpful {
		t.Fatalf("expected first note to have helpfulCount=1 and isHelpful=true, got count=%d helpful=%v", resp[0].HelpfulCount, resp[0].IsHelpful)
	}

	if resp[1].HelpfulCount != 0 || resp[1].IsHelpful {
		t.Fatalf("expected second note to have helpfulCount=0 and isHelpful=false, got count=%d helpful=%v", resp[1].HelpfulCount, resp[1].IsHelpful)
	}
}

func TestSaveNoteSuccess(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Math"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	author := models.User{Email: "author@example.com"}
	if err := config.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to seed author: %v", err)
	}

	saver := models.User{Email: "saver@example.com"}
	if err := config.DB.Create(&saver).Error; err != nil {
		t.Fatalf("failed to seed saver: %v", err)
	}

	note := models.Note{Title: "Algebra", Content: "Equations", CourseID: course.ID, UserID: author.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/notes/%d/save", note.ID), nil)
	ctx := context.WithValue(req.Context(), "userID", saver.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))
	rr := httptest.NewRecorder()

	SaveNote(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		IsSaved bool `json:"isSaved"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.IsSaved {
		t.Fatalf("expected isSaved true")
	}
}

func TestSaveNotePreventsDuplicateSaves(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Math"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	user := models.User{Email: "student@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	note := models.Note{Title: "Algebra", Content: "Equations", CourseID: course.ID, UserID: user.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	save := models.SavedNote{NoteID: note.ID, UserID: user.ID}
	if err := config.DB.Create(&save).Error; err != nil {
		t.Fatalf("failed to seed saved note: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/notes/%d/save", note.ID), nil)
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))
	rr := httptest.NewRecorder()

	SaveNote(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", rr.Code)
	}
}

func TestUnsaveNoteSuccess(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Math"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	author := models.User{Email: "author@example.com"}
	if err := config.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to seed author: %v", err)
	}

	saver := models.User{Email: "saver@example.com"}
	if err := config.DB.Create(&saver).Error; err != nil {
		t.Fatalf("failed to seed saver: %v", err)
	}

	note := models.Note{Title: "Algebra", Content: "Equations", CourseID: course.ID, UserID: author.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	save := models.SavedNote{NoteID: note.ID, UserID: saver.ID}
	if err := config.DB.Create(&save).Error; err != nil {
		t.Fatalf("failed to seed saved note: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/notes/%d/save", note.ID), nil)
	ctx := context.WithValue(req.Context(), "userID", saver.ID)
	req = req.Clone(context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{URLParams: chi.RouteParams{Keys: []string{"id"}, Values: []string{fmt.Sprintf("%d", note.ID)}}}))
	rr := httptest.NewRecorder()

	UnsaveNote(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		IsSaved bool `json:"isSaved"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.IsSaved {
		t.Fatalf("expected isSaved false")
	}
}

func TestGetSavedNotesSuccess(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Physics"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	author := models.User{Email: "author@example.com"}
	if err := config.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to seed author: %v", err)
	}

	saver := models.User{Email: "saver@example.com"}
	if err := config.DB.Create(&saver).Error; err != nil {
		t.Fatalf("failed to seed saver: %v", err)
	}

	otherUser := models.User{Email: "other@example.com"}
	if err := config.DB.Create(&otherUser).Error; err != nil {
		t.Fatalf("failed to seed other user: %v", err)
	}

	note1 := models.Note{Title: "Mechanics", Content: "Force", CourseID: course.ID, UserID: author.ID}
	note2 := models.Note{Title: "Optics", Content: "Lenses", CourseID: course.ID, UserID: author.ID}
	note3 := models.Note{Title: "Thermo", Content: "Entropy", CourseID: course.ID, UserID: author.ID}
	for _, note := range []models.Note{note1, note2, note3} {
		if err := config.DB.Create(&note).Error; err != nil {
			t.Fatalf("failed to seed note: %v", err)
		}
	}

	var seededNotes []models.Note
	if err := config.DB.Order("id asc").Find(&seededNotes).Error; err != nil {
		t.Fatalf("failed to fetch seeded notes: %v", err)
	}

	saves := []models.SavedNote{
		{NoteID: seededNotes[0].ID, UserID: saver.ID},
		{NoteID: seededNotes[1].ID, UserID: saver.ID},
		{NoteID: seededNotes[2].ID, UserID: otherUser.ID},
	}
	for _, save := range saves {
		if err := config.DB.Create(&save).Error; err != nil {
			t.Fatalf("failed to seed saved note: %v", err)
		}
	}

	vote := models.HelpfulVote{NoteID: seededNotes[0].ID, UserID: otherUser.ID}
	if err := config.DB.Create(&vote).Error; err != nil {
		t.Fatalf("failed to seed helpful vote: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/saved-notes", nil)
	ctx := context.WithValue(req.Context(), "userID", saver.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	GetSavedNotes(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []struct {
		Title        string `json:"title"`
		HelpfulCount int64  `json:"helpfulCount"`
		IsSaved      bool   `json:"isSaved"`
		Course       *struct {
			Name string `json:"name"`
		} `json:"course"`
		Author *struct {
			Email string `json:"email"`
		} `json:"author"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Fatalf("expected 2 saved notes, got %d", len(resp))
	}

	if !resp[0].IsSaved || !resp[1].IsSaved {
		t.Fatalf("expected saved notes to have isSaved true")
	}

	if resp[0].Course == nil || resp[0].Course.Name != "Physics" {
		t.Fatalf("expected course info in saved note response")
	}

	if resp[0].Author == nil || resp[0].Author.Email != "author@example.com" {
		t.Fatalf("expected author info in saved note response")
	}

	if resp[0].HelpfulCount != 1 {
		t.Fatalf("expected helpful count 1 for first saved note, got %d", resp[0].HelpfulCount)
	}
}

func TestGetSavedNotesEmpty(t *testing.T) {
	setupTestDB(t)

	user := models.User{Email: "saver@example.com"}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/saved-notes", nil)
	ctx := context.WithValue(req.Context(), "userID", user.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	GetSavedNotes(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []models.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 0 {
		t.Fatalf("expected 0 saved notes, got %d", len(resp))
	}
}

func TestGetNoteIncludesSavedMetadata(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Biology"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	author := models.User{Email: "author@example.com"}
	if err := config.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to seed author: %v", err)
	}

	saver := models.User{Email: "saver@example.com"}
	if err := config.DB.Create(&saver).Error; err != nil {
		t.Fatalf("failed to seed saver: %v", err)
	}

	note := models.Note{Title: "Cells", Content: "Membrane", CourseID: course.ID, UserID: author.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	save := models.SavedNote{NoteID: note.ID, UserID: saver.ID}
	if err := config.DB.Create(&save).Error; err != nil {
		t.Fatalf("failed to seed saved note: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notes/%d", note.ID), nil)
	ctx := context.WithValue(req.Context(), "userID", saver.ID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/notes/{id}", GetNote)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		IsSaved bool `json:"isSaved"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.IsSaved {
		t.Fatalf("expected isSaved true")
	}
}

func TestProtectedNoteEndpointsRequireAuthentication(t *testing.T) {
	setupTestDB(t)

	course := models.Course{Name: "Security"}
	if err := config.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	owner := models.User{Email: "owner@example.com"}
	if err := config.DB.Create(&owner).Error; err != nil {
		t.Fatalf("failed to seed owner: %v", err)
	}

	note := models.Note{Title: "Auth", Content: "JWT notes", CourseID: course.ID, UserID: owner.ID}
	if err := config.DB.Create(&note).Error; err != nil {
		t.Fatalf("failed to seed note: %v", err)
	}

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		handler    http.Handler
		wantStatus int
		wantBody   string
	}{
		{
			name:       "create note without auth",
			method:     http.MethodPost,
			path:       "/notes",
			body:       fmt.Sprintf(`{"title":"Protected","content":"Body","courseId":%d}`, course.ID),
			handler:    authMiddleware.ExtractUserFromToken(http.HandlerFunc(CreateNote)),
			wantStatus: http.StatusUnauthorized,
			wantBody:   "Missing authorization header",
		},
		{
			name:       "my notes without auth",
			method:     http.MethodGet,
			path:       "/my-notes",
			handler:    authMiddleware.ExtractUserFromToken(http.HandlerFunc(GetMyNotes)),
			wantStatus: http.StatusUnauthorized,
			wantBody:   "Missing authorization header",
		},
		{
			name:       "save note without auth",
			method:     http.MethodPost,
			path:       fmt.Sprintf("/notes/%d/save", note.ID),
			handler:    authMiddleware.ExtractUserFromToken(http.HandlerFunc(SaveNote)),
			wantStatus: http.StatusUnauthorized,
			wantBody:   "Missing authorization header",
		},
		{
			name:       "helpful vote without auth",
			method:     http.MethodPost,
			path:       fmt.Sprintf("/notes/%d/helpful", note.ID),
			handler:    authMiddleware.ExtractUserFromToken(http.HandlerFunc(MarkNoteHelpful)),
			wantStatus: http.StatusUnauthorized,
			wantBody:   "Missing authorization header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader *strings.Reader
			if tt.body == "" {
				bodyReader = strings.NewReader("")
			} else {
				bodyReader = strings.NewReader(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bodyReader)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			rr := httptest.NewRecorder()

			tt.handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}

			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Fatalf("expected response body to contain %q, got %q", tt.wantBody, rr.Body.String())
			}
		})
	}
}

func TestProtectedNoteEndpointsRejectInvalidToken(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/my-notes", nil)
	req.Header.Set("Authorization", "Bearer definitely-not-a-real-token")
	rr := httptest.NewRecorder()

	authMiddleware.ExtractUserFromToken(http.HandlerFunc(GetMyNotes)).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Invalid token") {
		t.Fatalf("expected invalid token message, got %q", rr.Body.String())
	}
}
