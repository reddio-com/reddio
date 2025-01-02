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

// StateDBWrapper provides a pending state for a transaction.
type StateDBWrapper struct {
	sync.RWMutex
	state *state.StateDB
}

func NewStateDBWrapper(db *state.StateDB) *StateDBWrapper {
	return &StateDBWrapper{
		state: db,
	}
}

func (s *StateDBWrapper) SetTxContext(txHash common.Hash, txIndex int) {
	s.Lock()
	defer s.Unlock()
	s.state.SetTxContext(txHash, txIndex)
}

func (s *StateDBWrapper) GetStateDB() *state.StateDB {
	s.RLock()
	defer s.RUnlock()
	return s.state
}

func (s *StateDBWrapper) CreateAccount(address common.Address) {
	s.Lock()
	defer s.Unlock()
	s.state.CreateAccount(address)
}

func (s *StateDBWrapper) SubBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	s.Lock()
	defer s.Unlock()
	s.state.SubBalance(address, u, reason)
}

func (s *StateDBWrapper) AddBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	s.Lock()
	defer s.Unlock()
	s.state.AddBalance(address, u, reason)
}

func (s *StateDBWrapper) GetBalance(address common.Address) *uint256.Int {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetBalance(address)
}

func (s *StateDBWrapper) GetNonce(address common.Address) uint64 {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetNonce(address)
}

func (s *StateDBWrapper) SetNonce(address common.Address, u uint64) {
	s.Lock()
	defer s.Unlock()
	s.state.SetNonce(address, u)
}

func (s *StateDBWrapper) GetCodeHash(address common.Address) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetCodeHash(address)
}

func (s *StateDBWrapper) GetCode(address common.Address) []byte {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetCode(address)
}

func (s *StateDBWrapper) SetCode(address common.Address, bytes []byte) {
	s.Lock()
	defer s.Unlock()
	s.state.SetCode(address, bytes)
}

func (s *StateDBWrapper) GetCodeSize(address common.Address) int {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetCodeSize(address)
}

func (s *StateDBWrapper) AddRefund(u uint64) {
	s.Lock()
	defer s.Unlock()
	s.state.AddRefund(u)
}

func (s *StateDBWrapper) SubRefund(u uint64) {
	s.Lock()
	defer s.Unlock()
	s.state.SubRefund(u)
}

func (s *StateDBWrapper) GetRefund() uint64 {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetRefund()
}

func (s *StateDBWrapper) GetCommittedState(address common.Address, hash common.Hash) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetCommittedState(address, hash)
}

func (s *StateDBWrapper) GetState(address common.Address, key common.Hash) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetState(address, key)
}

func (s *StateDBWrapper) SetState(address common.Address, key common.Hash, value common.Hash) {
	s.Lock()
	defer s.Unlock()
	s.state.SetState(address, key, value)
}

func (s *StateDBWrapper) GetStorageRoot(addr common.Address) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetStorageRoot(addr)
}

func (s *StateDBWrapper) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	s.RLock()
	defer s.RUnlock()
	return s.state.GetTransientState(addr, key)
}

func (s *StateDBWrapper) SetTransientState(addr common.Address, key, value common.Hash) {
	s.Lock()
	defer s.Unlock()
	s.state.SetTransientState(addr, key, value)
}

func (s *StateDBWrapper) SelfDestruct(address common.Address) {
	s.Lock()
	defer s.Unlock()
	s.state.SelfDestruct(address)
}

func (s *StateDBWrapper) HasSelfDestructed(address common.Address) bool {
	s.RLock()
	defer s.RUnlock()
	return s.state.HasSelfDestructed(address)
}

func (s *StateDBWrapper) Selfdestruct6780(address common.Address) {
	s.Lock()
	defer s.Unlock()
	s.state.Selfdestruct6780(address)
}

func (s *StateDBWrapper) Exist(address common.Address) bool {
	s.RLock()
	defer s.RUnlock()
	return s.state.Exist(address)
}

func (s *StateDBWrapper) Empty(address common.Address) bool {
	s.RLock()
	defer s.RUnlock()
	return s.state.Empty(address)
}

func (s *StateDBWrapper) AddressInAccessList(addr common.Address) bool {
	s.RLock()
	defer s.RUnlock()
	return s.state.AddressInAccessList(addr)
}

func (s *StateDBWrapper) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	s.RLock()
	defer s.RUnlock()
	return s.state.SlotInAccessList(addr, slot)
}

func (s *StateDBWrapper) AddAddressToAccessList(addr common.Address) {
	s.Lock()
	defer s.Unlock()
	s.state.AddAddressToAccessList(addr)
}

func (s *StateDBWrapper) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	s.Lock()
	defer s.Unlock()
	s.state.AddSlotToAccessList(addr, slot)
}

func (s *StateDBWrapper) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	s.Lock()
	defer s.Unlock()
	s.state.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
}

func (s *StateDBWrapper) RevertToSnapshot(i int) {
	s.Lock()
	defer s.Unlock()
	s.state.RevertToSnapshot(i)
}

func (s *StateDBWrapper) Snapshot() int {
	s.Lock()
	defer s.Unlock()
	return s.state.Snapshot()
}

func (s *StateDBWrapper) AddLog(log *types.Log) {
	s.Lock()
	defer s.Unlock()
	s.state.AddLog(log)
}

func (s *StateDBWrapper) AddPreimage(hash common.Hash, bytes []byte) {
	s.Lock()
	defer s.Unlock()
	s.state.AddPreimage(hash, bytes)
}

func (s *StateDBWrapper) AllPreimages() map[common.Hash][]byte {
	s.RLock()
	defer s.RUnlock()
	return s.state.Preimages()
}
