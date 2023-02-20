package routes

import (
	"api-enjor/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func NotesGroupApi(a *fiber.App) {
	// Create routes group.
	note := a.Group("/api/v1/note")
	note.Get("/", controllers.GetNotes)
	note.Post("/create", controllers.CreateNote)
}
