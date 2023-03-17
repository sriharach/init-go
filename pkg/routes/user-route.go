package routes

import (
	"api-enjor/app/controllers"
	"api-enjor/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserRoute interface {
	UserRoute(a *fiber.App)
}

type UserRouteOption struct {
	userController controllers.UserController
}

func NewRouteUser(uc controllers.UserController) UserRoute {
	return &UserRouteOption{
		userController: uc,
	}
}

func (uc *UserRouteOption) UserRoute(a *fiber.App) {
	router := a.Group("/api/user", middleware.DeserializeUser)
	router.Get("/", uc.userController.GetUser)
}
