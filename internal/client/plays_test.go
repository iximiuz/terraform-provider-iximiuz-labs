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

func TestListPlays(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/plays" {
			t.Errorf("expected path /plays, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[
			{"id":"play-1","status":{"stateEvents":[{"state":"RUNNING","at":"2025-01-01T00:00:00Z"}]},"playground":{"name":"pg-1"},"pageUrl":"https://example.com/plays/play-1"},
			{"id":"play-2","status":{"stateEvents":[{"state":"STOPPED","at":"2025-01-01T00:00:00Z"}]},"playground":{"name":"pg-2"},"pageUrl":"https://example.com/plays/play-2"}
		]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	plays, err := c.ListPlays(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plays) != 2 {
		t.Fatalf("expected 2 plays, got %d", len(plays))
	}
	if plays[0].ID != "play-1" {
		t.Errorf("expected ID %q, got %q", "play-1", plays[0].ID)
	}
	if plays[0].Status.CurrentState() != "RUNNING" {
		t.Errorf("expected state %q, got %q", "RUNNING", plays[0].Status.CurrentState())
	}
	if plays[1].Status.CurrentState() != "STOPPED" {
		t.Errorf("expected state %q, got %q", "STOPPED", plays[1].Status.CurrentState())
	}
}

func TestGetPlay(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/plays/play-123" {
			t.Errorf("expected path /plays/play-123, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "play-123",
			"createdAt": "2025-01-01T00:00:00Z",
			"updatedAt": "2025-01-01T00:00:00Z",
			"expiresIn": 3600,
			"status": {"stateEvents": [{"state": "RUNNING", "at": "2025-01-01T00:00:00Z"}]},
			"playground": {"name": "my-pg", "title": "My PG"},
			"machines": [{"name": "node-01"}],
			"pageUrl": "https://example.com/plays/play-123"
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	play, err := c.GetPlay(context.Background(), "play-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if play.ID != "play-123" {
		t.Errorf("expected ID %q, got %q", "play-123", play.ID)
	}
	if play.Status.CurrentState() != "RUNNING" {
		t.Errorf("expected state %q, got %q", "RUNNING", play.Status.CurrentState())
	}
	if play.Playground.Name != "my-pg" {
		t.Errorf("expected playground name %q, got %q", "my-pg", play.Playground.Name)
	}
	if play.ExpiresIn != 3600 {
		t.Errorf("expected expiresIn %d, got %d", 3600, play.ExpiresIn)
	}
	if len(play.Machines) != 1 {
		t.Fatalf("expected 1 machine, got %d", len(play.Machines))
	}
}

func TestGetPlay_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":"not found"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetPlay(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreatePlay(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/plays" {
			t.Errorf("expected path /plays, got %s", r.URL.Path)
		}

		var req CreatePlayRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Playground != "my-playground" {
			t.Errorf("expected playground %q, got %q", "my-playground", req.Playground)
		}
		if !req.SafetyDisclaimerConsent {
			t.Error("expected SafetyDisclaimerConsent to be true")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"id": "play-new",
			"status": {"stateEvents": [{"state": "CREATED", "at": "2025-01-01T00:00:00Z"}]},
			"playground": {"name": "my-playground"},
			"pageUrl": "https://example.com/plays/play-new"
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	play, err := c.CreatePlay(context.Background(), CreatePlayRequest{
		Playground:              "my-playground",
		SafetyDisclaimerConsent: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if play.ID != "play-new" {
		t.Errorf("expected ID %q, got %q", "play-new", play.ID)
	}
	if play.Status.CurrentState() != "CREATED" {
		t.Errorf("expected state %q, got %q", "CREATED", play.Status.CurrentState())
	}
}

func TestPlayAction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/plays/play-123/actions" {
			t.Errorf("expected path /plays/play-123/actions, got %s", r.URL.Path)
		}

		var req PlayActionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Action != "stop" {
			t.Errorf("expected action %q, got %q", "stop", req.Action)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.PlayAction(context.Background(), "play-123", "stop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDestroyPlay(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/plays/play-123/actions" {
			t.Errorf("expected path /plays/play-123/actions, got %s", r.URL.Path)
		}

		var req PlayActionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Action != "destroy" {
			t.Errorf("expected action %q, got %q", "destroy", req.Action)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.DestroyPlay(context.Background(), "play-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreatePlay_WithInitConditions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreatePlayRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.InitConditions["runtime"] != "docker" {
			t.Errorf("expected initConditions[runtime]=%q, got %q", "docker", req.InitConditions["runtime"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"play-cond","status":{"stateEvents":[{"state":"CREATED","at":"2025-01-01T00:00:00Z"}]},"playground":{"name":"pg"},"pageUrl":""}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.CreatePlay(context.Background(), CreatePlayRequest{
		Playground:     "pg",
		InitConditions: map[string]string{"runtime": "docker"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
