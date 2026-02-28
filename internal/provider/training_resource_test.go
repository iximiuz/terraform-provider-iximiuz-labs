// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrainingResource_basic(t *testing.T) {
	resourceName := "iximiuz-labs_training.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrainingResourceConfig("tf-acc-training-basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^tf-acc-training-basic`)),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-training-basic"),
					resource.TestCheckResourceAttrSet(resourceName, "title"),
					resource.TestCheckResourceAttrSet(resourceName, "page_url"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
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

func testAccTrainingResourceConfig(namePrefix string) string {
	return fmt.Sprintf(`
resource "iximiuz-labs_training" "test" {
  name_prefix = %[1]q
}
`, namePrefix)
}
