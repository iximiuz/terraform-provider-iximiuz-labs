// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// ListPlays lists all plays.
func (c *Client) ListPlays(ctx context.Context) ([]Play, error) {
	var plays []Play
	if err := c.GetJSON(ctx, "/plays", nil, &plays); err != nil {
		return nil, fmt.Errorf("listing plays: %w", err)
	}
	return plays, nil
}

// GetPlay gets a play by ID.
func (c *Client) GetPlay(ctx context.Context, id string) (*Play, error) {
	var play Play
	if err := c.GetJSON(ctx, "/plays/"+url.PathEscape(id), nil, &play); err != nil {
		return nil, fmt.Errorf("getting play %q: %w", id, err)
	}
	return &play, nil
}

// CreatePlay creates a new play (starts a playground instance).
func (c *Client) CreatePlay(ctx context.Context, req CreatePlayRequest) (*Play, error) {
	var play Play
	if err := c.PostJSON(ctx, "/plays", req, &play); err != nil {
		return nil, fmt.Errorf("creating play: %w", err)
	}
	return &play, nil
}

// PlayAction performs a lifecycle action on a play (stop, restart, destroy, make_persistent).
func (c *Client) PlayAction(ctx context.Context, id string, action string) error {
	req := PlayActionRequest{Action: action}
	if err := c.PostJSON(ctx, "/plays/"+url.PathEscape(id)+"/actions", req, nil); err != nil {
		return fmt.Errorf("play action %q on %q: %w", action, id, err)
	}
	return nil
}

// DestroyPlay destroys a play.
func (c *Client) DestroyPlay(ctx context.Context, id string) error {
	return c.PlayAction(ctx, id, "destroy")
}
