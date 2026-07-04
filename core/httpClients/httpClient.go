package httpClients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/pnaskardev/URL-Shortner-V1/core/helpers/utils"
)

type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

type RetryableHTTPClient struct {
	Client *http.Client
	Config RetryConfig
}

type DoWithRetryRequest struct {
	Ctx     context.Context
	Method  string
	Service utils.SERVICES
	Payload any
}

func (c *RetryableHTTPClient) DoWithRetry(target DoWithRetryRequest) (*http.Response, error) {

	url := utils.GlobalConstants[utils.SERVICES(target.Service)]

	// Serialize payload once before the loop to avoid re-marshaling on each retry.
	// A []byte payload is treated as an already-serialized body (e.g. protojson
	// output) and passed through untouched, since json.Marshal would corrupt it.
	var bodyBytes []byte
	if target.Payload != nil {
		if raw, ok := target.Payload.([]byte); ok {
			bodyBytes = raw
		} else {
			var err error
			bodyBytes, err = json.Marshal(target.Payload)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal payload: %w", err)
			}
		}
	}

	var (
		resp *http.Response
		err  error
	)

	for attempt := 0; attempt <= c.Config.MaxRetries; attempt++ {

		// Reconstruct the body reader on every attempt since io.Reader is consumed after first read
		var body io.Reader
		if bodyBytes != nil {
			body = bytes.NewBuffer(bodyBytes)
		}

		var req *http.Request
		req, err = http.NewRequestWithContext(target.Ctx, target.Method, url, body)
		if err != nil {
			return nil, fmt.Errorf("failed to build request: %w", err)
		}

		microserviceAuthToken, err := utils.GenerateMicroServiceAuthToken()
		if err != nil {
			slog.ErrorContext(target.Ctx, "HTTP request auth failed", "error", err)
			return nil, fmt.Errorf("HTTP request auth failed: %w", err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", microserviceAuthToken))

		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err = c.Client.Do(req)

		if err != nil {
			slog.ErrorContext(target.Ctx, "HTTP request failed",
				"attempt", attempt,
				"max_retries", c.Config.MaxRetries,
				"method", req.Method,
				"url", req.URL.String(),
				"error", err,
			)
			// Network-level error — fall through to retry logic below

		} else if resp.StatusCode >= 500 {
			// Server-side error — safe to retry
			resp.Body.Close()
			err = fmt.Errorf("server error: %d", resp.StatusCode)
			slog.WarnContext(target.Ctx, "HTTP request returned server error, will retry",
				"attempt", attempt,
				"max_retries", c.Config.MaxRetries,
				"method", req.Method,
				"url", req.URL.String(),
				"status_code", resp.StatusCode,
			)

		} else if resp.StatusCode >= 400 {
			// Client-side error — retrying won't help, return immediately
			resp.Body.Close()
			slog.WarnContext(target.Ctx, "HTTP request returned client error, aborting",
				"attempt", attempt,
				"method", req.Method,
				"url", req.URL.String(),
				"status_code", resp.StatusCode,
			)
			return nil, fmt.Errorf("client error: %d", resp.StatusCode)

		} else {
			// Success
			slog.InfoContext(target.Ctx, "HTTP request succeeded",
				"attempt", attempt,
				"method", req.Method,
				"url", req.URL.String(),
				"status_code", resp.StatusCode,
			)
			return resp, nil
		}

		// No more attempts left
		if attempt == c.Config.MaxRetries {
			slog.ErrorContext(target.Ctx, "max retries reached, giving up",
				"max_retries", c.Config.MaxRetries,
				"method", req.Method,
				"url", req.URL.String(),
				"error", err,
			)
			break
		}

		// Exponential backoff with jitter
		backoff := c.Config.BaseDelay * time.Duration(1<<attempt)
		if backoff > c.Config.MaxDelay {
			backoff = c.Config.MaxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
		sleep := backoff + jitter

		slog.InfoContext(target.Ctx, "retrying after backoff",
			"attempt", attempt,
			"next_attempt", attempt+1,
			"sleep", sleep.String(),
		)

		select {
		case <-time.After(sleep):
			continue
		case <-target.Ctx.Done():
			slog.WarnContext(target.Ctx, "context cancelled during retry backoff",
				"attempt", attempt,
				"error", target.Ctx.Err(),
			)
			return nil, target.Ctx.Err()
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.Config.MaxRetries, err)
}
