// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	maxRetries       = 5
	maxRetryDuration = 10 * time.Second
	baseRetryDelay   = 500 * time.Millisecond
)

// Client provides methods for interacting with the iximiuz Labs API.
type Client struct {
	baseURL     string
	apiBaseURL  string
	sessionID   string
	accessToken string
	userAgent   string
	httpClient  *http.Client
}

// ClientConfig holds configuration for creating a new Client.
type ClientConfig struct {
	BaseURL     string
	APIBaseURL  string
	SessionID   string
	AccessToken string
	UserAgent   string
}

// NewClient creates a new API client with the given configuration.
func NewClient(cfg ClientConfig) *Client {
	return &Client{
		baseURL:     cfg.BaseURL,
		apiBaseURL:  cfg.APIBaseURL,
		sessionID:   cfg.SessionID,
		accessToken: cfg.AccessToken,
		userAgent:   cfg.UserAgent,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// do executes an HTTP request with authentication, retry logic, and error handling.
func (c *Client) do(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	// Buffer the body so it can be re-read on retries.
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("reading request body: %w", err)
		}
	}

	var lastErr error
	start := time.Now()

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if time.Since(start) > maxRetryDuration {
			break
		}

		if attempt > 0 {
			delay := baseRetryDelay * time.Duration(math.Pow(2, float64(attempt-1)))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		// Set authentication header
		creds := base64.StdEncoding.EncodeToString([]byte(c.sessionID + ":" + c.accessToken))
		req.Header.Set("Authorization", "Basic "+creds)
		req.Header.Set("User-Agent", c.userAgent)
		if bodyBytes != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("executing request: %w", err)
			continue
		}

		// Check for retryable status codes
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			resp.Body.Close()
			lastErr = ErrRateLimitExceeded

			// Respect X-Ratelimit-Reset header
			if resetStr := resp.Header.Get("X-Ratelimit-Reset"); resetStr != "" {
				if resetTime, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
					waitDuration := time.Until(time.Unix(resetTime, 0))
					if waitDuration > 0 && waitDuration < maxRetryDuration {
						select {
						case <-ctx.Done():
							return nil, ctx.Err()
						case <-time.After(waitDuration):
						}
					}
				}
			}
			continue

		case http.StatusGatewayTimeout:
			resp.Body.Close()
			lastErr = ErrGatewayTimeout
			continue
		}

		return resp, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}

// checkResponse checks an HTTP response for errors and returns an appropriate error.
func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Body:       string(bodyBytes),
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		apiErr.Err = ErrAuthenticationRequired
	case http.StatusNotFound:
		apiErr.Err = ErrNotFound
	case http.StatusTooManyRequests:
		apiErr.Err = ErrRateLimitExceeded
	case http.StatusGatewayTimeout:
		apiErr.Err = ErrGatewayTimeout
	default:
		apiErr.Err = fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return apiErr
}

// Get performs a GET request to the given API path with optional query parameters.
func (c *Client) Get(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	u := c.apiBaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	return c.do(ctx, http.MethodGet, u, nil)
}

// GetJSON performs a GET request and decodes the JSON response into target.
func (c *Client) GetJSON(ctx context.Context, path string, query url.Values, target interface{}) error {
	resp, err := c.Get(ctx, path, query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// PostJSON performs a POST request with a JSON body and decodes the response.
func (c *Client) PostJSON(ctx context.Context, path string, body interface{}, target interface{}) error {
	return c.doJSON(ctx, http.MethodPost, path, body, target)
}

// PutJSON performs a PUT request with a JSON body and decodes the response.
func (c *Client) PutJSON(ctx context.Context, path string, body interface{}, target interface{}) error {
	return c.doJSON(ctx, http.MethodPut, path, body, target)
}

// PatchJSON performs a PATCH request with a JSON body and decodes the response.
func (c *Client) PatchJSON(ctx context.Context, path string, body interface{}, target interface{}) error {
	return c.doJSON(ctx, http.MethodPatch, path, body, target)
}

// Delete performs a DELETE request to the given API path.
func (c *Client) Delete(ctx context.Context, path string, query url.Values) error {
	u := c.apiBaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	resp, err := c.do(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkResponse(resp)
}

// doJSON is a helper that performs a JSON request (POST/PUT/PATCH) and optionally decodes the response.
func (c *Client) doJSON(ctx context.Context, method, path string, body interface{}, target interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	u := c.apiBaseURL + path
	resp, err := c.do(ctx, method, u, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return err
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}

	return nil
}
