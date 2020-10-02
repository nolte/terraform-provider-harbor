package harbor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/client/products"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func resourceWebhookPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"endpoint_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notify_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http",
			},
			"auth_header": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"skip_cert_verify": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"event_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		Create: resourceWebhookPolicyCreate,
		Read:   resourceWebhookPolicyRead,
		Update: resourceWebhookPolicyUpdate,
		Delete: resourceWebhookPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceWebhookPolicyCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	projectID := int64(d.Get("project_id").(int))
	webhookPolicyName := d.Get("name").(string)

	webhookTarget := &models.WebhookTargetObject{
		Type:           d.Get("notify_type").(string),
		AuthHeader:     d.Get("auth_header").(string),
		SkipCertVerify: d.Get("skip_cert_verify").(bool),
		Address:        d.Get("endpoint_url").(string),
	}

	webhookTargets := make([]*models.WebhookTargetObject, 1)
	webhookTargets[0] = webhookTarget

	body := products.NewPostProjectsProjectIDWebhookPoliciesParams().WithProjectID(projectID).WithPolicy(&models.WebhookPolicy{
		Name:        webhookPolicyName,
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		EventTypes:  expandWebhookPolicyEvents(d.Get("event_types").([]interface{})),
		Targets:     webhookTargets,
	})

	if _, err := apiClient.Products.PostProjectsProjectIDWebhookPolicies(body, nil); err != nil {
		return err
	}

	listResp, err := apiClient.Products.GetProjectsProjectIDWebhookPolicies(
		products.NewGetProjectsProjectIDWebhookPoliciesParams().
			WithProjectID(projectID),
		nil,
	)

	if err != nil {
		return fmt.Errorf("webhook policy not found")
	}

	var webhookPolicy *models.WebhookPolicy
	for _, policy := range listResp.Payload {
		if policy.Name == webhookPolicyName {
			webhookPolicy = policy
			break
		}
	}

	if webhookPolicy == nil {
		return fmt.Errorf("failed to find the webhook policy %s", webhookPolicyName)
	}

	d.SetId(fmt.Sprintf("%d/%d", webhookPolicy.ProjectID, webhookPolicy.ID))

	return resourceWebhookPolicyRead(d, m)
}

func resourceWebhookPolicyParseID(id string) (int64, int64, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid id %s", id)
	}

	projectID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	webhookPolicyID, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return int64(projectID), int64(webhookPolicyID), nil
}

func resourceWebhookPolicyRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if projectID, webhookPolicyID, err := resourceWebhookPolicyParseID(d.Id()); err == nil {

		resp, err := apiClient.Products.GetProjectsProjectIDWebhookPoliciesPolicyID(
			products.NewGetProjectsProjectIDWebhookPoliciesPolicyIDParams().
				WithProjectID(projectID).
				WithPolicyID(webhookPolicyID),
			nil,
		)

		if err != nil {
			return err
		}

		if err := setWebhookPolicySchema(d, resp.Payload); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("failed to find the webhook policy ")
}

func setWebhookPolicySchema(d *schema.ResourceData, webhookPolicy *models.WebhookPolicy) error {
	d.SetId(fmt.Sprintf("%d/%d", webhookPolicy.ProjectID, webhookPolicy.ID))

	if err := d.Set("name", webhookPolicy.Name); err != nil {
		return err
	}

	if err := d.Set("description", webhookPolicy.Description); err != nil {
		return err
	}

	if err := d.Set("project_id", webhookPolicy.ProjectID); err != nil {
		return err
	}

	if err := d.Set("event_types", flattenWebhookPolicyEvents(webhookPolicy.EventTypes)); err != nil {
		return err
	}

	if err := d.Set("enabled", webhookPolicy.Enabled); err != nil {
		return err
	}

	if err := d.Set("endpoint_url", webhookPolicy.Targets[0].Address); err != nil {
		return err
	}

	if err := d.Set("notify_type", webhookPolicy.Targets[0].Type); err != nil {
		return err
	}

	if err := d.Set("auth_header", webhookPolicy.Targets[0].AuthHeader); err != nil {
		return err
	}

	if err := d.Set("skip_cert_verify", webhookPolicy.Targets[0].SkipCertVerify); err != nil {
		return err
	}

	return nil
}

func resourceWebhookPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if projectID, webhookPolicyID, err := resourceWebhookPolicyParseID(d.Id()); err == nil {

		webhookTarget := &models.WebhookTargetObject{
			Type:           d.Get("notify_type").(string),
			AuthHeader:     d.Get("auth_header").(string),
			SkipCertVerify: d.Get("skip_cert_verify").(bool),
			Address:        d.Get("endpoint_url").(string),
		}

		webhookTargets := make([]*models.WebhookTargetObject, 1)
		webhookTargets[0] = webhookTarget

		body := products.NewPutProjectsProjectIDWebhookPoliciesPolicyIDParams().WithProjectID(projectID).WithPolicyID(webhookPolicyID).WithPolicy(&models.WebhookPolicy{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Enabled:     d.Get("enabled").(bool),
			EventTypes:  expandWebhookPolicyEvents(d.Get("event_types").([]interface{})),
			Targets:     webhookTargets,
		})

		_, err = apiClient.Products.PutProjectsProjectIDWebhookPoliciesPolicyID(body, nil)

		if err != nil {
			return fmt.Errorf("Webhook policy update failed")
		}

		return resourceWebhookPolicyRead(d, m)
	}

	return fmt.Errorf("failed to find the webhook policy")
}

func resourceWebhookPolicyDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if projectID, webhookPolicyID, err := resourceWebhookPolicyParseID(d.Id()); err == nil {

		delete := products.NewDeleteProjectsProjectIDWebhookPoliciesPolicyIDParams().WithProjectID(projectID).WithPolicyID(webhookPolicyID)
		if _, err := apiClient.Products.DeleteProjectsProjectIDWebhookPoliciesPolicyID(delete, nil); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("failed to find the webhook policy")
}

func expandWebhookPolicyEvents(eventTypes []interface{}) []string {
	events := make([]string, len(eventTypes))
	for i, event := range eventTypes {
		events[i] = event.(string)
	}
	return events
}

func flattenWebhookPolicyEvents(events []string) []interface{} {
	eventTypes := make([]interface{}, len(events))
	for i, eventType := range events {
		eventTypes[i] = eventType
	}
	return eventTypes
}
