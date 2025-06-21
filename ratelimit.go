package gohtb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RetryPolicy interface {
	ShouldRetry(resp *http.Response, err error) bool
	Wait(retries int) time.Duration
}

type RateLimiter struct {
	mu          sync.Mutex
	limit       RateLimitInfo
	lastRequest time.Time
	ctx         context.Context
	logger      Logger
}

type RateLimitInfo struct {
	Remaining int
	Limit     int
	Reset     time.Time
}

type APITransport struct {
	underlying  http.RoundTripper
	limiter     *RateLimiter
	retryConfig RetryConfig
	logger      Logger
}

func NewRateLimiter(ctx context.Context, logger Logger) *RateLimiter {
	if logger == nil {
		logger = NoopLogger{}
	}
	return &RateLimiter{ctx: ctx, logger: logger, limit: RateLimitInfo{Remaining: 10, Limit: 10}}
}

func NewAPITransport(underlying http.RoundTripper, limiter *RateLimiter, retryConfig RetryConfig, logger Logger) *APITransport {
	if underlying == nil {
		underlying = http.DefaultTransport
	}
	if logger == nil {
		logger = NoopLogger{}
	}
	// Provide a default retry policy if none is set
	if retryConfig.RetryPolicy == nil {
		retryConfig.RetryPolicy = &DefaultRetryPolicy{}
	}
	if retryConfig.MaxRetries <= 0 {
		retryConfig.MaxRetries = 3 // Default to 3 retries
	}

	return &APITransport{
		underlying:  underlying,
		limiter:     limiter,
		retryConfig: retryConfig,
		logger:      logger,
	}
}

func (r *RateLimiter) BeforeRequest() error {
	r.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(r.lastRequest)

	info := r.limit
	var wait time.Duration

	switch {
	case info.Remaining <= 1:
		wait = time.Until(info.Reset)
		if wait <= 0 {
			wait = 7 * time.Second
		}
		r.logger.Info("Rate limit hit (%d/%d), backing off for %v", info.Remaining, info.Limit, wait)

	case info.Remaining <= 3:
		wait = 5 * time.Second

	case info.Remaining <= info.Limit/2:
		wait = 2 * time.Second

	default:
		wait = 500 * time.Millisecond
	}

	if elapsed < wait {
		sleepFor := wait - elapsed
		r.lastRequest = now.Add(sleepFor)
		r.mu.Unlock()
		r.logger.Debug("Delaying %v due to rate pressure", sleepFor)
		return r.sleep(sleepFor)
	}

	r.lastRequest = now
	r.mu.Unlock()
	return nil
}

func (r *RateLimiter) AfterResponse(resp *http.Response) {
	remain, rrErr := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	limit, rlErr := strconv.Atoi(resp.Header.Get("X-Ratelimit-Limit"))
	resetUnix, reErr := strconv.ParseInt(resp.Header.Get("X-Ratelimit-Reset"), 10, 64)
	r.mu.Lock()
	defer r.mu.Unlock()

	reset := r.limit.Reset
	if reErr == nil {
		reset = time.Unix(resetUnix, 0)
	}

	if rrErr == nil && rlErr == nil {
		r.limit = RateLimitInfo{Remaining: remain, Limit: limit, Reset: reset}
	} else {
		r.logger.Debug("Rate limit headers missing or invalid; keeping existing values")
	}

	r.logger.Debug(fmt.Sprintf(
		"Rate limit headers — remaining: %d, limit: %d, reset: %v",
		remain, limit, reset,
	))
}

func (r *RateLimiter) Context() context.Context {
	return r.ctx
}

func (r *RateLimiter) Wrap(userCtx context.Context) context.Context {
	if userCtx == nil {
		return r.ctx
	}

	ctx, cancel := context.WithCancel(userCtx)
	go func() {
		select {
		case <-userCtx.Done():
		case <-r.ctx.Done():
			cancel()
		}
	}()
	return ctx
}

func (r *RateLimiter) sleep(d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-r.ctx.Done():
		return r.ctx.Err()
	case <-timer.C:
		return nil
	}
}

// DefaultRetryPolicy provides a basic retry strategy.
// It retries on 429 (Too Many Requests) and 5xx server errors.
type DefaultRetryPolicy struct{}

func isConnectionRefused(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}
	return false
}

