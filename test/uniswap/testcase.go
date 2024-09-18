package uniswap

import (
	"context"

	"github.com/reddio-com/reddio/test/pkg"
)

type TestCase interface {
	Prepare(ctx context.Context, m *pkg.WalletManager) error
	Run(ctx context.Context, m *pkg.WalletManager) error
	Name() string
}
