package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	CourseID     uint           `json:"courseId"`
	UserID       uint           `json:"userId"`
	Tags         datatypes.JSON `json:"tags" gorm:"type:json"`
	User         *User          `json:"author"`
	Course       *Course        `json:"course"`
	HelpfulVotes []HelpfulVote  `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	SavedNotes   []SavedNote    `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	HelpfulCount int64          `json:"helpfulCount" gorm:"-"`
	IsHelpful    bool           `json:"isHelpful" gorm:"-"`
	IsSaved      bool           `json:"isSaved" gorm:"-"`
}
