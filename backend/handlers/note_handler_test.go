package handlers

import (
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

	body := strings.NewReader(fmt.Sprintf(`{"title":"New","content":"Body","courseId":%d}`, course.ID))
	req := httptest.NewRequest(http.MethodPost, "/notes", body)
	req.Header.Set("Content-Type", "application/json")
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
}

func TestCreateNoteValidation(t *testing.T) {
	setupTestDB(t)

	body := strings.NewReader(`{"title":" ","content":"","courseId":0}`)
	req := httptest.NewRequest(http.MethodPost, "/notes", body)
	req.Header.Set("Content-Type", "application/json")
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
