package harbor

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

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
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pull",
				ForceNew: true,
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"robot_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Create: resourceRobotAccountCreate,
		Read:   resourceRobotAccountRead,
		Delete: resourceRobotAccountDelete,
	}
}

func resourceRobotAccountCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectid := int64(d.Get("project_id").(int))
	name := d.Get("name").(string)
	resource := "/project/" + strconv.FormatInt(projectid, 10) + "/repository"

	resp, err := apiClient.Client.Products.PostProjectsProjectIDRobots(products.NewPostProjectsProjectIDRobotsParams().WithProjectID(projectid).WithRobot(&models.RobotAccountCreate{
		Name:        name,
		Description: d.Get("description").(string),
		Access: []*models.RobotAccountAccess{
			{
				Action:   d.Get("action").(string),
				Resource: resource,
			},
		},
	}), nil)

	if err != nil {
		log.Fatal(err)
	}
	if err := d.Set("token", resp.Payload.Token); err != nil {
		return err
	}

	return resourceRobotAccountRead(d, m)
}

func resourceRobotAccountRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectid := int64(d.Get("project_id").(int))
	name := d.Get("name").(string)
	log.Printf("Load Robot Accounts from %i Project", projectid)
	resp, err := apiClient.Client.Products.GetProjectsProjectIDRobots(products.NewGetProjectsProjectIDRobotsParams().WithProjectID(projectid), nil)

	if err != nil {
		log.Fatal(err)
	}
	for _, v := range resp.Payload {
		if v.Name == "robot$"+name {
			d.Set("robot_id", int(v.ID))
			d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", v.Name)))
		}
	}
	return nil
}

func resourceRobotAccountDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	robotID := int64(d.Get("robot_id").(int))
	projectID := int64(d.Get("project_id").(int))

	_, err := apiClient.Client.Products.DeleteProjectsProjectIDRobotsRobotID(products.NewDeleteProjectsProjectIDRobotsRobotIDParams().WithRobotID(robotID).WithProjectID(projectID), nil)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
