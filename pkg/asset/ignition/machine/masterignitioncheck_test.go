package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// TestMasterIgnitionCheckGenerate tests generating the master ignition check asset.
func TestMasterIgnitionCheckGenerate(t *testing.T) {
	installConfig := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			Networking: &types.Networking{
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
			},
			Platform: types.Platform{
				AWS: &aws.Platform{
					Region: "us-east",
				},
			},
		},
	}

	rootCA := &tls.RootCA{}
	err := rootCA.Generate(nil)
	assert.NoError(t, err, "unexpected error generating root CA")

	parents := asset.Parents{}
	parents.Add(installConfig, rootCA)

	master := &Master{}
	err = master.Generate(parents)
	assert.NoError(t, err, "unexpected error generating master asset")

	parents.Add(master)
	masterIgnCheck := &MasterIgnitionCheck{}
	err = masterIgnCheck.Generate(parents)
	assert.NoError(t, err, "unexpected error generating master ignition check asset")

	actualFiles := masterIgnCheck.Files()
	assert.Equal(t, 1, len(actualFiles), "unexpected number of files in master state")
	assert.Equal(t, masterMachineConfigFileName, actualFiles[0].Filename, "unexpected name for master ignition config")
}
