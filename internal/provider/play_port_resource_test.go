// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlayPortResource_basic(t *testing.T) {
	resourceName := "iximiuz-labs_play_port.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlayPortResourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "play_id"),
					resource.TestCheckResourceAttr(resourceName, "machine", "docker-01"),
					resource.TestCheckResourceAttr(resourceName, "number", "8080"),
					resource.TestCheckResourceAttr(resourceName, "access", "public"),
					resource.TestCheckResourceAttr(resourceName, "tls", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
		},
	})
}

func TestAccPlayPortResource_privateDefault(t *testing.T) {
	resourceName := "iximiuz-labs_play_port.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlayPortResourceConfig_private(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "access", "private"),
				),
			},
		},
	})
}

func testAccPlayPortResourceConfig_basic() string {
	return testAccPlayResourceConfig_forSubresources() + `
resource "iximiuz-labs_play_port" "test" {
  play_id = iximiuz-labs_play.test.id
  machine = "docker-01"
  number  = 8080
  access  = "public"
  tls     = true
}
`
}

func testAccPlayPortResourceConfig_private() string {
	return testAccPlayResourceConfig_forSubresources() + `
resource "iximiuz-labs_play_port" "test" {
  play_id = iximiuz-labs_play.test.id
  machine = "docker-01"
  number  = 3000
}
`
}
