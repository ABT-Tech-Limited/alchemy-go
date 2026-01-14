package client

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/ABT-Tech-Limited/alchemy-go/errors"
)

// Retrier implements exponential backoff retry logic.
type Retrier struct {
	// MaxRetries is the maximum number of retry attempts.
	MaxRetries int
	// InitialDelay is the initial delay between retries.
	InitialDelay time.Duration
	// MaxDelay is the maximum delay between retries.
	MaxDelay time.Duration
	// Multiplier is the factor by which the delay increases.
	Multiplier float64
	// Jitter adds randomness to the delay (0.0 to 1.0).
	Jitter float64
}

// DefaultRetrier returns a Retrier with default settings.
func DefaultRetrier() *Retrier {
	return &Retrier{
		MaxRetries:   3,
		InitialDelay: time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		Jitter:       0.1,
	}
}

// Do executes the given function with retries.
func (r *Retrier) Do(ctx context.Context, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= r.MaxRetries; attempt++ {
		// Check context before attempting
		if err := ctx.Err(); err != nil {
			if err == context.Canceled {
				return errors.ErrContextCanceled
			}
			return errors.ErrContextDeadline
		}

		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should stop retrying
		if _, ok := err.(*stopRetry); ok {
			return err
		}

		// Check if the error is retryable
		if !r.ShouldRetry(err) {
			return err
		}

		// Don't sleep after the last attempt
		if attempt == r.MaxRetries {
			break
		}

		// Calculate delay with exponential backoff
		delay := r.calculateDelay(attempt)

		// Wait with context
		select {
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				return errors.ErrContextCanceled
			}
			return errors.ErrContextDeadline
		case <-time.After(delay):
		}
	}

	return lastErr
}

// calculateDelay calculates the delay for a given attempt with exponential backoff and jitter.
func (r *Retrier) calculateDelay(attempt int) time.Duration {
	// Calculate base delay with exponential backoff
	delay := float64(r.InitialDelay) * math.Pow(r.Multiplier, float64(attempt))

	// Apply maximum delay cap
	if delay > float64(r.MaxDelay) {
		delay = float64(r.MaxDelay)
	}

	// Apply jitter
	if r.Jitter > 0 {
		jitter := delay * r.Jitter * (rand.Float64()*2 - 1)
		delay = delay + jitter
	}

	// Ensure non-negative delay
	if delay < 0 {
		delay = 0
	}

	return time.Duration(delay)
}

// ShouldRetry determines if an error is retryable.
func (r *Retrier) ShouldRetry(err error) bool {
	return errors.IsRetryable(err)
}

// RetryableFunc is a function that can be retried.
type RetryableFunc func() error

// WithRetry wraps a function to be executed with retries.
func WithRetry(ctx context.Context, retrier *Retrier, fn RetryableFunc) error {
	if retrier == nil {
		retrier = DefaultRetrier()
	}
	return retrier.Do(ctx, fn)
}
