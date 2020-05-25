package harbor

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func robotAccountNamePrefix() string {
	return "robot$"
}

func resourceRobotAccount() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				ForceNew: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Create: resourceRobotAccountCreate,
		Read:   resourceRobotAccountRead,
		Delete: resourceRobotAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func factoryRobotDockerEndpoint(projectID int64) string {
	return "/project/" + strconv.FormatInt(projectID, 10) + "/repository"
}
func factoryRobotAccountAccessDockerRead(projectID int64) *models.RobotAccountAccess {
	return &models.RobotAccountAccess{
		Action:   "pull",
		Resource: factoryRobotDockerEndpoint(projectID),
	}
}
func factoryRobotAccountAccessDockerPush(projectID int64) *models.RobotAccountAccess {
	return &models.RobotAccountAccess{
		Action:   "push",
		Resource: factoryRobotDockerEndpoint(projectID),
	}
}

func factoryRobotAccountAccessHelmChartPush(projectID int64) *models.RobotAccountAccess {
	return &models.RobotAccountAccess{
		Action:   "create",
		Resource: "/project/" + strconv.FormatInt(projectID, 10) + "/helm-chart-version",
	}
}
func factoryRobotAccountAccessHelmChartRead(projectID int64) *models.RobotAccountAccess {
	return &models.RobotAccountAccess{
		Action:   "read",
		Resource: "/project/" + strconv.FormatInt(projectID, 10) + "/helm-chart",
	}
}

func resourceRobotAccountCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)
	projectid := int64(d.Get("project_id").(int))
	name := d.Get("name").(string)
	permissionsRoles := d.Get("actions").([]interface{})
	// create the permission list for the robot account
	robotAccountAccess := make([]*models.RobotAccountAccess, len(permissionsRoles))

	for i, role := range permissionsRoles {
		switch role.(string) {
		case "docker_read":
			robotAccountAccess[i] = (factoryRobotAccountAccessDockerRead(projectid))
		case "docker_write":
			robotAccountAccess[i] = (factoryRobotAccountAccessDockerPush(projectid))
		case "helm_read":
			robotAccountAccess[i] = (factoryRobotAccountAccessHelmChartRead(projectid))
		case "helm_write":
			robotAccountAccess[i] = (factoryRobotAccountAccessHelmChartPush(projectid))
		}
	}

	robotCreate := &models.RobotAccountCreate{
		Name:        name,
		Description: d.Get("description").(string),
		Access:      robotAccountAccess,
	}

	params := products.NewPostProjectsProjectIDRobotsParams().WithProjectID(projectid).WithRobot(robotCreate)

	resp, err := apiClient.Products.PostProjectsProjectIDRobots(params, nil)

	if err != nil {
		return err
	}

	if err := d.Set("token", resp.Payload.Token); err != nil {
		return err
	}

	robot, err := findRobotAccountByProjectAndName(d, m)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(int(robot.ID)))

	return resourceRobotAccountRead(d, m)
}

func findRobotAccountByProjectAndName(d *schema.ResourceData, m interface{}) (*models.RobotAccount, error) {
	apiClient := m.(*client.Harbor)
	name, nameOk := d.GetOk("name")
	projectID, projectIDOk := d.GetOk("project_id")

	if !nameOk || !projectIDOk {
		return &models.RobotAccount{}, errors.New("fail to get the name and/or project_id for robot account request")
	}

	query := products.NewGetProjectsProjectIDRobotsParams().WithProjectID(int64(projectID.(int)))

	resp, err := apiClient.Products.GetProjectsProjectIDRobots(query, nil)

	if err != nil {
		return &models.RobotAccount{}, err
	}

	for _, v := range resp.Payload {
		if v.Name == "robot$"+name.(string) {
			return v, nil
		}
	}

	return &models.RobotAccount{},
		fmt.Errorf("no Robot found for ProjectID %d, with Name %s", projectID.(int), name.(string))
}
func resourceRobotAccountRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)
	projectID := d.Get("project_id").(int)

	log.Printf("Load Robot Accounts from %v Project", projectID)

	if robotID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		query := products.NewGetProjectsProjectIDRobotsRobotIDParams().WithProjectID(int64(projectID)).WithRobotID(robotID)

		resp, err := apiClient.Products.GetProjectsProjectIDRobotsRobotID(query, nil)
		if err != nil {
			return err
		}

		if err := setRobotSchema(d, resp.Payload); err != nil {
			return err
		}
	}
	return nil
}
func setRobotSchema(d *schema.ResourceData, model *models.RobotAccount) error {
	d.SetId(strconv.Itoa(int(model.ID)))
	if err := d.Set("name", strings.Replace(model.Name, robotAccountNamePrefix(), "", 1)); err != nil {
		return err
	}
	if err := d.Set("description", model.Description); err != nil {
		return err
	}
	if err := d.Set("project_id", int(model.ProjectID)); err != nil {
		return err
	}

	if err := d.Set("disabled", model.Disabled); err != nil {
		return err
	}
	return nil
}

func resourceRobotAccountDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if robotID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		projectID := int64(d.Get("project_id").(int))
		params := products.NewDeleteProjectsProjectIDRobotsRobotIDParams().WithRobotID(robotID).WithProjectID(projectID)

		_, err := apiClient.Products.DeleteProjectsProjectIDRobotsRobotID(params, nil)
		if err != nil {
			return err
		}

		d.SetId("")
		return nil
	}
	return fmt.Errorf("Fail to Remove Robot Account")
}
