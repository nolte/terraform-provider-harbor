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

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	apiClient := m.(*client.Harbor)

	projectName := d.Get("name").(string)
	body := products.NewPostProjectsParams().WithProject(&models.ProjectReq{
		ProjectName: projectName,
		Metadata: &models.ProjectMetadata{
			AutoScan: d.Get("vulnerability_scanning").(string),
			Public:   d.Get("public").(string),
		},
	})

	if _, err := apiClient.Products.PostProjects(body, nil); err != nil {
		return err
	}
	project, err := findProjectByName(d, m)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(int(project.ProjectID)))
	return resourceProjectRead(d, m)

}

func findProjectByName(d *schema.ResourceData, m interface{}) (*models.Project, error) {
	apiClient := m.(*client.Harbor)
	if name, ok := d.GetOk("name"); ok {
		projectName := name.(string)
		query := products.NewGetProjectsParams().WithName(&projectName)
		resp, err := apiClient.Products.GetProjects(query, nil)
		if err != nil {
			d.SetId("")
			return &models.Project{}, err
		}
		if len(resp.Payload) < 1 {
			return &models.Project{}, fmt.Errorf("no project found with name %v", projectName)
		} else if resp.Payload[0].Name != projectName {
			return &models.Project{}, fmt.Errorf("Response Name %v not match with Expected Name %v", resp.Payload[0].Name, projectName)
		}
		return resp.Payload[0], nil
	}
	return &models.Project{}, fmt.Errorf("Fail to lookup project by Name")
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if projectID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		query := products.NewGetProjectsProjectIDParams().WithProjectID(projectID)
		resp, err := apiClient.Products.GetProjectsProjectID(query, nil)
		if err != nil {
			return err
		}
		if err := setProjectSchema(d, resp.Payload); err != nil {
			return err
		}
		return nil
	}
	fmt.Println(d.Id(), "is not an integer.")
	return fmt.Errorf("fail to load the project")
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if projectID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		body := products.NewPutProjectsProjectIDParams().WithProject(&models.ProjectReq{
			ProjectName: d.Get("name").(string),
			Metadata: &models.ProjectMetadata{
				AutoScan: d.Get("vulnerability_scanning").(string),
				Public:   d.Get("public").(string),
			},
		}).WithProjectID(projectID)

		_, err := apiClient.Products.PutProjectsProjectID(body, nil)
		if err != nil {
			log.Fatal(err)
		}

		return resourceProjectRead(d, m)
	}
	return fmt.Errorf("Project Id not a Integer")

}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)
	if projectID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		delete := products.NewDeleteProjectsProjectIDParams().WithProjectID(projectID)
		_, err := apiClient.Products.DeleteProjectsProjectID(delete, nil)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}
	return fmt.Errorf("Project Id not a Integer")
}

func setProjectSchema(data *schema.ResourceData, project *models.Project) error {
	data.SetId(strconv.Itoa(int(project.ProjectID)))

	if err := data.Set("name", project.Name); err != nil {
		return err
	}

	if err := data.Set("vulnerability_scanning", string(project.Metadata.AutoScan)); err != nil {
		return err
	}

	if err := data.Set("public", string(project.Metadata.Public)); err != nil {
		return err
	}

	return nil
}
