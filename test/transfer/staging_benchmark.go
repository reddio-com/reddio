package transfer

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/test/pkg"
)

type StagingBenchmark struct {
	initialCount uint64
	tm           *pkg.TransferManager
	wallets      []*pkg.EthWallet
	rm           *rate.Limiter
}

func (s *StagingBenchmark) Run(ctx context.Context, m *pkg.WalletManager) error {
	transferCase := s.tm.GenerateSwapTransferSteps(pkg.GenerateCaseWallets(s.initialCount, s.wallets))
	for i, step := range transferCase.Steps {
		if err := s.rm.Wait(ctx); err != nil {
			return err
		}
		if err := m.TransferEth(step.From, step.To, step.Count, uint64(i)+uint64(time.Now().UnixNano())); err != nil {
			logrus.Error("Failed to transfer step: from:%v, to:%v", step.From, step.To)
		}
	}
	return nil
}

func (s *StagingBenchmark) Name() string {
	return "staging_benchmark"
}

func NewStagingBenchmark(wallets []*pkg.EthWallet, rm *rate.Limiter) *StagingBenchmark {
	return &StagingBenchmark{
		initialCount: 2,
		tm:           pkg.NewTransferManager(),
		wallets:      wallets,
		rm:           rm,
	}
}
