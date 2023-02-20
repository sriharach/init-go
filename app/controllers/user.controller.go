package controllers

import (
	"api-enjor/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController interface {
	GetUser(c *fiber.Ctx) error
}

type DBgormUser struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) UserController {
	return &DBgormUser{
		DB: db,
	}
}

func (uc *DBgormUser) GetUser(c *fiber.Ctx) error {
	user := utils.ModuleUser(c)
	return c.Status(200).JSON(fiber.Map{
		"c.Locals": user,
	})
}
