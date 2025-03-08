package router

import (
	"simple-api-go/http/controller"
	m "simple-api-go/http/middleware"
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(api fiber.Router, p service.ProductService, t service.TokenService) {
	productController := controller.NewProductController(p)

	route := api.Group("/products")
	route.Get("/", productController.Index)
	route.Get("/:id", productController.Show)
	route.Post("/", m.Role(t), productController.Create)
	route.Put("/:id", m.Role(t), productController.Update)
	route.Delete("/:id", m.Role(t), productController.Delete)
}
