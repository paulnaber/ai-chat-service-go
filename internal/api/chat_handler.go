package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"ai-chat-service-go/internal/db"
	"ai-chat-service-go/internal/middleware"
	"ai-chat-service-go/internal/models"
)

// ChatHandler handles chat-related requests
type ChatHandler struct {
	store *db.Store
}

// NewChatHandler creates a new chat handler
func NewChatHandler(store *db.Store) *ChatHandler {
	return &ChatHandler{
		store: store,
	}
}

// CreateChat handles the creation of a new chat
func (h *ChatHandler) CreateChat(c *fiber.Ctx) error {
	// Get current user
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(models.NewUnauthorizedError(""))
	}

	// Parse request body
	var req models.CreateChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.NewValidationError("Invalid request body"))
	}

	// Validate request
	if req.Content == "" {
		return c.Status(http.StatusBadRequest).JSON(models.NewValidationError("Invalid request parameters", 
			models.ErrorDetail{Field: "content", Value: "Content cannot be empty"}))
	}

	// Generate IDs
	chatID := uuid.New().String()
	messageID := uuid.New().String()
	now := time.Now()

	// Execute the transaction
	var chat models.Chat
	var message models.Message

	err := h.store.ExecTx(context.Background(), func(q *db.Queries) error {
		// Create chat
		var err error
		result, err := q.CreateChat(context.Background(), db.CreateChatParams{
			ID:             chatID,
			Title:          req.Content, // Using the first message as title
			UserEmail:      user.Email,
			LastActiveDate: now,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
		if err != nil {
			return err
		}

		// Map to our domain model
		chat = models.Chat{
			ID:             result.ID,
			Title:          result.Title,
			UserEmail:      result.UserEmail,
			LastActiveDate: result.LastActiveDate,
			CreatedAt:      result.CreatedAt,
			UpdatedAt:      result.UpdatedAt,
		}

		// Create initial message
		msgResult, err := q.CreateMessage(context.Background(), db.CreateMessageParams{
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

		// Map to our domain model
		message = models.Message{
			ID:         msgResult.ID,
			Content:    msgResult.Content,
			SenderType: models.SenderType(msgResult.SenderType),
			ChatID:     msgResult.ChatID,
			CreatedAt:  msgResult.CreatedAt,
			UpdatedAt:  msgResult.UpdatedAt,
		}

		return nil
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.NewServerError("Failed to create chat"))
	}

	// Return response
	response := models.ChatWithMessageResponse{
		ChatResponse:   chat.ToChatResponse(),
		InitialMessage: message.ToMessageResponse(),
	}

	return c.Status(http.StatusOK).JSON(response)
}

// GetChats handles retrieving all chats for a user
func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
	// Get current user
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(models.NewUnauthorizedError(""))
	}

	// Query chats
	results, err := h.store.GetChatsByUserEmail(context.Background(), user.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.NewServerError("Failed to fetch chats"))
	}

	// Convert to response format
	var chats []models.ChatResponse
	for _, chat := range results {
		chats = append(chats, models.ChatResponse{
			ID:             chat.ID,
			Title:          chat.Title,
			LastActiveDate: chat.LastActiveDate,
		})
	}

	return c.Status(http.StatusOK).JSON(chats)
}