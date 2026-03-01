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

func TestListPlaygrounds(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/playgrounds" {
			t.Errorf("expected path /playgrounds, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":"pg1","name":"test-pg","title":"Test Playground"}]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	playgrounds, err := c.ListPlaygrounds(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(playgrounds) != 1 {
		t.Fatalf("expected 1 playground, got %d", len(playgrounds))
	}
	if playgrounds[0].Name != "test-pg" {
		t.Errorf("expected name %q, got %q", "test-pg", playgrounds[0].Name)
	}
	if playgrounds[0].Title != "Test Playground" {
		t.Errorf("expected title %q, got %q", "Test Playground", playgrounds[0].Title)
	}
}

func TestListPlaygrounds_WithFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("filter") != "docker" {
			t.Errorf("expected filter=docker, got %q", r.URL.Query().Get("filter"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	playgrounds, err := c.ListPlaygrounds(context.Background(), "docker")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(playgrounds) != 0 {
		t.Fatalf("expected 0 playgrounds, got %d", len(playgrounds))
	}
}

func TestGetPlayground(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/playgrounds/my-playground" {
			t.Errorf("expected path /playgrounds/my-playground, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("format") != "extended" {
			t.Errorf("expected format=extended, got %q", r.URL.Query().Get("format"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "pg-123",
			"name": "my-playground",
			"title": "My Playground",
			"description": "A test playground",
			"published": true,
			"pageUrl": "https://labs.iximiuz.com/playgrounds/my-playground",
			"machines": [{"name": "node-01"}],
			"networks": [{"name": "net-01", "subnet": "10.0.0.0/24"}]
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	pg, err := c.GetPlayground(context.Background(), "my-playground")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pg.ID != "pg-123" {
		t.Errorf("expected ID %q, got %q", "pg-123", pg.ID)
	}
	if pg.Name != "my-playground" {
		t.Errorf("expected name %q, got %q", "my-playground", pg.Name)
	}
	if pg.Title != "My Playground" {
		t.Errorf("expected title %q, got %q", "My Playground", pg.Title)
	}
	if !pg.Published {
		t.Error("expected published to be true")
	}
	if len(pg.Machines) != 1 {
		t.Fatalf("expected 1 machine, got %d", len(pg.Machines))
	}
	if pg.Machines[0].Name != "node-01" {
		t.Errorf("expected machine name %q, got %q", "node-01", pg.Machines[0].Name)
	}
}

func TestGetPlayground_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":"not found"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetPlayground(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreatePlayground(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/playgrounds" {
			t.Errorf("expected path /playgrounds, got %s", r.URL.Path)
		}

		var req CreatePlaygroundRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Name != "new-playground" {
			t.Errorf("expected name %q, got %q", "new-playground", req.Name)
		}
		if req.Title != "New Playground" {
			t.Errorf("expected title %q, got %q", "New Playground", req.Title)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"id": "pg-new",
			"name": "new-playground",
			"title": "New Playground",
			"description": "A new playground",
			"published": false,
			"pageUrl": "https://labs.iximiuz.com/playgrounds/new-playground"
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	pg, err := c.CreatePlayground(context.Background(), CreatePlaygroundRequest{
		Name:        "new-playground",
		Title:       "New Playground",
		Description: "A new playground",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pg.ID != "pg-new" {
		t.Errorf("expected ID %q, got %q", "pg-new", pg.ID)
	}
	if pg.Name != "new-playground" {
		t.Errorf("expected name %q, got %q", "new-playground", pg.Name)
	}
}

func TestUpdatePlayground(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/playgrounds/my-playground" {
			t.Errorf("expected path /playgrounds/my-playground, got %s", r.URL.Path)
		}

		var req UpdatePlaygroundRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Title != "Updated Title" {
			t.Errorf("expected title %q, got %q", "Updated Title", req.Title)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id": "pg-123",
			"name": "my-playground",
			"title": "Updated Title",
			"description": "Updated desc",
			"published": true,
			"pageUrl": "https://labs.iximiuz.com/playgrounds/my-playground"
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	pg, err := c.UpdatePlayground(context.Background(), "my-playground", UpdatePlaygroundRequest{
		Title:       "Updated Title",
		Description: "Updated desc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pg.Title != "Updated Title" {
		t.Errorf("expected title %q, got %q", "Updated Title", pg.Title)
	}
}

func TestDeletePlayground(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/playgrounds/my-playground" {
			t.Errorf("expected path /playgrounds/my-playground, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.DeletePlayground(context.Background(), "my-playground")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeletePlayground_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":"not found"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.DeletePlayground(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreatePlayground_WithMachinesAndNetworks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreatePlaygroundRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if len(req.Machines) != 1 {
			t.Fatalf("expected 1 machine, got %d", len(req.Machines))
		}
		if req.Machines[0].Name != "node-01" {
			t.Errorf("expected machine name %q, got %q", "node-01", req.Machines[0].Name)
		}
		if len(req.Networks) != 1 {
			t.Fatalf("expected 1 network, got %d", len(req.Networks))
		}
		if req.Networks[0].Subnet != "10.0.0.0/24" {
			t.Errorf("expected subnet %q, got %q", "10.0.0.0/24", req.Networks[0].Subnet)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"pg-complex","name":"complex-pg","title":"","description":"","published":false,"pageUrl":""}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.CreatePlayground(context.Background(), CreatePlaygroundRequest{
		Name: "complex-pg",
		Machines: []PlaygroundMachine{
			{Name: "node-01"},
		},
		Networks: []PlaygroundNetwork{
			{Name: "net-01", Subnet: "10.0.0.0/24"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
