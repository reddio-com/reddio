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
)

func init() {
	prometheus.MustRegister(EthApiBackendCounter)
}
