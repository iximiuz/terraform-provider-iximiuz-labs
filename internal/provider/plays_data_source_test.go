// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlaysDataSource_basic(t *testing.T) {
	dataSourceName := "data.iximiuz-labs_plays.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaysDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "plays.#"),
				),
			},
		},
	})
}

func testAccPlaysDataSourceConfig() string {
	return `
data "iximiuz-labs_plays" "test" {
}
`
}
