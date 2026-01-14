package errors

import (
	"encoding/json"
	"fmt"
)

// Standard JSON-RPC error codes.
const (
	// ParseError indicates invalid JSON was received.
	ParseError = -32700
	// InvalidRequest indicates the JSON is not a valid Request object.
	InvalidRequest = -32600
	// MethodNotFound indicates the method does not exist or is not available.
	MethodNotFound = -32601
	// InvalidParams indicates invalid method parameters.
	InvalidParams = -32602
	// InternalError indicates an internal JSON-RPC error.
	InternalError = -32603
	// ServerError indicates a server error (reserved range -32000 to -32099).
	ServerErrorStart = -32099
	ServerErrorEnd   = -32000
)

// JSONRPCError represents a JSON-RPC error response.
type JSONRPCError struct {
	// Code is the error code.
	Code int `json:"code"`
	// Message is the error message.
	Message string `json:"message"`
	// Data is optional additional data.
	Data json.RawMessage `json:"data,omitempty"`
}

// Error implements the error interface.
func (e *JSONRPCError) Error() string {
	if len(e.Data) > 0 {
		return fmt.Sprintf("JSON-RPC error %d: %s (data: %s)", e.Code, e.Message, string(e.Data))
	}
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}

// ErrorCode returns the error code as a string.
func (e *JSONRPCError) ErrorCode() string {
	return fmt.Sprintf("JSONRPC_%d", e.Code)
}

// Unwrap returns nil (JSONRPCError doesn't wrap other errors).
func (e *JSONRPCError) Unwrap() error {
	return nil
}

// IsParseError returns true if this is a parse error.
func (e *JSONRPCError) IsParseError() bool {
	return e.Code == ParseError
}

// IsInvalidRequest returns true if this is an invalid request error.
func (e *JSONRPCError) IsInvalidRequest() bool {
	return e.Code == InvalidRequest
}

// IsMethodNotFound returns true if this is a method not found error.
func (e *JSONRPCError) IsMethodNotFound() bool {
	return e.Code == MethodNotFound
}

// IsInvalidParams returns true if this is an invalid params error.
func (e *JSONRPCError) IsInvalidParams() bool {
	return e.Code == InvalidParams
}

// IsInternalError returns true if this is an internal error.
func (e *JSONRPCError) IsInternalError() bool {
	return e.Code == InternalError
}

// IsServerError returns true if this is a server error.
func (e *JSONRPCError) IsServerError() bool {
	return e.Code >= ServerErrorStart && e.Code <= ServerErrorEnd
}

// IsRetryable returns true if the error might be retryable.
func (e *JSONRPCError) IsRetryable() bool {
	// Server errors might be transient
	if e.IsServerError() || e.IsInternalError() {
		return true
	}
	return false
}

// DataAs unmarshals the Data field into the provided target.
func (e *JSONRPCError) DataAs(target interface{}) error {
	if len(e.Data) == 0 {
		return nil
	}
	return json.Unmarshal(e.Data, target)
}

// NewJSONRPCError creates a new JSONRPCError.
func NewJSONRPCError(code int, message string, data json.RawMessage) *JSONRPCError {
	return &JSONRPCError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
