// +build integration

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestKubernetesNginxNamespaceExists(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir: "../tf-acception-test",
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "harbor_project_id")
	assert.NotNil(t, output)

}

func TestHarborSystemConfig(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir: "../tf-project-only",
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

}
