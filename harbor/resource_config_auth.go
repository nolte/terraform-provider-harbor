package harbor

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func resourceConfigAuth() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auth_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"oidc_groups_claim": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_verify_cert": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Create: resourceConfigAuthUpdate,
		Read:   resourceConfigAuthRead,
		Update: resourceConfigAuthUpdate,
		Delete: resourceConfigAuthDelete,
	}
}

// dasdas
func resourceConfigAuthRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)
	resp, err := apiClient.Products.GetConfigurations(products.NewGetConfigurationsParams(), nil)

	if err != nil {
		log.Fatal(err)
	}

	if err := d.Set("auth_mode", resp.Payload.AuthMode.Value); err != nil {
		return err
	}

	if err := d.Set("oidc_name", resp.Payload.OidcName.Value); err != nil {
		return err
	}

	if err := d.Set("oidc_endpoint", resp.Payload.OidcEndpoint.Value); err != nil {
		return err
	}

	if err := d.Set("oidc_client_id", resp.Payload.OidcClientID.Value); err != nil {
		return err
	}

	if nil != resp.Payload.OidcGroupsClaim {
		if err := d.Set("oidc_groups_claim", resp.Payload.OidcGroupsClaim.Value); err != nil {
			return err
		}
	}

	if err := d.Set("oidc_scope", resp.Payload.OidcScope.Value); err != nil {
		return err
	}

	if err := d.Set("oidc_verify_cert", resp.Payload.OidcVerifyCert.Value); err != nil {
		return err
	}

	return nil
}

func resourceConfigAuthUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient, body := newAPIClient(d, m)

	_, err := apiClient.Products.PutConfigurations(products.NewPutConfigurationsParams().WithConfigurations(&body), nil)
	if err != nil {
		return err
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", d.Get("oidc_name").(string))))

	return resourceConfigAuthRead(d, m)
}

func resourceConfigAuthDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func newAPIClient(d *schema.ResourceData, m interface{}) (*client.Harbor, models.Configurations) {
	apiClient := m.(*client.Harbor)

	body := models.Configurations{
		AuthMode:         d.Get("auth_mode").(string),
		OidcName:         d.Get("oidc_name").(string),
		OidcEndpoint:     d.Get("oidc_endpoint").(string),
		OidcClientID:     d.Get("oidc_client_id").(string),
		OidcClientSecret: d.Get("oidc_client_secret").(string),
		OidcGroupsClaim:  d.Get("oidc_groups_claim").(string),
		OidcScope:        d.Get("oidc_scope").(string),
		OidcVerifyCert:   d.Get("oidc_verify_cert").(bool),
	}

	return apiClient, body
}
