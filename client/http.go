// Package client provides HTTP and JSON-RPC client implementations for the Alchemy SDK.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/ABT-Tech-Limited/alchemy-go/errors"
)

// HTTPClient is the HTTP client for making API requests.
type HTTPClient struct {
	baseURL     string
	apiKey      string
	httpClient  *http.Client
	middlewares []Middleware
	retrier     *Retrier
	debug       bool
}

// HTTPClientConfig holds configuration for HTTPClient.
type HTTPClientConfig struct {
	BaseURL       string
	APIKey        string
	Timeout       time.Duration
	MaxRetries    int
	RetryDelay    time.Duration
	RetryMaxDelay time.Duration
	HTTPClient    *http.Client
	Middlewares   []Middleware
	Debug         bool
}

// NewHTTPClient creates a new HTTPClient.
func NewHTTPClient(cfg HTTPClientConfig) *HTTPClient {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
		}
	}

	retrier := &Retrier{
		MaxRetries:   cfg.MaxRetries,
		InitialDelay: cfg.RetryDelay,
		MaxDelay:     cfg.RetryMaxDelay,
		Multiplier:   2.0,
	}

	return &HTTPClient{
		baseURL:     cfg.BaseURL,
		apiKey:      cfg.APIKey,
		httpClient:  httpClient,
		middlewares: cfg.Middlewares,
		retrier:     retrier,
		debug:       cfg.Debug,
	}
}

// BaseURL returns the base URL.
func (c *HTTPClient) BaseURL() string {
	return c.baseURL
}

// Do executes an HTTP request with retry and middleware support.
func (c *HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Build the handler chain with middlewares
	handler := c.doRequest
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i].Wrap(handler)
	}

	var resp *http.Response
	var lastErr error

	err := c.retrier.Do(ctx, func() error {
		var err error
		resp, err = handler(ctx, req)
		if err != nil {
			lastErr = err
			// Check if error is retryable
			if errors.IsRetryable(err) {
				return err
			}
			// Non-retryable error, stop retrying
			return &stopRetry{err: err}
		}

		// Check for retryable HTTP status codes
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = errors.NewHTTPError(resp.StatusCode, resp.Status, body)
			return lastErr
		}

		return nil
	})

	if err != nil {
		// Check if it's a stopRetry wrapper
		if sr, ok := err.(*stopRetry); ok {
			return nil, sr.err
		}
		if lastErr != nil {
			return nil, lastErr
		}
		return nil, err
	}

	return resp, nil
}

// doRequest executes a single HTTP request.
func (c *HTTPClient) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return c.httpClient.Do(req)
}

// Post makes a POST request with JSON body.
func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	url := c.baseURL + "/" + c.apiKey
	if path != "" {
		url = url + "/" + path
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "MARSHAL_ERROR", "failed to marshal request body")
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "REQUEST_ERROR", "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "READ_ERROR", "failed to read response body")
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.NewHTTPError(resp.StatusCode, resp.Status, respBody)
	}

	return respBody, nil
}

// Get makes a GET request.
func (c *HTTPClient) Get(ctx context.Context, path string) ([]byte, error) {
	url := c.baseURL + "/" + c.apiKey
	if path != "" {
		url = url + "/" + path
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "REQUEST_ERROR", "failed to create request")
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "READ_ERROR", "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.NewHTTPError(resp.StatusCode, resp.Status, respBody)
	}

	return respBody, nil
}

// GetWithQuery makes a GET request with query parameters.
func (c *HTTPClient) GetWithQuery(ctx context.Context, path string, query map[string]string) ([]byte, error) {
	url := c.baseURL + "/" + c.apiKey
	if path != "" {
		url = url + "/" + path
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "REQUEST_ERROR", "failed to create request")
	}

	if len(query) > 0 {
		q := req.URL.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "READ_ERROR", "failed to read response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.NewHTTPError(resp.StatusCode, resp.Status, respBody)
	}

	return respBody, nil
}

// stopRetry is used to signal that retrying should stop.
type stopRetry struct {
	err error
}

func (s *stopRetry) Error() string {
	return s.err.Error()
}

// RequestID is used to generate unique request IDs.
var requestIDCounter uint64

// NextRequestID returns the next request ID.
func NextRequestID() uint64 {
	return atomic.AddUint64(&requestIDCounter, 1)
}

