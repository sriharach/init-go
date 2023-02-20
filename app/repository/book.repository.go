package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid"`
	Name string    `json:"name"`
}

func NewBookTable(db *gorm.DB) {
	db.Table("books").AutoMigrate(&Book{})
}
