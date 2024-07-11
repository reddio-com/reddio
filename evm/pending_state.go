package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

type PendingState struct {
	reader *state.StateDB

	AccountWriteSet map[common.Address]*AccountOp
	AccountReadSet  map[common.Address]StateOpcode

	StateWriteSet map[StateKey]StateOp
	StateReadSet  map[StateKey]StateOpcode
}

func NewPendingState(prevStateDB *state.StateDB) (*PendingState, error) {
	return &PendingState{
		reader: prevStateDB,

		AccountWriteSet: make(map[common.Address]*AccountOp),
		AccountReadSet:  make(map[common.Address]StateOpcode),

		StateWriteSet: make(map[StateKey]StateOp),
		StateReadSet:  make(map[StateKey]StateOpcode),
	}, nil
}

func (s *PendingState) CreateAccount(address common.Address) {
	s.writeAccount(address, &AccountOp{
		Opcode: CreateAccount,
	})
}

func (s *PendingState) SubBalance(address common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) {
	s.writeAccount(address, &AccountOp{
		Opcode: SubBalance,
		Amount: amount,
		Reason: reason,
	})
}

func (s *PendingState) AddBalance(address common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) {
	s.writeAccount(address, &AccountOp{
		Opcode: AddBalance,
		Amount: amount,
		Reason: reason,
	})
}

func (s *PendingState) GetBalance(address common.Address) *uint256.Int {
	op := s.readAccount(GetBalance, address)
	if op != nil {
		return op.Amount
	}
	return s.reader.GetBalance(address)
}

func (s *PendingState) GetNonce(address common.Address) uint64 {
	s.readAccount(GetNonce, address)
	return s.reader.GetNonce(address)
}

func (s *PendingState) SetNonce(address common.Address, nonce uint64) {
	s.writeAccount(address, &AccountOp{
		Opcode: SetNonce,
		Nonce:  nonce,
	})
}

func (s *PendingState) GetCodeHash(address common.Address) common.Hash {
	s.readAccount(GetCodeHash, address)
	return s.reader.GetCodeHash(address)
}

func (s *PendingState) GetCode(address common.Address) []byte {
	s.readAccount(GetCode, address)
	return s.reader.GetCode(address)
}

func (s *PendingState) SetCode(address common.Address, code []byte) {
	s.writeAccount(address, &AccountOp{
		Opcode: SetCode,
		Code:   code,
	})
}

func (s *PendingState) GetCodeSize(address common.Address) int {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) AddRefund(u uint64) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) SubRefund(u uint64) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) GetRefund() uint64 {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) GetCommittedState(address common.Address, hash common.Hash) common.Hash {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) GetState(address common.Address, hash common.Hash) common.Hash {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) SetState(address common.Address, hash common.Hash, hash2 common.Hash) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) GetStorageRoot(addr common.Address) common.Hash {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) SetTransientState(addr common.Address, key, value common.Hash) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) SelfDestruct(address common.Address) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) HasSelfDestructed(address common.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) Selfdestruct6780(address common.Address) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) Exist(address common.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) Empty(address common.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) AddressInAccessList(addr common.Address) bool {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) AddAddressToAccessList(addr common.Address) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	s.reader.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
}

func (s *PendingState) RevertToSnapshot(i int) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) Snapshot() int {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) AddLog(log *types.Log) {
	//TODO implement me
	panic("implement me")
}

func (s *PendingState) AddPreimage(hash common.Hash, bytes []byte) {
	//TODO implement me
	panic("implement me")
}

type AccountOp struct {
	Opcode StateOpcode

	Amount *uint256.Int
	Reason tracing.BalanceChangeReason

	Nonce uint64

	Code []byte
}

func (s *PendingState) writeAccount(addr common.Address, op *AccountOp) {
	s.AccountWriteSet[addr] = op
	delete(s.AccountReadSet, addr)
}

func (s *PendingState) readAccount(opcode StateOpcode, addr common.Address) *AccountOp {
	s.AccountReadSet[addr] = opcode
	return s.AccountWriteSet[addr]
}

type StateKey struct {
	Addr common.Address
	Key  common.Hash
}

type StateOp struct {
	Opcode StateOpcode
	Key    common.Hash
	Value  common.Hash
}

func (s *PendingState) writeState(key StateKey, op StateOp) {
	s.StateWriteSet[key] = op
	delete(s.StateReadSet, key)
}

func (s *PendingState) readState(opcode StateOpcode, key StateKey) common.Hash {
	s.StateReadSet[key] = opcode
	if op, ok := s.StateWriteSet[key]; ok {
		return op.Value
	}
	return common.Hash{}
}

// stateDB.APPly (pending)
// pending.Wrset
