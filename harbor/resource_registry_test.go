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
	resource.AddTestSweepers("resource_harbor_registry", &resource.Sweeper{
		Name: "harbor_registry",
	})
}

func TestAccHarborRegistry_Basic(t *testing.T) {
	var registry models.Registry

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckRegistryResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckRegistryExists("harbor_registry.helmhub", &registry),
					resource.TestCheckResourceAttr(
						"harbor_registry.helmhub", "name", "acc-registry-test"),
					resource.TestCheckResourceAttr(
						"harbor_registry.helmhub", "url", "https://hub.helm.sh"),
					resource.TestCheckResourceAttr(
						"harbor_registry.helmhub", "description", "Helm Hub Registry"),
					resource.TestCheckResourceAttr(
						"harbor_registry.helmhub", "insecure", "false"),
				),
			},
		},
	})
}

func testAccHarborCheckRegistryResourceConfig() string {
	return `
resource "harbor_registry" "helmhub" {
    name        = "acc-registry-test"
    url         = "https://hub.helm.sh"
    type        = "helm-hub"
    description = "Helm Hub Registry"
    insecure    = false
}
`
}

func testAccHarborCheckRegistryExists(n string, registry *models.Registry) resource.TestCheckFunc {
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
			params := products.NewGetRegistriesIDParams().WithID(searchID)

			foundRegistry, err := client.Products.GetRegistriesID(params, nil)
			if err != nil {
				return err
			}

			if foundRegistry == nil {
				return fmt.Errorf("Record not found")
			}

			*registry = *foundRegistry.Payload
		}

		return nil
	}
}
