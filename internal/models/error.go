package models

// ErrorCode represents error types returned by the API
// swagger:enum ErrorCode
type ErrorCode string

const (
	// ValidationError indicates that validation failed
	ValidationError ErrorCode = "VALIDATION_ERROR"
	// UnauthorizedError indicates that authentication is required
	UnauthorizedError ErrorCode = "UNAUTHORIZED"
	// ForbiddenError indicates that the user doesn't have permission
	ForbiddenError ErrorCode = "FORBIDDEN"
	// ResourceNotFoundError indicates that the requested resource was not found
	ResourceNotFoundError ErrorCode = "RESOURCE_NOT_FOUND"
	// ServerError indicates an unexpected server error
	ServerError ErrorCode = "SERVER_ERROR"
)

// ErrorDetail represents details about validation errors
// swagger:model ErrorDetail
type ErrorDetail struct {
	// The name of the field with an error
	// example: content
	Field string `json:"field"`
	
	// The error message or problematic value for this field
	// example: Content cannot be empty
	Value string `json:"value"`
}

// ErrorResponse represents a standardized error response
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error code that identifies the error type
	// example: VALIDATION_ERROR
	Code ErrorCode `json:"code"`
	
	// Human-readable error description
	// example: The request contains invalid parameters
	Message string `json:"message"`
	
	// List of field-value pairs with additional error information
	Details []ErrorDetail `json:"details,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code ErrorCode, message string, details ...ErrorDetail) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewValidationError creates a validation error response
func NewValidationError(message string, details ...ErrorDetail) ErrorResponse {
	return NewErrorResponse(ValidationError, message, details...)
}

// NewUnauthorizedError creates an unauthorized error response
func NewUnauthorizedError(message string) ErrorResponse {
	if message == "" {
		message = "Authentication required"
	}
	return NewErrorResponse(UnauthorizedError, message)
}

// NewForbiddenError creates a forbidden error response
func NewForbiddenError(message string) ErrorResponse {
	if message == "" {
		message = "You do not have permission to access this resource"
	}
	return NewErrorResponse(ForbiddenError, message)
}

// NewResourceNotFoundError creates a not found error response
func NewResourceNotFoundError(message string, details ...ErrorDetail) ErrorResponse {
	if message == "" {
		message = "The requested resource could not be found"
	}
	return NewErrorResponse(ResourceNotFoundError, message, details...)
}

// NewServerError creates a server error response
func NewServerError(message string) ErrorResponse {
	if message == "" {
		message = "An unexpected error occurred"
	}
	return NewErrorResponse(ServerError, message)
}