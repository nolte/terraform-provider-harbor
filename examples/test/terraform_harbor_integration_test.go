// +build integration

package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestHarborBaseComponentsExists(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir:       "../tf-acception-test",
		Parallelism:        1,
		MaxRetries:         5,
		TimeBetweenRetries: time.Second + 10,
	}
	terraformOptionsPartTwo := &terraform.Options{
		TerraformDir:       "../tf-acception-test-part-2",
		Parallelism:        1,
		MaxRetries:         5,
		TimeBetweenRetries: time.Second + 10,
	}

	defer terraform.Destroy(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptionsPartTwo)

	terraform.InitAndApply(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptionsPartTwo)

	output := terraform.Output(t, terraformOptions, "harbor_project_id")
	assert.NotNil(t, output)

}

func TestHarborSystemConfig(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir: "../tf-project-only",
		Parallelism:  1,
		MaxRetries:   5,
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

}
