// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlaygroundResource_basic(t *testing.T) {
	resourceName := "iximiuz-labs_playground.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPlaygroundResourceConfig_basic("tf-acc-pg-basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^tf-acc-pg-basic`)),
					resource.TestCheckResourceAttr(resourceName, "title", "TF Acc Test Playground"),
					resource.TestCheckResourceAttr(resourceName, "description", "Acceptance test playground"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "page_url"),
				),
			},
			// ImportState testing
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "network", "machine", "tab", "access_control"},
				ImportStateIdFunc:       testAccContentImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccPlaygroundResource_full(t *testing.T) {
	resourceName := "iximiuz-labs_playground.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaygroundResourceConfig_full("tf-acc-pg-full"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^tf-acc-pg-full`)),
					resource.TestCheckResourceAttr(resourceName, "title", "Full TF Acc Test"),
					resource.TestCheckResourceAttr(resourceName, "description", "Full acceptance test"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "page_url"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "network", "machine", "tab", "access_control"},
				ImportStateIdFunc:       testAccContentImportStateIdFunc(resourceName),
			},
		},
	})
}

// testAccPlaygroundResourceConfig_basic creates a minimal playground extending the
// docker base. No machine block is needed — the base's machine is inherited.
func testAccPlaygroundResourceConfig_basic(namePrefix string) string {
	return fmt.Sprintf(`
resource "iximiuz-labs_playground" "test" {
  name_prefix = %[1]q
  base        = "docker"
  title       = "TF Acc Test Playground"
  description = "Acceptance test playground"
}
`, namePrefix)
}

func testAccPlaygroundResourceConfig_full(namePrefix string) string {
	return fmt.Sprintf(`
resource "iximiuz-labs_playground" "test" {
  name_prefix = %[1]q
  base        = "docker"
  title       = "Full TF Acc Test"
  description = "Full acceptance test"

  tab {
    kind    = "terminal"
    name    = "Extra Terminal"
    machine = "docker-01"
  }
}
`, namePrefix)
}
