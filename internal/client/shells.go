// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// ListShells lists all exposed shells for a play.
func (c *Client) ListShells(ctx context.Context, playID string) ([]Shell, error) {
	var shells []Shell
	path := "/plays/" + url.PathEscape(playID) + "/shells"
	if err := c.GetJSON(ctx, path, nil, &shells); err != nil {
		return nil, fmt.Errorf("listing shells for play %q: %w", playID, err)
	}
	return shells, nil
}

// ExposeShell exposes a shell on a play machine.
func (c *Client) ExposeShell(ctx context.Context, playID string, req ExposeShellRequest) (*Shell, error) {
	var shell Shell
	path := "/plays/" + url.PathEscape(playID) + "/shells"
	if err := c.PostJSON(ctx, path, req, &shell); err != nil {
		return nil, fmt.Errorf("exposing shell on play %q: %w", playID, err)
	}
	return &shell, nil
}

// UnexposeShell removes an exposed shell.
func (c *Client) UnexposeShell(ctx context.Context, playID string, shellID string) error {
	path := "/plays/" + url.PathEscape(playID) + "/shells/" + url.PathEscape(shellID)
	if err := c.Delete(ctx, path, nil); err != nil {
		return fmt.Errorf("unexposing shell %q on play %q: %w", shellID, playID, err)
	}
	return nil
}
