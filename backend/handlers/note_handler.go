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
	config.DB.Preload("User").Where("course_id = ?", courseID).Find(&notes)

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
	result := config.DB.Preload("User").First(&note, id)

	if result.Error != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(uint)

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

	// Set note owner to authenticated user
	note.UserID = userID

	config.DB.Create(&note)

	// Reload note with author information
	config.DB.Preload("User").First(&note, note.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Note created successfully",
		"note":    note,
	})
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uint)

	var note models.Note
	result := config.DB.First(&note, id)

	if result.Error != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	if note.UserID != userID {
		http.Error(w, "Unauthorized: you can only delete your own notes", http.StatusForbidden)
		return
	}

	config.DB.Delete(&note)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Note deleted successfully",
	})
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uint)

	var updateData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	missingFields := []string{}
	if strings.TrimSpace(updateData.Title) == "" {
		missingFields = append(missingFields, "title")
	}
	if strings.TrimSpace(updateData.Content) == "" {
		missingFields = append(missingFields, "content")
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

	var note models.Note
	result := config.DB.First(&note, id)

	if result.Error != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	if note.UserID != userID {
		http.Error(w, "Unauthorized: you can only update your own notes", http.StatusForbidden)
		return
	}

	// Update the note
	config.DB.Model(&note).Updates(models.Note{
		Title:   updateData.Title,
		Content: updateData.Content,
	})

	// Reload note with author information
	config.DB.Preload("User").First(&note, note.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Note updated successfully",
		"note":    note,
	})
}
