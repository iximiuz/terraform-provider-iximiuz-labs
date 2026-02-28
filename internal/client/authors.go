// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
)

// GetAuthor gets the current user's author profile.
func (c *Client) GetAuthor(ctx context.Context) (*Author, error) {
	var a Author
	if err := c.GetJSON(ctx, "/author", nil, &a); err != nil {
		return nil, fmt.Errorf("getting author profile: %w", err)
	}
	return &a, nil
}

// CreateAuthor creates a new author profile.
func (c *Client) CreateAuthor(ctx context.Context, req CreateAuthorRequest) (*Author, error) {
	var a Author
	if err := c.PostJSON(ctx, "/authors", req, &a); err != nil {
		return nil, fmt.Errorf("creating author: %w", err)
	}
	return &a, nil
}
