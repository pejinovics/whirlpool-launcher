package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	ProbeChecks = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wh_probe_checks_total",
			Help: "Number of probe checks performed.",
		},
		[]string{"probe", "target", "container", "probe_type"},
	)

	ProbeSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wh_probe_success_total",
			Help: "Number of successful probe checks.",
		},
		[]string{"probe", "target", "container", "probe_type"},
	)

	ProbeFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wh_probe_failures_total",
			Help: "Number of failed probe checks.",
		},
		[]string{"probe", "target", "container", "probe_type"},
	)

	ContainerRestarts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wh_container_restarts_total",
			Help: "Number of container restarts triggered by failed probes.",
		},
		[]string{"probe", "target", "container", "probe_type"},
	)
)

func Register(reg prometheus.Registerer) {
	reg.MustRegister(ProbeChecks, ProbeSuccess, ProbeFailures, ContainerRestarts)
}

func BuildLabels(probe, target, container, probeType string) prometheus.Labels {
	return prometheus.Labels{
		"probe":      probe,
		"target":     target,
		"container":  container,
		"probe_type": probeType,
	}
}

func ObserveResult(labels prometheus.Labels, success bool, dur time.Duration, consecFails int) {
	ProbeChecks.With(labels).Inc()
	if success {
		ProbeSuccess.With(labels).Inc()
	} else {
		ProbeFailures.With(labels).Inc()
	}
}

func IncRestart(labels prometheus.Labels) {
	ContainerRestarts.With(labels).Inc()
}
