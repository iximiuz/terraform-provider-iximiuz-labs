// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTutorialResource_basic(t *testing.T) {
	resourceName := "iximiuz-labs_tutorial.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTutorialResourceConfig("tf-acc-tutorial-basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^tf-acc-tutorial-basic`)),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-tutorial-basic"),
					resource.TestCheckResourceAttrSet(resourceName, "title"),
					resource.TestCheckResourceAttrSet(resourceName, "page_url"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"name_prefix"},
				ImportStateIdFunc:                    testAccContentImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccTutorialResourceConfig(namePrefix string) string {
	return fmt.Sprintf(`
resource "iximiuz-labs_tutorial" "test" {
  name_prefix = %[1]q
}
`, namePrefix)
}
