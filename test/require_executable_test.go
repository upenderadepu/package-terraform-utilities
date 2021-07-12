package test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func TestRequireExecutableWorksWithExistingExecutable(t *testing.T) {
	t.Parallel()

	randomString := random.UniqueId()
	testFolder := test_structure.CopyTerraformFolderToTemp(t, "..", "examples")
	terraformModulePath := filepath.Join(testFolder, "require-executable")
	terratestOptions := createBaseTerratestOptions(t, terraformModulePath)
	terratestOptions.Vars = map[string]interface{}{
		"required_executables": []string{"go"},
		"error_message":        randomString,
	}
	defer terraform.Destroy(t, terratestOptions)

	out := terraform.InitAndApply(t, terratestOptions)
	assert.False(t, strings.Contains(out, randomString))
}

func TestRequireExecutableFailsForMissingExecutable(t *testing.T) {
	t.Parallel()

	randomString := random.UniqueId()
	testFolder := test_structure.CopyTerraformFolderToTemp(t, "..", "examples")
	terraformModulePath := filepath.Join(testFolder, "require-executable")
	terratestOptions := createBaseTerratestOptions(t, terraformModulePath)
	terratestOptions.Vars = map[string]interface{}{
		"required_executables": []string{"this-should-not-exist"},
		"error_message":        randomString,
	}
	// Swallow the error by using the E variant of Destroy, because the data source is expected to fail on both Apply
	// and Destroy
	defer terraform.DestroyE(t, terratestOptions)

	out, err := terraform.InitAndApplyE(t, terratestOptions)
	assert.Error(t, err)
	assert.True(t, strings.Contains(out, randomString))
}

func TestConditionalRequireExecutable(t *testing.T) {
	t.Parallel()

	randomString := random.UniqueId()
	testFolder := test_structure.CopyTerraformFolderToTemp(t, "..", "examples")
	terraformModulePath := filepath.Join(testFolder, "require-executable")
	terratestOptions := createBaseTerratestOptions(t, terraformModulePath)
	terratestOptions.Vars = map[string]interface{}{
		"required_executables":         []string{},
		"error_message":                "",
		"validate_bad_executable":      "1",
		"bad_executable_error_message": randomString,
	}
	defer terraform.Destroy(t, terratestOptions)

	out, err := terraform.InitAndApplyE(t, terratestOptions)
	assert.Error(t, err)
	assert.True(t, strings.Contains(out, randomString))

	terratestOptions.Vars["validate_bad_executable"] = "0"
	out = terraform.InitAndApply(t, terratestOptions)
	assert.False(t, strings.Contains(out, randomString))
}
