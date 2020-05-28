package harbor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func init() {
	resource.AddTestSweepers("data_source_project", &resource.Sweeper{
		Name: "harbor_project",
	})
}
func TestAccHarborDataSourceProject(t *testing.T) {
	var project models.Project

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckGlobalProjectDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckProjectExists("harbor_project.project", &project),
					resource.TestCheckResourceAttr(
						"data.harbor_project.ds_1", "name", "acc-project-ds-test"),
					resource.TestCheckResourceAttr(
						"data.harbor_project.ds_1", "vulnerability_scanning", "false"),
					resource.TestCheckResourceAttr(
						"data.harbor_project.ds_1", "public", "false"),
					resource.TestCheckResourceAttr(
						"data.harbor_project.ds_2", "name", "acc-project-ds-test"),
					resource.TestCheckResourceAttr(
						"data.harbor_project.ds_2", "vulnerability_scanning", "false"),
					resource.TestCheckResourceAttr(
						"data.harbor_project.ds_2", "public", "false"),
				),
			},
		},
	})
}

func testAccHarborCheckGlobalProjectDataSourceConfig() string {
	return `
resource "harbor_project" "project" {
    name                   = "acc-project-ds-test"
    public                 = false
    vulnerability_scanning = false
}

data "harbor_project" "ds_1" {
    id = harbor_project.project.id
}

data "harbor_project" "ds_2" {
    name = harbor_project.project.name
}
`
}
