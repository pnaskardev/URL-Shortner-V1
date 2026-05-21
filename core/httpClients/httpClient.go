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

func (c *RetryableHTTPClient) DoWithRetry(ctx context.Context, method string, service utils.SERVICES, payload *any) (*http.Response, error) {

	url := utils.GlobalConstants[utils.SERVICES(service)]

	// Serialize payload once before the loop to avoid re-marshaling on each retry
	var bodyBytes []byte
	if payload != nil {
		var err error
		bodyBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
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
		req, err = http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, fmt.Errorf("failed to build request: %w", err)
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err = c.Client.Do(req)

		if err != nil {
			slog.ErrorContext(ctx, "HTTP request failed",
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
			slog.WarnContext(ctx, "HTTP request returned server error, will retry",
				"attempt", attempt,
				"max_retries", c.Config.MaxRetries,
				"method", req.Method,
				"url", req.URL.String(),
				"status_code", resp.StatusCode,
			)

		} else if resp.StatusCode >= 400 {
			// Client-side error — retrying won't help, return immediately
			resp.Body.Close()
			slog.WarnContext(ctx, "HTTP request returned client error, aborting",
				"attempt", attempt,
				"method", req.Method,
				"url", req.URL.String(),
				"status_code", resp.StatusCode,
			)
			return nil, fmt.Errorf("client error: %d", resp.StatusCode)

		} else {
			// Success
			slog.InfoContext(ctx, "HTTP request succeeded",
				"attempt", attempt,
				"method", req.Method,
				"url", req.URL.String(),
				"status_code", resp.StatusCode,
			)
			return resp, nil
		}

		// No more attempts left
		if attempt == c.Config.MaxRetries {
			slog.ErrorContext(ctx, "max retries reached, giving up",
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

		slog.InfoContext(ctx, "retrying after backoff",
			"attempt", attempt,
			"next_attempt", attempt+1,
			"sleep", sleep.String(),
		)

		select {
		case <-time.After(sleep):
			continue
		case <-ctx.Done():
			slog.WarnContext(ctx, "context cancelled during retry backoff",
				"attempt", attempt,
				"error", ctx.Err(),
			)
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.Config.MaxRetries, err)
}
