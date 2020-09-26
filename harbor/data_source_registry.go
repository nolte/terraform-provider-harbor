package harbor

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
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
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"credential_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Read: dataSourceRegistryRead,
	}
}

func dataSourceRegistryRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if _, ok := d.GetOk("name"); ok {
		registry, err := findRegistryByName(d, m)
		if err != nil {
			return err
		}

		if err := setRegistrySchema(d, registry); err != nil {
			return err
		}

		return nil
	}

	if id, ok := d.GetOk("id"); ok {
		resp, err := apiClient.Products.GetRegistriesID(products.NewGetRegistriesIDParams().WithID(int64(id.(int))), nil)
		if err != nil {
			return err
		}

		if err := setRegistrySchema(d, resp.Payload); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("please specify a name to lookup for a registries")
}
