package parallel

import (
	"time"
	
	"github.com/sirupsen/logrus"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/metrics"
)

type BlockTxnStatManager struct {
	TxnCount           int
	TxnBatchCount      int
	TxnBatchRedoCount  int
	ConflictCount      int
	ExecuteDuration    time.Duration
	ExecuteTxnDuration time.Duration
	PrepareDuration    time.Duration
	CommitDuration     time.Duration
	CopyDuration       time.Duration
}

func (stat *BlockTxnStatManager) UpdateMetrics() {
	metrics.BlockExecuteTxnCountGauge.WithLabelValues().Set(float64(stat.TxnCount))
	metrics.BlockExecuteTxnDurationGauge.WithLabelValues().Set(float64(stat.ExecuteDuration.Seconds()))
	metrics.BlockTxnAllExecuteDurationGauge.WithLabelValues().Set(float64(stat.ExecuteTxnDuration.Seconds()))
	metrics.BlockTxnPrepareDurationGauge.WithLabelValues().Set(float64(stat.PrepareDuration.Seconds()))
	metrics.BlockTxnCommitDurationGauge.WithLabelValues().Set(float64(stat.CommitDuration.Seconds()))
	if config.GlobalConfig.IsBenchmarkMode {
		logrus.Infof("execute %v txn, total:%v, execute cost:%v, prepare:%v, copy:%v, commit:%v, txnBatch:%v, conflict:%v, redoBatch:%v",
			stat.TxnCount, stat.ExecuteDuration.String(), stat.ExecuteTxnDuration.String(),
			stat.PrepareDuration.String(), stat.CopyDuration.String(), stat.CommitDuration.String(), stat.TxnBatchCount, stat.ConflictCount, stat.TxnBatchRedoCount)
	}
}
