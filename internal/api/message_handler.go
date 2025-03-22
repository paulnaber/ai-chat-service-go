package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"ai-chat-service-go/internal/db"
	"ai-chat-service-go/internal/middleware"
	"ai-chat-service-go/internal/models"
	"ai-chat-service-go/internal/services"
)

// MessageHandler handles message-related requests
type MessageHandler struct {
	store *db.Store
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(store *db.Store) *MessageHandler {
	return &MessageHandler{
		store: store,
	}
}

// CreateMessage handles the creation of a new message
func (h *MessageHandler) CreateMessage(c *fiber.Ctx) error {
	// Get current user
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(models.NewUnauthorizedError(""))
	}

	// Get chat ID from path
	chatID := c.Params("chatId")
	if chatID == "" {
		return c.Status(http.StatusBadRequest).JSON(models.NewValidationError("Invalid chat ID"))
	}

	// Parse request body
	var req models.CreateMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.NewValidationError("Invalid request body"))
	}

	// Validate request
	if req.Content == "" {
		return c.Status(http.StatusBadRequest).JSON(models.NewValidationError("Invalid request parameters",
			models.ErrorDetail{Field: "content", Value: "Content cannot be empty"}))
	}

	// Check if chat exists and belongs to user
	chat, err := h.store.GetChat(context.Background(), chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(models.NewResourceNotFoundError("The requested chat could not be found",
				models.ErrorDetail{Field: "chatId", Value: chatID}))
		}
		return c.Status(http.StatusInternalServerError).JSON(models.NewServerError(""))
	}

	// Check if chat belongs to user
	if chat.UserEmail != user.Email {
		return c.Status(http.StatusForbidden).JSON(models.NewForbiddenError("You do not have permission to access this chat"))
	}

	// Generate message ID
	messageID := uuid.New().String()
	now := time.Now()

	// Create user message
	// var userMessage models.Message
	err = h.store.ExecTx(context.Background(), func(q *db.Queries) error {
		// Create message
		_, err := q.CreateMessage(context.Background(), db.CreateMessageParams{
			ID:         messageID,
			Content:    req.Content,
			SenderType: string(models.SenderTypeUser),
			ChatID:     chatID,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
		if err != nil {
			return err
		}

		// Update chat last active date
		err = q.UpdateChatLastActive(context.Background(), db.UpdateChatLastActiveParams{
			ID:             chatID,
			LastActiveDate: now,
			UpdatedAt:      now,
		})
		if err != nil {
			return err
		}

		// Map to our domain model
		// userMessage = models.Message{
		// 	ID:         result.ID,
		// 	Content:    result.Content,
		// 	SenderType: models.SenderType(result.SenderType),
		// 	ChatID:     result.ChatID,
		// 	CreatedAt:  result.CreatedAt,
		// 	UpdatedAt:  result.UpdatedAt,
		// }

		return nil
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.NewServerError("Failed to create message"))
	}

	// Generate AI response (in a real app, you'd call an AI service here)
	aiResponse, err := services.GenerateAIResponse(req.Content)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.NewServerError("Failed to generate AI response"))
	}

	// Create AI message
	aiMessageID := uuid.New().String()
	now = time.Now()
	var aiMessage models.Message

	err = h.store.ExecTx(context.Background(), func(q *db.Queries) error {
		// Create AI message
		result, err := q.CreateMessage(context.Background(), db.CreateMessageParams{
			ID:         aiMessageID,
			Content:    aiResponse,
			SenderType: string(models.SenderTypeLLM),
			ChatID:     chatID,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
		if err != nil {
			return err
		}

		// Update chat last active date again
		err = q.UpdateChatLastActive(context.Background(), db.UpdateChatLastActiveParams{
			ID:             chatID,
			LastActiveDate: now,
			UpdatedAt:      now,
		})
		if err != nil {
			return err
		}

		// Map to our domain model
		aiMessage = models.Message{
			ID:         result.ID,
			Content:    result.Content,
			SenderType: models.SenderType(result.SenderType),
			ChatID:     result.ChatID,
			CreatedAt:  result.CreatedAt,
			UpdatedAt:  result.UpdatedAt,
		}

		return nil
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.NewServerError("Failed to create AI message"))
	}

	// Return AI message response
	return c.Status(http.StatusOK).JSON(aiMessage.ToMessageResponse())
}

// GetMessages handles retrieving all messages for a chat
func (h *MessageHandler) GetMessages(c *fiber.Ctx) error {
	// Get current user
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(models.NewUnauthorizedError(""))
	}

	// Get chat ID from path
	chatID := c.Params("chatId")
	if chatID == "" {
		return c.Status(http.StatusBadRequest).JSON(models.NewValidationError("Invalid chat ID",
			models.ErrorDetail{Field: "chatId", Value: "Invalid format"}))
	}

	// Check if chat exists and belongs to user
	chat, err := h.store.GetChat(context.Background(), chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(models.NewResourceNotFoundError("The requested chat could not be found",
				models.ErrorDetail{Field: "chatId", Value: chatID}))
		}
		return c.Status(http.StatusInternalServerError).JSON(models.NewServerError(""))
	}

	// Check if chat belongs to user
	if chat.UserEmail != user.Email {
		return c.Status(http.StatusForbidden).JSON(models.NewForbiddenError("You do not have permission to access this chat"))
	}

	// Query messages
	results, err := h.store.GetMessagesByChatID(context.Background(), chatID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.NewServerError("Failed to fetch messages"))
	}

	// Convert to response format
	var messages []models.MessageResponse
	for _, msg := range results {
		messages = append(messages, models.MessageResponse{
			ID:         msg.ID,
			Content:    msg.Content,
			SenderType: models.SenderType(msg.SenderType),
			CreatedAt:  msg.CreatedAt,
			ChatID:     msg.ChatID,
		})
	}

	return c.Status(http.StatusOK).JSON(messages)
}
