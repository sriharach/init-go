package main

import (
	"api-enjor/app/controllers"
	"api-enjor/app/repository"
	"api-enjor/pkg/configs"
	"api-enjor/pkg/middleware"
	"api-enjor/pkg/routes"
	"api-enjor/pkg/utils"

	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	configs.Godotenv()
	config := configs.FiberConfig()
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.
	var (
		//start connect database
		postgresDB = configs.InitDatabase()

		//connected database to connect controllers
		authenController = controllers.NewAuthController(postgresDB)
		userController   = controllers.NewUserController(postgresDB)

		//and connect router for first process
		authenRouter = routes.NewRouteAuth(authenController)
		userRouter   = routes.NewRouteUser(userController)
	)

	repository.NewUsersTable(postgresDB)

	authenRouter.SigninUserRoute(app)
	userRouter.UserRoute(app)

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
