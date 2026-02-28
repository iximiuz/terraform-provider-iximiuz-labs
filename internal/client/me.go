// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// GetMe gets the current authenticated user's info.
func (c *Client) GetMe(ctx context.Context) (*Me, error) {
	query := url.Values{}
	query.Set("authenticate", "true")

	var me Me
	if err := c.GetJSON(ctx, "/auth/me", query, &me); err != nil {
		return nil, fmt.Errorf("getting current user: %w", err)
	}
	return &me, nil
}
