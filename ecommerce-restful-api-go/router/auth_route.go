package router

import (
	"simple-api-go/http/controller"
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(
	api fiber.Router, a service.AuthService, u service.UserService,
	t service.TokenService,
) {
	authController := controller.NewAuthController(a, u, t)

	auth := api.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
}
