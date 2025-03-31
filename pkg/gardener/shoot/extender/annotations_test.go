package extender

import (
	"github.com/kyma-project/infrastructure-manager/pkg/gardener/shoot/extender/testutils"
	"testing"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAnnotationsExtender(t *testing.T) {
	licenceType := "licence"

	for _, testCase := range []struct {
		name                string
		runtime             imv1.Runtime
		expectedAnnotations map[string]string
	}{
		{
			name: "Create with basic annotations",
			runtime: imv1.Runtime{
				ObjectMeta: v1.ObjectMeta{
					Name:      "runtime",
					Namespace: "namespace",
					Labels: map[string]string{
						"kyma-project.io/runtime-id": "runtime-id",
					},
					Generation: 100,
				},
			},
			expectedAnnotations: map[string]string{
				"infrastructuremanager.kyma-project.io/runtime-id":         "runtime-id",
				"infrastructuremanager.kyma-project.io/runtime-generation": "100"},
		},
		{
			name: "Create licence type annotation",
			runtime: imv1.Runtime{
				ObjectMeta: v1.ObjectMeta{
					Name:      "runtime",
					Namespace: "namespace",
					Labels: map[string]string{
						"kyma-project.io/runtime-id": "runtime-id",
					},
				},
				Spec: imv1.RuntimeSpec{
					Shoot: imv1.RuntimeShoot{
						LicenceType: &licenceType,
					},
				},
			},
			expectedAnnotations: map[string]string{
				"infrastructuremanager.kyma-project.io/runtime-id":         "runtime-id",
				"infrastructuremanager.kyma-project.io/licence-type":       "licence",
				"infrastructuremanager.kyma-project.io/runtime-generation": "0"},
		},
	} {
		// given
		shoot := testutils.FixEmptyGardenerShoot("shoot", "kcp-system")

		// when
		err := ExtendWithAnnotations(testCase.runtime, &shoot)
		require.NoError(t, err)

		// then
		assert.Equal(t, testCase.expectedAnnotations, shoot.Annotations)
	}
}
