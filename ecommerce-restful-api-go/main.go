package main

import (
	"fmt"
	"log"
	"simple-api-go/config"
	"simple-api-go/database"
	"simple-api-go/router"
	"simple-api-go/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func main() {
	app := setupFiberApp()
	db := setupDatabase()
	defer closeDatabase(db)
	setupRoutes(app, db)

	address := fmt.Sprintf("%s:%d", config.AppHost, config.AppPort)
	log.Fatal(app.Listen(address))
}

func setupFiberApp() *fiber.App {
	app := fiber.New(config.FiberConfig())

	// Middleware setup
	app.Use("/api/auth", config.LimiterConfig())
	app.Use(config.LoggerConfig())
	app.Use(cors.New())
	app.Use(config.RecoverConfig())
	return app
}

func setupDatabase() *gorm.DB {
	db := database.Connect(config.DBHost, config.DBName)
	// Add any additional database setup if needed
	return db
}

func setupRoutes(app *fiber.App, db *gorm.DB) {
	router.Routes(app, db)
	app.Use(util.NotFoundHandler)
}

func closeDatabase(db *gorm.DB) {
	sqlDB, errDB := db.DB()
	if errDB != nil {
		util.Log.Errorf("Error getting database instance: %v", errDB)
		return
	}

	if err := sqlDB.Close(); err != nil {
		util.Log.Errorf("Error closing database connection: %v", err)
	} else {
		util.Log.Info("Database connection closed successfully")
	}
}
