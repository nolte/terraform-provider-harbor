package harbor

import (
	"fmt"
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
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"vulnerability_scanning": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	projectName := d.Get("name").(string)
	body := products.NewPostProjectsParams().WithProject(&models.ProjectReq{
		ProjectName: projectName,
		Metadata: &models.ProjectMetadata{
			AutoScan:             strconv.FormatBool(d.Get("vulnerability_scanning").(bool)),
			Public:               strconv.FormatBool(d.Get("public").(bool)),
			ReuseSysCveWhitelist: strconv.FormatBool(d.Get("reuse_sys_cve_whitelist").(bool)),
		},
		CveWhitelist: &models.CVEWhitelist{
			Items: expandCveWhitelist(d.Get("cve_whitelist").([]interface{})),
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

	if len(d.Get("cve_whitelist").([]interface{})) > 0 {
		// There is a bug in the Harbor API where the CVE Whitelist does not persist during Project creation
		err = resourceProjectUpdate(d, m)
		if err != nil {
			return err
		}
	}

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
			return &models.Project{},
				fmt.Errorf("response Name %v not match with Expected Name %v", resp.Payload[0].Name, projectName)
		}

		return resp.Payload[0], nil
	}

	return &models.Project{}, fmt.Errorf("fail to lookup project by Name")
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

	return fmt.Errorf("fail to load the project")
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if projectID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		body := products.NewPutProjectsProjectIDParams().WithProject(&models.ProjectReq{
			ProjectName: d.Get("name").(string),
			Metadata: &models.ProjectMetadata{
				AutoScan:             strconv.FormatBool(d.Get("vulnerability_scanning").(bool)),
				Public:               strconv.FormatBool(d.Get("public").(bool)),
				ReuseSysCveWhitelist: strconv.FormatBool(d.Get("reuse_sys_cve_whitelist").(bool)),
			},
			CveWhitelist: &models.CVEWhitelist{
				Items: expandCveWhitelist(d.Get("cve_whitelist").([]interface{})),
			},
		}).WithProjectID(projectID)

		if _, err := apiClient.Products.PutProjectsProjectID(body, nil); err != nil {
			return err
		}

		return resourceProjectRead(d, m)
	}

	return fmt.Errorf("project Id not a Integer")
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if projectID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		delete := products.NewDeleteProjectsProjectIDParams().WithProjectID(projectID)
		if _, err := apiClient.Products.DeleteProjectsProjectID(delete, nil); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("project Id not a Integer")
}

func setProjectSchema(data *schema.ResourceData, project *models.Project) error {
	data.SetId(strconv.Itoa(int(project.ProjectID)))

	if err := data.Set("name", project.Name); err != nil {
		return err
	}

	// prevent errors where auto_scan is unset
	if project.Metadata.AutoScan != "" {
		autoScan, err := strconv.ParseBool(project.Metadata.AutoScan)
		if err != nil {
			return err
		}

		if err := data.Set("vulnerability_scanning", autoScan); err != nil {
			return err
		}
	}

	public, err := strconv.ParseBool(project.Metadata.Public)
	if err != nil {
		return err
	}

	if err := data.Set("public", public); err != nil {
		return err
	}

	if project.Metadata.ReuseSysCveWhitelist != "" {
		reuseSysCveWhitelist, err := strconv.ParseBool(project.Metadata.ReuseSysCveWhitelist)
		if err != nil {
			return err
		}

		if err := data.Set("reuse_sys_cve_whitelist", reuseSysCveWhitelist); err != nil {
			return err
		}
	}

	if err := data.Set("cve_whitelist", flattenCveWhitelist(data, project.CveWhitelist)); err != nil {
		return err
	}

	return nil
}

func expandCveWhitelist(cveWhitelist []interface{}) []*models.CVEWhitelistItem {
	var cveWhitelistItems []*models.CVEWhitelistItem

	for _, cve := range cveWhitelist {
		cveWhitelistItem := &models.CVEWhitelistItem{
			CveID: cve.(string),
		}
		cveWhitelistItems = append(cveWhitelistItems, cveWhitelistItem)
	}

	return cveWhitelistItems
}

func flattenCveWhitelist(data *schema.ResourceData, CveWhitelist *models.CVEWhitelist) []interface{} {
	var cveWhitelist []interface{}

	for _, cve := range CveWhitelist.Items {
		if cve.CveID != "" {
			cveWhitelist = append(cveWhitelist, cve.CveID)
		}
	}

	return cveWhitelist
}
