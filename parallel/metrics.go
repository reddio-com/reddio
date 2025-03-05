package parallel

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/reddio-com/reddio/metrics"
)

var (
	BlockTxnCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "reddio",
		Subsystem: "block_txn",
		Name:      "counter",
		Help:      "counter of block txn",
	}, []string{metrics.TypeLbl})
)

func init() {
	prometheus.MustRegister(BlockTxnCounter)
}
