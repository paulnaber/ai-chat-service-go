package api

import (
	"ai-chat-service-go/internal/database"
	"ai-chat-service-go/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Ensure ChatServer implements ServerInterface
var _ ServerInterface = (*ChatServer)(nil)

type ChatServer struct {
	Store *database.Queries
}

func (s *ChatServer) GetChats(fiberContext *fiber.Ctx) error {
	// mocked user for now
	user := middleware.MockUserInfo("paul.naber@gmail.com", []string{"admin"})
	chats, err := s.Store.GetChatsByUserEmail(fiberContext.Context(), user.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch chats")
	}
	return fiberContext.JSON(chats)
}

func (s *ChatServer) CreateChat(fiberContext *fiber.Ctx) error {
	// TODO I guess I need a mapper here
	// because the sqlc generated structs for me to use for querying
	// oapi-codegen generates the routes and structs for the endpoints
	// var body CreateChatJSONBody                                            // this is the struct of the openapi spec
	var mappedBody database.CreateChatParams = database.CreateChatParams{} // this is the struct of the generated sqlc code
	chat, err := s.Store.CreateChat(fiberContext.Context(), mappedBody)    // Assuming CreateChat(content string) exists
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create chat")
	}

	return fiberContext.JSON(chat)
}

func (s *ChatServer) CreateMessage(c *fiber.Ctx, chatId int) error {
	return fiber.NewError(fiber.StatusNotImplemented, "CreateMessage not implemented yet")
}

func (s *ChatServer) GetMessages(c *fiber.Ctx, chatId int) error {
	return fiber.NewError(fiber.StatusNotImplemented, "GetMessages not implemented yet")
}
