package pending_state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

type PendingStateWrapper struct {
	PS    *PendingState
	TxnID int64
	sCtx  *StateContext
	logs  []*types.Log
}

func NewPendingStateWrapper(ps *PendingState, TxnID int64) *PendingStateWrapper {
	psw := &PendingStateWrapper{
		PS:    ps,
		TxnID: TxnID,
		sCtx:  NewStateContext(),
		logs:  make([]*types.Log, 0),
	}
	return psw
}

func (psw *PendingStateWrapper) CreateAccount(address common.Address) {
	psw.sCtx.WriteAccount(address)
	psw.PS.CreateAccount(address)
}

func (psw *PendingStateWrapper) SubBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	psw.sCtx.WriteBalance(address)
	psw.PS.SubBalance(address, u, reason)
}

func (psw *PendingStateWrapper) AddBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	psw.sCtx.WriteBalance(address)
	psw.PS.AddBalance(address, u, reason)
}

func (psw *PendingStateWrapper) GetBalance(address common.Address) *uint256.Int {
	psw.sCtx.ReadBalance(address)
	return psw.PS.GetBalance(address)
}

func (psw *PendingStateWrapper) GetNonce(address common.Address) uint64 {
	return psw.PS.GetNonce(address)
}

func (psw *PendingStateWrapper) SetNonce(address common.Address, u uint64) {
	psw.PS.SetNonce(address, u)
}

func (psw *PendingStateWrapper) GetCodeHash(address common.Address) common.Hash {
	psw.sCtx.ReadCode(address)
	return psw.PS.GetCodeHash(address)
}

func (psw *PendingStateWrapper) GetCode(address common.Address) []byte {
	psw.sCtx.ReadCode(address)
	return psw.PS.GetCode(address)
}

func (psw *PendingStateWrapper) SetCode(address common.Address, bytes []byte) {
	psw.sCtx.WriteCode(address)
	psw.PS.SetCode(address, bytes)
}

func (psw *PendingStateWrapper) GetCodeSize(address common.Address) int {
	psw.sCtx.ReadCode(address)
	return psw.PS.GetCodeSize(address)
}

func (psw *PendingStateWrapper) AddRefund(u uint64) {
	psw.PS.AddRefund(u)
}

func (psw *PendingStateWrapper) SubRefund(u uint64) {
	psw.PS.SubRefund(u)
}

func (psw *PendingStateWrapper) GetRefund() uint64 {
	return psw.PS.GetRefund()
}

func (psw *PendingStateWrapper) GetCommittedState(address common.Address, hash common.Hash) common.Hash {
	return psw.PS.GetCommittedState(address, hash)
}

func (psw *PendingStateWrapper) GetState(address common.Address, hash common.Hash) common.Hash {
	psw.sCtx.ReadState(address, hash)
	return psw.PS.GetState(address, hash)
}

func (psw *PendingStateWrapper) SetState(address common.Address, hash common.Hash, hash2 common.Hash) {
	psw.sCtx.WriteState(address, hash)
	psw.PS.SetState(address, hash, hash2)
}

func (psw *PendingStateWrapper) GetStorageRoot(addr common.Address) common.Hash {
	return psw.PS.GetStorageRoot(addr)
}

func (psw *PendingStateWrapper) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	return psw.PS.GetTransientState(addr, key)
}

func (psw *PendingStateWrapper) SetTransientState(addr common.Address, key, value common.Hash) {
	psw.PS.SetTransientState(addr, key, value)
}

func (psw *PendingStateWrapper) SelfDestruct(address common.Address) {
	psw.sCtx.SelfDestruct(address)
	psw.PS.SelfDestruct(address)
}

func (psw *PendingStateWrapper) HasSelfDestructed(address common.Address) bool {
	return psw.PS.HasSelfDestructed(address)
}

func (psw *PendingStateWrapper) Selfdestruct6780(address common.Address) {
	psw.sCtx.SelfDestruct(address)
	psw.PS.Selfdestruct6780(address)
}

func (psw *PendingStateWrapper) Exist(address common.Address) bool {
	return psw.PS.Exist(address)
}

func (psw *PendingStateWrapper) Empty(address common.Address) bool {
	return psw.PS.Empty(address)
}

func (psw *PendingStateWrapper) AddressInAccessList(addr common.Address) bool {
	return psw.PS.AddressInAccessList(addr)
}

func (psw *PendingStateWrapper) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	return psw.PS.SlotInAccessList(addr, slot)
}

func (psw *PendingStateWrapper) AddAddressToAccessList(addr common.Address) {
	psw.sCtx.AddAddressToList(addr)
	psw.PS.AddAddressToAccessList(addr)
}

func (psw *PendingStateWrapper) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	psw.sCtx.AddSlot2Address(slotToAddress{
		addr: addr,
		slot: slot,
	})
	psw.PS.AddSlotToAccessList(addr, slot)
}

func (psw *PendingStateWrapper) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	psw.sCtx.SetPrepare(rules, sender, coinbase, dest, precompiles, txAccesses)
	psw.PS.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
}

func (psw *PendingStateWrapper) RevertToSnapshot(i int) {
	psw.PS.RevertToSnapshot(i)
}

func (psw *PendingStateWrapper) Snapshot() int {
	return psw.PS.Snapshot()
}

func (psw *PendingStateWrapper) AddLog(log *types.Log) {
	psw.logs = append(psw.logs, log)
	psw.PS.AddLog(log)
}

func (psw *PendingStateWrapper) AddPreimage(hash common.Hash, bytes []byte) {
	psw.PS.AddPreimage(hash, bytes)
}

func (psw *PendingStateWrapper) SetTxContext(txHash common.Hash, txIndex int) {
	psw.PS.SetTxContext(txHash, txIndex)
}

func (psw *PendingStateWrapper) GetStateDB() *state.StateDB {
	return psw.PS.GetStateDB()
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
		if psw.PS.state.Exist(addr) {
			stateDB.CreateAccount(addr)
		}
	}

	stateDB.SetNonce(sender, psw.PS.state.GetNonce(sender))

	for addr := range psw.sCtx.Write.Balance {
		stateDB.SetBalance(addr, psw.PS.state.GetBalance(addr), tracing.BalanceChangeTransfer)
	}
	for addr := range psw.sCtx.Write.Code {
		stateDB.SetCode(addr, psw.PS.state.GetCode(addr))
	}
	for addr, keys := range psw.sCtx.Write.State {
		for key := range keys {
			stateDB.SetState(addr, key, psw.PS.state.GetState(addr, key))
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
	for hash, bs := range psw.PS.AllPreimages() {
		stateDB.AddPreimage(hash, bs)
	}
}
