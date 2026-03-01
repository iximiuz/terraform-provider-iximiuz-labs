// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/iximiuz/terraform-provider-iximiuz-labs/internal/client"
)

// waitForPlayState polls the play API until the play reaches one of the target states.
func waitForPlayState(ctx context.Context, c *client.Client, playID string, targetStates []string, timeout time.Duration) (*client.Play, error) {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		play, err := c.GetPlay(ctx, playID)
		if err != nil {
			return nil, fmt.Errorf("polling play %q: %w", playID, err)
		}

		if play.Status != nil {
			currentState := play.Status.CurrentState()
			for _, target := range targetStates {
				if currentState == target {
					return play, nil
				}
			}

			// Check for terminal failure state
			if currentState == client.PlayStateFailed {
				return play, fmt.Errorf("play %q reached FAILED state", playID)
			}
		}

		if time.Now().After(deadline) {
			currentState := "unknown"
			if play.Status != nil {
				currentState = play.Status.CurrentState()
			}
			return play, fmt.Errorf("timeout waiting for play %q to reach state %v (current: %s)", playID, targetStates, currentState)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}
