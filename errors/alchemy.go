package errors

import (
	"fmt"
)

// Alchemy API error types.
const (
	ErrTypeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	ErrTypeInvalidAPIKey     = "INVALID_API_KEY"
	ErrTypeInvalidParams     = "INVALID_PARAMS"
	ErrTypeResourceNotFound  = "RESOURCE_NOT_FOUND"
	ErrTypeInternalError     = "INTERNAL_ERROR"
	ErrTypeUnsupportedMethod = "UNSUPPORTED_METHOD"
	ErrTypeUnknown           = "UNKNOWN_ERROR"
)

// APIError represents an Alchemy API-specific error.
type APIError struct {
	// Type is the error type (e.g., RATE_LIMIT_EXCEEDED).
	Type string `json:"type"`
	// Message is the error message.
	Message string `json:"message"`
	// Details contains additional error details.
	Details string `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("Alchemy API error [%s]: %s (details: %s)", e.Type, e.Message, e.Details)
	}
	return fmt.Sprintf("Alchemy API error [%s]: %s", e.Type, e.Message)
}

// Code returns the error type as the code.
func (e *APIError) Code() string {
	return e.Type
}

// Unwrap returns nil (APIError doesn't wrap other errors).
func (e *APIError) Unwrap() error {
	return nil
}

// IsRateLimitExceeded returns true if this is a rate limit error.
func (e *APIError) IsRateLimitExceeded() bool {
	return e.Type == ErrTypeRateLimitExceeded
}

// IsInvalidAPIKey returns true if this is an invalid API key error.
func (e *APIError) IsInvalidAPIKey() bool {
	return e.Type == ErrTypeInvalidAPIKey
}

// IsInvalidParams returns true if this is an invalid params error.
func (e *APIError) IsInvalidParams() bool {
	return e.Type == ErrTypeInvalidParams
}

// IsResourceNotFound returns true if this is a resource not found error.
func (e *APIError) IsResourceNotFound() bool {
	return e.Type == ErrTypeResourceNotFound
}

// IsRetryable returns true if the error might be retryable.
func (e *APIError) IsRetryable() bool {
	switch e.Type {
	case ErrTypeRateLimitExceeded, ErrTypeInternalError:
		return true
	default:
		return false
	}
}

// NewAPIError creates a new APIError.
func NewAPIError(errType, message string) *APIError {
	return &APIError{
		Type:    errType,
		Message: message,
	}
}

// NewAPIErrorWithDetails creates a new APIError with details.
func NewAPIErrorWithDetails(errType, message, details string) *APIError {
	return &APIError{
		Type:    errType,
		Message: message,
		Details: details,
	}
}
