package transfer

import (
	"context"

	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/test/pkg"
)

type RandomBenchmarkTest struct {
	CaseName     string
	walletCount  int
	initialCount uint64
	steps        int
	tm           *pkg.TransferManager
	wallets      []*pkg.EthWallet
	rm           *rate.Limiter
}

func NewRandomBenchmarkTest(name string, count int, initial uint64, steps int, wallets []*pkg.EthWallet, rm *rate.Limiter) *RandomBenchmarkTest {
	return &RandomBenchmarkTest{
		CaseName:     name,
		walletCount:  count,
		initialCount: initial,
		steps:        steps,
		tm:           pkg.NewTransferManager(),
		wallets:      wallets,
		rm:           rm,
	}
}

func (tc *RandomBenchmarkTest) Name() string {
	return tc.CaseName
}

func (tc *RandomBenchmarkTest) Run(ctx context.Context, m *pkg.WalletManager) error {
	transferCase := tc.tm.GenerateRandomTransferSteps(tc.steps, pkg.GenerateCaseWallets(tc.initialCount, tc.wallets))
	for _, step := range transferCase.Steps {
		if err := tc.rm.Wait(ctx); err != nil {
			return err
		}
		m.TransferEth(step.From, step.To, step.Count)
	}
	return nil
}
