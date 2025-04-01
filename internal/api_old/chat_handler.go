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

type ChatHandler struct {
	store *database.Queries
}

func NewChatHandler(store *database.Queries) *ChatHandler {
	return &ChatHandler{store: store}
}

func (h *ChatHandler) CreateChat(c *fiber.Ctx) error {
	var req models.CreateChatRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// user := middleware.GetCurrentUser(c)
	// if user == nil {
	// 	return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	// }
	user := middleware.MockUserInfo("paul.naber@gmail.com", []string{"admin"})

	chatID := uuid.New()
	now := time.Now().UTC()

	chat, err := h.store.CreateChat(c.Context(), database.CreateChatParams{
		ID:             chatID,
		Title:          req.Content, // Use the first message as the title
		UserEmail:      user.Email,
		LastActiveDate: now,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create chat")
	}

	// Create the initial message
	message, err := h.store.CreateMessage(c.Context(), database.CreateMessageParams{
		ID:         uuid.New(),
		Content:    req.Content,
		SenderType: string(models.SenderTypeUser),
		ChatID:     chatID,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create initial message")
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

	response := models.ChatWithMessageResponse{
		ChatResponse: models.ChatResponse{
			ID:             chat.ID,
			Title:          chat.Title,
			LastActiveDate: chat.LastActiveDate,
		},
		InitialMessage: models.MessageResponse{
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

func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
	// user := middleware.GetCurrentUser(c)
	// if user == nil {
	// 	return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	// }
	user := middleware.MockUserInfo("paul.naber@gmail.com", []string{"admin"})

	chats, err := h.store.GetChatsByUserEmail(c.Context(), user.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch chats")
	}

	var response []models.ChatResponse
	for _, chat := range chats {
		response = append(response, models.ChatResponse{
			ID:             chat.ID,
			Title:          chat.Title,
			LastActiveDate: chat.LastActiveDate,
		})
	}

	return c.JSON(response)
}
