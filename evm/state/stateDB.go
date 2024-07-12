package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/holiman/uint256"
)

type StateDBWrapper struct {
	sCtx  *StateContext
	state *state.StateDB
}

func NewStateDB(db *state.StateDB) *StateDBWrapper {
	return &StateDBWrapper{
		sCtx:  NewStateContext(),
		state: db,
	}
}

func (s *StateDBWrapper) CreateAccount(address common.Address) {
	s.state.CreateAccount(address)
}

func (s *StateDBWrapper) SubBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	s.sCtx.WriteBalance(address)
	s.state.SubBalance(address, u, reason)
}

func (s *StateDBWrapper) AddBalance(address common.Address, u *uint256.Int, reason tracing.BalanceChangeReason) {
	s.sCtx.WriteBalance(address)
	s.state.AddBalance(address, u, reason)
}

func (s *StateDBWrapper) GetBalance(address common.Address) *uint256.Int {
	s.sCtx.ReadBalance(address)
	return s.state.GetBalance(address)
}

func (s *StateDBWrapper) GetNonce(address common.Address) uint64 {
	return s.state.GetNonce(address)
}

func (s *StateDBWrapper) SetNonce(address common.Address, u uint64) {
	s.state.SetNonce(address, u)
}

func (s *StateDBWrapper) GetCodeHash(address common.Address) common.Hash {
	s.sCtx.ReadCode(address)
	return s.state.GetCodeHash(address)
}

func (s *StateDBWrapper) GetCode(address common.Address) []byte {
	s.sCtx.ReadCode(address)
	return s.state.GetCode(address)
}

func (s *StateDBWrapper) SetCode(address common.Address, bytes []byte) {
	s.sCtx.WriteCode(address)
	s.state.SetCode(address, bytes)
}

func (s *StateDBWrapper) GetCodeSize(address common.Address) int {
	s.sCtx.ReadCode(address)
	return s.GetCodeSize(address)
}

func (s *StateDBWrapper) AddRefund(u uint64) {
	s.state.AddRefund(u)
}

func (s *StateDBWrapper) SubRefund(u uint64) {
	s.state.SubRefund(u)
}

func (s *StateDBWrapper) GetRefund() uint64 {
	return s.state.GetRefund()
}

func (s *StateDBWrapper) GetCommittedState(address common.Address, hash common.Hash) common.Hash {
	return s.state.GetCommittedState(address, hash)
}

func (s *StateDBWrapper) GetState(address common.Address, hash common.Hash) common.Hash {
	s.sCtx.ReadState(address)
	return s.state.GetState(address, hash)
}

func (s *StateDBWrapper) SetState(address common.Address, hash common.Hash, hash2 common.Hash) {

	s.sCtx.WriteState(address)
	s.state.SetState(address, hash, hash2)
}

func (s *StateDBWrapper) GetStorageRoot(addr common.Address) common.Hash {
	return s.state.GetStorageRoot(addr)
}

func (s *StateDBWrapper) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	return s.state.GetTransientState(addr, key)
}

func (s *StateDBWrapper) SetTransientState(addr common.Address, key, value common.Hash) {
	s.state.SetTransientState(addr, key, value)
}

func (s *StateDBWrapper) SelfDestruct(address common.Address) {
	s.SelfDestruct(address)
	s.state.SelfDestruct(address)
}

func (s *StateDBWrapper) HasSelfDestructed(address common.Address) bool {

	return s.state.HasSelfDestructed(address)
}

func (s *StateDBWrapper) Selfdestruct6780(address common.Address) {
	s.SelfDestruct(address)
	s.state.Selfdestruct6780(address)
}

func (s *StateDBWrapper) Exist(address common.Address) bool {

	return s.state.Exist(address)
}

func (s *StateDBWrapper) Empty(address common.Address) bool {

	return s.state.Empty(address)
}

func (s *StateDBWrapper) AddressInAccessList(addr common.Address) bool {
	return s.state.AddressInAccessList(addr)
}

func (s *StateDBWrapper) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	return s.state.SlotInAccessList(addr, slot)
}

func (s *StateDBWrapper) AddAddressToAccessList(addr common.Address) {
	s.sCtx.AddAddressToList(addr)
	s.state.AddAddressToAccessList(addr)
}

func (s *StateDBWrapper) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	s.sCtx.AddSlot2Address(slotToAddress{
		addr: addr,
		slot: slot,
	})
	s.state.AddSlotToAccessList(addr, slot)
}

func (s *StateDBWrapper) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	s.sCtx.SetPrepare(rules, sender, coinbase, dest, precompiles, txAccesses)
	s.state.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
}

func (s *StateDBWrapper) RevertToSnapshot(i int) {
	s.sCtx.HasRevertToSnapshot = true
	s.state.RevertToSnapshot(i)
}

func (s *StateDBWrapper) Snapshot() int {
	return s.state.Snapshot()
}

func (s *StateDBWrapper) AddLog(log *types.Log) {
	s.sCtx.AddLog(log)
	s.state.AddLog(log)
}

func (s *StateDBWrapper) AddPreimage(hash common.Hash, bytes []byte) {
	s.sCtx.AddPreImage(preImage{
		hash: hash,
		b:    bytes,
	})
	s.state.AddPreimage(hash, bytes)
}

func (s *StateDBWrapper) StopPrefetcher() {
	s.state.StopPrefetcher()
}

func (s *StateDBWrapper) Commit(block uint64, deleteEmptyObjects bool) (common.Hash, error) {
	return s.state.Commit(block, deleteEmptyObjects)
}

func (s *StateDBWrapper) SetBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) {
	s.state.SetBalance(addr, amount, reason)
}

func (s *StateDBWrapper) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	s.state.SetStorage(addr, storage)
}

func (s *StateDBWrapper) Error() error {
	return s.state.Error()
}

func (s *StateDBWrapper) TrieDB() *triedb.Database {
	return s.state.Database().TrieDB()
}

func (s *StateDBWrapper) Finalise(deleteEmptyObjects bool) {
	s.state.Finalise(deleteEmptyObjects)
}

func (s *StateDBWrapper) Internal() *state.StateDB {
	return s.state
}

func (s *StateDBWrapper) Copy() *StateDBWrapper {
	return NewStateDB(s.state.Copy())
}
