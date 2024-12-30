package pending_state

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

type PendingStateWrapper struct {
	sync.RWMutex
	statedb *state.StateDB
	TxnID   int64
	sCtx    *StateContext
	logs    []*types.Log
}

func NewPendingStateWrapper(statedb *state.StateDB, TxnID int64) *PendingStateWrapper {
	psw := &PendingStateWrapper{
		statedb: statedb,
		TxnID:   TxnID,
		sCtx:    NewStateContext(),
		logs:    make([]*types.Log, 0),
	}
	return psw
}

func (psw *PendingStateWrapper) CreateAccount(address common.Address) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.WriteAccount(address, psw.TxnID)
	psw.statedb.CreateAccount(address)
}

func (psw *PendingStateWrapper) SubBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.WriteBalance(address, psw.TxnID)
	psw.statedb.SubBalance(address, u, reason)
}

func (psw *PendingStateWrapper) AddBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.WriteBalance(address, psw.TxnID)
	psw.statedb.AddBalance(address, u, reason)
}

func (psw *PendingStateWrapper) GetBalance(address common.Address) *uint256.Int {
	psw.RLock()
	defer psw.RUnlock()
	psw.sCtx.ReadBalance(address, psw.TxnID)
	return psw.statedb.GetBalance(address)
}

func (psw *PendingStateWrapper) GetNonce(address common.Address) uint64 {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.GetNonce(address)
}

func (psw *PendingStateWrapper) SetNonce(address common.Address, u uint64) {
	psw.Lock()
	defer psw.Unlock()
	psw.statedb.SetNonce(address, u)
}

func (psw *PendingStateWrapper) GetCodeHash(address common.Address) common.Hash {
	psw.RLock()
	defer psw.RUnlock()
	psw.sCtx.ReadCode(address, psw.TxnID)
	return psw.statedb.GetCodeHash(address)
}

func (psw *PendingStateWrapper) GetCode(address common.Address) []byte {
	psw.RLock()
	defer psw.RUnlock()
	psw.sCtx.ReadCode(address, psw.TxnID)
	return psw.statedb.GetCode(address)
}

func (psw *PendingStateWrapper) SetCode(address common.Address, bytes []byte) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.WriteCode(address, psw.TxnID)
	psw.statedb.SetCode(address, bytes)
}

func (psw *PendingStateWrapper) GetCodeSize(address common.Address) int {
	psw.RLock()
	defer psw.RUnlock()
	psw.sCtx.ReadCode(address, psw.TxnID)
	return psw.statedb.GetCodeSize(address)
}

func (psw *PendingStateWrapper) AddRefund(u uint64) {
	psw.Lock()
	defer psw.Unlock()
	psw.statedb.AddRefund(u)
}

func (psw *PendingStateWrapper) SubRefund(u uint64) {
	psw.Lock()
	defer psw.Unlock()
	psw.statedb.SubRefund(u)
}

func (psw *PendingStateWrapper) GetRefund() uint64 {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.GetRefund()
}

func (psw *PendingStateWrapper) GetCommittedState(address common.Address, hash common.Hash) common.Hash {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.GetCommittedState(address, hash)
}

func (psw *PendingStateWrapper) GetState(address common.Address, hash common.Hash) common.Hash {
	psw.RLock()
	defer psw.RUnlock()
	psw.sCtx.ReadState(address, hash, psw.TxnID)
	return psw.statedb.GetState(address, hash)
}

func (psw *PendingStateWrapper) SetState(address common.Address, hash common.Hash, hash2 common.Hash) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.WriteState(address, hash, psw.TxnID)
	psw.statedb.SetState(address, hash, hash2)
}

func (psw *PendingStateWrapper) GetStorageRoot(addr common.Address) common.Hash {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.GetStorageRoot(addr)
}

func (psw *PendingStateWrapper) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.GetTransientState(addr, key)
}

func (psw *PendingStateWrapper) SetTransientState(addr common.Address, key, value common.Hash) {
	psw.Lock()
	defer psw.Unlock()
	psw.statedb.SetTransientState(addr, key, value)
}

func (psw *PendingStateWrapper) SelfDestruct(address common.Address) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.SelfDestruct(address, psw.TxnID)
	psw.statedb.SelfDestruct(address)
}

