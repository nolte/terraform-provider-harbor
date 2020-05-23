package harbor

import (
	"fmt"
	"log"
	"strconv"

	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"

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
	if name, ok := d.GetOk("name"); ok {
		query := products.NewGetProjectsParams().WithName(name.(*string))
		resp, err := apiClient.Products.GetProjects(query, nil)

		if err != nil {
			d.SetId("")
			log.Fatal(err)
		}
		if len(resp.Payload) < 1 || resp.Payload[0].Name == name.(string) {
			return fmt.Errorf("no project found with name %v", name)
		}
		setProjectSchema(d, resp.Payload[0])
		return nil
	}

	return fmt.Errorf("please specify a name to lookup for a project")
}

func setProjectSchema(data *schema.ResourceData, project *models.Project) {
	data.SetId(strconv.Itoa(int(project.ProjectID)))
	data.Set("name", project.Name)
}
