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

func resourceLabel() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"label_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "green",
			},
			"deleted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "g",
				ForceNew: true,
				// TODO validator for p & g
			},
		},
		Create: resourceLabelCreate,
		Read:   resourceLabelRead,
		Update: resourceLabelUpdate,
		Delete: resourceLabelDelete,
	}
}
func buildLabel(d *schema.ResourceData) *models.Label {
	return &models.Label{
		Color:       d.Get("color").(string),
		Scope:       d.Get("scope").(string),
		ProjectID:   int64(d.Get("project_id").(int)),
		Deleted:     d.Get("deleted").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}
func resourceLabelCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := products.NewPostReplicationPoliciesParams().WithPolicy(&models.ReplicationPolicy{})

	_, err := apiClient.Client.Products.PostReplicationPolicies(body, nil)
	if err != nil {
		log.Fatal(err)
	}
	//d.Set("Label_id", resp.Payload[0].LabelID)
	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", d.Get("name").(string))))
	return resourceLabelRead(d, m)
}

func resourceLabelRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	LabelName := d.Get("name").(string)
	query := products.NewGetLabelsParams().WithScope(d.Get("scope").(string)).WithName(&LabelName)
	resp, err := apiClient.Client.Products.GetLabels(query, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := d.Set("label_id", int(resp.Payload[0].ID)); err != nil {
		return err
	}

	if err := d.Set("name", string(resp.Payload[0].Name)); err != nil {
		return err
	}

	if err := d.Set("description", resp.Payload[0].Description); err != nil {
		return err
	}

	if err := d.Set("deleted", resp.Payload[0].Deleted); err != nil {
		return err
	}
	if err := d.Set("color", resp.Payload[0].Color); err != nil {
		return err
	}
	if err := d.Set("scope", resp.Payload[0].Scope); err != nil {
		return err
	}
	if err := d.Set("project_id", int(resp.Payload[0].ProjectID)); err != nil {
		return err
	}

	return nil
}

func resourceLabelUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := products.NewPutLabelsIDParams().WithLabel(buildLabel(d))

	_, err := apiClient.Client.Products.PutLabelsID(body, nil)
	if err != nil {
		log.Fatal(err)
	}

	return resourceLabelRead(d, m)
}

func resourceLabelDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	labelID := d.Get("label_id").(int)

	delete := products.NewDeleteLabelsIDParams().WithID(int64(labelID))
	_, err := apiClient.Client.Products.DeleteLabelsID(delete, nil)
	if err != nil {
		log.Fatal(err)
	}
	d.SetId("")
	return nil
}
