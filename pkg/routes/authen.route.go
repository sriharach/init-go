package routes

import (
	"api-enjor/app/controllers"
	"api-enjor/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type AuthenRouteController struct {
	authenController controllers.AuthController
}

func NewRouteAuth(ac controllers.AuthController) AuthenRouteController {
	return AuthenRouteController{ac}
}

func (ac *AuthenRouteController) SigninUserRoute(a *fiber.App) {
	router := a.Group("/api/auth")
	router.Post("/sign", ac.authenController.SigninUserController)
	router.Post("/login", ac.authenController.LoginUserControlles)
	router.Get("/refreshtoken", middleware.DeserializeUser, ac.authenController.RefreshTokenControlles)

	router.Get("/oauth/login", ac.authenController.UserOauthController)
	router.Post("/oauth/callback", ac.authenController.CallbackUserController)
}
