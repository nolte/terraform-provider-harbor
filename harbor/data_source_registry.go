package harbor

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
)

func dataSourceRegistry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Read: dataSourceRegistryRead,
	}
}

func dataSourceRegistryRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	query := products.NewGetRegistriesParams().WithName(d.Get("name").(*string))
	resp, err := apiClient.Products.GetRegistries(query, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	d.Set("repository_id", resp.Payload[0].ID)

	return nil
}
