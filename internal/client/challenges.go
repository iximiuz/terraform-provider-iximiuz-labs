// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// GetChallenge gets a challenge by name.
func (c *Client) GetChallenge(ctx context.Context, name string) (*Challenge, error) {
	var ch Challenge
	if err := c.GetJSON(ctx, "/challenges/"+url.PathEscape(name), nil, &ch); err != nil {
		return nil, fmt.Errorf("getting challenge %q: %w", name, err)
	}
	return &ch, nil
}

// ListChallenges lists all challenges.
func (c *Client) ListChallenges(ctx context.Context) ([]Challenge, error) {
	var challenges []Challenge
	if err := c.GetJSON(ctx, "/challenges", nil, &challenges); err != nil {
		return nil, fmt.Errorf("listing challenges: %w", err)
	}
	return challenges, nil
}

// CreateChallenge creates a new challenge.
func (c *Client) CreateChallenge(ctx context.Context, name string) (*Challenge, error) {
	req := CreateContentRequest{Name: name}
	var ch Challenge
	if err := c.PostJSON(ctx, "/challenges", req, &ch); err != nil {
		return nil, fmt.Errorf("creating challenge %q: %w", name, err)
	}
	return &ch, nil
}

// DeleteChallenge deletes a challenge by name.
func (c *Client) DeleteChallenge(ctx context.Context, name string) error {
	if err := c.Delete(ctx, "/challenges/"+url.PathEscape(name), nil); err != nil {
		return fmt.Errorf("deleting challenge %q: %w", name, err)
	}
	return nil
}
