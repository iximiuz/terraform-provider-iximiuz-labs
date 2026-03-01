// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Me Tests ---

func TestGetMe(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/auth/me" {
			t.Errorf("expected path /auth/me, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("authenticate") != "true" {
			t.Errorf("expected authenticate=true, got %q", r.URL.Query().Get("authenticate"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "user-123",
			"githubProfileId": "gh-456",
			"premiumAccess": {
				"until": "2026-01-01T00:00:00Z",
				"lifetime": false,
				"trial": false
			}
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	me, err := c.GetMe(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if me.ID != "user-123" {
		t.Errorf("expected ID %q, got %q", "user-123", me.ID)
	}
	if me.GithubProfileID != "gh-456" {
		t.Errorf("expected GithubProfileID %q, got %q", "gh-456", me.GithubProfileID)
	}
	if me.PremiumAccess == nil {
		t.Fatal("expected PremiumAccess to be non-nil")
	}
	if me.PremiumAccess.Until != "2026-01-01T00:00:00Z" {
		t.Errorf("expected until %q, got %q", "2026-01-01T00:00:00Z", me.PremiumAccess.Until)
	}
}

func TestGetMe_NoPremium(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"user-free","githubProfileId":"gh-789"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	me, err := c.GetMe(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if me.PremiumAccess != nil {
		t.Error("expected PremiumAccess to be nil for free user")
	}
}

func TestGetMe_Unauthorized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error":"unauthorized"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetMe(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Port Tests ---

func TestListPorts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/plays/play-123/ports" {
			t.Errorf("expected path /plays/play-123/ports, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{"id":"port-1","playId":"play-123","machine":"node-01","number":8080,"access":"public","url":"https://h1.example.com"},
			{"id":"port-2","playId":"play-123","machine":"node-01","number":3000,"access":"private"}
		]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	ports, err := c.ListPorts(context.Background(), "play-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(ports))
	}
	if ports[0].ID != "port-1" {
		t.Errorf("expected ID %q, got %q", "port-1", ports[0].ID)
	}
	if ports[0].Number != 8080 {
		t.Errorf("expected number %d, got %d", 8080, ports[0].Number)
	}
	if ports[0].AccessMode != "public" {
		t.Errorf("expected accessMode %q, got %q", "public", ports[0].AccessMode)
	}
}

func TestExposePort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/plays/play-123/ports" {
			t.Errorf("expected path /plays/play-123/ports, got %s", r.URL.Path)
		}

		var req ExposePortRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Machine != "node-01" {
			t.Errorf("expected machine %q, got %q", "node-01", req.Machine)
		}
		if req.Number != 8080 {
			t.Errorf("expected number %d, got %d", 8080, req.Number)
		}
		if req.Access != "public" {
			t.Errorf("expected access %q, got %q", "public", req.Access)
		}
		if !req.TLS {
			t.Error("expected TLS to be true")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"id": "port-new",
			"playId": "play-123",
			"machine": "node-01",
			"number": 8080,
			"access": "public",
			"tls": true,
			"url": "https://h.example.com"
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	port, err := c.ExposePort(context.Background(), "play-123", ExposePortRequest{
		Machine: "node-01",
		Number:  8080,
		Access:  "public",
		TLS:     true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port.ID != "port-new" {
		t.Errorf("expected ID %q, got %q", "port-new", port.ID)
	}
	if port.URL != "https://h.example.com" {
		t.Errorf("expected URL %q, got %q", "https://h.example.com", port.URL)
	}
}

func TestUnexposePort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/plays/play-123/ports/port-456" {
			t.Errorf("expected path /plays/play-123/ports/port-456, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.UnexposePort(context.Background(), "play-123", "port-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Shell Tests ---

func TestListShells(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/plays/play-123/shells" {
			t.Errorf("expected path /plays/play-123/shells, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{"id":"shell-1","playId":"play-123","machine":"node-01","user":"root","access":"public","url":"https://s1.example.com"}
		]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	shells, err := c.ListShells(context.Background(), "play-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(shells) != 1 {
		t.Fatalf("expected 1 shell, got %d", len(shells))
	}
	if shells[0].ID != "shell-1" {
		t.Errorf("expected ID %q, got %q", "shell-1", shells[0].ID)
	}
	if shells[0].User != "root" {
		t.Errorf("expected user %q, got %q", "root", shells[0].User)
	}
}

func TestExposeShell(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/plays/play-123/shells" {
			t.Errorf("expected path /plays/play-123/shells, got %s", r.URL.Path)
		}

		var req ExposeShellRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Machine != "node-01" {
			t.Errorf("expected machine %q, got %q", "node-01", req.Machine)
		}
		if req.User != "root" {
			t.Errorf("expected user %q, got %q", "root", req.User)
		}
		if req.Access != "public" {
			t.Errorf("expected access %q, got %q", "public", req.Access)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"id": "shell-new",
			"playId": "play-123",
			"machine": "node-01",
			"user": "root",
			"access": "public",
			"url": "https://s.example.com"
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	shell, err := c.ExposeShell(context.Background(), "play-123", ExposeShellRequest{
		Machine: "node-01",
		User:    "root",
		Access:  "public",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shell.ID != "shell-new" {
		t.Errorf("expected ID %q, got %q", "shell-new", shell.ID)
	}
	if shell.URL != "https://s.example.com" {
		t.Errorf("expected URL %q, got %q", "https://s.example.com", shell.URL)
	}
}

func TestUnexposeShell(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/plays/play-123/shells/shell-456" {
			t.Errorf("expected path /plays/play-123/shells/shell-456, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.UnexposeShell(context.Background(), "play-123", "shell-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
