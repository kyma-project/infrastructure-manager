package metrics

import (
	"fmt"
	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	ctrlMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	runtimeId = "runtimeId"
	state     = "state"
)

var (

	//nolint:godox //TODO: test custom metric, remove when done with https://github.com/kyma-project/infrastructure-manager/issues/11
	playgroundTotalReconciliationLoopsStarted = prometheus.NewCounter( //nolint:gochecknoglobals
		prometheus.CounterOpts{
			Name: "im_playground_reconciliation_loops_started_total",
			Help: "Number of times reconciliation loop was started",
		},
	)

	metricGardenerClustersState = prometheus.NewGaugeVec( //nolint:gochecknoglobals
		prometheus.GaugeOpts{ //nolint:gochecknoglobals
			Subsystem: "infrastructure_manager",
			Name:      "im_gardener_clusters_state",
			Help:      "Indicates the Status.state for GardenerCluster CRs",
		}, []string{runtimeId, state})
)

func init() {
	ctrlMetrics.Registry.MustRegister(playgroundTotalReconciliationLoopsStarted, metricGardenerClustersState)
}

func IncrementReconciliationLoopsStarted() {
	playgroundTotalReconciliationLoopsStarted.Inc()
}

func SetGardenerClusterStates(cluster v1.GardenerCluster) {
	metricGardenerClustersState.WithLabelValues(cluster.Name, string(cluster.Status.State)).Set(1)
}

func UnSetGardenerClusterStates(secret corev1.Secret) {
	var runtimeId = secret.GetLabels()["kyma-project.io/runtime-id"]
	var deletedReady = metricGardenerClustersState.DeleteLabelValues(runtimeId, "Ready")
	var deletedError = metricGardenerClustersState.DeleteLabelValues(runtimeId, "Error")

	if deletedReady || deletedError {
		fmt.Printf("GardenerClusterStates deleted value for %v", runtimeId)
	}
}
