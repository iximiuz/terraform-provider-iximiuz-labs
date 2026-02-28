// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient creates a Client pointing at the given test server.
func newTestClient(ts *httptest.Server) *Client {
	return &Client{
		baseURL:     ts.URL,
		apiBaseURL:  ts.URL,
		sessionID:   "test-session",
		accessToken: "test-token",
		userAgent:   "test-agent/1.0",
		httpClient:  ts.Client(),
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(ClientConfig{
		BaseURL:     "https://example.com",
		APIBaseURL:  "https://example.com/api",
		SessionID:   "sess123",
		AccessToken: "tok456",
		UserAgent:   "my-agent/1.0",
	})

	if c.baseURL != "https://example.com" {
		t.Errorf("expected baseURL %q, got %q", "https://example.com", c.baseURL)
	}
	if c.apiBaseURL != "https://example.com/api" {
		t.Errorf("expected apiBaseURL %q, got %q", "https://example.com/api", c.apiBaseURL)
	}
	if c.sessionID != "sess123" {
		t.Errorf("expected sessionID %q, got %q", "sess123", c.sessionID)
	}
	if c.accessToken != "tok456" {
		t.Errorf("expected accessToken %q, got %q", "tok456", c.accessToken)
	}
	if c.userAgent != "my-agent/1.0" {
		t.Errorf("expected userAgent %q, got %q", "my-agent/1.0", c.userAgent)
	}
	if c.httpClient == nil {
		t.Error("expected httpClient to be non-nil")
	}
}

func TestAuthHeader(t *testing.T) {
	var receivedAuth string
	var receivedUA string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		receivedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.Get(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCreds := base64.StdEncoding.EncodeToString([]byte("test-session:test-token"))
	expectedAuth := "Basic " + expectedCreds
	if receivedAuth != expectedAuth {
		t.Errorf("expected Authorization %q, got %q", expectedAuth, receivedAuth)
	}
	if receivedUA != "test-agent/1.0" {
		t.Errorf("expected User-Agent %q, got %q", "test-agent/1.0", receivedUA)
	}
}

func TestContentTypeHeader(t *testing.T) {
	var receivedContentType string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.PostJSON(context.Background(), "/test", map[string]string{"key": "value"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedContentType != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", receivedContentType)
	}
}

func TestGetJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/test-path" {
			t.Errorf("expected path /test-path, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"name":"hello","value":42}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)

	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	err := c.GetJSON(context.Background(), "/test-path", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "hello" {
		t.Errorf("expected name %q, got %q", "hello", result.Name)
	}
	if result.Value != 42 {
		t.Errorf("expected value %d, got %d", 42, result.Value)
	}
}

func TestGetWithQueryParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("filter") != "docker" {
			t.Errorf("expected filter=docker, got %q", r.URL.Query().Get("filter"))
		}
		if r.URL.Query().Get("format") != "extended" {
			t.Errorf("expected format=extended, got %q", r.URL.Query().Get("format"))
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	resp, err := c.Get(context.Background(), "/test", map[string][]string{
		"filter": {"docker"},
		"format": {"extended"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if body["name"] != "test-name" {
			t.Errorf("expected body name %q, got %q", "test-name", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"created-123"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)

	var result struct {
		ID string `json:"id"`
	}
	err := c.PostJSON(context.Background(), "/test", map[string]string{"name": "test-name"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "created-123" {
		t.Errorf("expected id %q, got %q", "created-123", result.ID)
	}
}

func TestPutJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"updated":true}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	var result struct {
		Updated bool `json:"updated"`
	}
	err := c.PutJSON(context.Background(), "/test", map[string]string{}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Updated {
		t.Error("expected updated to be true")
	}
}

func TestPatchJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"patched":true}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	var result struct {
		Patched bool `json:"patched"`
	}
	err := c.PatchJSON(context.Background(), "/test", map[string]string{}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Patched {
		t.Error("expected patched to be true")
	}
}

func TestDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/test-path" {
			t.Errorf("expected path /test-path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.Delete(context.Background(), "/test-path", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckResponse_Success(t *testing.T) {
	for _, code := range []int{200, 201, 204, 299} {
		resp := &http.Response{
			StatusCode: code,
			Body:       http.NoBody,
		}
		if err := checkResponse(resp); err != nil {
			t.Errorf("expected no error for status %d, got: %v", code, err)
		}
	}
}

func TestCheckResponse_NotFound(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"error":"not found"}`)),
	}
	err := checkResponse(resp)
	if err == nil {
		t.Fatal("expected error for 404")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true")
	}
}

func TestCheckResponse_Unauthorized(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(strings.NewReader(`{"error":"unauthorized"}`)),
	}
	err := checkResponse(resp)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if !errors.Is(err, ErrAuthenticationRequired) {
		t.Error("expected ErrAuthenticationRequired")
	}
}

func TestCheckResponse_RateLimit(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	err := checkResponse(resp)
	if err == nil {
		t.Fatal("expected error for 429")
	}
	if !errors.Is(err, ErrRateLimitExceeded) {
		t.Error("expected ErrRateLimitExceeded")
	}
}

func TestCheckResponse_GatewayTimeout(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusGatewayTimeout,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	err := checkResponse(resp)
	if err == nil {
		t.Fatal("expected error for 504")
	}
	if !errors.Is(err, ErrGatewayTimeout) {
		t.Error("expected ErrGatewayTimeout")
	}
}

func TestCheckResponse_ServerError(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`internal server error`)),
	}
	err := checkResponse(resp)
	if err == nil {
		t.Fatal("expected error for 500")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", apiErr.StatusCode)
	}
	if apiErr.Body != "internal server error" {
		t.Errorf("expected body %q, got %q", "internal server error", apiErr.Body)
	}
}

func TestRetryOnGatewayTimeout(t *testing.T) {
	var attempts int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"ok":true}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	// Override the http client timeout so retries don't wait forever
	c.httpClient.Timeout = 30 * time.Second

	var result struct {
		OK bool `json:"ok"`
	}
	err := c.GetJSON(context.Background(), "/test", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.OK {
		t.Error("expected ok to be true")
	}
	if n := atomic.LoadInt32(&attempts); n < 3 {
		t.Errorf("expected at least 3 attempts, got %d", n)
	}
}

func TestRetryOnRateLimit(t *testing.T) {
	var attempts int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"ok":true}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)

	var result struct {
		OK bool `json:"ok"`
	}
	err := c.GetJSON(context.Background(), "/test", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.OK {
		t.Error("expected ok to be true")
	}
	if n := atomic.LoadInt32(&attempts); n < 2 {
		t.Errorf("expected at least 2 attempts, got %d", n)
	}
}

func TestContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := newTestClient(ts)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := c.Get(ctx, "/test", nil)
	if err == nil {
		t.Fatal("expected error due to context cancellation")
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name:     "with body",
			err:      &APIError{StatusCode: 404, Body: `{"msg":"not found"}`, Err: ErrNotFound},
			expected: `not found (HTTP 404): {"msg":"not found"}`,
		},
		{
			name:     "without body",
			err:      &APIError{StatusCode: 401, Err: ErrAuthenticationRequired},
			expected: "authentication required (HTTP 401)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	apiErr := &APIError{StatusCode: 404, Err: ErrNotFound}
	if !errors.Is(apiErr, ErrNotFound) {
		t.Error("expected Unwrap to return ErrNotFound")
	}
}

func TestIsNotFound(t *testing.T) {
	if IsNotFound(fmt.Errorf("some error")) {
		t.Error("expected IsNotFound to return false for non-NotFound error")
	}
	if !IsNotFound(ErrNotFound) {
		t.Error("expected IsNotFound to return true for ErrNotFound")
	}
	apiErr := &APIError{StatusCode: 404, Err: ErrNotFound}
	if !IsNotFound(apiErr) {
		t.Error("expected IsNotFound to return true for wrapped ErrNotFound")
	}
}
