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
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"vulnerability_scanning": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"reuse_sys_cve_whitelist": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"cve_whitelist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
		params := products.NewGetProjectsProjectIDParams().WithProjectID(int64(id.(int)))

		resp, err := apiClient.Products.GetProjectsProjectID(params, nil)
		if err != nil {
			return err
		}

		if err = setProjectSchema(d, resp.Payload); err != nil {
			return err
		}

		return nil
	}

	d.SetId("")

	return fmt.Errorf("please specify a name to lookup for a project")
}