// JSONRPCRequest represents a JSON-RPC request.
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
	ID      uint64        `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC response.
type JSONRPCResponse struct {
	JSONRPC string               `json:"jsonrpc"`
	ID      uint64               `json:"id"`
	Result  json.RawMessage      `json:"result,omitempty"`
	Error   *errors.JSONRPCError `json:"error,omitempty"`
}

// JSONRPCClient is a client for making JSON-RPC calls.
type JSONRPCClient struct {
	httpClient *HTTPClient
}

// NewJSONRPCClient creates a new JSONRPCClient.
func NewJSONRPCClient(httpClient *HTTPClient) *JSONRPCClient {
	return &JSONRPCClient{
		httpClient: httpClient,
	}
}

// Call makes a JSON-RPC call and unmarshals the result.
func (c *JSONRPCClient) Call(ctx context.Context, method string, params []interface{}, result interface{}) error {
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      NextRequestID(),
	}

	respBody, err := c.httpClient.Post(ctx, "", req)
	if err != nil {
		return err
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return errors.Wrap(err, "UNMARSHAL_ERROR", "failed to unmarshal JSON-RPC response")
	}

	if resp.Error != nil {
		return resp.Error
	}

	if result != nil && len(resp.Result) > 0 {
		if err := json.Unmarshal(resp.Result, result); err != nil {
			return errors.Wrap(err, "UNMARSHAL_ERROR", "failed to unmarshal result")
		}
	}

	return nil
}

// CallRaw makes a JSON-RPC call and returns the raw result.
func (c *JSONRPCClient) CallRaw(ctx context.Context, method string, params []interface{}) (json.RawMessage, error) {
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      NextRequestID(),
	}

	respBody, err := c.httpClient.Post(ctx, "", req)
	if err != nil {
		return nil, err
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, errors.Wrap(err, "UNMARSHAL_ERROR", "failed to unmarshal JSON-RPC response")
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Result, nil
}

// BatchCall represents a single call in a batch request.
type BatchCall struct {
	Method string
	Params []interface{}
	Result interface{} // pointer to result type
}

// BatchResult represents the result of a single call in a batch request.
type BatchResult struct {
	Error  error
	Result interface{}
}

// BatchCallResponse holds the response for a batch call item.
type BatchCallResponse struct {
	JSONRPC string               `json:"jsonrpc"`
	ID      uint64               `json:"id"`
	Result  json.RawMessage      `json:"result,omitempty"`
	Error   *errors.JSONRPCError `json:"error,omitempty"`
}

// BatchCall makes multiple JSON-RPC calls in a single HTTP request.
func (c *JSONRPCClient) BatchCall(ctx context.Context, calls []BatchCall) ([]BatchResult, error) {
	if len(calls) == 0 {
		return nil, nil
	}

	// Build batch request
	requests := make([]JSONRPCRequest, len(calls))
	for i, call := range calls {
		requests[i] = JSONRPCRequest{
			JSONRPC: "2.0",
			Method:  call.Method,
			Params:  call.Params,
			ID:      uint64(i + 1),
		}
	}

	respBody, err := c.httpClient.Post(ctx, "", requests)
	if err != nil {
		return nil, err
	}

	// Parse batch response
	var responses []BatchCallResponse
	if err := json.Unmarshal(respBody, &responses); err != nil {
		return nil, errors.Wrap(err, "UNMARSHAL_ERROR", "failed to unmarshal batch response")
	}

	// Create a map of responses by ID for easier lookup
	responseMap := make(map[uint64]*BatchCallResponse)
	for i := range responses {
		responseMap[responses[i].ID] = &responses[i]
	}

	// Process results in order
	results := make([]BatchResult, len(calls))
	for i, call := range calls {
		resp, ok := responseMap[uint64(i+1)]
		if !ok {
			results[i] = BatchResult{
				Error: fmt.Errorf("missing response for call %d", i),
			}
			continue
		}

		if resp.Error != nil {
			results[i] = BatchResult{
				Error: resp.Error,
			}
			continue
		}

		if call.Result != nil && len(resp.Result) > 0 {
			if err := json.Unmarshal(resp.Result, call.Result); err != nil {
				results[i] = BatchResult{
					Error: errors.Wrap(err, "UNMARSHAL_ERROR", "failed to unmarshal result"),
				}
				continue
			}
			results[i] = BatchResult{
				Result: call.Result,
			}
		}
	}

	return results, nil
}
