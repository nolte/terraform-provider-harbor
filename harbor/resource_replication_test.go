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
	resource.AddTestSweepers("resource_harbor_replicationPull", &resource.Sweeper{
		Name: "harbor_replicationPull",
	})
}

func TestAccHarborReplicationPull_Basic(t *testing.T) {
	var project models.Project

	var registry models.Registry

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckReplicationPullResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckProjectExists("harbor_project.project_replica", &project),
					testAccHarborCheckRegistryExists("harbor_registry.registry_replica_helm_hub", &registry),
					resource.TestCheckResourceAttr(
						"harbor_replication_pull.pull_helm_chart", "name", "acc-helm-prometheus-operator-test"),
					resource.TestCheckResourceAttr(
						"harbor_replication_pull.pull_helm_chart", "source_registry_filter_name", "stable/prometheus-operator"),
					resource.TestCheckResourceAttr(
						"harbor_replication_pull.pull_helm_chart", "description", "Prometheus Operator Replica"),
					resource.TestCheckResourceAttr(
						"harbor_replication_pull.pull_helm_chart", "destination_namespace", "acc-project-replica-test"),
				),
			},
		},
	})
}

func testAccHarborCheckReplicationPullResourceConfig() string {
	return `
resource "harbor_project" "project_replica" {
    name                   = "acc-project-replica-test"
    public                 = false
    vulnerability_scanning = false
}
resource "harbor_registry" "registry_replica_helm_hub" {
  name        = "acc-registry-replica-test"
  url         = "https://hub.helm.sh"
  type        = "helm-hub"
  description = "Helm Hub Registry"
  insecure    = false
}
resource "harbor_replication_pull" "pull_helm_chart" {
  name                        = "acc-helm-prometheus-operator-test"
  description                 = "Prometheus Operator Replica"
  source_registry_id          = harbor_registry.registry_replica_helm_hub.id
  source_registry_filter_name = "stable/prometheus-operator"
  source_registry_filter_tag  = "**"
  destination_namespace       = harbor_project.project_replica.name
}
`
}

func testAccHarborCheckReplicationPullExists(n string, replicationPull *models.ReplicationPolicy) resource.TestCheckFunc {
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
			params := products.NewGetReplicationPoliciesIDParams().WithID(searchID)

			foundReplicationPull, err := client.Products.GetReplicationPoliciesID(params, nil)
			if err != nil {
				return err
			}

			if foundReplicationPull == nil {
				return fmt.Errorf("Record not found")
			}

			*replicationPull = *foundReplicationPull.Payload
		}

		return nil
	}
}
