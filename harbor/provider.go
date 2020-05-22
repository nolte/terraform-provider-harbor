package harbor

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/nolte/terraform-provider-harbor/client"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"basepath": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/api",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"harbor_project":       resourceProject(),
			"harbor_registry":      resourceRegistry(),
			"harbor_config_email":  resourceConfigEmail(),
			"harbor_config_auth":   resourceConfigAuth(),
			"harbor_config_system": resourceConfigSystem(),
			"harbor_robot_account": resourceRobotAccount(),
			"harbor_tasks":         resourceTasks(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"harbor_project": dataSourceProject(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	insecure := d.Get("insecure").(bool)
	basepath := d.Get("basepath").(string)

	return client.NewClient(url, username, password, insecure, basepath), nil
}
