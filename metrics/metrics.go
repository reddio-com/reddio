package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	TypeLbl      = "type"
	TypeCountLbl = "count"
)

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
			Buckets:   TxnBuckets,
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

	BlockTxnCommitDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "reddio",
			Subsystem: "block_txn",
			Name:      "commit_duration_seconds",
			Help:      "txn commit duration seconds",
		},
		[]string{},
	)

	BatchTxnDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "batch_txn",
			Name:      "execute_duration_seconds",
			Help:      "txn execute duration distribution.",
			Buckets:   TxnBuckets,
		},
		[]string{TypeLbl},
	)

	BatchTxnSplitCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reddio",
			Subsystem: "batch_txn",
			Name:      "split_txn_count",
			Help:      "split sub batch txn count",
		},
		[]string{TypeCountLbl},
	)

	BlockExecuteTxnDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "reddio",
			Subsystem: "block",
			Name:      "execute_duration_seconds",
			Help:      "block execute txn duration",
		},
		[]string{},
	)

	BlockExecuteTxnCountGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "reddio",
			Subsystem: "block",
			Name:      "execute_txn_count",
			Help:      "txn count for each block",
		}, []string{})

	BlockTxnPrepareDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "reddio",
			Subsystem: "block_txn",
			Name:      "prepare_txn_duration_seconds",
			Help:      "split batch txn duration",
		},
		[]string{},
	)

	BlockTxnAllExecuteDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "reddio",
			Subsystem: "block_txn",
			Name:      "execute_all_duration_seconds",
			Help:      "split batch txn duration",
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(TxnCounter)
	prometheus.MustRegister(TxnDuration)

	prometheus.MustRegister(BlockExecuteTxnCountGauge)
	prometheus.MustRegister(BlockTxnPrepareDurationGauge)
	prometheus.MustRegister(BlockTxnAllExecuteDurationGauge)
	prometheus.MustRegister(BlockTxnCommitDurationGauge)
	prometheus.MustRegister(BlockExecuteTxnDurationGauge)

	prometheus.MustRegister(BatchTxnCounter)
	prometheus.MustRegister(BatchTxnSplitCounter)
	prometheus.MustRegister(BatchTxnDuration)
}

var TxnBuckets = []float64{.00005, .0001, .00025, .0005, .001, .0025, .005, 0.01, 0.025, 0.05, 0.1}
