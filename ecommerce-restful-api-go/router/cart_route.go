package router

import (
	"simple-api-go/http/controller"
	m "simple-api-go/http/middleware"
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
)

func CartRoutes(api fiber.Router, c service.CartService, t service.TokenService) {
	cartController := controller.NewCartController(c)

	route := api.Group("/carts")

	route.Get("/", m.Auth(t), cartController.Index)
	route.Post("/add", m.Auth(t), cartController.CreateOrUpdate)
	route.Post("/remove/:productId", m.Auth(t), cartController.Delete)
	route.Post("/pay", m.Auth(t), cartController.Pay)

}
