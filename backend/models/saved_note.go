package models

import "gorm.io/gorm"

type SavedNote struct {
	gorm.Model
	NoteID uint `json:"noteId" gorm:"uniqueIndex:idx_note_user_save"`
	UserID uint `json:"userId" gorm:"uniqueIndex:idx_note_user_save"`
	Note   Note `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	User   User `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}
