// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlayShellResource_basic(t *testing.T) {
	resourceName := "iximiuz-labs_play_shell.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlayShellResourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "play_id"),
					resource.TestCheckResourceAttr(resourceName, "machine", "docker-01"),
					resource.TestCheckResourceAttr(resourceName, "user", "root"),
					resource.TestCheckResourceAttr(resourceName, "access", "public"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
		},
	})
}

func testAccPlayShellResourceConfig_basic() string {
	return testAccPlayResourceConfig_forSubresources() + `
resource "iximiuz-labs_play_shell" "test" {
  play_id = iximiuz-labs_play.test.id
  machine = "docker-01"
  user    = "root"
  access  = "public"
}
`
}
