package harbor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("data_source_label", &resource.Sweeper{
		Name: "harbor_label",
	})
}
func TestAccHarborDataSourceLabel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckGlobalLabelDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_1", "name", "global-acc-ds-golang"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_1", "scope", "g"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_1", "description", "Golang Global Acc Test Label"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_1", "color", "#333333"),
				),
			},
			//{
			//	Config: testAccHarborCheckProjectLabelDataSourceConfig(),
			//	Check: resource.ComposeTestCheckFunc(
			//		//testAccHarborCheckProjectExists("harbor_project.project_label", &project),
			//		//testAccCheckLabelExists("harbor_label.project_label", &label),
			//		resource.TestCheckResourceAttr(
			//			"data.harbor_label.ds_2", "name", "project-acc-golang"),
			//		resource.TestCheckResourceAttr(
			//			"data.harbor_label.ds_2", "scope", "p"),
			//		resource.TestCheckResourceAttr(
			//			"data.harbor_label.ds_2", "description", "Golang Project Acc Test Label"),
			//		resource.TestCheckResourceAttr(
			//			"data.harbor_label.ds_2", "color", "#666666"),
			//	),
			//},
		},
	})
}

func testAccHarborCheckGlobalLabelDataSourceConfig() string {
	return `
resource "harbor_label" "main_project" {
    name        = "global-acc-ds-golang"
    description = "Golang Global Acc Test Label"
    color       = "#333333"
    scope       = "g"
}

data "harbor_label" "ds_1" {
    name = harbor_label.main_project.name
    scope = "g"
  }
`
}

func testAccHarborCheckProjectLabelDataSourceConfig() string {
	return `
resource "harbor_project" "project" {
  name                   = "acc-label-test"
  public                 = false
  vulnerability_scanning = false
}
resource "harbor_label" "project_label" {
    name        = "project-acc-golang"
    description = "Golang Project Acc Test Label"
    color       = "#666666"
    scope       = "p"
    project_id = harbor_project.project.id
}
data "harbor_label" "ds_2" {
  name = harbor_label.project_label.name
  scope = "p"
}
`
}
