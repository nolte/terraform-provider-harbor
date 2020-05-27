package harbor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("data_source_datacenter", &resource.Sweeper{
		Name: "harbor_label",
	})
}
func TestAccHarborDataSourceLabel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckDatacenterDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_1", "id", "27"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_1", "name", "testlabel2"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_1", "scope", "g"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_1", "description", "Test Label"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_2", "id", "20"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_2", "name", "foooba"),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_2", "description", ""),
					resource.TestCheckResourceAttr(
						"data.harbor_label.ds_2", "scope", "g"),
				),
			},
		},
	})
}

func testAccHarborCheckDatacenterDataSourceConfig() string {
	return `
data "harbor_label" "ds_1" {
  name = "testlabel2"
  scope = "g"
}
data "harbor_label" "ds_2" {
  id = 20
}
`
}
