package harbor

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/nolte/terraform-provider-harbor/client"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_ENDPOINT", nil),
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https",
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
			"harbor_label":         resourceLabel(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"harbor_project":  dataSourceProject(),
			"harbor_registry": dataSourceRegistry(),
			"harbor_label":    dataSourceLabel(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	insecure := d.Get("insecure").(bool)
	basepath := d.Get("basepath").(string)
	schema := d.Get("schema").(string)

	return client.NewClient(host, username, password, insecure, basepath, schema), nil
}
