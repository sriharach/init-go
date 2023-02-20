package routes

import (
	"api-enjor/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func BooksGroupApi(a *fiber.App) {
	book := a.Group("/api/v2")
	book.Get("/book", controllers.GetBooks)
	// book.Post("/book/create", controllers.CreateBook)
}
