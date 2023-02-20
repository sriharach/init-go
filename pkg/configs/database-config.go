package configs

import (
	"api-enjor/pkg/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	dsn, _ := utils.ConnectionURLBuilder("POSTGRES")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	fmt.Println("Database connected!")

	return db
}
