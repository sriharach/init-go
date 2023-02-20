package repository

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID         uuid.UUID `json:"id" gorm:"type:uuid"`
	User_name  string    `json:"user_name" gorm:"type:varchar(50);"`
	Password   string    `json:"password" gorm:"type:varchar;"`
	E_mail     string    `json:"e_mail" gorm:"type:varchar(50);unique;not null"`
	First_name string    `json:"first_name" gorm:"type:varchar(50);"`
	Last_name  string    `json:"last_name" gorm:"type:varchar(50);"`
	Activate   uint8     `json:"activate" gorm:"type:uint;default:1"` //uint8 smailint postgres
}

func NewUsersTable(db *gorm.DB) {
	err := db.Table("users").AutoMigrate(&Users{})
	if err != nil {
		log.Fatal("fail migrate:", err.Error())
	}
}
