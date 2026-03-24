package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"courseshare/config"
	"courseshare/models"
)

func TestGetCoursesReturnsCourses(t *testing.T) {
	setupTestDB(t)

	courses := []models.Course{{Name: "Math"}, {Name: "Science"}}
	for _, course := range courses {
		if err := config.DB.Create(&course).Error; err != nil {
			t.Fatalf("failed to seed courses: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/courses", nil)
	rr := httptest.NewRecorder()

	GetCourses(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp []models.Course
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 2 {
		t.Fatalf("expected 2 courses, got %d", len(resp))
	}
}

func TestCreateCourseSuccess(t *testing.T) {
	setupTestDB(t)

	body := strings.NewReader(`{"name":"Biology"}`)
	req := httptest.NewRequest(http.MethodPost, "/courses", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateCourse(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var resp models.Course
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID == 0 {
		t.Fatalf("expected course ID to be set")
	}

	if resp.Name != "Biology" {
		t.Fatalf("expected course name Biology, got %s", resp.Name)
	}
}

func TestCreateCourseMissingName(t *testing.T) {
	setupTestDB(t)

	body := strings.NewReader(`{"description":"no name"}`)
	req := httptest.NewRequest(http.MethodPost, "/courses", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateCourse(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}