func (psw *PendingStateWrapper) HasSelfDestructed(address common.Address) bool {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.HasSelfDestructed(address)
}

func (psw *PendingStateWrapper) Selfdestruct6780(address common.Address) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.SelfDestruct(address, psw.TxnID)
	psw.statedb.Selfdestruct6780(address)
}

func (psw *PendingStateWrapper) Exist(address common.Address) bool {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.Exist(address)
}

func (psw *PendingStateWrapper) Empty(address common.Address) bool {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.Empty(address)
}

func (psw *PendingStateWrapper) AddressInAccessList(addr common.Address) bool {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.AddressInAccessList(addr)
}

func (psw *PendingStateWrapper) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb.SlotInAccessList(addr, slot)
}

func (psw *PendingStateWrapper) AddAddressToAccessList(addr common.Address) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.AddAddressToList(addr)
	psw.statedb.AddAddressToAccessList(addr)
}

func (psw *PendingStateWrapper) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.AddSlot2Address(slotToAddress{
		addr: addr,
		slot: slot,
	})
	psw.statedb.AddSlotToAccessList(addr, slot)
}

func (psw *PendingStateWrapper) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	psw.Lock()
	defer psw.Unlock()
	psw.sCtx.SetPrepare(rules, sender, coinbase, dest, precompiles, txAccesses)
	psw.statedb.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
}

func (psw *PendingStateWrapper) RevertToSnapshot(i int) {
	psw.Lock()
	defer psw.Unlock()
	psw.statedb.RevertToSnapshot(i)
}

func (psw *PendingStateWrapper) Snapshot() int {
	psw.Lock()
	defer psw.Unlock()
	return psw.statedb.Snapshot()
}

func (psw *PendingStateWrapper) AddLog(log *types.Log) {
	psw.Lock()
	defer psw.Unlock()
	psw.logs = append(psw.logs, log)
	psw.statedb.AddLog(log)
}

func (psw *PendingStateWrapper) AddPreimage(hash common.Hash, bytes []byte) {
	psw.Lock()
	defer psw.Unlock()
	psw.statedb.AddPreimage(hash, bytes)
}

func (psw *PendingStateWrapper) SetTxContext(txHash common.Hash, txIndex int) {
	psw.Lock()
	defer psw.Unlock()
	psw.statedb.SetTxContext(txHash, txIndex)
}

func (psw *PendingStateWrapper) GetStateDB() *state.StateDB {
	psw.RLock()
	defer psw.RUnlock()
	return psw.statedb
}

func (psw *PendingStateWrapper) GetCtx() *StateContext {
	return psw.sCtx
}

func (psw *PendingStateWrapper) AllLogs() []*types.Log {
	return psw.logs
}

func (psw *PendingStateWrapper) MergeInto(stateDB *state.StateDB, sender common.Address) {
	if psw.sCtx.prepareParams != nil {
		pre := psw.sCtx.prepareParams
		stateDB.Prepare(pre.rules, pre.sender, pre.coinbase, pre.dest, pre.precompiles, pre.txAccesses)
	}
	for addr := range psw.sCtx.Write.Account {
		if psw.statedb.Exist(addr) {
			stateDB.CreateAccount(addr)
		}
	}

	stateDB.SetNonce(sender, psw.statedb.GetNonce(sender))

	for addr := range psw.sCtx.Write.Balance {
		stateDB.SetBalance(addr, psw.statedb.GetBalance(addr), tracing.BalanceChangeTransfer)
	}
	for addr := range psw.sCtx.Write.Code {
		stateDB.SetCode(addr, psw.statedb.GetCode(addr))
	}
	for addr, keys := range psw.sCtx.Write.State {
		for key := range keys {
			stateDB.SetState(addr, key, psw.statedb.GetState(addr, key))
		}
	}
	for _, addr := range psw.sCtx.addAddressToList {
		stateDB.AddAddressToAccessList(addr)
	}
	for _, sd := range psw.sCtx.addSlotToAddress {
		stateDB.AddSlotToAccessList(sd.addr, sd.slot)
	}
	for _, log := range psw.AllLogs() {
		stateDB.AddLog(log)
	}
	for hash, bs := range psw.statedb.Preimages() {
		stateDB.AddPreimage(hash, bs)
	}
}
