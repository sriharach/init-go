package controllers

import (
	"api-enjor/app/repository"
	"api-enjor/pkg/configs"

	"github.com/gofiber/fiber/v2"
)

func GetBooks(c *fiber.Ctx) error {
	db := configs.InitDatabase()
	var books []repository.Book

	db.Find(&books)

	if len(books) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No notes present", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Notes Found", "data": books})
}

// func CreateBook(c *fiber.Ctx) error {

// }
