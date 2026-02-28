// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// GetRoadmap gets a roadmap by name.
func (c *Client) GetRoadmap(ctx context.Context, name string) (*Roadmap, error) {
	var r Roadmap
	if err := c.GetJSON(ctx, "/roadmaps/"+url.PathEscape(name), nil, &r); err != nil {
		return nil, fmt.Errorf("getting roadmap %q: %w", name, err)
	}
	return &r, nil
}

// ListRoadmaps lists all roadmaps.
func (c *Client) ListRoadmaps(ctx context.Context) ([]Roadmap, error) {
	var roadmaps []Roadmap
	if err := c.GetJSON(ctx, "/roadmaps", nil, &roadmaps); err != nil {
		return nil, fmt.Errorf("listing roadmaps: %w", err)
	}
	return roadmaps, nil
}

// CreateRoadmap creates a new roadmap.
func (c *Client) CreateRoadmap(ctx context.Context, name string) (*Roadmap, error) {
	req := CreateContentRequest{Name: name}
	var r Roadmap
	if err := c.PostJSON(ctx, "/roadmaps", req, &r); err != nil {
		return nil, fmt.Errorf("creating roadmap %q: %w", name, err)
	}
	return &r, nil
}

// DeleteRoadmap deletes a roadmap by name.
func (c *Client) DeleteRoadmap(ctx context.Context, name string) error {
	if err := c.Delete(ctx, "/roadmaps/"+url.PathEscape(name), nil); err != nil {
		return fmt.Errorf("deleting roadmap %q: %w", name, err)
	}
	return nil
}
