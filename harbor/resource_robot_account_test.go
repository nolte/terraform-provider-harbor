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
	resource.AddTestSweepers("resource_robot_account", &resource.Sweeper{
		Name: "harbor_robot_account",
	})
}

func TestAccHarborRobot_Basic(t *testing.T) {
	var robot models.RobotAccount
	var project models.Project

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckRobotResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckProjectExists("harbor_project.main", &project),
					testAccHarborCheckRobotExists("harbor_robot_account.master_robot", &robot),
					resource.TestCheckResourceAttr(
						"harbor_robot_account.master_robot", "name", "acc-robot-test"),
					resource.TestCheckResourceAttr(
						"harbor_robot_account.master_robot", "description", "Robot account used created by gloang acc"),
					testAccHarborCheckRobotHasGeneratedTokenExists(
						"harbor_robot_account.master_robot"),
				),
			},
		},
	})
}

func testAccHarborCheckRobotHasGeneratedTokenExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if rs.Primary.Attributes["token"] == "" {
			return fmt.Errorf("No token generated")
		}
		return nil
	}
}

func testAccHarborCheckRobotResourceConfig() string {
	return `
    resource "harbor_project" "main" {
        name = "acc-robot-project"
    }

    resource "harbor_robot_account" "master_robot" {
      name        = "acc-robot-test"
      description = "Robot account used created by gloang acc"
      project_id  = harbor_project.main.id
      actions     = ["docker_read", "docker_write", "helm_read", "helm_write"]
    }
`
}

func testAccHarborCheckRobotExists(n string, robot *models.RobotAccount) resource.TestCheckFunc {
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
			if projectId, err := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64); err == nil {
				query := products.NewGetProjectsProjectIDRobotsRobotIDParams().WithProjectID(projectId).WithRobotID(searchID)
				foundRobot, err := client.Products.GetProjectsProjectIDRobotsRobotID(query, nil)
				if err != nil {
					return err
				}

				if foundRobot == nil {
					return fmt.Errorf("Record not found")
				}

				*robot = *foundRobot.Payload
			}
		}

		return nil
	}
}
