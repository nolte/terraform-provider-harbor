package harbor

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func resourceProject() *schema.Resource {
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
			"public": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  false,
			},
			"vulnerability_scanning": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
			},
		},
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := products.NewPostProjectsParams().WithProject(&models.ProjectReq{
		ProjectName: d.Get("name").(string),
		Metadata: &models.ProjectMetadata{
			AutoScan: d.Get("vulnerability_scanning").(string),
			Public:   d.Get("public").(string),
		},
	})

	_, err := apiClient.Client.Products.PostProjects(body, nil)
	if err != nil {
		log.Fatal(err)
	}
	//d.Set("project_id", resp.Payload[0].ProjectID)
	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", d.Get("name").(string))))
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectName := d.Get("name").(string)
	query := products.NewGetProjectsParams().WithName(&projectName)
	resp, err := apiClient.Client.Products.GetProjects(query, nil)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Payload) < 1 {
		d.SetId("")
		return nil
	}
	if err := d.Set("project_id", int(resp.Payload[0].ProjectID)); err != nil {
		return err
	}

	if err := d.Set("name", string(resp.Payload[0].Name)); err != nil {
		return err
	}

	if err := d.Set("vulnerability_scanning", string(resp.Payload[0].Metadata.AutoScan)); err != nil {
		return err
	}

	if err := d.Set("public", string(resp.Payload[0].Metadata.Public)); err != nil {
		return err
	}
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := products.NewPutProjectsProjectIDParams().WithProject(&models.ProjectReq{
		ProjectName: d.Get("name").(string),
		Metadata: &models.ProjectMetadata{
			AutoScan: d.Get("vulnerability_scanning").(string),
			Public:   d.Get("public").(string),
		},
	}).WithProjectID(int64(d.Get("project_id").(int)))

	_, err := apiClient.Client.Products.PutProjectsProjectID(body, nil)
	if err != nil {
		log.Fatal(err)
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectId := d.Get("project_id").(int)

	delete := products.NewDeleteProjectsProjectIDParams().WithProjectID(int64(projectId))
	_, err := apiClient.Client.Products.DeleteProjectsProjectID(delete, nil)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
