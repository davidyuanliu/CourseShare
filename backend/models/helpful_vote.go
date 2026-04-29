package models

import "gorm.io/gorm"

type HelpfulVote struct {
	gorm.Model
	NoteID uint `json:"noteId" gorm:"uniqueIndex:idx_note_user_vote"`
	UserID uint `json:"userId" gorm:"uniqueIndex:idx_note_user_vote"`
	Note   Note `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	User   User `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}
