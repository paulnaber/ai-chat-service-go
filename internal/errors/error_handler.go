package errors

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
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
	var response ErrorResponse
	switch code {
	case fiber.StatusBadRequest:
		response = NewValidationError(err.Error())
	case fiber.StatusUnauthorized:
		response = NewUnauthorizedError(err.Error())
	case fiber.StatusForbidden:
		response = NewForbiddenError(err.Error())
	case fiber.StatusNotFound:
		response = NewResourceNotFoundError(err.Error())
	default:
		// Log unexpected errors
		log.Printf("Unexpected error: %v", err)
		response = NewServerError("An unexpected error occurred")
	}

	// Set the content type and return the error
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(code).JSON(response)
}
