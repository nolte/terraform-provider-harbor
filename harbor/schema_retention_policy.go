package harbor

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func retentionPolicyRuleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"priority": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"disabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"action": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"template": {
			Type:     schema.TypeString,
			Required: true,
		},
		"params": {
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
			Optional: true,
			Default:  map[string]int{},
		},
		"tag_selectors": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: retentionPolicyRuleRetentionSelectorsFields(),
			},
		},
		"scope_selectors": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: retentionPolicyRuleScopeSelectorsFields(),
			},
		},
	}

	return s
}

func retentionPolicyRuleScopeSelectorsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"repository": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: retentionPolicyRuleRetentionSelectorsFields(),
			},
		},
	}

	return s
}

func retentionPolicyRuleRetentionSelectorsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"kind": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "doublestar",
		},
		"extras": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"decoration": {
			Type:     schema.TypeString,
			Required: true,
		},
		"pattern": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	return s
}

func retentionPolicyScopeFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"level": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "project",
		},
		"ref": {
			Type:     schema.TypeInt,
			Required: true,
		},
	}

	return s
}

func retentionPolicyTriggerSettingsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cron": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
	}

	return s
}

func retentionPolicyTriggerReferencesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"job_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}

	return s
}

func retentionPolicyTriggerFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"kind": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "Schedule",
		},
		"references": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: retentionPolicyTriggerReferencesFields(),
			},
			Computed: true,
		},
		"settings": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: retentionPolicyTriggerSettingsFields(),
			},
		},
	}

	return s
}
