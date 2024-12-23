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

// PendingState provides a pending state for a transaction.
type PendingState struct {
	sync.RWMutex
	state *state.StateDB
}

func NewPendingState(db *state.StateDB) *PendingState {
	return &PendingState{
		state: db,
	}
}

func (s *PendingState) SetTxContext(txHash common.Hash, txIndex int) {
	s.Lock()
	defer s.Unlock()
	s.state.SetTxContext(txHash, txIndex)
}

func (s *PendingState) GetStateDB() *state.StateDB {
	s.RLock()
	defer s.RUnlock()
	return s.state
}

func (s *PendingState) CreateAccount(address common.Address) {
	s.Lock()
	defer s.Unlock()
	s.state.CreateAccount(address)
}

func (s *PendingState) SubBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	s.Lock()
	defer s.Unlock()
	s.state.SubBalance(address, u, reason)
}

func (s *PendingState) AddBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	s.Lock()
	defer s.Unlock()
	s.state.AddBalance(address, u, reason)
}

func (s *PendingState) GetBalance(address common.Address) *uint256.Int {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetBalance(address)
}

func (s *PendingState) GetNonce(address common.Address) uint64 {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetNonce(address)
}

func (s *PendingState) SetNonce(address common.Address, u uint64) {
	s.Lock()
	defer s.Unlock()
	s.state.SetNonce(address, u)
}

func (s *PendingState) GetCodeHash(address common.Address) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetCodeHash(address)
}

func (s *PendingState) GetCode(address common.Address) []byte {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetCode(address)
}

func (s *PendingState) SetCode(address common.Address, bytes []byte) {
	s.Lock()
	defer s.Unlock()
	s.state.SetCode(address, bytes)
}

func (s *PendingState) GetCodeSize(address common.Address) int {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetCodeSize(address)
}

func (s *PendingState) AddRefund(u uint64) {
	s.Lock()
	defer s.Unlock()
	s.state.AddRefund(u)
}

func (s *PendingState) SubRefund(u uint64) {
	s.Lock()
	defer s.Unlock()
	s.state.SubRefund(u)
}

func (s *PendingState) GetRefund() uint64 {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetRefund()
}

func (s *PendingState) GetCommittedState(address common.Address, hash common.Hash) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetCommittedState(address, hash)
}

func (s *PendingState) GetState(address common.Address, key common.Hash) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetState(address, key)
}

func (s *PendingState) SetState(address common.Address, key common.Hash, value common.Hash) {
	s.Lock()
	defer s.Unlock()
	s.state.SetState(address, key, value)
}

func (s *PendingState) GetStorageRoot(addr common.Address) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetStorageRoot(addr)
}

func (s *PendingState) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetTransientState(addr, key)
}

func (s *PendingState) SetTransientState(addr common.Address, key, value common.Hash) {
	s.Lock()
	defer s.Unlock()
	s.state.SetTransientState(addr, key, value)
}

func (s *PendingState) SelfDestruct(address common.Address) {
	s.Lock()
	defer s.Unlock()
	s.state.SelfDestruct(address)
}

func (s *PendingState) HasSelfDestructed(address common.Address) bool {
	s.RLock()
	defer s.RUnlock()
	return s.state.HasSelfDestructed(address)
}

func (s *PendingState) Selfdestruct6780(address common.Address) {
	s.Lock()
	defer s.Unlock()
	s.state.Selfdestruct6780(address)
}

func (s *PendingState) Exist(address common.Address) bool {
	s.RLock()
	defer s.RUnlock()
	return s.state.Exist(address)
}

func (s *PendingState) Empty(address common.Address) bool {
	s.RLock()
	defer s.RUnlock()
	return s.state.Empty(address)
}

func (s *PendingState) AddressInAccessList(addr common.Address) bool {
	s.RLock()
	defer s.RUnlock()
	return s.state.AddressInAccessList(addr)
}

func (s *PendingState) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	s.RLock()
	defer s.RUnlock()
	return s.state.SlotInAccessList(addr, slot)
}

func (s *PendingState) AddAddressToAccessList(addr common.Address) {
	s.Lock()
	defer s.Unlock()
	s.state.AddAddressToAccessList(addr)
}

func (s *PendingState) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	s.Lock()
	defer s.Unlock()
	s.state.AddSlotToAccessList(addr, slot)
}

func (s *PendingState) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	s.Lock()
	defer s.Unlock()
	s.state.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
}

func (s *PendingState) RevertToSnapshot(i int) {
	s.Lock()
	defer s.Unlock()
	s.state.RevertToSnapshot(i)
}

func (s *PendingState) Snapshot() int {
	s.Lock()
	defer s.Unlock()
	return s.state.Snapshot()
}

func (s *PendingState) AddLog(log *types.Log) {
	s.Lock()
	defer s.Unlock()
	s.state.AddLog(log)
}

func (s *PendingState) AddPreimage(hash common.Hash, bytes []byte) {
	s.Lock()
	defer s.Unlock()
	s.state.AddPreimage(hash, bytes)
}

func (s *PendingState) AllPreimages() map[common.Hash][]byte {
	s.RLock()
	defer s.RUnlock()
	return s.state.Preimages()
}
