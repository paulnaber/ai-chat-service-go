package models

import (
	"time"
)

// SenderType defines the type of message sender
// swagger:enum SenderType
type SenderType string

const (
	// SenderTypeUser represents a message sent by the end user
	SenderTypeUser SenderType = "user"
	// SenderTypeBackend represents a message sent by the backend system
	SenderTypeBackend SenderType = "backend"
	// SenderTypeLLM represents a message sent by the AI/LLM model
	SenderTypeLLM SenderType = "llm"
)

// Message represents a message in a chat
// swagger:model Message
type Message struct {
	// The unique identifier for the message
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id" db:"id"`

	// The content of the message
	// example: How do I configure my device?
	Content string `json:"content" db:"content"`

	// The type of sender (user, backend, or llm)
	// example: user
	SenderType SenderType `json:"senderType" db:"sender_type"`

	// The ID of the chat this message belongs to
	// example: 123e4567-e89b-12d3-a456-426614174000
	ChatID string `json:"chatId" db:"chat_id"`

	// When the message was created
	// example: 2023-07-15T14:32:21Z
	CreatedAt time.Time `json:"createdAt" db:"created_at"`

	// When the message was last updated
	// example: 2023-07-15T14:32:21Z
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// MessageResponse is the data transfer object for messages in API responses
// swagger:model MessageResponse
type MessageResponse struct {
	// The unique identifier for the message
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id"`

	// The content of the message
	// example: How do I configure my device?
	Content string `json:"content"`

	// The type of sender (user, backend, or llm)
	// example: user
	SenderType SenderType `json:"senderType"`

	// When the message was created
	// example: 2023-07-15T14:32:21Z
	CreatedAt time.Time `json:"createdAt"`

	// The ID of the chat this message belongs to
	// example: 123e4567-e89b-12d3-a456-426614174000
	ChatID string `json:"chatId"`
}

// ToMessageResponse converts a Message model to a MessageResponse
func (m Message) ToMessageResponse() MessageResponse {
	return MessageResponse{
		ID:         m.ID,
		Content:    m.Content,
		SenderType: m.SenderType,
		CreatedAt:  m.CreatedAt,
		ChatID:     m.ChatID,
	}
}

// CreateMessageRequest is the request body for creating a message
// swagger:model CreateMessageRequest
type CreateMessageRequest struct {
	// The content of the message
	// required: true
	// example: How do I configure my device?
	Content string `json:"content" validate:"required"`
}