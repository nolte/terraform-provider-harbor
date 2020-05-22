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

func resourceRegistry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Create: resourceRegistryCreate,
		Read:   resourceRegistryRead,
		Update: resourceRegistryUpdate,
		Delete: resourceRegistryDelete,
	}
}

func resourceRegistryCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := products.NewPostRegistriesParams().WithRegistry(&models.Registry{
		Description: d.Get("description").(string),
		Insecure:    d.Get("insecure").(bool),
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		URL:         d.Get("url").(string),
	})
	_, err := apiClient.Client.Products.PostRegistries(body, nil)
	if err != nil {
		log.Fatal(err)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", d.Get("name").(string))))
	return resourceRegistryRead(d, m)
}

func resourceRegistryRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	projectName := d.Get("name").(string)

	query := products.NewGetRegistriesParams().WithName(&projectName)
	resp, err := apiClient.Client.Products.GetRegistries(query, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := d.Set("repository_id", int(resp.Payload[0].ID)); err != nil {
		return err
	}
	if err := d.Set("description", resp.Payload[0].Description); err != nil {
		return err
	}
	if err := d.Set("insecure", resp.Payload[0].Insecure); err != nil {
		return err
	}
	if err := d.Set("name", resp.Payload[0].Name); err != nil {
		return err
	}
	if err := d.Set("type", resp.Payload[0].Type); err != nil {
		return err
	}
	if err := d.Set("url", resp.Payload[0].URL); err != nil {
		return err
	}
	return nil
}

func resourceRegistryUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	_, err := apiClient.Client.Products.PutRegistriesID(products.NewPutRegistriesIDParams().WithRepoTarget(&models.PutRegistry{
		Description: d.Get("description").(string),
		Insecure:    d.Get("insecure").(bool),
		Name:        d.Get("name").(string),
		URL:         d.Get("url").(string),
	}), nil)

	if err != nil {
		log.Fatal(err)
	}
	return resourceRegistryRead(d, m)
}

func resourceRegistryDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, err := apiClient.Client.Products.DeleteRegistriesID(products.NewDeleteRegistriesIDParams().WithID(int64(d.Get("repository_id").(int))), nil)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
