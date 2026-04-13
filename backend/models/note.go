package models

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	CourseID uint   `json:"courseId"`
	UserID   uint   `json:"userId"`
	User     *User  `json:"author"`
}
