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
	resource.AddTestSweepers("harbor_project", &resource.Sweeper{
		Name: "harbor_project",
	})
}

func TestAccHarborProject_Basic(t *testing.T) {
	var project models.Project

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckGlobalProjectResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckProjectExists("harbor_project.main", &project),
					resource.TestCheckResourceAttr(
						"harbor_project.main", "name", "acc-project-test"),
					resource.TestCheckResourceAttr(
						"harbor_project.main", "public", "false"),
					resource.TestCheckResourceAttr(
						"harbor_project.main", "vulnerability_scanning", "false"),
				),
			},
			{
				Config: testAccHarborCheckProjectActiveAllResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckProjectExists("harbor_project.main", &project),
					resource.TestCheckResourceAttr(
						"harbor_project.main", "name", "acc-project-minimal-test"),
					resource.TestCheckResourceAttr(
						"harbor_project.main", "public", "true"),
					resource.TestCheckResourceAttr(
						"harbor_project.main", "vulnerability_scanning", "true"),
				),
			},
		},
	})
}

func testAccHarborCheckGlobalProjectResourceConfig() string {
	return `
resource "harbor_project" "main" {
    name                   = "acc-project-test"
    public                 = false
    vulnerability_scanning = false
}
`
}
func testAccHarborCheckProjectActiveAllResourceConfig() string {
	return `
resource "harbor_project" "main" {
    name                   = "acc-project-minimal-test"
    public                 = true
    vulnerability_scanning = true
}
`
}

func testAccHarborCheckProjectExists(n string, project *models.Project) resource.TestCheckFunc {
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
			params := products.NewGetProjectsProjectIDParams().WithProjectID(searchID)

			foundProject, err := client.Products.GetProjectsProjectID(params, nil)
			if err != nil {
				return err
			}

			if foundProject == nil {
				return fmt.Errorf("Record not found")
			}

			*project = *foundProject.Payload
		}

		return nil
	}
}
