// Package errors provides error types for the Alchemy SDK.
package errors

import (
	"errors"
	"fmt"
)

// Common sentinel errors.
var (
	ErrNilResponse      = errors.New("nil response")
	ErrInvalidResponse  = errors.New("invalid response")
	ErrContextCanceled  = errors.New("context canceled")
	ErrContextDeadline  = errors.New("context deadline exceeded")
	ErrInvalidAPIKey    = errors.New("invalid API key")
	ErrRateLimited      = errors.New("rate limited")
	ErrNetworkNotFound  = errors.New("network not found")
	ErrInvalidAddress   = errors.New("invalid address")
	ErrInvalidHash      = errors.New("invalid hash")
	ErrInvalidParameter = errors.New("invalid parameter")
)

// Error is the interface for all SDK errors.
type Error interface {
	error
	// Code returns the error code.
	Code() string
	// Unwrap returns the underlying error.
	Unwrap() error
}

// AlchemyError is the base error type for the SDK.
type AlchemyError struct {
	code    string
	message string
	cause   error
}

// Error implements the error interface.
func (e *AlchemyError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.code, e.message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

// Code returns the error code.
func (e *AlchemyError) Code() string {
	return e.code
}

// Unwrap returns the underlying error.
func (e *AlchemyError) Unwrap() error {
	return e.cause
}

// New creates a new AlchemyError.
func New(code, message string) *AlchemyError {
	return &AlchemyError{
		code:    code,
		message: message,
	}
}

// Wrap wraps an error with additional context.
func Wrap(err error, code, message string) *AlchemyError {
	return &AlchemyError{
		code:    code,
		message: message,
		cause:   err,
	}
}

// Wrapf wraps an error with formatted message.
func Wrapf(err error, code, format string, args ...interface{}) *AlchemyError {
	return &AlchemyError{
		code:    code,
		message: fmt.Sprintf(format, args...),
		cause:   err,
	}
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// IsRetryable returns true if the error is retryable.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for HTTP errors
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.IsRetryable()
	}

	// Check for rate limit errors
	if errors.Is(err, ErrRateLimited) {
		return true
	}

	// Check for context errors (not retryable)
	if errors.Is(err, ErrContextCanceled) || errors.Is(err, ErrContextDeadline) {
		return false
	}

	return false
}

// IsAuthError returns true if the error is an authentication error.
func IsAuthError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrInvalidAPIKey) {
		return true
	}

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == 401 || httpErr.StatusCode == 403
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Type == ErrTypeInvalidAPIKey
	}

	return false
}

// IsRateLimitError returns true if the error is a rate limit error.
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrRateLimited) {
		return true
	}

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == 429
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Type == ErrTypeRateLimitExceeded
	}

	return false
}
