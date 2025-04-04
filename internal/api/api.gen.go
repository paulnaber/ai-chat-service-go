// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime"
)

// Defines values for ErrorMessageCode.
const (
	FORBIDDEN        ErrorMessageCode = "FORBIDDEN"
	RESOURCENOTFOUND ErrorMessageCode = "RESOURCE_NOT_FOUND"
	SERVERERROR      ErrorMessageCode = "SERVER_ERROR"
	UNAUTHORIZED     ErrorMessageCode = "UNAUTHORIZED"
	VALIDATIONERROR  ErrorMessageCode = "VALIDATION_ERROR"
)

// Defines values for SenderType.
const (
	BACKEND SenderType = "BACKEND"
	LLM     SenderType = "LLM"
	USER    SenderType = "USER"
)

// ChatDTO defines model for ChatDTO.
type ChatDTO struct {
	// Id Unique identifier for the chat (auto-generated)
	Id *int `json:"id,omitempty"`

	// LastActiveDate Date when the chat was last active
	LastActiveDate *LocalDateTime `json:"lastActiveDate,omitempty"`

	// Title Name of the chat (derived from first message)
	Title *string `json:"title,omitempty"`
}

// ErrorMessage defines model for ErrorMessage.
type ErrorMessage struct {
	// Code Error code that identifies the error type
	Code ErrorMessageCode `json:"code"`

	// Details List of field-value pairs with additional error information
	Details *[]struct {
		// Field The name of the field with an error
		Field *string `json:"field,omitempty"`

		// Value The error message or problematic value for this field
		Value *string `json:"value,omitempty"`
	} `json:"details,omitempty"`

	// Message Human-readable error description
	Message string `json:"message"`
}

// ErrorMessageCode Error code that identifies the error type
type ErrorMessageCode string

// LocalDateTime defines model for LocalDateTime.
type LocalDateTime = time.Time

// MessageDTO defines model for MessageDTO.
type MessageDTO struct {
	// ChatId Reference to the chat this message belongs to (auto-generated)
	ChatId *int `json:"chatId,omitempty"`

	// Content Content of the message
	Content *string `json:"content,omitempty"`

	// CreatedAt Date when the message was created (auto-generated)
	CreatedAt *LocalDateTime `json:"createdAt,omitempty"`

	// Id Unique identifier for the message (auto-generated)
	Id *int `json:"id,omitempty"`

	// SenderType Type of sender (automatically set to 'user' for user messages)
	SenderType *SenderType `json:"senderType,omitempty"`
}

// SenderType Type of sender (automatically set to 'user' for user messages)
type SenderType string

// CreateChatJSONBody defines parameters for CreateChat.
type CreateChatJSONBody struct {
	// Content Content of the first message to start the chat with
	Content string `json:"content"`
}

// CreateMessageJSONBody defines parameters for CreateMessage.
type CreateMessageJSONBody struct {
	// Content Content of the message
	Content string `json:"content"`
}

// CreateChatJSONRequestBody defines body for CreateChat for application/json ContentType.
type CreateChatJSONRequestBody CreateChatJSONBody

// CreateMessageJSONRequestBody defines body for CreateMessage for application/json ContentType.
type CreateMessageJSONRequestBody CreateMessageJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all chats for a user
	// (GET /v1/chats)
	GetChats(c *fiber.Ctx) error
	// Create a new chat
	// (POST /v1/chats)
	CreateChat(c *fiber.Ctx) error
	// Get all messages for a chat
	// (GET /v1/chats/{chatId}/messages)
	GetMessages(c *fiber.Ctx, chatId int) error
	// Create a new message
	// (POST /v1/chats/{chatId}/messages)
	CreateMessage(c *fiber.Ctx, chatId int) error
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

type MiddlewareFunc fiber.Handler

// GetChats operation middleware
func (siw *ServerInterfaceWrapper) GetChats(c *fiber.Ctx) error {

	return siw.Handler.GetChats(c)
}

// CreateChat operation middleware
func (siw *ServerInterfaceWrapper) CreateChat(c *fiber.Ctx) error {

	return siw.Handler.CreateChat(c)
}

// GetMessages operation middleware
func (siw *ServerInterfaceWrapper) GetMessages(c *fiber.Ctx) error {

	var err error

	// ------------- Path parameter "chatId" -------------
	var chatId int

	err = runtime.BindStyledParameterWithOptions("simple", "chatId", c.Params("chatId"), &chatId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter chatId: %w", err).Error())
	}

	return siw.Handler.GetMessages(c, chatId)
}

// CreateMessage operation middleware
func (siw *ServerInterfaceWrapper) CreateMessage(c *fiber.Ctx) error {

	var err error

	// ------------- Path parameter "chatId" -------------
	var chatId int

	err = runtime.BindStyledParameterWithOptions("simple", "chatId", c.Params("chatId"), &chatId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter chatId: %w", err).Error())
	}

	return siw.Handler.CreateMessage(c, chatId)
}

// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	for _, m := range options.Middlewares {
		router.Use(fiber.Handler(m))
	}

	router.Get(options.BaseURL+"/v1/chats", wrapper.GetChats)

	router.Post(options.BaseURL+"/v1/chats", wrapper.CreateChat)

	router.Get(options.BaseURL+"/v1/chats/:chatId/messages", wrapper.GetMessages)

	router.Post(options.BaseURL+"/v1/chats/:chatId/messages", wrapper.CreateMessage)

}
