package router

import (
	"simple-api-go/http/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	userService := service.NewUserService(db)
	tokenService := service.NewTokenService(db)
	authService := service.NewAuthService(db, userService, tokenService)
	productService := service.NewProductService(db)
	cartService := service.NewCartService(db)
	purchaseService := service.NewPurchaseService(db)

	api := app.Group("/api")

	AuthRoutes(api, authService, userService, tokenService)
	UserRoutes(api, userService, tokenService)
	ProductRoutes(api, productService, tokenService)
	CartRoutes(api, cartService, tokenService)
	PurchaseRoutes(api, purchaseService, tokenService)
}
