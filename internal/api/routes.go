package api

import (
	"github.com/gofiber/fiber/v2"

	"ai-chat-service-go/internal/config"
	"ai-chat-service-go/internal/database"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, store *database.Queries, cfg *config.Config) {
	// Create API group with version prefix
	api := app.Group("/v1")

	// Create handlers
	chatHandler := NewChatHandler(store)
	messageHandler := NewMessageHandler(store)

	// Chat routes
	chats := api.Group("/chats")

	// TODO skip auth for now
	// chats.Use(middleware.Auth(cfg.Auth))

	// swagger:route POST /v1/chats Chats createChat
	// Creates a new chat with the initial message content.
	// User identity (email) is extracted from JWT token.
	// responses:
	//   200: ChatWithMessageResponse
	//   400: ErrorResponse
	//   401: ErrorResponse
	//   500: ErrorResponse
	chats.Post("/", chatHandler.CreateChat)

	// swagger:route GET /v1/chats Chats getChats
	// Returns an array of chats owned by the user.
	// User identity (email) is extracted from JWT token.
	// responses:
	//   200: []ChatResponse
	//   401: ErrorResponse
	//   500: ErrorResponse
	chats.Get("/", chatHandler.GetChats)

	// Message routes
	// swagger:route POST /v1/chats/{chatId}/messages Messages createMessage
	// Creates a new message in a specific chat.
	// User identity (email) is extracted from JWT token.
	// parameters:
	//   + name: chatId
	//     in: path
	//     required: true
	//     type: string
	// responses:
	//   200: MessageResponse
	//   400: ErrorResponse
	//   401: ErrorResponse
	//   403: ErrorResponse
	//   404: ErrorResponse
	//   500: ErrorResponse
	chats.Post("/:chatId/messages", messageHandler.CreateMessage)

	// swagger:route GET /v1/chats/{chatId}/messages Messages getMessages
	// Returns an array of messages for the specified chat.
	// User identity (email) is extracted from JWT token.
	// parameters:
	//   + name: chatId
	//     in: path
	//     required: true
	//     type: string
	// responses:
	//   200: []MessageResponse
	//   400: ErrorResponse
	//   401: ErrorResponse
	//   403: ErrorResponse
	//   404: ErrorResponse
	//   500: ErrorResponse
	chats.Get("/:chatId/messages", messageHandler.GetMessages)
}
