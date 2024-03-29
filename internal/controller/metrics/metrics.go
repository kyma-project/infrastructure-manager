package metrics

import (
	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/prometheus/client_golang/prometheus"
	ctrlMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	runtimeIDKeyName               = "runtimeId"
	state                          = "state"
	reason                         = "reason"
	componentName                  = "infrastructure_manager"
	RuntimeIDLabel                 = "kyma-project.io/runtime-id"
	GardenerClusterStateMetricName = "im_gardener_clusters_state"
)

type Metrics struct {
	gardenerClustersStateGaugeVec *prometheus.GaugeVec
}

func NewMetrics() Metrics {
	m := Metrics{
		gardenerClustersStateGaugeVec: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: componentName,
				Name:      GardenerClusterStateMetricName,
				Help:      "Indicates the Status.state for GardenerCluster CRs",
			}, []string{runtimeIDKeyName, state, reason}),
	}
	ctrlMetrics.Registry.MustRegister(m.gardenerClustersStateGaugeVec)
	return m
}

func (m Metrics) SetGardenerClusterStates(cluster v1.GardenerCluster) {
	var runtimeID = cluster.GetLabels()[RuntimeIDLabel]

	if runtimeID != "" {
		if len(cluster.Status.Conditions) != 0 {
			var reason = cluster.Status.Conditions[0].Reason

			// first clean the old metric
			m.CleanUpGardenerClusterGauge(runtimeID)
			m.gardenerClustersStateGaugeVec.WithLabelValues(runtimeID, string(cluster.Status.State), reason).Set(1)
		}
	}
}

func (m Metrics) CleanUpGardenerClusterGauge(runtimeID string) {
	m.gardenerClustersStateGaugeVec.DeletePartialMatch(prometheus.Labels{
		runtimeIDKeyName: runtimeID,
	})
}
