package e2e

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	test_helper "github.com/Azure/terraform-module-test-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestExamplesBasic(t *testing.T) {
	t.Skip()
	createPublicIp := []bool{
		false, true,
	}
	for _, publicIp := range createPublicIp {
		t.Run(fmt.Sprintf("createPublicIp-%t", publicIp), func(t *testing.T) {
			vars := map[string]interface{}{
				"create_public_ip": publicIp,
			}
			managedIdentityId := os.Getenv("MSI_ID")
			if managedIdentityId != "" {
				vars["managed_identity_principal_id"] = managedIdentityId
			}
			test_helper.RunE2ETest(t, "../../", "examples/basic", terraform.Options{
				Upgrade: true,
				Vars:    vars,
			}, func(t *testing.T, output test_helper.TerraformOutput) {
				vmIdRegex := `/subscriptions/.+/resourceGroups/.+/providers/Microsoft.Compute/virtualMachines/.+`
				linuxVmId, ok := output["linux_vm_id"]
				assert.True(t, ok)
				assert.Regexp(t, vmIdRegex, linuxVmId)
				windowsVmId, ok := output["windows_vm_id"]
				assert.True(t, ok)
				assert.Regexp(t, vmIdRegex, windowsVmId)
				if publicIp {
					linuxPublicIp, ok := output["linux_public_ip"].(string)
					assert.True(t, ok)
					windowsPublicIp, ok := output["windows_public_ip"].(string)
					assert.True(t, ok)
					ipRegex := `((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}`
					assert.Regexp(t, ipRegex, linuxPublicIp)
					assert.Regexp(t, ipRegex, windowsPublicIp)
					nsgIdRegex := `/subscriptions/.+/resourceGroups/.+/providers/Microsoft.Network/networkSecurityGroups/.+`
					linuxNsgId, ok := output["linux_network_security_group_id"].(string)
					assert.True(t, ok)
					assert.Regexp(t, nsgIdRegex, linuxNsgId)
					windowsNsgId, ok := output["windows_network_security_group_id"].(string)
					assert.True(t, ok)
					assert.Regexp(t, nsgIdRegex, windowsNsgId)
				}
			})
		})
	}
}

func TestExamplesVmss(t *testing.T) {
	test_helper.RunE2ETest(t, "../../", "examples/vmss", terraform.Options{
		Upgrade: true,
	}, func(t *testing.T, output test_helper.TerraformOutput) {
		vmssIdRegex := `/subscriptions/.+/resourceGroups/.+/providers/Microsoft.Compute/virtualMachineScaleSets/.+`
		vmssId, ok := output["linux_vm_vmss_id"]
		require.True(t, ok)
		require.Regexp(t, vmssIdRegex, vmssId)
	})
}
