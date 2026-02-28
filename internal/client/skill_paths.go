// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// GetSkillPath gets a skill path by name.
func (c *Client) GetSkillPath(ctx context.Context, name string) (*SkillPath, error) {
	var sp SkillPath
	if err := c.GetJSON(ctx, "/skill-paths/"+url.PathEscape(name), nil, &sp); err != nil {
		return nil, fmt.Errorf("getting skill path %q: %w", name, err)
	}
	return &sp, nil
}

// ListSkillPaths lists all skill paths.
func (c *Client) ListSkillPaths(ctx context.Context) ([]SkillPath, error) {
	var paths []SkillPath
	if err := c.GetJSON(ctx, "/skill-paths", nil, &paths); err != nil {
		return nil, fmt.Errorf("listing skill paths: %w", err)
	}
	return paths, nil
}

// CreateSkillPath creates a new skill path.
func (c *Client) CreateSkillPath(ctx context.Context, name string) (*SkillPath, error) {
	req := CreateContentRequest{Name: name}
	var sp SkillPath
	if err := c.PostJSON(ctx, "/skill-paths", req, &sp); err != nil {
		return nil, fmt.Errorf("creating skill path %q: %w", name, err)
	}
	return &sp, nil
}

// DeleteSkillPath deletes a skill path by name.
func (c *Client) DeleteSkillPath(ctx context.Context, name string) error {
	if err := c.Delete(ctx, "/skill-paths/"+url.PathEscape(name), nil); err != nil {
		return fmt.Errorf("deleting skill path %q: %w", name, err)
	}
	return nil
}
