// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccChallengeResource_basic(t *testing.T) {
	resourceName := "iximiuz-labs_challenge.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccChallengeResourceConfig("tf-acc-challenge-basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^tf-acc-challenge-basic`)),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-challenge-basic"),
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

func TestAccChallengeResource_disappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChallengeResourceConfig("tf-acc-challenge-disappears"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("iximiuz-labs_challenge.test", "name", regexp.MustCompile(`^tf-acc-challenge-disappears`)),
				),
			},
		},
	})
}

func testAccChallengeResourceConfig(namePrefix string) string {
	return fmt.Sprintf(`
resource "iximiuz-labs_challenge" "test" {
  name_prefix = %[1]q
}
`, namePrefix)
}

// testAccContentImportStateIdFunc returns an ImportStateIdFunc that extracts
// the "name" attribute from state for content resources that use name as
// the import identifier (no separate "id" attribute).
func testAccContentImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		name := rs.Primary.Attributes["name"]
		if name == "" {
			return "", fmt.Errorf("name attribute not set for %s", resourceName)
		}
		return name, nil
	}
}
