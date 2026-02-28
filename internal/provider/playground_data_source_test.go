// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlaygroundDataSource_basic(t *testing.T) {
	dataSourceName := "data.iximiuz-labs_playground.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaygroundDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "name", regexp.MustCompile(`^tf-acc-ds-pg`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "title"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "page_url"),
				),
			},
		},
	})
}

func testAccPlaygroundDataSourceConfig() string {
	return `
resource "iximiuz-labs_playground" "test" {
  name_prefix = "tf-acc-ds-pg"
  base        = "docker"
  title       = "Data Source Test Playground"
  description = "Playground for data source acceptance tests"
}

data "iximiuz-labs_playground" "test" {
  name = iximiuz-labs_playground.test.name
}
`
}
