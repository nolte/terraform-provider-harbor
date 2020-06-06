package harbor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func init() {
	resource.AddTestSweepers("resource_harbor_replication", &resource.Sweeper{
		Name: "harbor_replication",
	})
}

func TestAccHarborReplication_Basic(t *testing.T) {
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
						"harbor_replication.pull_helm_chart", "name", "acc-helm-prometheus-operator-test"),
					resource.TestCheckResourceAttr(
						"harbor_replication.pull_helm_chart", "source_registry_filter_name", "stable/prometheus-operator"),
					resource.TestCheckResourceAttr(
						"harbor_replication.pull_helm_chart", "description", "Prometheus Operator Replica"),
					resource.TestCheckResourceAttr(
						"harbor_replication.pull_helm_chart", "destination_namespace", "acc-project-replica-test"),
				),
			}, {
				Config: testAccHarborCheckReplicationPushResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckRegistryExists("harbor_registry.registry_replica_push_helm_hub", &registry),
					resource.TestCheckResourceAttr(
						"harbor_replication.push_helm_chart", "name", "acc-push-test"),
					resource.TestCheckResourceAttr(
						"harbor_replication.push_helm_chart", "source_registry_filter_name", "stable/prometheus-operator"),
					resource.TestCheckResourceAttr(
						"harbor_replication.push_helm_chart", "description", "Push Replica"),
					resource.TestCheckResourceAttr(
						"harbor_replication.push_helm_chart", "destination_namespace", "notexistingtest"),
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
resource "harbor_replication" "pull_helm_chart" {
  name                        = "acc-helm-prometheus-operator-test"
  description                 = "Prometheus Operator Replica"
  source_registry_id          = harbor_registry.registry_replica_helm_hub.id
  source_registry_filter_name = "stable/prometheus-operator"
  source_registry_filter_tag  = "**"
  destination_namespace       = harbor_project.project_replica.name
}
`
}

func testAccHarborCheckReplicationPushResourceConfig() string {
	return `
resource "harbor_registry" "registry_replica_push_helm_hub" {
  name        = "acc-registry-push-replica-test"
  url         = "https://hub.helm.sh"
  type        = "helm-hub"
  description = "Helm Hub Registry"
  insecure    = false
}
resource "harbor_label" "main_push" {
    name        = "global-acc-push-golang"
    description = "Golang Acc Test Label"
    color       = "#61717D"
}
resource "harbor_replication" "push_helm_chart" {
  name                        = "acc-push-test"
  description                 = "Push Replica"
  source_registry_filter_name = "stable/prometheus-operator"
  source_registry_filter_tag  = "**"
  destination_registry_id = harbor_registry.registry_replica_push_helm_hub.id
  destination_namespace       = "notexistingtest"
}
`
}
