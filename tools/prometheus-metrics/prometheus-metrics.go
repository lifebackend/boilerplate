package prometheusmetrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	lock sync.Mutex
	m    *Metrics
)

type Metrics struct {
	apiResponseDuration *prometheus.HistogramVec

	version *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	lock.Lock()
	defer lock.Unlock()

	if m != nil {
		return m
	}

	m = &Metrics{
		apiResponseDuration: prometheus.NewHistogramVec(
			// nolint:exhaustivestruct
			prometheus.HistogramOpts{
				Name: "http_api_response_duration_seconds",
				Help: "HTTP API response duration",
			},
			[]string{"resource"},
		),
		version: prometheus.NewCounterVec(
			// nolint:exhaustivestruct,promlinter
			prometheus.CounterOpts{
				Name: "wallet_pod_version_count",
				Help: "Wallet pod version count",
			},
			[]string{"version"},
		),
	}

	prometheus.MustRegister(
		m.apiResponseDuration,
		m.version,
	)

	return m
}

func (m *Metrics) AddAPIResponseDuration(resource string, duration time.Duration) {
	m.apiResponseDuration.With(prometheus.Labels{"resource": resource}).Observe(duration.Seconds())
}

func (m *Metrics) SetPodVersion(version string) {
	m.version.With(prometheus.Labels{"version": version}).Inc()
}