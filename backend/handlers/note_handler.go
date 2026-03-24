package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"courseshare/config"
	"courseshare/models"
)

func GetNotesByCourse(w http.ResponseWriter, r *http.Request) {
	courseIDParam := chi.URLParam(r, "id")

	courseID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var notes []models.Note
	config.DB.Where("course_id = ?", courseID).Find(&notes)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var note models.Note
	result := config.DB.First(&note, id)

	if result.Error != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	missingFields := []string{}
	if strings.TrimSpace(note.Title) == "" {
		missingFields = append(missingFields, "title")
	}
	if strings.TrimSpace(note.Content) == "" {
		missingFields = append(missingFields, "content")
	}
	if note.CourseID == 0 {
		missingFields = append(missingFields, "courseId")
	}

	if len(missingFields) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":         "Missing required fields",
			"missingFields": missingFields,
		})
		return
	}

	config.DB.Create(&note)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Note created successfully",
		"note":    note,
	})
}
