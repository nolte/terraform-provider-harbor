package harbor_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/nolte/terraform-provider-harbor/harbor"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = harbor.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"harbor": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := harbor.Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = harbor.Provider()
}

func testAccHarborPreCheck(t *testing.T) {
	if v := os.Getenv("HARBOR_ENDPOINT"); v == "" {
		t.Fatal("HARBOR_ENDPOINT must be set for acceptance tests")
	}

	if v := os.Getenv("HARBOR_USERNAME"); v == "" {
		t.Fatal("HARBOR_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("HARBOR_PASSWORD"); v == "" {
		t.Fatal("HARBOR_PASSWORD must be set for acceptance tests")
	}
}
