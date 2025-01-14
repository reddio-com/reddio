package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	SolidityCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reddio",
			Subsystem: "solidity",
			Name:      "op_counter",
			Help:      "Total Operator number of counter",
		},
		[]string{TypeLbl, TypeStatusLbl},
	)

	SolidityHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "reddio",
		Subsystem: "solidity",
		Name:      "op_execute_hist",
		Help:      "solidity operation execute duration distribution.",
		Buckets:   prometheus.ExponentialBuckets(10, 2, 20), // 10us ~ 5s
	}, []string{TypeLbl})
)

func init() {
	prometheus.MustRegister(SolidityCounter)
	prometheus.MustRegister(SolidityHist)
}
