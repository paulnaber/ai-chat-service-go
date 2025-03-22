package models

import (
	"time"
)

// Chat represents a user's chat session
// swagger:model Chat
type Chat struct {
	// The unique identifier for the chat
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id" db:"id"`

	// The title of the chat, derived from the first message
	// example: How do I configure my device?
	Title string `json:"title" db:"title"`

	// The email of the user who owns the chat
	// example: user@example.com
	UserEmail string `json:"userEmail" db:"user_email"`

	// When the chat was last active
	// example: 2023-07-15T14:32:21Z
	LastActiveDate time.Time `json:"lastActiveDate" db:"last_active_date"`

	// When the chat was created
	// example: 2023-07-15T14:32:21Z
	CreatedAt time.Time `json:"createdAt" db:"created_at"`

	// When the chat was last updated
	// example: 2023-07-15T14:35:42Z
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// ChatResponse is the data transfer object for chats in API responses
// swagger:model ChatResponse
type ChatResponse struct {
	// The unique identifier for the chat
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id"`

	// The title of the chat, derived from the first message
	// example: How do I configure my device?
	Title string `json:"title"`

	// When the chat was last active
	// example: 2023-07-15T14:32:21Z
	LastActiveDate time.Time `json:"lastActiveDate"`
}

// ToChatResponse converts a Chat model to a ChatResponse
func (c Chat) ToChatResponse() ChatResponse {
	return ChatResponse{
		ID:             c.ID,
		Title:          c.Title,
		LastActiveDate: c.LastActiveDate,
	}
}

// CreateChatRequest is the request body for creating a chat
// swagger:model CreateChatRequest
type CreateChatRequest struct {
	// The content of the first message to start the chat with
	// required: true
	// example: How do I configure my device?
	Content string `json:"content" validate:"required"`
}

// ChatWithMessageResponse is the response for creating a chat
// swagger:model ChatWithMessageResponse
type ChatWithMessageResponse struct {
	// The chat details
	ChatResponse
	
	// The initial message in the chat
	// required: true
	InitialMessage MessageResponse `json:"initialMessage"`
}