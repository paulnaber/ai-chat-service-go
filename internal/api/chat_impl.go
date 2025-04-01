package api

import (
	"ai-chat-service-go/internal/database"

	"github.com/gofiber/fiber/v2"
)

type ChatServer struct {
	store *database.Queries
}

func (s *ChatServer) GetChats(c *fiber.Ctx) error {
	chats, err := s.store.GetChatsByUserEmail()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch chats")
	}
	return c.JSON(chats)
}

func (s *ChatServer) CreateChat(c *fiber.Ctx) error {
	var body api.CreateChatJSONBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	chat, err := s.store.CreateChat(body.Content) // Assuming CreateChat(content string) exists
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create chat")
	}

	return c.JSON(chat)
}
