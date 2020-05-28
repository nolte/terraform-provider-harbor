package harbor

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func resourceReplicationPull() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_registry_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"source_registry_filter_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_registry_filter_tag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_namespace": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Destination namespace Specify the destination namespace.
                    If empty, the resources will be put under the same namespace as the source.`,
			},
			"trigger_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "manual",
			},
			"trigger_cron": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"override": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Create: resourceReplicationPullCreate,
		Read:   resourceReplicationPullRead,
		Update: resourceReplicationUpdate,
		Delete: resourceReplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func buildReplicationPolicy(d *schema.ResourceData, m interface{}) *models.ReplicationPolicy {
	filters := make([]*models.ReplicationFilter, 2)
	filters[0] = &models.ReplicationFilter{
		Type:  "name",
		Value: d.Get("source_registry_filter_name").(string),
	}
	filters[1] = &models.ReplicationFilter{
		Type:  "tag",
		Value: d.Get("source_registry_filter_tag").(string),
	}
	return &models.ReplicationPolicy{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		DestNamespace: d.Get("destination_namespace").(string),
		SrcRegistry: &models.Registry{
			ID: int64(d.Get("source_registry_id").(int)),
		},
		Filters:  filters,
		Enabled:  d.Get("enabled").(bool),
		Override: d.Get("override").(bool),
		Trigger: &models.ReplicationTrigger{
			Type: d.Get("trigger_mode").(string),
			TriggerSettings: &models.TriggerSettings{
				Cron: d.Get("trigger_cron").(string),
			},
		},
	}
}

func resourceReplicationPullCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	params := products.NewPostReplicationPoliciesParams().WithPolicy(buildReplicationPolicy(d, m))

	_, err := apiClient.Products.PostReplicationPolicies(params, nil)
	if err != nil {
		return err
	}

	registry, err := findReplicationByName(d, m)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(int(registry.ID)))

	return resourceReplicationPullRead(d, m)
}

func findReplicationByName(d *schema.ResourceData, m interface{}) (*models.ReplicationPolicy, error) {
	apiClient := m.(*client.Harbor)

	if name, ok := d.GetOk("name"); ok {
		registryName := name.(string)
		query := products.NewGetReplicationPoliciesParams().WithName(&registryName)

		resp, err := apiClient.Products.GetReplicationPolicies(query, nil)
		if err != nil {
			d.SetId("")
			return &models.ReplicationPolicy{}, err
		}

		if len(resp.Payload) < 1 {
			return &models.ReplicationPolicy{}, fmt.Errorf("no Replication found with name %v", registryName)
		} else if resp.Payload[0].Name != registryName {
			return &models.ReplicationPolicy{},
				fmt.Errorf("response Name %v not match with Expected Name %v", resp.Payload[0].Name, registryName)
		}

		return resp.Payload[0], nil
	}

	return &models.ReplicationPolicy{}, fmt.Errorf("fail to lookup Replication by Name")
}

func resourceReplicationPullRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if registryID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		params := products.NewGetReplicationPoliciesIDParams().WithID(registryID)

		resp, err := apiClient.Products.GetReplicationPoliciesID(params, nil)
		if err != nil {
			return err
		}

		if err = setReplicationSchema(d, resp.Payload); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("replication Id not a Integer currently: '%s'", d.Id())
}

func resourceReplicationUpdate(d *schema.ResourceData, m interface{}) error {

	apiClient := m.(*client.Harbor)

	if registryID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		params := products.NewPutReplicationPoliciesIDParams().WithPolicy(buildReplicationPolicy(d, m)).WithID(registryID)
		if _, err := apiClient.Products.PutReplicationPoliciesID(params, nil); err != nil {
			return err
		}

		return resourceReplicationPullRead(d, m)
	}

	return fmt.Errorf("replication Id not a Integer")

}

func resourceReplicationDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if registryID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		params := products.NewDeleteReplicationPoliciesIDParams().WithID(registryID)

		if _, err := apiClient.Products.DeleteReplicationPoliciesID(params, nil); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("replication Id not a Integer")
}

func setReplicationSchema(d *schema.ResourceData, registry *models.ReplicationPolicy) error {
	d.SetId(strconv.Itoa(int(registry.ID)))

	if err := d.Set("description", registry.Description); err != nil {
		return err
	}

	if err := d.Set("name", registry.Name); err != nil {
		return err
	}

	if err := d.Set("destination_namespace", registry.DestNamespace); err != nil {
		return err
	}

	if err := d.Set("trigger_mode", registry.Trigger.Type); err != nil {
		return err
	}

	if err := d.Set("trigger_cron", registry.Trigger.TriggerSettings.Cron); err != nil {
		return err
	}

	if err := d.Set("enabled", registry.Enabled); err != nil {
		return err
	}

	if err := d.Set("override", registry.Override); err != nil {
		return err
	}

	if err := d.Set("source_registry_id", registry.SrcRegistry.ID); err != nil {
		return err
	}

	return nil
}
