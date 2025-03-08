package router

import (
	"simple-api-go/http/controller"
	"simple-api-go/http/middleware"
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(api fiber.Router, s service.UserService, t service.TokenService) {
	userController := controller.NewUserController(s, t)

	user := api.Group("/users")
	needAuth := user.Use(middleware.Role(t))
	needAuth.Get("/", userController.Index)
	needAuth.Post("/", userController.Create)
	needAuth.Get("/:id", userController.Show)
	needAuth.Put("/:id", userController.Update)
	needAuth.Delete("/:id", userController.Delete)
}
