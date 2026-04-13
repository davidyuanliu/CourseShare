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

	"courseshare/config"
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
		ID     uint `json:"id"`
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
		ID     uint `json:"id"`
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
