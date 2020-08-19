// +build integration

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func testsBaseDirectory() string {
	accPath := os.Getenv("HARBOR_ACC_BASEPATH")
	if accPath == "" {
		accPath = "tf-v13"
	}
	return accPath
}

func TestHarborBaseComponentsExists(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir: filepath.Join(testsBaseDirectory(), "tf-acception-test"),
		//Parallelism:        1,
		//MaxRetries:         5,
		//TimeBetweenRetries: time.Second + 10,
	}
	terraformOptionsPartTwo := &terraform.Options{
		TerraformDir: filepath.Join(testsBaseDirectory(), "tf-acception-test-part-2"),
		//Parallelism:        1,
		//MaxRetries:         5,
		//TimeBetweenRetries: time.Second + 10,
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
		TerraformDir: filepath.Join(testsBaseDirectory(), "tf-project-only"),
		//Parallelism:  1,
		//MaxRetries:   5,
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
}
