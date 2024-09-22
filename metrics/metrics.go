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

	BlockTxnCommitDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "block_txn",
			Name:      "commit_duration_seconds",
			Help:      "txn commit duration distribution.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{},
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

	BatchTxnStatedbCopyDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "batch_txn",
			Name:      "statedb_copy_seconds",
			Help:      "stateDB copy duration per block distribution.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{TypeCountLbl},
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

	BlockExecuteTxnDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "block",
			Name:      "execute_duration_seconds",
			Help:      "block execute txn duration",
			Buckets:   prometheus.DefBuckets,
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

	BlockTxnSplitDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "block_txn",
			Name:      "split_txn_duration_seconds",
			Help:      "split batch txn duration",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{},
	)

	BlockTxnPrepareDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "block_txn",
			Name:      "prepare_txn_duration_seconds",
			Help:      "split batch txn duration",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{},
	)

	BatchTxnAllExecuteDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "reddio",
			Subsystem: "batch_txn",
			Name:      "all_execute_duration_seconds",
			Help:      "split batch txn duration",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(TxnCounter)
	prometheus.MustRegister(TxnDuration)

	prometheus.MustRegister(BlockExecuteTxnCountGauge)
	prometheus.MustRegister(BlockTxnPrepareDuration)
	prometheus.MustRegister(BlockTxnSplitDuration)
	prometheus.MustRegister(BatchTxnAllExecuteDuration)
	prometheus.MustRegister(BlockTxnCommitDuration)

	prometheus.MustRegister(BlockExecuteTxnDuration)

	prometheus.MustRegister(BatchTxnCounter)
	prometheus.MustRegister(BatchTxnSplitCounter)
	prometheus.MustRegister(BatchTxnDuration)
	prometheus.MustRegister(BatchTxnStatedbCopyDuration)
}
