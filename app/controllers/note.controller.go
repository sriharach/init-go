package controllers

import (
	"api-enjor/app/repository"
	"api-enjor/pkg/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetNotes(c *fiber.Ctx) error {
	db := configs.InitDatabase()
	var notes []repository.Note

	// find all notes in the database
	db.Find(&notes)

	// If no note is present return an error
	if len(notes) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No notes present", "data": notes})
	}

	// Else return notes
	return c.JSON(fiber.Map{"status": "success", "message": "Notes Found", "data": notes})
}

func CreateNote(c *fiber.Ctx) error {
	db := configs.InitDatabase()
	note := new(repository.Note)

	err := c.BodyParser(note)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	note.ID = uuid.New()

	err = db.Create(&note).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create note", "data": err})
	}
	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Created Note", "data": note})
}
