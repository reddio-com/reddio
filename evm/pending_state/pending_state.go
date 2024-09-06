package pending_state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// PendingState provides a pending state for a transaction.
type PendingState struct {
	sCtx   *StateContext
	state  *state.StateDB
	sender common.Address
	logs   []*types.Log
}

func NewPendingState(sender common.Address, db *state.StateDB) *PendingState {
	return &PendingState{
		sCtx:   NewStateContext(),
		state:  db,
		sender: sender,
		logs:   make([]*types.Log, 0),
	}
}

func (s *PendingState) SetTxContext(txHash common.Hash, txIndex int) {
	s.state.SetTxContext(txHash, txIndex)
}

func (s *PendingState) GetStateDB() *state.StateDB {
	return s.state
}

func (s *PendingState) GetCtx() *StateContext {
	return s.sCtx
}

func (s *PendingState) CreateAccount(address common.Address) {
	s.sCtx.WriteAccount(address)
	s.state.CreateAccount(address)
}

func (s *PendingState) SubBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	s.sCtx.WriteBalance(address)
	s.state.SubBalance(address, u, reason)
}

func (s *PendingState) AddBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	s.sCtx.WriteBalance(address)
	s.state.AddBalance(address, u, reason)
}

func (s *PendingState) GetBalance(address common.Address) *uint256.Int {
	s.sCtx.ReadBalance(address)
	return s.state.GetBalance(address)
}

func (s *PendingState) GetNonce(address common.Address) uint64 {
	return s.state.GetNonce(address)
}

func (s *PendingState) SetNonce(address common.Address, u uint64) {
	s.state.SetNonce(address, u)
}

func (s *PendingState) GetCodeHash(address common.Address) common.Hash {
	s.sCtx.ReadCode(address)
	return s.state.GetCodeHash(address)
}

func (s *PendingState) GetCode(address common.Address) []byte {
	s.sCtx.ReadCode(address)
	return s.state.GetCode(address)
}

func (s *PendingState) SetCode(address common.Address, bytes []byte) {
	s.sCtx.WriteCode(address)
	s.state.SetCode(address, bytes)
}

func (s *PendingState) GetCodeSize(address common.Address) int {
	s.sCtx.ReadCode(address)
	return s.state.GetCodeSize(address)
}

func (s *PendingState) AddRefund(u uint64) {
	s.state.AddRefund(u)
}

func (s *PendingState) SubRefund(u uint64) {
	s.state.SubRefund(u)
}

func (s *PendingState) GetRefund() uint64 {
	return s.state.GetRefund()
}

func (s *PendingState) GetCommittedState(address common.Address, hash common.Hash) common.Hash {
	return s.state.GetCommittedState(address, hash)
}

func (s *PendingState) GetState(address common.Address, key common.Hash) common.Hash {
	s.sCtx.ReadState(address, key)
	return s.state.GetState(address, key)
}

func (s *PendingState) SetState(address common.Address, key common.Hash, value common.Hash) {
	s.sCtx.WriteState(address, key)
	s.state.SetState(address, key, value)
}

func (s *PendingState) GetStorageRoot(addr common.Address) common.Hash {
	return s.state.GetStorageRoot(addr)
}

func (s *PendingState) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	return s.state.GetTransientState(addr, key)
}

func (s *PendingState) SetTransientState(addr common.Address, key, value common.Hash) {
	s.state.SetTransientState(addr, key, value)
}

func (s *PendingState) SelfDestruct(address common.Address) {
	s.sCtx.SelfDestruct(address)
	s.state.SelfDestruct(address)
}

func (s *PendingState) HasSelfDestructed(address common.Address) bool {
	return s.state.HasSelfDestructed(address)
}

func (s *PendingState) Selfdestruct6780(address common.Address) {
	s.sCtx.SelfDestruct(address)
	s.state.Selfdestruct6780(address)
}

func (s *PendingState) Exist(address common.Address) bool {

	return s.state.Exist(address)
}

func (s *PendingState) Empty(address common.Address) bool {

	return s.state.Empty(address)
}

func (s *PendingState) AddressInAccessList(addr common.Address) bool {
	return s.state.AddressInAccessList(addr)
}

func (s *PendingState) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	return s.state.SlotInAccessList(addr, slot)
}

func (s *PendingState) AddAddressToAccessList(addr common.Address) {
	s.sCtx.AddAddressToList(addr)
	s.state.AddAddressToAccessList(addr)
}

func (s *PendingState) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	s.sCtx.AddSlot2Address(slotToAddress{
		addr: addr,
		slot: slot,
	})
	s.state.AddSlotToAccessList(addr, slot)
}

func (s *PendingState) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	s.sCtx.SetPrepare(rules, sender, coinbase, dest, precompiles, txAccesses)
	s.state.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
}

func (s *PendingState) RevertToSnapshot(i int) {
	s.state.RevertToSnapshot(i)
}

func (s *PendingState) Snapshot() int {
	return s.state.Snapshot()
}

func (s *PendingState) AddLog(log *types.Log) {
	s.logs = append(s.logs, log)
	s.state.AddLog(log)
}

func (s *PendingState) AddPreimage(hash common.Hash, bytes []byte) {
	s.state.AddPreimage(hash, bytes)
}

func (s *PendingState) AllLogs() []*types.Log {
	return s.logs
}

func (s *PendingState) AllPreimages() map[common.Hash][]byte {
	return s.state.Preimages()
}

func (s *PendingState) MergeInto(stateDB *state.StateDB) {
	if s.sCtx.prepareParams != nil {
		pre := s.sCtx.prepareParams
		stateDB.Prepare(pre.rules, pre.sender, pre.coinbase, pre.dest, pre.precompiles, pre.txAccesses)
	}
	for addr := range s.sCtx.Write.Account {
		if s.state.Exist(addr) {
			stateDB.CreateAccount(addr)
		}
	}

	stateDB.SetNonce(s.sender, s.state.GetNonce(s.sender))

	for addr := range s.sCtx.Write.Balance {
		stateDB.SetBalance(addr, s.state.GetBalance(addr), tracing.BalanceChangeTransfer)
	}
	for addr := range s.sCtx.Write.Code {
		stateDB.SetCode(addr, s.state.GetCode(addr))
	}
	for addr, keys := range s.sCtx.Write.State {
		for key := range keys {
			stateDB.SetState(addr, key, s.state.GetState(addr, key))
		}
	}
	for _, addr := range s.sCtx.addAddressToList {
		stateDB.AddAddressToAccessList(addr)
	}
	for _, sd := range s.sCtx.addSlotToAddress {
		stateDB.AddSlotToAccessList(sd.addr, sd.slot)
	}
	for _, log := range s.AllLogs() {
		stateDB.AddLog(log)
	}
	for hash, bs := range s.AllPreimages() {
		stateDB.AddPreimage(hash, bs)
	}
}
