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

func resourceConfigSystem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_creation_restriction": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "adminonly",
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"robot_token_expiration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},
		Create: resourceConfigSystemUpdate,
		Read:   resourceConfigSystemRead,
		Update: resourceConfigSystemUpdate,
		Delete: resourceConfigSystemDelete,
	}
}

func resourceConfigSystemRead(d *schema.ResourceData, m interface{}) error {

	apiClient := m.(*client.Client)
	resp, err := apiClient.Client.Products.GetConfigurations(products.NewGetConfigurationsParams(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// some internal harbor process will be change the value to project_creation_restriction = "oidc_auth"
	//if err := d.Set("project_creation_restriction", resp.Payload.AuthMode.Value); err != nil {
	//	return err
	//}
	if nil != resp.Payload.ReadOnly {
		if err := d.Set("read_only", resp.Payload.ReadOnly.Value); err != nil {
			return err
		}
	} else {
		d.Set("read_only", false)
	}
	if err := d.Set("robot_token_expiration", resp.Payload.TokenExpiration.Value); err != nil {
		return err
	}
	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", "systemconfig")))
	return nil
}

func resourceConfigSystemUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	_, err := apiClient.Client.Products.PutConfigurations(products.NewPutConfigurationsParams().WithConfigurations(&models.Configurations{
		ProjectCreationRestriction: d.Get("project_creation_restriction").(string),
		TokenExpiration:            int64(d.Get("robot_token_expiration").(int)),
		ReadOnly:                   d.Get("read_only").(bool),
	}), nil)
	if err != nil {
		log.Fatal(err)
	}

	return resourceConfigSystemRead(d, m)
}

func resourceConfigSystemDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
