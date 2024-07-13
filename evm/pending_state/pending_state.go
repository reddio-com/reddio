package pending_state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/holiman/uint256"
)

type PendingState struct {
	sCtx  *StateContext
	state *state.StateDB
}

func NewPendingState(db *state.StateDB) *PendingState {
	return &PendingState{
		sCtx:  NewStateContext(),
		state: db,
	}
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

func (s *PendingState) GetState(address common.Address, hash common.Hash) common.Hash {
	s.sCtx.ReadState(address)
	return s.state.GetState(address, hash)
}

func (s *PendingState) SetState(address common.Address, hash common.Hash, hash2 common.Hash) {
	s.sCtx.WriteState(address)
	s.state.SetState(address, hash, hash2)
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
	s.sCtx.HasRevertToSnapshot = true
	s.state.RevertToSnapshot(i)
}

func (s *PendingState) Snapshot() int {
	return s.state.Snapshot()
}

func (s *PendingState) AddLog(log *types.Log) {
	s.state.AddLog(log)
}

func (s *PendingState) AddPreimage(hash common.Hash, bytes []byte) {
	s.state.AddPreimage(hash, bytes)
}

func (s *PendingState) StopPrefetcher() {
	s.state.StopPrefetcher()
}

func (s *PendingState) Commit(block uint64, deleteEmptyObjects bool) (common.Hash, error) {
	return s.state.Commit(block, deleteEmptyObjects)
}

func (s *PendingState) SetBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) {
	s.state.SetBalance(addr, amount, reason)
}

func (s *PendingState) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	s.state.SetStorage(addr, storage)
}

func (s *PendingState) Error() error {
	return s.state.Error()
}

func (s *PendingState) TrieDB() *triedb.Database {
	return s.state.Database().TrieDB()
}

func (s *PendingState) Finalise(deleteEmptyObjects bool) {
	s.state.Finalise(deleteEmptyObjects)
}

func (s *PendingState) Internal() *state.StateDB {
	return s.state
}

func (s *PendingState) Copy() *PendingState {
	return NewPendingState(s.state.Copy())
}
