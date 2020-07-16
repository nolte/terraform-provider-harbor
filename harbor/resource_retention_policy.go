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

func resourceRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "or",
			},
			"scope": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: retentionPolicyScopeFields(),
				},
			},
			// is it optional? 0 item = null array?
			"rule": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: retentionPolicyRuleFields(),
				},
			},
			"trigger": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: retentionPolicyTriggerFields(),
				},
			},
		},
		Create: resourceRetentionPolicyCreate,
		Read:   resourceRetentionPolicyRead,
		Update: resourceRetentionPolicyUpdate,
		Delete: resourceRetentionPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRetentionPolicyCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	projectID := int64(d.Get("scope.0.ref").(int))
	body := products.NewPostRetentionsParams().WithPolicy(&models.RetentionPolicy{
		Algorithm: d.Get("algorithm").(string),
		Scope: &models.RetentionPolicyScope{
			Level: d.Get("scope.0.level").(string),
			Ref:   projectID,
		},
		Rules: expandRetentionPolicyRules(d.Get("rule").([]interface{})),
		Trigger: &models.RetentionRuleTrigger{
			Kind:       d.Get("trigger.0.kind").(string),
			References: map[string]int{}, // default reference is empty map
			Settings: map[string]string{
				"cron": d.Get("trigger.0.settings.0.cron").(string),
			},
		},
	})

	if _, err := apiClient.Products.PostRetentions(body, nil); err != nil {
		return err
	}

	retentionID, err := findRetentionIDByProject(d, m, projectID)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(int(retentionID)))

	return resourceRetentionPolicyRead(d, m)
}

func findRetentionIDByProject(d *schema.ResourceData, m interface{}, projectID int64) (int64, error) {
	apiClient := m.(*client.Harbor)

	query := products.NewGetProjectsProjectIDParams().WithProjectID(projectID)

	resp, err := apiClient.Products.GetProjectsProjectID(query, nil)
	if err != nil {
		d.SetId("")
		return 0, err
	}

	retentionID := resp.Payload.Metadata.RetentionID
	if retentionID == "" {
		return 0, fmt.Errorf("retention ID is null or empty for project ID %v", projectID)
	}

	retentionIDInt, err := strconv.Atoi(retentionID)
	if err != nil {
		return 0, fmt.Errorf("couldn't parse retention id %v from project ID %v: %v", retentionID, projectID, err)
	}

	return int64(retentionIDInt), nil
}

func resourceRetentionPolicyRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if retentionID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		query := products.NewGetRetentionsIDParams().WithID(retentionID)

		resp, err := apiClient.Products.GetRetentionsID(query, nil)
		if err != nil {
			return err
		}

		if err := setRetentionPolicySchema(d, resp.Payload); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("fail to load the retention policy")
}

func resourceRetentionPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Harbor)

	if retentionID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
		body := products.NewPutRetentionsIDParams().WithPolicy(&models.RetentionPolicy{
			Algorithm: d.Get("algorithm").(string),
			Scope: &models.RetentionPolicyScope{
				Level: d.Get("scope.0.level").(string),
				Ref:   int64(d.Get("scope.0.ref").(int)),
			},
			Rules: expandRetentionPolicyRules(d.Get("rule").([]interface{})),
			Trigger: &models.RetentionRuleTrigger{
				Kind:       d.Get("trigger.0.kind").(string),
				References: expandRetentionPolicyTriggerReferencesUpdate(d.Get("trigger.0.references.0").(map[string]interface{})),
				Settings: map[string]string{
					"cron": d.Get("trigger.0.settings.0.cron").(string),
				},
			},
		}).WithID(retentionID)

		if _, err := apiClient.Products.PutRetentionsID(body, nil); err != nil {
			return err
		}

		return resourceRetentionPolicyRead(d, m)
	}

	return fmt.Errorf("retention policy Id not a Integer")
}

// Delete is not supported by API
// a good way to implement it would be:
//     - to empty the policy (but the id is still existent)
// and - to use "put" for creation if empty policy exists
// and - to prevent creation if non-empty pocity exists
//func resourceRetentionPolicyDelete(d *schema.ResourceData, m interface{}) error {
//	apiClient := m.(*client.Harbor)
//
//	if retentionID, err := strconv.ParseInt(d.Id(), 10, 64); err == nil {
//		delete := products.NewDeleteRetentionsIDParams().WithID(retentionID)
//		if _, err := apiClient.Products.DeleteRetentionsID(delete, nil); err != nil {
//			return err
//		}
//
//		return nil
//	}
//
//	return fmt.Errorf("Retention policy Id not a Integer")
//}

func resourceRetentionPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Print("not supported at the moment ... the problem is set empty values to the retention policy")
	return nil
}

func setRetentionPolicySchema(data *schema.ResourceData, policy *models.RetentionPolicy) error {
	data.SetId(strconv.Itoa(int(policy.ID)))

	if err := data.Set("algorithm", policy.Algorithm); err != nil {
		return err
	}

	scope, err := flattenRetentionPolicyScope(policy.Scope, data)
	if err != nil {
		return err
	}

	if err := data.Set("scope", scope); err != nil {
		return err
	}

	trigger, err := flattenRetentionPolicyTrigger(policy.Trigger, data)
	if err != nil {
		return err
	}

	if err := data.Set("trigger", trigger); err != nil {
		return err
	}

	rule, err := flattenRetentionPolicyRules(policy.Rules, data)
	if err != nil {
		return err
	}

	if err := data.Set("rule", rule); err != nil {
		return err
	}

	return nil
}
