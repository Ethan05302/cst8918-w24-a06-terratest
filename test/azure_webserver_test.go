package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAzureLinuxVMCreation(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix":    "ZheZhang",
			"admin_username": "azureuser", // 根据你实际设置填写
		},
		Upgrade: true,
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Retry mechanism to wait for Public IP to be available
	maxRetries := 10
	timeBetweenRetries := 15 * time.Second

	publicIP := ""
	_, err := retry.DoWithRetryE(t, "Get Public IP", maxRetries, timeBetweenRetries, func() (string, error) {
		publicIP = terraform.Output(t, terraformOptions, "public_ip")
		if publicIP == "" {
			return "", assert.AnError
		}
		return "OK", nil
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, publicIP)

	// Also check VM name and resource group
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroup := terraform.Output(t, terraformOptions, "resource_group_name")

	assert.NotEmpty(t, vmName)
	assert.NotEmpty(t, resourceGroup)

	t.Logf("VM Name: %s", vmName)
	t.Logf("Resource Group: %s", resourceGroup)
	t.Logf("Public IP: %s", publicIP)
}
