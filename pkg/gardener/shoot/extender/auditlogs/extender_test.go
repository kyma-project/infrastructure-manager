package auditlogs

import (
	"testing"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/require"
)

func Test_AuditlogExtender(t *testing.T) {
	var zero imv1.Runtime
	for _, tc := range []struct {
		shoot               gardener.Shoot
		data                AuditLogData
		policyConfigmapName string
	}{
		{
			shoot: gardener.Shoot{},
			data: AuditLogData{
				TenantID:   "tenant-id",
				ServiceURL: "testme",
			},
		},
	} {
		// given
		extendWithAuditlogs := NewAuditlogExtenderForCreate(tc.policyConfigmapName, tc.data)

		// when
		err := extendWithAuditlogs(zero, &tc.shoot)

		// then
		require.NoError(t, err)
	}
}
