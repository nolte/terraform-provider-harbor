package harbor

import (
	"fmt"

	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceProject() *schema.Resource {
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

		Read: dataSourceProjectRead,
	}
}

func dataSourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)
	if _, ok := d.GetOk("name"); ok {
		project, err := findProjectByName(d, m)
		if err != nil {
			return err
		}
		if err = setProjectSchema(d, project); err != nil {
			return err
		}
		return nil
	}
	if id, ok := d.GetOk("id"); ok {
		if resp, err := apiClient.Products.GetProjectsProjectID(products.NewGetProjectsProjectIDParams().WithProjectID(int64(id.(int))), nil); err != nil {
			if err = setProjectSchema(d, resp.Payload); err != nil {
				return err
			}
			return nil
		}
	}
	d.SetId("")
	return fmt.Errorf("please specify a name to lookup for a project")
}
