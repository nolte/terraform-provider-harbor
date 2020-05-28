package harbor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func init() {
	resource.AddTestSweepers("data_source_label", &resource.Sweeper{
		Name: "harbor_label",
	})
}
func TestAccHarborDataSourceLabel(t *testing.T) {
	var project models.Project

	var label models.Label

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckGlobalLabelDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLabelExists("harbor_label.label_global", &label),
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
			{
				Config: testAccHarborCheckProjectLabelDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(

					testAccHarborCheckProjectExists("harbor_project.project", &project),
					testAccCheckLabelExists("harbor_label.project_label", &label),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_p_1", "name", "project-label-acc-golang"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_p_1", "scope", "p"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_p_1", "description", "Golang Project Acc Test Label"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_p_1", "color", "#333333"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_p_2", "name", "project-label-acc-golang"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_p_2", "scope", "p"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_p_2", "description", "Golang Project Acc Test Label"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_p_2", "color", "#333333"),
				),
			},
		},
	})
}

func testAccHarborCheckGlobalLabelDataSourceConfig() string {
	return `
resource "harbor_label" "label_global" {
    name        = "global-acc-ds-golang"
    description = "Golang Global Acc Test Label"
    color       = "#333333"
    scope       = "g"
}

data "harbor_label" "ds_1" {
    name = harbor_label.label_global.name
    scope = "g"
  }
`
}

func testAccHarborCheckProjectLabelDataSourceConfig() string {
	return `
resource "harbor_project" "project" {
    name                   = "acc-project-label-test"
    public                 = false
    vulnerability_scanning = false
}

resource "harbor_label" "project_label" {
    name        = "project-label-acc-golang"
    description = "Golang Project Acc Test Label"
    color       = "#333333"
    scope       = "p"
    project_id = harbor_project.project.id
}
data "harbor_label" "ds_p_1" {
    id = harbor_label.project_label.id
}

data "harbor_label" "ds_p_2" {
    project_id = harbor_project.project.id
    name = harbor_label.project_label.name
    scope       = "p"
}
`
}
