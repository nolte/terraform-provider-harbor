package harbor

import (
	"log"

	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},

		Read: dataSourceProjectRead,
	}
}

func dataSourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)
	query := products.NewGetProjectsParams().WithName(d.Get("name").(*string))
	resp, err := apiClient.Products.GetProjects(query, nil)
	if err != nil {
		log.Fatal(err)
	}

	d.Set("project_id", resp.Payload[0].ProjectID)
	return nil
}
