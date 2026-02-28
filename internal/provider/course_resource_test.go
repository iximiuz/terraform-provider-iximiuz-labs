// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCourseResource_basic(t *testing.T) {
	resourceName := "iximiuz-labs_course.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCourseResourceConfig_basic("tf-acc-course-basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^tf-acc-course-basic`)),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-course-basic"),
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

func TestAccCourseResource_full(t *testing.T) {
	resourceName := "iximiuz-labs_course.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with all attributes
			{
				Config: testAccCourseResourceConfig_full("tf-acc-course-full", "modular", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^tf-acc-course-full`)),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-course-full"),
					resource.TestCheckResourceAttr(resourceName, "variant", "modular"),
					resource.TestCheckResourceAttr(resourceName, "sample", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "title"),
					resource.TestCheckResourceAttrSet(resourceName, "page_url"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			// ImportState testing -- variant and sample are write-only (not returned by API read),
			// so we skip verifying them on import.
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"name_prefix", "variant", "sample"},
				ImportStateIdFunc:                    testAccContentImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCourseResourceConfig_basic(namePrefix string) string {
	return fmt.Sprintf(`
resource "iximiuz-labs_course" "test" {
  name_prefix = %[1]q
}
`, namePrefix)
}

func testAccCourseResourceConfig_full(namePrefix, variant string, sample bool) string {
	return fmt.Sprintf(`
resource "iximiuz-labs_course" "test" {
  name_prefix = %[1]q
  variant     = %[2]q
  sample      = %[3]t
}
`, namePrefix, variant, sample)
}
