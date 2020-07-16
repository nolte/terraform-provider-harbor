package harbor

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func init() {
	resource.AddTestSweepers("resource_retention_policy", &resource.Sweeper{
		Name: "harbor_retention_policy",
	})
}

func flattenRetentionPolicyScope(scope *models.RetentionPolicyScope, data *schema.ResourceData) ([]interface{}, error) {
	att := map[string]interface{}{
		"level": scope.Level,
		"ref":   scope.Ref,
	}

	return []interface{}{att}, nil
}

func flattenRetentionPolicyTrigger(trigger *models.RetentionRuleTrigger, data *schema.ResourceData) ([]interface{}, error) {
	att := map[string]interface{}{
		"kind":       trigger.Kind,
		"references": []interface{}{trigger.References},
		"settings":   []interface{}{trigger.Settings},
	}

	return []interface{}{att}, nil
}

func flattenRetentionPolicyTriggerSettings(settings interface{}, data *schema.ResourceData) ([]interface{}, error) {
	settingsStruct, ok := settings.(models.TriggerSettings)
	if !ok {
		return nil, fmt.Errorf("Can't cast settings value %v to TriggerSettings", settings)
	}

	att := map[string]interface{}{
		"cron": settingsStruct.Cron,
	}

	return []interface{}{att}, nil
}

func flattenRetentionPolicyRulesRetentionSelectors(selector *models.RetentionSelector, data *schema.ResourceData) (interface{}, error) {
	att := map[string]interface{}{
		"decoration": selector.Decoration,
		"extras":     selector.Extras,
		"kind":       selector.Kind,
		"pattern":    selector.Pattern,
	}

	return att, nil
}

func flattenRetentionPolicyRulesTagSelectors(tagSelectors []*models.RetentionSelector, data *schema.ResourceData) ([]interface{}, error) {
	att := make([]interface{}, len(tagSelectors))

	for i, selector := range tagSelectors {
		selectorAtt, err := flattenRetentionPolicyRulesRetentionSelectors(selector, data)
		if err != nil {
			return []interface{}{}, err
		}

		att[i] = selectorAtt
	}

	return att, nil
}

func flattenRetentionPolicyRulesScopeSelectors(scopeSelectors []models.RetentionSelector, data *schema.ResourceData) ([]interface{}, error) {
	att := make([]interface{}, len(scopeSelectors))

	for i, selector := range scopeSelectors {
		selectorAtt, err := flattenRetentionPolicyRulesRetentionSelectors(&selector, data)
		if err != nil {
			return []interface{}{}, err
		}

		att[i] = selectorAtt
	}

	return att, nil
}

func flattenRetentionPolicyRulesScopesSelectors(scopesSelectors map[string][]models.RetentionSelector, data *schema.ResourceData) ([]interface{}, error) {
	att := make([]interface{}, len(scopesSelectors))

	i := 0
	for scope, selector := range scopesSelectors {
		if scope != "repository" {
			return []interface{}{}, fmt.Errorf("scope %s is unsupported", scope)
		}
		selectorAtt, err := flattenRetentionPolicyRulesScopeSelectors(selector, data)
		if err != nil {
			return []interface{}{}, err
		}

		att[i] = map[string]interface{}{
			scope: selectorAtt,
		}

		i++
	}

	return att, nil
}

func flattenRetentionPolicyRules(rules []*models.RetentionRule, data *schema.ResourceData) ([]interface{}, error) {
	att := make([]interface{}, len(rules))

	for i, rule := range rules {
		tagSelectorsAtt, err := flattenRetentionPolicyRulesTagSelectors(rule.TagSelectors, data)
		if err != nil {
			return []interface{}{}, err
		}
		scopeSelectorsAtt, err := flattenRetentionPolicyRulesScopesSelectors(rule.ScopeSelectors, data)
		if err != nil {
			return []interface{}{}, err
		}

		ruleAtt := map[string]interface{}{
			"id":              rule.ID,
			"priority":        rule.Priority,
			"disabled":        rule.Disabled,
			"action":          rule.Action,
			"template":        rule.Template,
			"params":          rule.Params,
			"tag_selectors":   tagSelectorsAtt,
			"scope_selectors": scopeSelectorsAtt,
		}

		att[i] = ruleAtt
	}

	return att, nil
}

