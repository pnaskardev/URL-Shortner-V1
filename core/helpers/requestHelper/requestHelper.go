package requesthelper

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

func CreateHTTPRequest[T any](ctx context.Context, method string, url string, payload *T) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil

}

func (c *RetryableHTTPClient) DoWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.Config.MaxRetries; attempt++ {
		clonedReq := req.Clone(ctx)

		resp, err = c.Client.Do(clonedReq)

		if err != nil {
			slog.ErrorContext(ctx, "HTTP request failed",
				"attempt", attempt,
				"max_retries", c.Config.MaxRetries,
				"method", req.Method,
				"url", req.URL.String(),
				"error", err,
			)
		} else if resp.StatusCode >= 400 {
			slog.WarnContext(ctx, "HTTP request returned server error",
				"attempt", attempt,
				"max_retries", c.Config.MaxRetries,
				"method", req.Method,
				"url", req.URL.String(),
				"status_code", resp.StatusCode,
			)
			return nil, fmt.Errorf("REQUEST FAILED WITH - ", resp.StatusCode)
		} else {
			slog.InfoContext(ctx, "HTTP request succeeded",
				"attempt", attempt,
				"method", req.Method,
				"url", req.URL.String(),
				"status_code", resp.StatusCode,
			)
			return resp, nil
		}

		if attempt == c.Config.MaxRetries {
			slog.ErrorContext(ctx, "max retries reached, giving up",
				"max_retries", c.Config.MaxRetries,
				"method", req.Method,
				"url", req.URL.String(),
				"error", err,
			)
			break
		}

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

	return resp, err
}
