package harbor

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func dataSourceRegistry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
		Read: dataSourceRegistryRead,
	}
}

func dataSourceRegistryRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if name, ok := d.GetOk("name"); ok {
		query := products.NewGetRegistriesParams().WithName(name.(*string))
		resp, err := apiClient.Products.GetRegistries(query, nil)

		if err != nil {
			d.SetId("")
			log.Fatal(err)
		}
		if len(resp.Payload) < 1 || resp.Payload[0].Name == name.(string) {
			return fmt.Errorf("no registries found with name %v", name)
		}
		setRegistrySchema(d, resp.Payload[0])
		return nil
	}
	if id, ok := d.GetOk("id"); ok {
		resp, err := apiClient.Products.GetRegistriesID(products.NewGetRegistriesIDParams().WithID(int64(id.(int))), nil)
		if err != nil {
			d.SetId("")
			log.Fatal(err)
		}
		setRegistrySchema(d, resp.Payload)
	}

	return fmt.Errorf("please specify a name to lookup for a registries")
}
func setRegistrySchema(data *schema.ResourceData, project *models.Registry) {
	data.SetId(strconv.Itoa(int(project.ID)))
	data.Set("name", project.Name)
}
