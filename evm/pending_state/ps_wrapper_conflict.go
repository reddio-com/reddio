package pending_state

import "github.com/ethereum/go-ethereum/common"

func (sctx *StateContext) WriteConflict(addr common.Address, txnID int64) bool {
	wa := sctx.Write.Address[addr]
	ra := sctx.Read.Address[addr]
	if checkVisitedTxnConflict(txnID, wa) || checkVisitedTxnConflict(txnID, ra) {
		return true
	}
	return false
}

func (sctx *StateContext) ReadConflict(addr common.Address, txnID int64) bool {
	ra := sctx.Read.Address[addr]
	return checkVisitedTxnConflict(txnID, ra)
}

func checkVisitedTxnConflict(txnID int64, w VisitTxnID) bool {
	if len(w) < 1 {
		return false
	}
	if len(w) == 1 {
		_, ok := w[txnID]
		if ok {
			return false
		}
	}
	return true
}
