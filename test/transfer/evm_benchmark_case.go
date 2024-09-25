package transfer

import (
	"context"
	"github.com/reddio-com/reddio/test/pkg"
	"golang.org/x/time/rate"
)

type RandomEVMBenchmarkTest struct {
	CaseName     string
	walletCount  int
	initialCount uint64
	steps        int
	tm           *pkg.TransferManager
	wallets      []*pkg.EthWallet
	rm           *rate.Limiter
}

func NewRandomEVMBenchmarkTest(name string, count int, initial uint64, steps int, wallets []*pkg.EthWallet, rm *rate.Limiter) *RandomEVMBenchmarkTest {
	return &RandomEVMBenchmarkTest{
		CaseName:     name,
		walletCount:  count,
		initialCount: initial,
		steps:        steps,
		tm:           pkg.NewTransferManager(),
		wallets:      wallets,
		rm:           rm,
	}
}

func (tc *RandomEVMBenchmarkTest) Name() string {
	return tc.CaseName
}

func (tc *RandomEVMBenchmarkTest) Run(ctx context.Context, m *pkg.WalletManager) error {
	transferCase := tc.tm.GenerateRandomTransferSteps(tc.steps, pkg.GenerateCaseWallets(tc.initialCount, tc.wallets))
	if err := tc.rm.WaitN(ctx, tc.steps); err != nil {
		return err
	}
	return m.TransferEthForEVM(transferCase.Steps)
}

func (tc *RandomEVMBenchmarkTest) BatchRun(ctx context.Context, m *pkg.WalletManager) error {
	panic("implement me")
}
