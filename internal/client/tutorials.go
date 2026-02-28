// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// GetTutorial gets a tutorial by name.
func (c *Client) GetTutorial(ctx context.Context, name string) (*Tutorial, error) {
	var t Tutorial
	if err := c.GetJSON(ctx, "/tutorials/"+url.PathEscape(name), nil, &t); err != nil {
		return nil, fmt.Errorf("getting tutorial %q: %w", name, err)
	}
	return &t, nil
}

// ListTutorials lists all tutorials.
func (c *Client) ListTutorials(ctx context.Context) ([]Tutorial, error) {
	var tutorials []Tutorial
	if err := c.GetJSON(ctx, "/tutorials", nil, &tutorials); err != nil {
		return nil, fmt.Errorf("listing tutorials: %w", err)
	}
	return tutorials, nil
}

// CreateTutorial creates a new tutorial.
func (c *Client) CreateTutorial(ctx context.Context, name string) (*Tutorial, error) {
	req := CreateContentRequest{Name: name}
	var t Tutorial
	if err := c.PostJSON(ctx, "/tutorials", req, &t); err != nil {
		return nil, fmt.Errorf("creating tutorial %q: %w", name, err)
	}
	return &t, nil
}

// DeleteTutorial deletes a tutorial by name.
func (c *Client) DeleteTutorial(ctx context.Context, name string) error {
	if err := c.Delete(ctx, "/tutorials/"+url.PathEscape(name), nil); err != nil {
		return fmt.Errorf("deleting tutorial %q: %w", name, err)
	}
	return nil
}
