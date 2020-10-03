package harbor

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

//nolint:funlen
func resourceReplication() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name from the Replication Policy",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Will be displayed in harbor",
			},
			"source_registry_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Pull the resources from the remote registry to the local Harbor.`,
			},
			"source_registry_filter_name": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `
                    Filter the name of the resource.
                    Leave empty or use '**' to match all. 'library/**' only matches resources under 'library'.
                    For more patterns, please refer to the user guide.
                `,
			},
			"source_registry_filter_tag": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `
                    Filter the tag/version part of the resources.
                    Leave empty or use '**' to match all. '1.0*' only matches the tags that starts with '1.0'.
                    For more patterns, please refer to the user guide.
                `,
			},
			//  not supported for the moment swagger client problems
			//  "source_registry_filter_labels": {
			//  	Type:     schema.TypeList,
			//  	Optional: true,
			//  	Elem:     &schema.Schema{Type: schema.TypeString},
			//  },
			"destination_namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Destination namespace Specify the destination namespace.
                    If empty, the resources will be put under the same namespace as the source.`,
			},
			"destination_registry_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The Id from the destination registry used for push-based replication policies",
			},
			"trigger_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "manual",
				Description: "Can be manual,scheduled and for push-based addition event_based",
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
				Description: `Specify whether to override the resources at the destination
                if a resource with the same name exists.`,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Create: resourceReplicationCreate,
		Read:   resourceReplicationRead,
		Update: resourceReplicationUpdate,
		Delete: resourceReplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func buildReplicationPolicy(d *schema.ResourceData) *models.ReplicationPolicy {
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
		DestRegistry: &models.Registry{
			ID: int64(d.Get("destination_registry_id").(int)),
		},
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

func resourceReplicationCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	params := products.NewPostReplicationPoliciesParams().WithPolicy(buildReplicationPolicy(d))

	_, err := apiClient.Products.PostReplicationPolicies(params, nil)
	if err != nil {
		return err
	}

	registry, err := findReplicationByName(d, m)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(int(registry.ID)))

	return resourceReplicationRead(d, m)
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
		for _, element := range resp.Payload {
			if element.Name == registryName {
				return element, nil
			}
		}
		return &models.ReplicationPolicy{}, fmt.Errorf("no Replication found with name %v", registryName)
	}

	return &models.ReplicationPolicy{}, fmt.Errorf("fail to lookup Replication by Name")
}

func resourceReplicationRead(d *schema.ResourceData, m interface{}) error {
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
		params := products.NewPutReplicationPoliciesIDParams().WithPolicy(buildReplicationPolicy(d)).WithID(registryID)
		if _, err := apiClient.Products.PutReplicationPoliciesID(params, nil); err != nil {
			return err
		}

		return resourceReplicationRead(d, m)
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

	if err := d.Set("destination_registry_id", registry.DestRegistry.ID); err != nil {
		return err
	}

	return nil
}
