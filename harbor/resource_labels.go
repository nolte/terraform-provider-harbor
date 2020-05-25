package harbor

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func resourceLabel() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "#61717D",
			},
			"deleted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "g",
				ForceNew:     true,
				ValidateFunc: validateLabelGroup,
			},
		},
		Create: resourceLabelCreate,
		Read:   resourceLabelRead,
		Update: resourceLabelUpdate,
		Delete: resourceLabelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func validateLabelGroup(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value != "g" && value != "p" {
		errors = append(errors, fmt.Errorf("label group must be 'p' or 'g' %s", k))
		return
	}

	return
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
	apiClient := m.(*client.Harbor)

	body := products.NewPostLabelsParams().WithLabel(buildLabel(d))

	_, err := apiClient.Products.PostLabels(body, nil)
	if err != nil {
		log.Fatal(err)
	}
	//d.Set("Label_id", resp.Payload[0].LabelID)
	label, err := findLabelByNameAndScope(d, m)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(int(label.ID)))

	return resourceLabelRead(d, m)
}

func findLabelByNameAndScope(d *schema.ResourceData, m interface{}) (*models.Label, error) {
	apiClient := m.(*client.Harbor)

	if name, ok := d.GetOk("name"); ok {
		if scope, ok := d.GetOk("scope"); ok {
			searchName := name.(string)
			scopeName := scope.(string)

			query := products.NewGetLabelsParams().WithScope(scopeName).WithName(&searchName)

			if scopeName == "p" {
				projectID := int64(d.Get("project_id").(int))
				query = query.WithProjectID(&projectID)
			}

			resp, err := apiClient.Products.GetLabels(query, nil)
			if err != nil {
				d.SetId("")
				return &models.Label{}, err
			}

			if len(resp.Payload) < 1 {
				return &models.Label{}, fmt.Errorf("no label found with name %v", searchName)
			} else if resp.Payload[0].Name != searchName {
				return &models.Label{},
					fmt.Errorf("response Name %v not match with Expected Name %v", resp.Payload[0].Name, searchName)
			}

			return resp.Payload[0], nil
		}
	}

	return &models.Label{}, fmt.Errorf("fail to lookup label by Name and Scope")
}

func findLabelByID(d *schema.ResourceData, m interface{}) (*models.Label, error) {
	apiClient := m.(*client.Harbor)

	if searchID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		query := products.NewGetLabelsIDParams().WithID(searchID)

		resp, err := apiClient.Products.GetLabelsID(query, nil)
		if err != nil {
			return &models.Label{}, err
		}

		return resp.Payload, nil
	}

	return &models.Label{}, fmt.Errorf("fail to find the label")
}

func setLabelSchema(d *schema.ResourceData, resp *models.Label) error {
	d.SetId(strconv.Itoa(int(resp.ID)))

	if err := d.Set("name", resp.Name); err != nil {
		return err
	}

	if err := d.Set("description", resp.Description); err != nil {
		return err
	}

	if err := d.Set("deleted", resp.Deleted); err != nil {
		return err
	}

	if err := d.Set("color", resp.Color); err != nil {
		return err
	}

	if err := d.Set("scope", resp.Scope); err != nil {
		return err
	}

	if err := d.Set("project_id", int(resp.ProjectID)); err != nil {
		return err
	}

	return nil
}
func resourceLabelRead(d *schema.ResourceData, m interface{}) error {
	label, err := findLabelByID(d, m)
	if err != nil {
		return err
	}

	if err := setLabelSchema(d, label); err != nil {
		return err
	}

	return nil
}

func resourceLabelUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if resourceID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		body := products.NewPutLabelsIDParams().WithLabel(buildLabel(d)).WithID(resourceID)
		if _, err := apiClient.Products.PutLabelsID(body, nil); err != nil {
			return err
		}

		return resourceLabelRead(d, m)
	}

	return fmt.Errorf("label Id not a Integer")
}

func resourceLabelDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if resourceID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		delete := products.NewDeleteLabelsIDParams().WithID(resourceID)
		if _, err := apiClient.Products.DeleteLabelsID(delete, nil); err != nil {
			return err
		}

		d.SetId("")

		return nil
	}

	return fmt.Errorf("label Id not a Integer")
}
