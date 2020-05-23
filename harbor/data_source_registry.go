package harbor

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/client"
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
		Read: dataSourceProjectRead,
	}
}

func dataSourceRegistryRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	query := products.NewGetRegistriesParams().WithName(d.Get("name").(*string))
	resp, err := apiClient.Client.Products.GetRegistries(query, nil)
	if err != nil {
		log.Fatal(err)
	}

	d.Set("repository_id", resp.Payload[0].ID)

	return nil
}
