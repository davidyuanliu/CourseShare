package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/davidyuanliu/CourseShare/backend/config"
	"github.com/davidyuanliu/CourseShare/backend/models"
)

func GetCourses(w http.ResponseWriter, r *http.Request) {
	var courses []models.Course

	config.DB.Find(&courses)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if course.Name == "" {
		http.Error(w, "Course name is required", http.StatusBadRequest)
		return
	}

	config.DB.Create(&course)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}
