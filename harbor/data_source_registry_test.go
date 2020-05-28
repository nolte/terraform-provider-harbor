package harbor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func init() {
	resource.AddTestSweepers("data_source_registry", &resource.Sweeper{
		Name: "harbor_registry",
	})
}
func TestAccHarborDataSourceRegistry(t *testing.T) {
	var registry models.Registry
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckGlobalRegistryDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckRegistryExists("harbor_registry.registry", &registry),
					resource.TestCheckResourceAttr(
						"data.harbor_registry.ds_1", "name", "acc-registry-ds-test"),
					resource.TestCheckResourceAttr(
						"data.harbor_registry.ds_1", "url", "https://hub.helm.sh"),
					resource.TestCheckResourceAttr(
						"data.harbor_registry.ds_1", "type", "helm-hub"),
					resource.TestCheckResourceAttr(
						"data.harbor_registry.ds_1", "insecure", "false"),
					resource.TestCheckResourceAttr(
						"data.harbor_registry.ds_2", "name", "acc-registry-ds-test"),
					resource.TestCheckResourceAttr(
						"data.harbor_registry.ds_2", "url", "https://hub.helm.sh"),
					resource.TestCheckResourceAttr(
						"data.harbor_registry.ds_2", "type", "helm-hub"),
					resource.TestCheckResourceAttr(
						"data.harbor_registry.ds_2", "insecure", "false"),
				),
			},
		},
	})
}

func testAccHarborCheckGlobalRegistryDataSourceConfig() string {
	return `
resource "harbor_registry" "registry" {
    name        = "acc-registry-ds-test"
    url         = "https://hub.helm.sh"
    type        = "helm-hub"
    description = "Helm Hub Registry"
    insecure    = false
}

data "harbor_registry" "ds_1" {
    id = harbor_registry.registry.id
}

data "harbor_registry" "ds_2" {
    name = harbor_registry.registry.name
}

`
}
