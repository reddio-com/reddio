package api

import (
	"sync"

	"gorm.io/gorm"
)

var (

	// L2UnclaimedWithdrawalsByAddressCtl the L2UnclaimedWithdrawalsByAddressController instance
	L2UnclaimedWithdrawalsByAddressCtl *L2UnclaimedWithdrawalsByAddressController
	// TxsByAddressCtl the TxsByAddressController instance
	TxsByAddressCtl *TxsByAddressController

	// L2WithdrawalsByAddressCtl the L2WithdrawalsByAddressController instance
	initControllerOnce sync.Once
)

// InitController inits Controller with database
func InitController(db *gorm.DB) {
	initControllerOnce.Do(func() {
		TxsByAddressCtl = NewTxsByAddressController(db)
		L2UnclaimedWithdrawalsByAddressCtl = NewL2UnclaimedWithdrawalsByAddressController(db)

	})
}
