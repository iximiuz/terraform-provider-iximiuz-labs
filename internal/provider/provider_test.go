// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"iximiuz-labs": providerserver.NewProtocol6WithError(New("test")()),
}

// testAccPreCheck validates that required environment variables or credentials
// are available before running acceptance tests. It checks for env vars first,
// then falls back to the labctl config file presence.
func testAccPreCheck(t *testing.T) {
	t.Helper()

	// Check for env-based credentials first
	sessionID := os.Getenv("IXIMIUZ_SESSION_ID")
	accessToken := os.Getenv("IXIMIUZ_ACCESS_TOKEN")

	if sessionID != "" && accessToken != "" {
		return // Env-based credentials are available
	}

	// Fall back: check if labctl config file exists
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Cannot determine home directory for labctl config fallback: %s", err)
	}

	configPath := home + "/.iximiuz/labctl/config.yaml"
	if _, err := os.Stat(configPath); err != nil {
		t.Skip("Acceptance tests require IXIMIUZ_SESSION_ID and IXIMIUZ_ACCESS_TOKEN environment variables, " +
			"or a labctl config file at ~/.iximiuz/labctl/config.yaml (run `labctl auth login`). Skipping.")
	}
}
