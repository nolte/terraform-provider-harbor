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

func resourceConfigEmail() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"email_host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"email_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"email_from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// "email_verify_cert": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
		},
		Create: resourceConfigEmailCreate,
		Read:   resourceConfigEmailRead,
		Update: resourceConfigEmailCreate,
		Delete: resourceConfigEmailDelete,
	}
}

func resourceConfigEmailCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	_, err := apiClient.Client.Products.PutConfigurations(products.NewPutConfigurationsParams().WithConfigurations(&models.Configurations{
		EmailHost:     d.Get("email_host").(string),
		EmailPort:     int64(d.Get("email_port").(int)),
		EmailUsername: d.Get("email_username").(string),
		EmailPassword: d.Get("email_password").(string),
		EmailFrom:     d.Get("email_from").(string),
		EmailSsl:      d.Get("email_ssl").(bool),
	}), nil)
	if err != nil {
		log.Fatal(err)
	}
	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", d.Get("email_host").(string))))
	return resourceConfigEmailRead(d, m)
}

func resourceConfigEmailRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, err := apiClient.Client.Products.GetConfigurations(products.NewGetConfigurationsParams(), nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := d.Set("email_host", resp.Payload.EmailHost.Value); err != nil {
		return err
	}
	if err := d.Set("email_port", int(resp.Payload.EmailPort.Value)); err != nil {
		return err
	}
	if err := d.Set("email_username", resp.Payload.EmailUsername.Value); err != nil {
		return err
	}
	if err := d.Set("email_from", resp.Payload.EmailFrom.Value); err != nil {
		return err
	}
	if err := d.Set("email_ssl", resp.Payload.EmailSsl.Value); err != nil {
		return err
	}
	return nil
}

func resourceConfigEmailDelete(d *schema.ResourceData, m interface{}) error {
	log.Print("not supported at the moment ... the problem is set empty values to the config")
	return nil
}
