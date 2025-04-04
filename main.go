package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "ai-chat-service-go/docs"
	"ai-chat-service-go/internal/api"
	"ai-chat-service-go/internal/config"
	"ai-chat-service-go/internal/database"
	"ai-chat-service-go/internal/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "github.com/lib/pq"
)

// Implement the server interface
type ChatServer struct{}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbConn, err := sql.Open("postgres", cfg.Database.DatabaseUrl)
	fmt.Println(cfg.Database.DatabaseUrl)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	// connection is up, possible to do querries here
	queries := database.New(dbConn)

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: errors.ErrorHandler,
	})

	// Setup middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Configure CORS
	corsConfig := cors.Config{
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}
	if cfg.CORS.AllowedOrigins == "*" {
		corsConfig.AllowOrigins = "*"
		corsConfig.AllowCredentials = false
	} else {
		corsConfig.AllowOrigins = cfg.CORS.AllowedOrigins // Expected to be comma-separated
		corsConfig.AllowCredentials = true
	}
	app.Use(cors.New(corsConfig))

	// serving the static file is needed in order to serve the Swagger UI
	app.Static("/api.yml", "./api.yml") // Serve your api.yml file here

	// Serve the Swagger UI at '/swagger' and point it to '/api.yml'
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/api.yml",
	}))

	// Setup routes
	chatServer := &api.ChatServer{Store: queries}
	api.RegisterHandlers(app, chatServer)

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
