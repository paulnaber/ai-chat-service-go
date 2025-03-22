package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateChatRequest struct {
	Content string `json:"content" validate:"required"`
}

type ChatResponse struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	LastActiveDate time.Time `json:"lastActiveDate"`
}

type ChatWithMessageResponse struct {
	ChatResponse
	InitialMessage MessageResponse `json:"initialMessage"`
	AIResponse     MessageResponse `json:"aiResponse"`
}

type CreateMessageRequest struct {
	Content string `json:"content" validate:"required"`
}

type CreateMessageResponse struct {
	UserMessage MessageResponse `json:"userMessage"`
	AIResponse  MessageResponse `json:"aiResponse"`
}

type MessageResponse struct {
	ID         uuid.UUID  `json:"id"`
	Content    string     `json:"content"`
	SenderType SenderType `json:"senderType"`
	CreatedAt  time.Time  `json:"createdAt"`
	ChatID     uuid.UUID  `json:"chatId"`
}

type SenderType string

const (
	SenderTypeUser    SenderType = "user"
	SenderTypeBackend SenderType = "backend"
	SenderTypeLLM     SenderType = "llm"
)
