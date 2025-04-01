package api_old

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"ai-chat-service-go/internal/database"
	"ai-chat-service-go/internal/middleware"
	"ai-chat-service-go/internal/models"
	"ai-chat-service-go/internal/services"
)

type MessageHandler struct {
	store *database.Queries
}

func NewMessageHandler(store *database.Queries) *MessageHandler {
	return &MessageHandler{store: store}
}

func (h *MessageHandler) CreateMessage(c *fiber.Ctx) error {
	var req models.CreateMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// user := middleware.GetCurrentUser(c)
	// if user == nil {
	// 	return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	// }
	user := middleware.MockUserInfo("paul.naber@gmail.com", []string{"admin"})

	chatID, err := uuid.Parse(c.Params("chatId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid chat ID")
	}

	// Check if the chat belongs to the user
	chat, err := h.store.GetChat(c.Context(), chatID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Chat not found")
	}
	if chat.UserEmail != user.Email {
		return fiber.NewError(fiber.StatusForbidden, "You don't have access to this chat")
	}

	now := time.Now().UTC()

	// Create user message
	message, err := h.store.CreateMessage(c.Context(), database.CreateMessageParams{
		ID:         uuid.New(),
		Content:    req.Content,
		SenderType: string(models.SenderTypeUser),
		ChatID:     chatID,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create message")
	}

	// Update chat's last active date
	err = h.store.UpdateChatLastActive(c.Context(), database.UpdateChatLastActiveParams{
		ID:             chatID,
		LastActiveDate: now,
		UpdatedAt:      now,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update chat")
	}

	// Generate AI response
	aiResponse, err := services.GenerateAIResponse(req.Content)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate AI response")
	}

	// Create AI response message
	aiMessage, err := h.store.CreateMessage(c.Context(), database.CreateMessageParams{
		ID:         uuid.New(),
		Content:    aiResponse,
		SenderType: string(models.SenderTypeLLM),
		ChatID:     chatID,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create AI response message")
	}

	response := models.CreateMessageResponse{
		UserMessage: models.MessageResponse{
			ID:         message.ID,
			Content:    message.Content,
			SenderType: models.SenderType(message.SenderType),
			CreatedAt:  message.CreatedAt,
			ChatID:     message.ChatID,
		},
		AIResponse: models.MessageResponse{
			ID:         aiMessage.ID,
			Content:    aiMessage.Content,
			SenderType: models.SenderType(aiMessage.SenderType),
			CreatedAt:  aiMessage.CreatedAt,
			ChatID:     aiMessage.ChatID,
		},
	}

	return c.JSON(response)
}

func (h *MessageHandler) GetMessages(c *fiber.Ctx) error {
	// user := middleware.GetCurrentUser(c)
	// if user == nil {
	// 	return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	// }
	user := middleware.MockUserInfo("paul.naber@gmail.com", []string{"admin"})

	chatID, err := uuid.Parse(c.Params("chatId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid chat ID")
	}

	// Check if the chat belongs to the user
	chat, err := h.store.GetChat(c.Context(), chatID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Chat not found")
	}
	if chat.UserEmail != user.Email {
		return fiber.NewError(fiber.StatusForbidden, "You don't have access to this chat")
	}

	messages, err := h.store.GetMessagesByChatID(c.Context(), chatID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch messages")
	}

	var response []models.MessageResponse
	for _, msg := range messages {
		response = append(response, models.MessageResponse{
			ID:         msg.ID,
			Content:    msg.Content,
			SenderType: models.SenderType(msg.SenderType),
			CreatedAt:  msg.CreatedAt,
			ChatID:     msg.ChatID,
		})
	}

	return c.JSON(response)
}
