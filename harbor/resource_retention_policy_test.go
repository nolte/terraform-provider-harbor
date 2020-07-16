package harbor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/nolte/terraform-provider-harbor/gen/harborctl/models"
)

func init() {
	resource.AddTestSweepers("resource_harbor_retention_policy", &resource.Sweeper{
		Name: "harbor_retention_policy",
	})
}

func TestAccHarborRetentionPolicy_Basic(t *testing.T) {
	var project models.Project

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHarborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHarborCheckRetentionPolicyMinimalResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckProjectExists("harbor_project.default", &project),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "algorithm", "or"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "scope.0.level", "project"),
					// add check on scope.0.ref = project_id?
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "trigger.0.settings.0.cron", ""),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.id", "0"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.priority", "0"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.action", "retain"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.template", "always"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.tag_selectors.0.kind", "doublestar"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.tag_selectors.0.decoration", "matches"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.tag_selectors.0.pattern", "**"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.tag_selectors.0.extras", ""),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.scope_selectors.0.repository.0.kind", "doublestar"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.scope_selectors.0.repository.0.decoration", "repoMatches"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.scope_selectors.0.repository.0.pattern", "**"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.scope_selectors.0.repository.0.extras", ""),
				),
			},
			{
				Config: testAccHarborCheckRetentionPolicyUpdateResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccHarborCheckProjectExists("harbor_project.default", &project),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "algorithm", "or"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "scope.0.level", "project"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "trigger.0.settings.0.cron", "0 0 0 * * *"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.id", "0"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.priority", "0"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.action", "retain"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.template", "nDaysSinceLastPush"),
					//resource.TestCheckResourceAttr(
					//	"harbor_retention_policy.default", "rule.0.params", "{\"nDaysSinceLastPush\": 1}"),
					//resource.TestCheckResourceAttr(
					//	"harbor_retention_policy.default", "rule.0.params", map[string]int{"nDaysSinceLastPush": 1}),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.tag_selectors.0.kind", "doublestar"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.tag_selectors.0.decoration", "matches"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.tag_selectors.0.pattern", "master"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.tag_selectors.0.extras", "{\"untagged\":true}"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.scope_selectors.0.repository.0.kind", "doublestar"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.scope_selectors.0.repository.0.decoration", "repoMatches"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.scope_selectors.0.repository.0.pattern", "**"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.0.scope_selectors.0.repository.0.extras", ""),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.id", "0"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.priority", "0"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.disabled", "true"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.action", "retain"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.template", "latestPulledN"),
					//resource.TestCheckResourceAttr(
					//	"harbor_retention_policy.default", "rule.1.params", "{\"latestPulledN\": 15}"),
					//resource.TestCheckResourceAttr(
					//	"harbor_retention_policy.default", "rule.1.params", map[string]int{"latestPulledN": 15}),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.tag_selectors.0.kind", "doublestar"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.tag_selectors.0.decoration", "excludes"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.tag_selectors.0.pattern", "master"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.tag_selectors.0.extras", "{\"untagged\":false}"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.scope_selectors.0.repository.0.kind", "doublestar"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.scope_selectors.0.repository.0.decoration", "repoExcludes"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.scope_selectors.0.repository.0.pattern", "nginx-*"),
					resource.TestCheckResourceAttr(
						"harbor_retention_policy.default", "rule.1.scope_selectors.0.repository.0.extras", ""),
				),
			},
		},
	})
}

func testAccHarborCheckRetentionPolicyMinimalResourceConfig() string {
	return `
resource "harbor_project" "default" {
    name                   = "acc-retention-policy-test"
    public                 = false
    vulnerability_scanning = false
}
resource "harbor_retention_policy" "default" {
    scope {
        ref = harbor_project.default.id
	}
	rule {
		template = "always"
		tag_selectors {
			decoration = "matches"
			pattern    = "**"
		}
		scope_selectors {
			repository {
				decoration = "repoMatches"
				pattern    = "**"
			}
		}
	}
	trigger {
		settings {
		}
	}
}
`
}

func testAccHarborCheckRetentionPolicyUpdateResourceConfig() string {
	return `
resource "harbor_project" "default" {
    name                   = "acc-retention-policy-test"
    public                 = false
    vulnerability_scanning = false
}
resource "harbor_retention_policy" "default" {
    scope {
        ref = harbor_project.default.id
	}
	rule {
		template = "nDaysSinceLastPush"
		params = {
			"nDaysSinceLastPush" = 1
		}
		tag_selectors {
			decoration = "matches"
			pattern    = "master"
			extras     = jsonencode({
				untagged: true
			})
		}
		scope_selectors {
			repository {
				decoration = "repoMatches"
				pattern    = "**"
			}
		}
	}
	rule {
		disabled = true
		template = "latestPulledN"
		params = {
			"latestPulledN"      = 15
			"nDaysSinceLastPush" = 7
		}
		tag_selectors {
			kind       = "doublestar"
			decoration = "excludes"
			pattern    = "master"
			extras     = jsonencode({
				untagged: false
			})
		}
		scope_selectors {
			repository {
				kind       = "doublestar"
				decoration = "repoExcludes"
				pattern    = "nginx-*"
			}
		}
	}
	trigger {
		references {
		}
		settings {
			cron = "0 0 0 * * *"
		}
	}
}
`
}
