package harbor

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func resourceRegistry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"credential_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceRegistryCreate,
		Read:   resourceRegistryRead,
		Update: resourceRegistryUpdate,
		Delete: resourceRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRegistryCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	body := products.NewPostRegistriesParams().WithRegistry(&models.Registry{
		Description: d.Get("description").(string),
		Insecure:    d.Get("insecure").(bool),
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		URL:         d.Get("url").(string),
		Credential: &models.RegistryCredential{
			AccessKey:    d.Get("access_key").(string),
			AccessSecret: d.Get("access_secret").(string),
			Type:         d.Get("credential_type").(string),
		},
	})

	_, err := apiClient.Products.PostRegistries(body, nil)
	if err != nil {
		return err
	}

	registry, err := findRegistryByName(d, m)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(int(registry.ID)))

	return resourceRegistryRead(d, m)
}
func findRegistryByName(d *schema.ResourceData, m interface{}) (*models.Registry, error) {
	apiClient := m.(*client.Harbor)

	if name, ok := d.GetOk("name"); ok {
		registryName := name.(string)
		query := products.NewGetRegistriesParams().WithName(&registryName)

		resp, err := apiClient.Products.GetRegistries(query, nil)
		if err != nil {
			d.SetId("")
			return &models.Registry{}, err
		}

		if len(resp.Payload) < 1 {
			return &models.Registry{}, fmt.Errorf("no Registry found with name %v", registryName)
		} else if resp.Payload[0].Name != registryName {
			return &models.Registry{},
				fmt.Errorf("response Name %v not match with Expected Name %v", resp.Payload[0].Name, registryName)
		}

		return resp.Payload[0], nil
	}

	return &models.Registry{}, fmt.Errorf("fail to lookup Registry by Name")
}
func resourceRegistryRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if registryID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		resp, err := apiClient.Products.GetRegistriesID(products.NewGetRegistriesIDParams().WithID(registryID), nil)

		if err != nil {
			return err
		}

		if err = setRegistrySchema(d, resp.Payload); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("registry Id not a Integer currently: '%s'", d.Id())
}

func resourceRegistryUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if registryID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		params := products.NewPutRegistriesIDParams().WithID(registryID).WithRepoTarget(&models.PutRegistry{
			Description:    d.Get("description").(string),
			Insecure:       d.Get("insecure").(bool),
			Name:           d.Get("name").(string),
			URL:            d.Get("url").(string),
			AccessKey:      d.Get("access_key").(string),
			AccessSecret:   d.Get("access_secret").(string),
			CredentialType: d.Get("credential_type").(string),
		})
		if _, err := apiClient.Products.PutRegistriesID(params, nil); err != nil {
			return err
		}

		return resourceRegistryRead(d, m)
	}

	return fmt.Errorf("registry Id not a Integer")
}

func resourceRegistryDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if registryID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		params := products.NewDeleteRegistriesIDParams().WithID(registryID)

		if _, err := apiClient.Products.DeleteRegistriesID(params, nil); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("registry Id not a Integer")
}

func setRegistrySchema(d *schema.ResourceData, registry *models.Registry) error {
	d.SetId(strconv.Itoa(int(registry.ID)))

	if err := d.Set("description", registry.Description); err != nil {
		return err
	}

	if err := d.Set("insecure", registry.Insecure); err != nil {
		return err
	}

	if err := d.Set("name", registry.Name); err != nil {
		return err
	}

	if err := d.Set("type", registry.Type); err != nil {
		return err
	}

	if err := d.Set("url", registry.URL); err != nil {
		return err
	}

	if err := d.Set("access_key", registry.Credential.AccessKey); err != nil {
		return err
	}

	if err := d.Set("credential_type", registry.Credential.Type); err != nil {
		return err
	}

	return nil
}