func expandRetentionPolicyTriggerReferencesUpdate(selectorsAtt map[string]interface{}) map[string]int {
	references := map[string]int{}

	if referencesJobId, ok := selectorsAtt["job_id"]; ok {
		references["job_id"] = referencesJobId.(int)
	}

	return references
}

func expandRetentionPolicyRulesRetentionSelectors(selectorsAtt interface{}) models.RetentionSelector {
	retentionSelectors := models.RetentionSelector{}
	selectorsAttMap := selectorsAtt.(map[string]interface{})

	if decoration, ok := selectorsAttMap["decoration"]; ok {
		retentionSelectors.Decoration = decoration.(string)
	}

	if extras, ok := selectorsAttMap["extras"]; ok {
		retentionSelectors.Extras = extras.(string)
	}

	if kind, ok := selectorsAttMap["kind"]; ok {
		retentionSelectors.Kind = kind.(string)
	}

	if pattern, ok := selectorsAttMap["pattern"]; ok {
		retentionSelectors.Pattern = pattern.(string)
	}

	return retentionSelectors
}

func expandRetentionPolicyRulesTagSelectors(tagSelectorsAtt []interface{}) []*models.RetentionSelector {
	tagSelectors := make([]*models.RetentionSelector, len(tagSelectorsAtt))

	for i, selectorAtt := range tagSelectorsAtt {
		tagSelector := expandRetentionPolicyRulesRetentionSelectors(selectorAtt)

		tagSelectors[i] = &tagSelector
	}

	return tagSelectors
}

func expandRetentionPolicyRulesScopeSelectors(scopeSelectorsAtt []interface{}) []models.RetentionSelector {
	scopeSelectors := make([]models.RetentionSelector, len(scopeSelectorsAtt))

	for i, selectorAtt := range scopeSelectorsAtt {
		scopeSelector := expandRetentionPolicyRulesRetentionSelectors(selectorAtt)

		scopeSelectors[i] = scopeSelector
	}

	return scopeSelectors
}

func expandRetentionPolicyRulesScopesSelectors(scopesSelectorsAtt []interface{}) map[string][]models.RetentionSelector {
	scopesSelectors := make(map[string][]models.RetentionSelector, len(scopesSelectorsAtt))

	for _, scopeSelectorsAtt := range scopesSelectorsAtt {
		scopeSelectorsMap := scopeSelectorsAtt.(map[string]interface{})

		if repositoryScopeSelectorsAtt, ok := scopeSelectorsMap["repository"]; ok {
			scopesSelectors["repository"] = expandRetentionPolicyRulesScopeSelectors(repositoryScopeSelectorsAtt.([]interface{}))
		}
	}

	return scopesSelectors
}

func expandRetentionPolicyRules(rulesAtt []interface{}) []*models.RetentionRule {
	rules := make([]*models.RetentionRule, len(rulesAtt))

	for i, ruleAtt := range rulesAtt {
		ruleAttMap := ruleAtt.(map[string]interface{})
		var rule models.RetentionRule

		if disabled, ok := ruleAttMap["disabled"]; ok {
			rule.Disabled = disabled.(bool)
		}

		// computed value. Always as "retain".
		rule.Action = "retain"

		if template, ok := ruleAttMap["template"]; ok {
			rule.Template = template.(string)
		}

		if params, ok := ruleAttMap["params"]; ok {
			rule.Params = params.(map[string]interface{})
		}

		if tagSelectors, ok := ruleAttMap["tag_selectors"]; ok {
			rule.TagSelectors = expandRetentionPolicyRulesTagSelectors(tagSelectors.([]interface{}))
		}

		if scopesSelectors, ok := ruleAttMap["scope_selectors"]; ok {
			rule.ScopeSelectors = expandRetentionPolicyRulesScopesSelectors(scopesSelectors.([]interface{}))
		}

		rules[i] = &rule
	}

	return rules
}
