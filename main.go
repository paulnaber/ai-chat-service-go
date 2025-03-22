// AI Chat Service API
//
// HTTP API for the AI Chat Service
//
//	Schemes: http, https
//	Host: localhost:3000
//	BasePath: /
//	Version: 1.0.0
//	License: MIT http://opensource.org/licenses/MIT
//	Contact: API Support <support@example.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	SecurityDefinitions:
//	  BearerAuth:
//	    type: apiKey
//	    in: header
//	    name: Authorization
//
// swagger:meta
package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "ai-chat-service-go/docs/swagger" // Import swagger docs
	"ai-chat-service-go/internal/api"
	"ai-chat-service-go/internal/config"
	"ai-chat-service-go/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbConn, err := sql.Open("postgres", cfg.Database.GetConnectionString())
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	store := database.New(dbConn)

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: api.ErrorHandler,
	})

	// Setup middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Configure CORS
	corsConfig := cors.Config{
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}

	// Handle AllowCredentials and AllowOrigins correctly
	if cfg.CORS.AllowedOrigins == "*" {
		corsConfig.AllowOrigins = "*"
		corsConfig.AllowCredentials = false
	} else {
		corsConfig.AllowOrigins = cfg.CORS.AllowedOrigins // Expected to be comma-separated
		corsConfig.AllowCredentials = true
	}

	app.Use(cors.New(corsConfig))

	// Setup Swagger
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: true,
	}))

	// Setup routes
	api.SetupRoutes(app, store, cfg)

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
