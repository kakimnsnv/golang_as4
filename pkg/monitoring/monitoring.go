package monitoring

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "https_requests_total", Help: "Total number of HTTPS requests",
		},
		[]string{"path", "method", "status"},
	)

	RequestDuration = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "https_request_duration_seconds",
			Help: "Duration of HTTPS requests",
		},
		[]string{"path"},
	)

	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "https_errors_total", Help: "Total number of HTTPS errors",
		},
		[]string{"path", "status"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(ErrorCounter)
}
