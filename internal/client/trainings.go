// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// GetTraining gets a training by name.
func (c *Client) GetTraining(ctx context.Context, name string) (*Training, error) {
	var t Training
	if err := c.GetJSON(ctx, "/trainings/"+url.PathEscape(name), nil, &t); err != nil {
		return nil, fmt.Errorf("getting training %q: %w", name, err)
	}
	return &t, nil
}

// ListTrainings lists all trainings.
func (c *Client) ListTrainings(ctx context.Context) ([]Training, error) {
	var trainings []Training
	if err := c.GetJSON(ctx, "/trainings", nil, &trainings); err != nil {
		return nil, fmt.Errorf("listing trainings: %w", err)
	}
	return trainings, nil
}

// CreateTraining creates a new training.
func (c *Client) CreateTraining(ctx context.Context, name string) (*Training, error) {
	req := CreateContentRequest{Name: name}
	var t Training
	if err := c.PostJSON(ctx, "/trainings", req, &t); err != nil {
		return nil, fmt.Errorf("creating training %q: %w", name, err)
	}
	return &t, nil
}

// DeleteTraining deletes a training by name.
func (c *Client) DeleteTraining(ctx context.Context, name string) error {
	if err := c.Delete(ctx, "/trainings/"+url.PathEscape(name), nil); err != nil {
		return fmt.Errorf("deleting training %q: %w", name, err)
	}
	return nil
}
