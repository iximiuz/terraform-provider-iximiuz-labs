// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// ListPorts lists all exposed ports for a play.
func (c *Client) ListPorts(ctx context.Context, playID string) ([]Port, error) {
	var ports []Port
	path := "/plays/" + url.PathEscape(playID) + "/ports"
	if err := c.GetJSON(ctx, path, nil, &ports); err != nil {
		return nil, fmt.Errorf("listing ports for play %q: %w", playID, err)
	}
	return ports, nil
}

// ExposePort exposes a port on a play machine.
func (c *Client) ExposePort(ctx context.Context, playID string, req ExposePortRequest) (*Port, error) {
	var port Port
	path := "/plays/" + url.PathEscape(playID) + "/ports"
	if err := c.PostJSON(ctx, path, req, &port); err != nil {
		return nil, fmt.Errorf("exposing port on play %q: %w", playID, err)
	}
	return &port, nil
}

// UnexposePort removes an exposed port.
func (c *Client) UnexposePort(ctx context.Context, playID string, portID string) error {
	path := "/plays/" + url.PathEscape(playID) + "/ports/" + url.PathEscape(portID)
	if err := c.Delete(ctx, path, nil); err != nil {
		return fmt.Errorf("unexposing port %q on play %q: %w", portID, playID, err)
	}
	return nil
}
