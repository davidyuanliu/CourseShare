package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"courseshare/config"
	"courseshare/models"
)

var errInvalidTags = errors.New("invalid tags")

// validateTags validates tag input
func validateTags(tags interface{}) (datatypes.JSON, error) {
	if tags == nil {
		return datatypes.JSON([]byte("[]")), nil
	}

	switch v := tags.(type) {
	case []interface{}:
		normalizedTags := make([]string, 0, len(v))

		// Validate each tag
		for _, tag := range v {
			tagStr, ok := tag.(string)
			if !ok {
				return nil, errInvalidTags
			}
			tagStr = strings.TrimSpace(tagStr)
			if tagStr == "" {
				return nil, errInvalidTags
			}
			if len(tagStr) > 50 {
				return nil, errInvalidTags
			}
			normalizedTags = append(normalizedTags, tagStr)
		}
		tagBytes, _ := json.Marshal(normalizedTags)
		return datatypes.JSON(tagBytes), nil
	default:
		return nil, errInvalidTags
	}
}

func noteHasTag(note models.Note, filterTag string) bool {
	if len(note.Tags) == 0 {
		return false
	}

	var tags []string
	if err := json.Unmarshal(note.Tags, &tags); err != nil {
		return false
	}

	for _, tag := range tags {
		if strings.EqualFold(strings.TrimSpace(tag), filterTag) {
			return true
		}
	}

	return false
}

func getOptionalUserIDFromRequest(r *http.Request) *uint {
	if userID, ok := r.Context().Value("userID").(uint); ok {
		return &userID
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-key-please-change-in-production"
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return nil
	}

	userID := uint(userIDFloat)
	return &userID
}

func applyHelpfulMetadata(notes []models.Note, userID *uint) error {
	if len(notes) == 0 {
		return nil
	}

	noteIDs := make([]uint, 0, len(notes))
	for _, note := range notes {
		noteIDs = append(noteIDs, note.ID)
	}

	var countRows []struct {
		NoteID uint
		Count  int64
	}
	if err := config.DB.Model(&models.HelpfulVote{}).
		Select("note_id, COUNT(*) as count").
		Where("note_id IN ?", noteIDs).
		Group("note_id").
		Scan(&countRows).Error; err != nil {
		return err
	}

	countsByNoteID := make(map[uint]int64, len(countRows))
	for _, row := range countRows {
		countsByNoteID[row.NoteID] = row.Count
	}

	helpfulByNoteID := map[uint]bool{}
	if userID != nil {
		var votes []models.HelpfulVote
		if err := config.DB.Where("user_id = ? AND note_id IN ?", *userID, noteIDs).Find(&votes).Error; err != nil {
			return err
		}
		for _, vote := range votes {
			helpfulByNoteID[vote.NoteID] = true
		}
	}

	for i := range notes {
		notes[i].HelpfulCount = countsByNoteID[notes[i].ID]
		notes[i].IsHelpful = helpfulByNoteID[notes[i].ID]
	}

	return nil
}

func applySavedMetadata(notes []models.Note, userID *uint) error {
	if len(notes) == 0 || userID == nil {
		return nil
	}

	noteIDs := make([]uint, 0, len(notes))
	for _, note := range notes {
		noteIDs = append(noteIDs, note.ID)
	}

	var saves []models.SavedNote
	if err := config.DB.Where("user_id = ? AND note_id IN ?", *userID, noteIDs).Find(&saves).Error; err != nil {
		return err
	}

	savedByNoteID := make(map[uint]bool, len(saves))
	for _, save := range saves {
		savedByNoteID[save.NoteID] = true
	}

	for i := range notes {
		notes[i].IsSaved = savedByNoteID[notes[i].ID]
	}

	return nil
}

func applyNoteMetadata(notes []models.Note, userID *uint) error {
	if err := applyHelpfulMetadata(notes, userID); err != nil {
		return err
	}
	if err := applySavedMetadata(notes, userID); err != nil {
		return err
	}
	return nil
}

