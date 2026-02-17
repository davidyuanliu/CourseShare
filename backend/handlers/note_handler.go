package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/davidyuanliu/CourseShare/backend/config"
	"github.com/davidyuanliu/CourseShare/backend/models"
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

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if note.Title == "" || note.Content == "" || note.CourseID == 0 {
		http.Error(w, "Title, content, and courseId are required", http.StatusBadRequest)
		return
	}

	config.DB.Create(&note)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}
