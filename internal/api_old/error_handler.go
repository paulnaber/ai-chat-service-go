package api_old

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"

	"ai-chat-service-go/internal/models"
)

// ErrorHandler is a custom Fiber error handler
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Get the status code from the error, default to 500
	code := fiber.StatusInternalServerError

	// Check for custom errors
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Handle specific errors with appropriate status codes
	var response models.ErrorResponse
	switch code {
	case fiber.StatusBadRequest:
		response = models.NewValidationError(err.Error())
	case fiber.StatusUnauthorized:
		response = models.NewUnauthorizedError(err.Error())
	case fiber.StatusForbidden:
		response = models.NewForbiddenError(err.Error())
	case fiber.StatusNotFound:
		response = models.NewResourceNotFoundError(err.Error())
	default:
		// Log unexpected errors
		log.Printf("Unexpected error: %v", err)
		response = models.NewServerError("An unexpected error occurred")
	}

	// Set the content type and return the error
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(code).JSON(response)
}
