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
	resource.AddTestSweepers("resource_harbor_usergroup", &resource.Sweeper{
		Name: "harbor_usergroup",
	})
}

func TestAccHarborUsergroup_Basic(t *testing.T) {
	var usergroup models.UserGroup

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckUsergroupResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckUsergroupExists("harbor_usergroup.main", &usergroup),
					resource.TestCheckResourceAttr(
						"harbor_usergroup.main", "name", "acc-usergroup-test"),
					resource.TestCheckResourceAttr(
						"harbor_usergroup.main", "ldap_dn", ""),
					resource.TestCheckResourceAttr(
						"harbor_usergroup.main", "type", "http"),
				),
			},
			// Can't run the Ldap test without a LDAP server configured.
			//			{
			//				Config: testAccHarborCheckUsergroupLdapResourceConfig(),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccHarborCheckUsergroupExists("harbor_usergroup.main", &usergroup),
			//					resource.TestCheckResourceAttr(
			//						"harbor_usergroup.main", "name", "acc-usergroup-test-2"),
			//					resource.TestCheckResourceAttr(
			//						"harbor_usergroup.main", "ldap_dn", "cn=developers,ou=groups,dc=example,dc=com"),
			//					resource.TestCheckResourceAttr(
			//						"harbor_usergroup.main", "type", "http"),
			//				),
			//			},
		},
	})
}

func testAccHarborCheckUsergroupResourceConfig() string {
	return `
resource "harbor_usergroup" "main" {
    name = "acc-usergroup-test"
    type = "http"
}
`
}

// Can't test LDAP type without an LDAP server or mock
func testAccHarborCheckUsergroupLdapResourceConfig() string {
	return `
resource "harbor_usergroup" "main" {
    name    = "acc-usergroup-test"
	ldap_dn = "cn=developers,ou=groups,dc=example,dc=com"
    type    = "ldap"
}
`
}

func testAccHarborCheckUsergroupExists(n string, usergroup *models.UserGroup) resource.TestCheckFunc {
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
			params := products.NewGetUsergroupsGroupIDParams().WithGroupID(searchID)

			foundUsergroup, err := client.Products.GetUsergroupsGroupID(params, nil)
			if err != nil {
				return err
			}

			if foundUsergroup == nil {
				return fmt.Errorf("Record not found")
			}

			*usergroup = *foundUsergroup.Payload
		}

		return nil
	}
}
