package client

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Handler is a function that processes an HTTP request.
type Handler func(ctx context.Context, req *http.Request) (*http.Response, error)

// Middleware is an interface for HTTP middleware.
type Middleware interface {
	// Wrap wraps a handler with middleware logic.
	Wrap(next Handler) Handler
}

// MiddlewareFunc is a function that implements Middleware.
type MiddlewareFunc func(next Handler) Handler

// Wrap implements Middleware.
func (f MiddlewareFunc) Wrap(next Handler) Handler {
	return f(next)
}

// LoggingMiddleware logs HTTP requests and responses.
type LoggingMiddleware struct {
	Logger *slog.Logger
}

// NewLoggingMiddleware creates a new LoggingMiddleware.
func NewLoggingMiddleware(logger *slog.Logger) *LoggingMiddleware {
	if logger == nil {
		logger = slog.Default()
	}
	return &LoggingMiddleware{Logger: logger}
}

// Wrap implements Middleware.
func (m *LoggingMiddleware) Wrap(next Handler) Handler {
	return func(ctx context.Context, req *http.Request) (*http.Response, error) {
		start := time.Now()

		m.Logger.Debug("HTTP request",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
		)

		resp, err := next(ctx, req)
		duration := time.Since(start)

		if err != nil {
			m.Logger.Error("HTTP request failed",
				slog.String("method", req.Method),
				slog.String("url", req.URL.String()),
				slog.Duration("duration", duration),
				slog.String("error", err.Error()),
			)
			return nil, err
		}

		m.Logger.Debug("HTTP response",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
			slog.Int("status", resp.StatusCode),
			slog.Duration("duration", duration),
		)

		return resp, nil
	}
}

// HeaderMiddleware adds custom headers to requests.
type HeaderMiddleware struct {
	Headers map[string]string
}

// NewHeaderMiddleware creates a new HeaderMiddleware.
func NewHeaderMiddleware(headers map[string]string) *HeaderMiddleware {
	return &HeaderMiddleware{Headers: headers}
}

// Wrap implements Middleware.
func (m *HeaderMiddleware) Wrap(next Handler) Handler {
	return func(ctx context.Context, req *http.Request) (*http.Response, error) {
		for k, v := range m.Headers {
			req.Header.Set(k, v)
		}
		return next(ctx, req)
	}
}

// UserAgentMiddleware sets the User-Agent header.
type UserAgentMiddleware struct {
	UserAgent string
}

// NewUserAgentMiddleware creates a new UserAgentMiddleware.
func NewUserAgentMiddleware(userAgent string) *UserAgentMiddleware {
	return &UserAgentMiddleware{UserAgent: userAgent}
}

// Wrap implements Middleware.
func (m *UserAgentMiddleware) Wrap(next Handler) Handler {
	return func(ctx context.Context, req *http.Request) (*http.Response, error) {
		req.Header.Set("User-Agent", m.UserAgent)
		return next(ctx, req)
	}
}

// MetricsMiddleware collects request metrics.
type MetricsMiddleware struct {
	// OnRequest is called before a request is made.
	OnRequest func(method, url string)
	// OnResponse is called after a response is received.
	OnResponse func(method, url string, statusCode int, duration time.Duration, err error)
}

// NewMetricsMiddleware creates a new MetricsMiddleware.
func NewMetricsMiddleware(onRequest func(method, url string), onResponse func(method, url string, statusCode int, duration time.Duration, err error)) *MetricsMiddleware {
	return &MetricsMiddleware{
		OnRequest:  onRequest,
		OnResponse: onResponse,
	}
}

// Wrap implements Middleware.
func (m *MetricsMiddleware) Wrap(next Handler) Handler {
	return func(ctx context.Context, req *http.Request) (*http.Response, error) {
		if m.OnRequest != nil {
			m.OnRequest(req.Method, req.URL.String())
		}

		start := time.Now()
		resp, err := next(ctx, req)
		duration := time.Since(start)

		if m.OnResponse != nil {
			statusCode := 0
			if resp != nil {
				statusCode = resp.StatusCode
			}
			m.OnResponse(req.Method, req.URL.String(), statusCode, duration, err)
		}

		return resp, err
	}
}

// Chain chains multiple middlewares together.
func Chain(middlewares ...Middleware) Middleware {
	return MiddlewareFunc(func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i].Wrap(next)
		}
		return next
	})
}
