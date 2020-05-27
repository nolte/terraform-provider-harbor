package harbor_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func init() {
	resource.AddTestSweepers("harbor_label", &resource.Sweeper{
		Name: "harbor_label",
	})
}

func TestAccHcloudLabel_Basic(t *testing.T) {
	var label models.Label

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckGlobalLabelResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHcloudCheckLabelExists("harbor_label.main", &label),
					resource.TestCheckResourceAttr(
						"harbor_label.main", "name", "global-acc-golang"),
					resource.TestCheckResourceAttr(
						"harbor_label.main", "description", "Golang Acc Test Label"),
					resource.TestCheckResourceAttr(
						"harbor_label.main", "scope", "g"),
					resource.TestCheckResourceAttr(
						"harbor_label.main", "color", "#61717D"),
					resource.TestCheckResourceAttr(
						"harbor_label.main", "description", "Golang Acc Test Label"),
				),
			},
			{
				Config: testAccHarborCheckProjectLabelResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHcloudCheckLabelExists("harbor_label.main_project", &label),
					resource.TestCheckResourceAttr(
						"harbor_label.main_project", "name", "project-acc-golang"),
					resource.TestCheckResourceAttr(
						"harbor_label.main_project", "scope", "p"),
					resource.TestCheckResourceAttr(
						"harbor_label.main_project", "color", "#333333"),
					resource.TestCheckResourceAttr(
						"harbor_label.main_project", "description", "Golang Project Acc Test Label"),
				),
			},
		},
	})
}

func testAccHarborCheckGlobalLabelResourceConfig() string {
	return `
resource "harbor_label" "main" {
    name        = "global-acc-golang"
    description = "Golang Acc Test Label"
    color       = "#61717D"
}
`
}

func testAccHarborCheckProjectLabelResourceConfig() string {
	return `
resource "harbor_project" "main" {
  name                   = "acc-label-test"
  public                 = false
  vulnerability_scanning = false
}

resource "harbor_label" "main_project" {
    name        = "project-acc-golang"
    description = "Golang Project Acc Test Label"
    color       = "#333333"
    scope       = "p"
    project_id = harbor_project.main.id
}
`
}

func testAccHcloudCheckLabelExists(n string, label *models.Label) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*client.Harbor)

		if searchID, err := strconv.ParseInt(rs.Primary.ID, 10, 64); err == nil {
			foundLabel, err := client.Products.GetLabelsID(products.NewGetLabelsIDParams().WithID(searchID), nil)
			if err != nil {
				return err
			}

			if foundLabel == nil {
				return fmt.Errorf("Record not found")
			}

			*label = *foundLabel.Payload
		}

		return nil
	}
}
