package transfer

import (
	"fmt"
	"log"

	"github.com/reddio-com/reddio/test/pkg"
)

type ConflictTransfer struct {
	CaseName     string
	walletCount  int
	initialCount uint64
	steps        int
	tm           *pkg.TransferManager
}

func (c *ConflictTransfer) Run(m *pkg.WalletManager) error {
	wallets, err := m.GenerateRandomWallet(c.walletCount, c.initialCount)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("%s create wallets finish", c.CaseName))
	cwallets := pkg.GenerateCaseWallets(c.initialCount, wallets)
	transferCase := c.tm.GenerateSameTargetTransferSteps(c.steps, cwallets, cwallets[0])
	return runAndAssert(transferCase, m, wallets)
}

func (c *ConflictTransfer) Name() string {
	return c.CaseName
}

func NewConflictTest(name string, count int, initial uint64, steps int) *ConflictTransfer {
	return &ConflictTransfer{
		CaseName:     name,
		walletCount:  count,
		initialCount: initial,
		steps:        steps,
		tm:           pkg.NewTransferManager(),
	}
}
