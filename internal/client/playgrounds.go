// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
)

// ListPlaygrounds lists all playgrounds, optionally filtered.
func (c *Client) ListPlaygrounds(ctx context.Context, filter string) ([]Playground, error) {
	query := url.Values{}
	if filter != "" {
		query.Set("filter", filter)
	}

	var playgrounds []Playground
	if err := c.GetJSON(ctx, "/playgrounds", query, &playgrounds); err != nil {
		return nil, fmt.Errorf("listing playgrounds: %w", err)
	}
	return playgrounds, nil
}

// GetPlayground gets a playground by name with extended format.
func (c *Client) GetPlayground(ctx context.Context, name string) (*Playground, error) {
	query := url.Values{}
	query.Set("format", "extended")

	var pg Playground
	if err := c.GetJSON(ctx, "/playgrounds/"+url.PathEscape(name), query, &pg); err != nil {
		return nil, fmt.Errorf("getting playground %q: %w", name, err)
	}
	return &pg, nil
}

// CreatePlayground creates a new playground.
func (c *Client) CreatePlayground(ctx context.Context, req CreatePlaygroundRequest) (*Playground, error) {
	var pg Playground
	if err := c.PostJSON(ctx, "/playgrounds", req, &pg); err != nil {
		return nil, fmt.Errorf("creating playground: %w", err)
	}
	return &pg, nil
}

// UpdatePlayground updates an existing playground.
func (c *Client) UpdatePlayground(ctx context.Context, name string, req UpdatePlaygroundRequest) (*Playground, error) {
	var pg Playground
	if err := c.PutJSON(ctx, "/playgrounds/"+url.PathEscape(name), req, &pg); err != nil {
		return nil, fmt.Errorf("updating playground %q: %w", name, err)
	}
	return &pg, nil
}

// DeletePlayground deletes a playground by name.
// If the API returns 400 (playground still has an active play), it retries up
// to 3 times with exponential backoff (1 s, 2 s, 4 s) before giving up.
func (c *Client) DeletePlayground(ctx context.Context, name string) error {
	const maxDeleteRetries = 3

	var lastErr error
	for attempt := 0; attempt <= maxDeleteRetries; attempt++ {
		if attempt > 0 {
			delay := time.Second * time.Duration(math.Pow(2, float64(attempt-1)))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		err := c.Delete(ctx, "/playgrounds/"+url.PathEscape(name), nil)
		if err == nil {
			return nil
		}

		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusBadRequest {
			lastErr = err
			continue
		}

		return fmt.Errorf("deleting playground %q: %w", name, err)
	}

	return fmt.Errorf("deleting playground %q: %w", name, lastErr)
}
