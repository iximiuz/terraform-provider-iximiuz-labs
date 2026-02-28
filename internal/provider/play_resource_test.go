// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlayResource_basic(t *testing.T) {
	resourceName := "iximiuz-labs_play.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPlayResourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_in"),
					resource.TestCheckResourceAttrSet(resourceName, "page_url"),
				),
			},
			// ImportState testing
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"safety_disclaimer_consent", "expires_in", "updated_at"},
			},
		},
	})
}

func TestAccPlayResource_withConsent(t *testing.T) {
	resourceName := "iximiuz-labs_play.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlayResourceConfig_withConsent(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
				),
			},
		},
	})
}

// testAccPlayResourceConfig_basic creates a play using the built-in docker playground.
// We use the docker playground directly instead of creating a custom one to avoid
// cleanup issues (API won't let you delete a playground while a play is running).
func testAccPlayResourceConfig_basic() string {
	return `
resource "iximiuz-labs_play" "test" {
  playground = "docker"
}
`
}

func testAccPlayResourceConfig_withConsent() string {
	return `
resource "iximiuz-labs_play" "test" {
  playground                  = "docker"
  safety_disclaimer_consent   = true
}
`
}

// testAccPlayResourceConfig_forSubresources creates a play for
// sub-resource tests (port, shell). Uses the built-in docker playground
// to avoid cleanup issues.
func testAccPlayResourceConfig_forSubresources() string {
	return `
resource "iximiuz-labs_play" "test" {
  playground = "docker"
}
`
}
