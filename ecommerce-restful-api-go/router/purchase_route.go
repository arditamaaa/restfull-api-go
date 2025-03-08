package router

import (
	"simple-api-go/http/controller"
	m "simple-api-go/http/middleware"
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
)

func PurchaseRoutes(api fiber.Router, s service.PurchaseService, t service.TokenService) {
	purchaseController := controller.NewPurchaseController(s)

	route := api.Group("/purchases")
	route.Get("/", m.Auth(t), purchaseController.Index)
	route.Get("/:id", m.Auth(t), purchaseController.Show)
	route.Post("/", m.Auth(t), purchaseController.Create)
	route.Put("/:id", m.Auth(t), purchaseController.Update)
	route.Delete("/:id", m.Auth(t), purchaseController.Delete)
}
