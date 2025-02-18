package ethrpc

import "github.com/prometheus/client_golang/prometheus"

const (
	TypeLbl       = "type"
	TypeCountLbl  = "count"
	TypeStatusLbl = "status"
)

var (
	EthApiBackendCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reddio",
			Subsystem: "eth_api_backend",
			Name:      "op_counter",
			Help:      "Total Operator number of counter",
		},
		[]string{TypeLbl},
	)

	EthApiBackendDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "eth_api_backend",
			Name:      "op_duration_microseconds",
			Help:      "Operation Duration",
			Buckets:   prometheus.ExponentialBuckets(10, 2, 20), // 10us ~ 5s
		},
		[]string{TypeLbl},
	)
)

func init() {
	prometheus.MustRegister(EthApiBackendCounter)
	prometheus.MustRegister(EthApiBackendDuration)
}