// ShouldRetry determines if a request should be retried based on the response or error.
func (p *DefaultRetryPolicy) ShouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return true
		}

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true // Retry on timeouts
		}

		// Optionally retry on other connection-level errors
		if isConnectionRefused(err) {
			return true
		}

		return false
	}

	// Retry on specific HTTP status codes
	if resp.StatusCode == http.StatusTooManyRequests { // 429
		return true
	}
	if resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented && resp.StatusCode != http.StatusHTTPVersionNotSupported { // 5xx errors except 501 and 505
		return true
	}

	return false
}

// Wait calculates the duration to wait before the next retry attempt.
// Uses exponential backoff with jitter.
func (p *DefaultRetryPolicy) Wait(retries int) time.Duration {
	// Simple exponential backoff: 1s, 2s, 4s, ...
	baseDelay := time.Second
	// Calculate delay: 1s, 2s, 4s, 8s... (capped potentially later)
	delay := baseDelay * time.Duration(1<<(retries-1)) // retries starts from 1 for Wait

	// Cap the delay to avoid excessively long waits (e.g., max 30 seconds)
	maxDelay := 30 * time.Second
	if delay > maxDelay {
		delay = maxDelay
	}

	// Add jitter: +/- 10% of the delay
	// rand.Float64() returns a pseudo-random float64 in [0.0, 1.0)
	// We shift it to [-0.1, 0.1) by doing (rand.Float64() - 0.5) * 0.2
	jitterFraction := (rand.Float64() - 0.5) * 0.2 // Range roughly -0.1 to +0.1
	jitter := time.Duration(float64(delay) * jitterFraction)

	waitDuration := delay + jitter
	// Ensure wait duration is not negative
	if waitDuration < 0 {
		waitDuration = 0
	}

	return waitDuration
}

func (t *APITransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	var reqBodyBytes []byte

	// Read the body only if it exists and is not nil.
	// This allows retries for requests like GET that might have a nil body.
	if req.Body != nil {
		reqBodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body for retry: %w", err)
		}
		// It's crucial to close the original body after reading.
		req.Body.Close()
	}

	for retries := 0; ; retries++ {
		// --- Rate Limiter Check ---
		// Check rate limit *before* each attempt.
		if err := t.limiter.BeforeRequest(); err != nil {
			// If context is canceled during wait, return the context error.
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				t.logger.Warn("Context cancelled or deadline exceeded before request", "error", err)
				return nil, err // Return the context error
			}
			// Log other limiter errors but potentially allow retry logic to handle them if applicable
			t.logger.Error("Rate limiter pre-request check failed", "error", err)
			// Depending on the error, you might want to return immediately or let retry logic decide.
			// For now, we'll let the retry policy check this error.
		}

		// --- Prepare Request for Attempt ---
		// For each attempt, reset the request body if it exists.
		if reqBodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(reqBodyBytes))
		}

		// --- Make the HTTP Request ---
		currentResp, currentErr := t.underlying.RoundTrip(req)

		// --- Update Rate Limiter Info ---
		// Update rate limit info *after* each attempt, even if it failed,
		// as some APIs might return rate limit headers on error responses (e.g., 429).
		if currentResp != nil {
			t.limiter.AfterResponse(currentResp)
		}

		// --- Check if Retry is Needed ---
		// Use the latest response and error for the retry decision.
		resp = currentResp
		err = currentErr
		shouldRetry := t.retryConfig.RetryPolicy.ShouldRetry(resp, err)

		// --- Decide to Break or Continue ---
		if !shouldRetry || retries >= t.retryConfig.MaxRetries {
			// If we shouldn't retry, or we've exhausted retries, break the loop.
			break
		}

		// --- Wait Before Retrying ---
		waitTime := t.retryConfig.RetryPolicy.Wait(retries + 1) // Pass the *next* retry attempt number

		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if secs, err := strconv.Atoi(retryAfter); err == nil {
					waitTime = time.Duration(secs) * time.Second
				}
			}
		}

		t.logger.Debug("Retrying request",
			"attempt", retries+1,
			"max_retries", t.retryConfig.MaxRetries,
			"wait_duration", waitTime,
			"url", req.URL.String(),
			"error", err, // Log the error that triggered the retry
			"status_code", func() int { // Log status code if available
				if resp != nil {
					return resp.StatusCode
				}
				return 0
			}())

		select {
		case <-req.Context().Done():
			t.logger.Warn("Request context cancelled during retry wait", "error", req.Context().Err())
			// Return the latest response/error along with the context error
			// It might be more informative than just the context error alone.
			if err == nil { // If the last attempt had no error, return the context error
				err = req.Context().Err()
			}
			return resp, err // Return last known state + context error
		case <-time.After(waitTime):
			// Continue to the next iteration after waiting.
		}
	}

	// Return the response and error from the last attempt.
	return resp, err
}
