package errors

import (
	"fmt"
)

// HTTPError represents an HTTP error response.
type HTTPError struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// Status is the HTTP status text.
	Status string
	// Body is the response body.
	Body []byte
	// URL is the request URL.
	URL string
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	if len(e.Body) > 0 && len(e.Body) < 200 {
		return fmt.Sprintf("HTTP %d %s: %s", e.StatusCode, e.Status, string(e.Body))
	}
	return fmt.Sprintf("HTTP %d %s", e.StatusCode, e.Status)
}

// Code returns the error code.
func (e *HTTPError) Code() string {
	return fmt.Sprintf("HTTP_%d", e.StatusCode)
}

// Unwrap returns nil (HTTPError doesn't wrap other errors).
func (e *HTTPError) Unwrap() error {
	return nil
}

// IsRetryable returns true if the HTTP error is retryable.
func (e *HTTPError) IsRetryable() bool {
	// 429 Too Many Requests (rate limit)
	if e.StatusCode == 429 {
		return true
	}
	// 5xx Server Errors
	if e.StatusCode >= 500 && e.StatusCode < 600 {
		return true
	}
	// 408 Request Timeout
	if e.StatusCode == 408 {
		return true
	}
	return false
}

// IsClientError returns true if this is a 4xx error.
func (e *HTTPError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError returns true if this is a 5xx error.
func (e *HTTPError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// IsRateLimited returns true if this is a 429 error.
func (e *HTTPError) IsRateLimited() bool {
	return e.StatusCode == 429
}

// IsUnauthorized returns true if this is a 401 error.
func (e *HTTPError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsForbidden returns true if this is a 403 error.
func (e *HTTPError) IsForbidden() bool {
	return e.StatusCode == 403
}

// IsNotFound returns true if this is a 404 error.
func (e *HTTPError) IsNotFound() bool {
	return e.StatusCode == 404
}

// NewHTTPError creates a new HTTPError.
func NewHTTPError(statusCode int, status string, body []byte) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Status:     status,
		Body:       body,
	}
}
