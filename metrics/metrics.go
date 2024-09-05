package metrics

import "github.com/prometheus/client_golang/prometheus"

const TypeLbl = "type"

var (
	TxnCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reddio",
			Subsystem: "txn",
			Name:      "total_count",
			Help:      "Total number of count for txn",
		},
		[]string{TypeLbl},
	)

	TxnDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "txn",
			Name:      "execute_duration_seconds",
			Help:      "txn execute duration distribution.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{},
	)

	BatchTxnCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reddio",
			Subsystem: "batch_txn",
			Name:      "total_count",
			Help:      "Total number of redo count for batch txn",
		},
		[]string{TypeLbl},
	)

	BatchTxnDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "batch_txn",
			Name:      "execute_duration_seconds",
			Help:      "txn execute duration distribution.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{TypeLbl},
	)

	StatedbCopyDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "batch_txn",
			Name:      "statedb_copy",
			Help:      "stateDB copy duration distribution.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{},
	)

	InternalTPS = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "reddio",
		Subsystem: "batch_txn",
		Name:      "internal_TPS",
		Help:      "EVM invocation TPS",
	})
)

func init() {
	prometheus.MustRegister(TxnCounter)
	prometheus.MustRegister(TxnDuration)
	prometheus.MustRegister(BatchTxnCounter)
	prometheus.MustRegister(BatchTxnDuration)
	prometheus.MustRegister(StatedbCopyDuration)
	prometheus.MustRegister(InternalTPS)
}
