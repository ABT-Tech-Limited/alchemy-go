package alchemy

import (
	"net/http"
	"time"
)

// Config holds the configuration for the Alchemy client.
type Config struct {
	// APIKey is the Alchemy API key (required).
	APIKey string

	// Network is the blockchain network to connect to (default: EthMainnet).
	Network Network

	// BaseURL overrides the default API endpoint.
	// If empty, the endpoint is derived from Network.
	BaseURL string

	// Timeout is the request timeout (default: 30s).
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts (default: 3).
	MaxRetries int

	// RetryDelay is the initial delay between retries (default: 1s).
	RetryDelay time.Duration

	// RetryMaxDelay is the maximum delay between retries (default: 30s).
	RetryMaxDelay time.Duration

	// HTTPClient is a custom HTTP client to use.
	// If nil, a default client is created.
	HTTPClient *http.Client

	// Debug enables debug logging.
	Debug bool
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() Config {
	return Config{
		Network:       EthMainnet,
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryDelay:    1 * time.Second,
		RetryMaxDelay: 30 * time.Second,
	}
}

// Validate validates the configuration and returns an error if invalid.
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}
	return nil
}

// WithDefaults returns a copy of the config with default values applied
// for any zero-valued fields.
func (c Config) WithDefaults() Config {
	defaults := DefaultConfig()

	if c.Network == "" {
		c.Network = defaults.Network
	}
	if c.Timeout == 0 {
		c.Timeout = defaults.Timeout
	}
	if c.MaxRetries == 0 {
		c.MaxRetries = defaults.MaxRetries
	}
	if c.RetryDelay == 0 {
		c.RetryDelay = defaults.RetryDelay
	}
	if c.RetryMaxDelay == 0 {
		c.RetryMaxDelay = defaults.RetryMaxDelay
	}

	return c
}

// GetBaseURL returns the base URL for API requests.
func (c *Config) GetBaseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return c.Network.BaseURL()
}

// GetHTTPClient returns the HTTP client to use.
func (c *Config) GetHTTPClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{
		Timeout: c.Timeout,
	}
}

// Common configuration errors.
var (
	ErrMissingAPIKey = &ConfigError{Message: "API key is required"}
)

// ConfigError represents a configuration error.
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Message
}
