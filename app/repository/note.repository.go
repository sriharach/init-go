package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"type:uuid"`
	Completed bool      `json:"completed"`
	Title     string    `json:"title"`
	Sub_title string    `json:"sub_title"`
	Text      string    `json:"text"`
}

func NewNoteTable(db *gorm.DB) {
	db.Table("notes").AutoMigrate(&Note{})
}