func applyNoteMetadataToNote(note *models.Note, userID *uint) error {
	notes := []models.Note{*note}
	if err := applyNoteMetadata(notes, userID); err != nil {
		return err
	}
	*note = notes[0]
	return nil
}

func GetNotesByCourse(w http.ResponseWriter, r *http.Request) {
	courseIDParam := chi.URLParam(r, "id")

	courseID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var notes []models.Note
	config.DB.Preload("User").Where("course_id = ?", courseID).Find(&notes)

	tagFilter := strings.TrimSpace(r.URL.Query().Get("tag"))
	if tagFilter != "" {
		filteredNotes := make([]models.Note, 0)
		for _, note := range notes {
			if noteHasTag(note, tagFilter) {
				filteredNotes = append(filteredNotes, note)
			}
		}
		notes = filteredNotes
	}

	if err := applyNoteMetadata(notes, getOptionalUserIDFromRequest(r)); err != nil {
		http.Error(w, "Failed to retrieve note metadata", http.StatusInternalServerError)
		return
	}

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

	if err := applyNoteMetadataToNote(&note, getOptionalUserIDFromRequest(r)); err != nil {
		http.Error(w, "Failed to retrieve note metadata", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(uint)

	var noteRequest struct {
		Title      string        `json:"title"`
		Content    string        `json:"content"`
		CourseID   uint          `json:"courseId"`
		AuthorName string        `json:"authorName"`
		Tags       []interface{} `json:"tags"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&noteRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	missingFields := []string{}
	if strings.TrimSpace(noteRequest.Title) == "" {
		missingFields = append(missingFields, "title")
	}
	if strings.TrimSpace(noteRequest.Content) == "" {
		missingFields = append(missingFields, "content")
	}
	if noteRequest.CourseID == 0 {
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

	// Validate and process tags
	tags, err := validateTags(noteRequest.Tags)
	if err != nil {
		http.Error(w, "Invalid tags provided", http.StatusBadRequest)
		return
	}

	note := models.Note{
		Title:      noteRequest.Title,
		Content:    noteRequest.Content,
		CourseID:   noteRequest.CourseID,
		UserID:     userID,
		AuthorName: noteRequest.AuthorName,
		Tags:       tags,
	}

	config.DB.Create(&note)

	// Reload note with author information
	config.DB.Preload("User").Preload("Course").First(&note, note.ID)
	if err := applyNoteMetadataToNote(&note, &userID); err != nil {
		http.Error(w, "Failed to retrieve note metadata", http.StatusInternalServerError)
		return
	}

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
		Title      string        `json:"title"`
		Content    string        `json:"content"`
		AuthorName string        `json:"authorName"`
		Tags       []interface{} `json:"tags"`
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

	// Validate and process tags
	tags, err := validateTags(updateData.Tags)
	if err != nil {
		http.Error(w, "Invalid tags provided", http.StatusBadRequest)
		return
	}

	// Update the note
	config.DB.Model(&note).Updates(models.Note{
		Title:      updateData.Title,
		Content:    updateData.Content,
		AuthorName: updateData.AuthorName,
		Tags:       tags,
	})

	// Reload note with author information
	config.DB.Preload("User").Preload("Course").First(&note, note.ID)
	if err := applyNoteMetadataToNote(&note, &userID); err != nil {
		http.Error(w, "Failed to retrieve note metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Note updated successfully",
		"note":    note,
	})
}

func GetMyNotes(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Unauthorized: no user ID found", http.StatusUnauthorized)
		return
	}

	var notes []models.Note
	result := config.DB.Preload("User").Preload("Course").Where("user_id = ?", userID).Find(&notes)

	if result.Error != nil {
		http.Error(w, "Failed to retrieve notes", http.StatusInternalServerError)
		return
	}

	// Return empty array instead of null if no notes
	if notes == nil {
		notes = []models.Note{}
	}

	userIDUint, _ := userID.(uint)
	if err := applyNoteMetadata(notes, &userIDUint); err != nil {
		http.Error(w, "Failed to retrieve note metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}

func MarkNoteHelpful(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized: no user ID found", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var note models.Note
	if err := config.DB.First(&note, id).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	var existingVote models.HelpfulVote
	if err := config.DB.Where("note_id = ? AND user_id = ?", note.ID, userID).First(&existingVote).Error; err == nil {
		http.Error(w, "You have already marked this note as helpful", http.StatusConflict)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Failed to retrieve helpful vote data", http.StatusInternalServerError)
		return
	}

	vote := models.HelpfulVote{
		NoteID: note.ID,
		UserID: userID,
	}
	if err := config.DB.Create(&vote).Error; err != nil {
		http.Error(w, "Failed to mark note as helpful", http.StatusInternalServerError)
		return
	}

	var helpfulCount int64
	if err := config.DB.Model(&models.HelpfulVote{}).Where("note_id = ?", note.ID).Count(&helpfulCount).Error; err != nil {
		http.Error(w, "Failed to retrieve helpful vote data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Note marked as helpful",
		"helpfulCount": helpfulCount,
		"isHelpful":    true,
	})
}

func UnmarkNoteHelpful(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized: no user ID found", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var note models.Note
	if err := config.DB.First(&note, id).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	result := config.DB.Unscoped().Where("note_id = ? AND user_id = ?", note.ID, userID).Delete(&models.HelpfulVote{})
	if result.Error != nil {
		http.Error(w, "Failed to remove helpful vote", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Helpful vote not found", http.StatusNotFound)
		return
	}

	var helpfulCount int64
	if err := config.DB.Model(&models.HelpfulVote{}).Where("note_id = ?", note.ID).Count(&helpfulCount).Error; err != nil {
		http.Error(w, "Failed to retrieve helpful vote data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Helpful vote removed",
		"helpfulCount": helpfulCount,
		"isHelpful":    false,
	})
}

func SaveNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized: no user ID found", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var note models.Note
	if err := config.DB.First(&note, id).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	var existingSave models.SavedNote
	if err := config.DB.Where("note_id = ? AND user_id = ?", note.ID, userID).First(&existingSave).Error; err == nil {
		http.Error(w, "You have already saved this note", http.StatusConflict)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Failed to retrieve saved note data", http.StatusInternalServerError)
		return
	}

	save := models.SavedNote{
		NoteID: note.ID,
		UserID: userID,
	}
	if err := config.DB.Create(&save).Error; err != nil {
		http.Error(w, "Failed to save note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Note saved successfully",
		"isSaved": true,
	})
}

func UnsaveNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized: no user ID found", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var note models.Note
	if err := config.DB.First(&note, id).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	result := config.DB.Unscoped().Where("note_id = ? AND user_id = ?", note.ID, userID).Delete(&models.SavedNote{})
	if result.Error != nil {
		http.Error(w, "Failed to unsave note", http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		http.Error(w, "Saved note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Note unsaved successfully",
		"isSaved": false,
	})
}

func GetSavedNotes(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "Unauthorized: no user ID found", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var savedRecords []models.SavedNote
	if err := config.DB.Preload("Note.User").Preload("Note.Course").Where("user_id = ?", userID).Find(&savedRecords).Error; err != nil {
		http.Error(w, "Failed to retrieve saved notes", http.StatusInternalServerError)
		return
	}

	notes := make([]models.Note, 0, len(savedRecords))
	for _, record := range savedRecords {
		note := record.Note
		notes = append(notes, note)
	}

	if notes == nil {
		notes = []models.Note{}
	}

	if err := applyNoteMetadata(notes, &userID); err != nil {
		http.Error(w, "Failed to retrieve note metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}
