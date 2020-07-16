package harbor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func init() {
	resource.AddTestSweepers("resource_harbor_project_member", &resource.Sweeper{
		Name: "harbor_project_member",
	})
}

func TestAccHarborProjectMember_Basic(t *testing.T) {
	var project models.Project

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckProjectMemberResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckProjectExists("harbor_project.project_member", &project),
					resource.TestCheckResourceAttr(
						"harbor_project_member.group_member", "role", "limited_guest"),
					resource.TestCheckResourceAttr(
						"harbor_project_member.group_member", "group_type", "http"),
					resource.TestCheckResourceAttr(
						"harbor_project_member.group_member", "group_name", "example_group"),
				),
			},
		},
	})
}

func testAccHarborCheckProjectMemberResourceConfig() string {
	return `
resource "harbor_project" "project_member" {
    name                   = "acc-project-member-test"
    public                 = false
    vulnerability_scanning = false
}
resource "harbor_usergroup" "project_group" {
	name = "example_group"
	type = "http"
}
resource "harbor_project_member" "group_member" {
  project_id = harbor_project.project_member.id
  role       = "limited_guest"
  group_type = "http"
  group_name = harbor_usergroup.project_group.name
}
`
}
