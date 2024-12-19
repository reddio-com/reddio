package pending_state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

type VisitTxnID map[int64]struct{}

type StateContext struct {
	Read  *VisitedAddress
	Write *VisitedAddress

	addSlotToAddress []slotToAddress
	addAddressToList []common.Address
	prepareParams    *prepareParams
}

func NewStateContext() *StateContext {
	sctx := &StateContext{
		Read:  NewVisitedAddress(),
		Write: NewVisitedAddress(),
	}
	return sctx
}

func (sctx *StateContext) IsConflict(tar *StateContext) bool {
	if len(sctx.GetWriteState()) < 1 && len(sctx.GetWriteAddress()) < 1 &&
		len(sctx.GetReadState()) < 1 && len(sctx.GetReadAddress()) < 1 {
		return false
	}
	// write/write conflict of Address
	if IsAddressConflict(sctx.GetWriteAddress(), tar.GetWriteAddress()) {
		return true
	}
	// read/write conflict of Address
	if IsAddressConflict(sctx.GetReadAddress(), tar.GetWriteAddress()) {
		return true
	}
	// read/write conflict of Address
	if IsAddressConflict(sctx.GetWriteAddress(), tar.GetReadAddress()) {
		return true
	}

	// write/write conflict of (Address, StateKey)
	if IsStateConflict(sctx.GetWriteState(), tar.GetWriteState()) {
		return true
	}
	// read/write conflict of (Address, StateKey)
	if IsStateConflict(sctx.GetReadState(), tar.GetWriteState()) {
		return true
	}
	// read/write conflict of (Address, StateKey)
	if IsStateConflict(sctx.GetWriteState(), tar.GetReadState()) {
		return true
	}
	return false
}

func (sctx *StateContext) GetReadState() map[common.Address]map[common.Hash]VisitTxnID {
	return sctx.Read.State
}

func (sctx *StateContext) GetWriteState() map[common.Address]map[common.Hash]VisitTxnID {
	return sctx.Write.State
}

func (sctx *StateContext) ReadAddress() map[common.Address]VisitTxnID {
	return sctx.Read.Address
}

func (sctx *StateContext) GetWriteAddress() map[common.Address]VisitTxnID {
	return sctx.Write.Address
}

func (sctx *StateContext) GetReadAddress() map[common.Address]VisitTxnID {
	return sctx.Read.Address
}

func (sctx *StateContext) AddSlot2Address(slot slotToAddress) {
	sctx.addSlotToAddress = append(sctx.addSlotToAddress, slot)
}

func (sctx *StateContext) AddAddressToList(address common.Address) {
	sctx.addAddressToList = append(sctx.addAddressToList, address)
}

func (sctx *StateContext) SetPrepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	sctx.prepareParams = &prepareParams{
		rules:       rules,
		sender:      sender,
		coinbase:    coinbase,
		dest:        dest,
		precompiles: precompiles,
		txAccesses:  txAccesses,
	}
}

func (sctx *StateContext) SelfDestruct(addr common.Address, txnID int64) {
	sctx.Write.VisitDestruct(addr, txnID)
}

func (sctx *StateContext) WriteAccount(addr common.Address, txnID int64) {
	sctx.Write.VisitAccount(addr, txnID)
}

func (sctx *StateContext) WriteBalance(addr common.Address, txnID int64) {
	sctx.Write.VisitBalance(addr, txnID)
}

func (sctx *StateContext) ReadBalance(addr common.Address, txnID int64) {
	sctx.Read.VisitBalance(addr, txnID)
}

func (sctx *StateContext) WriteCode(addr common.Address, txnID int64) {
	sctx.Write.VisitCode(addr, txnID)
}

func (sctx *StateContext) ReadCode(addr common.Address, txnID int64) {
	sctx.Read.VisitCode(addr, txnID)
}

func (sctx *StateContext) WriteState(addr common.Address, key common.Hash, txnID int64) {
	sctx.Write.VisitState(addr, key, txnID)
}

func (sctx *StateContext) ReadState(addr common.Address, key common.Hash, txnID int64) {
	sctx.Read.VisitState(addr, key, txnID)
}

type VisitedAddress struct {
	// address -> TxnID
	Address  map[common.Address]VisitTxnID
	Account  map[common.Address]VisitTxnID
	Destruct map[common.Address]VisitTxnID
	Balance  map[common.Address]VisitTxnID
	Code     map[common.Address]VisitTxnID
	State    map[common.Address]map[common.Hash]VisitTxnID
}

func NewVisitedAddress() *VisitedAddress {
	return &VisitedAddress{
		Address:  make(map[common.Address]VisitTxnID),
		Account:  make(map[common.Address]VisitTxnID),
		Destruct: make(map[common.Address]VisitTxnID),
		Balance:  make(map[common.Address]VisitTxnID),
		Code:     make(map[common.Address]VisitTxnID),
		State:    make(map[common.Address]map[common.Hash]VisitTxnID),
	}
}

func (v *VisitedAddress) VisitAccount(addr common.Address, txnID int64) {
	v.Address = txnVisitAddrMap(v.Address, addr, txnID)
	v.Account = txnVisitAddrMap(v.Account, addr, txnID)
}

func (v *VisitedAddress) VisitBalance(addr common.Address, txnID int64) {
	v.Address = txnVisitAddrMap(v.Address, addr, txnID)
	v.Balance = txnVisitAddrMap(v.Balance, addr, txnID)
}

func (v *VisitedAddress) VisitCode(addr common.Address, txnID int64) {
	v.Address = txnVisitAddrMap(v.Address, addr, txnID)
	v.Code = txnVisitAddrMap(v.Code, addr, txnID)
}

func (v *VisitedAddress) VisitState(addr common.Address, key common.Hash, txnID int64) {
	v1, ok := v.State[addr]
	if !ok {
		v1 = make(map[common.Hash]VisitTxnID)
	}
	v1 = txnVisitHashMap(v1, key, txnID)
	v.State[addr] = v1
}

func (v *VisitedAddress) VisitDestruct(addr common.Address, txnID int64) {
	v.Address = txnVisitAddrMap(v.Address, addr, txnID)
	v.Destruct = txnVisitAddrMap(v.Destruct, addr, txnID)
}

type prepareParams struct {
	rules       params.Rules
	sender      common.Address
	coinbase    common.Address
	dest        *common.Address
	precompiles []common.Address
	txAccesses  types.AccessList
}

type slotToAddress struct {
	addr common.Address
	slot common.Hash
}

func IsAddressConflict(a1, a2 map[common.Address]VisitTxnID) bool {
	for k := range a1 {
		_, ok := a2[k]
		if ok {
			return true
		}
	}
	for k := range a2 {
		_, ok := a1[k]
		if ok {
			return true
		}
	}
	return false
}

func IsStateConflict(a1, a2 map[common.Address]map[common.Hash]VisitTxnID) bool {
	for addr, v1 := range a1 {
		v2, ok := a2[addr]
		if !ok {
			continue
		}
		for k := range v1 {
			_, ok := v2[k]
			if ok {
				return true
			}
		}
	}
	for addr, v2 := range a2 {
		v1, ok := a1[addr]
		if !ok {
			continue
		}
		for k := range v2 {
			_, ok := v1[k]
			if ok {
				return true
			}
		}
	}
	return false
}

func txnVisitAddrMap(m map[common.Address]VisitTxnID, addr common.Address, txnID int64) map[common.Address]VisitTxnID {
	v, ok := m[addr]
	if ok {
		v[txnID] = struct{}{}
		return m
	}
	m[addr] = make(VisitTxnID)
	m[addr][txnID] = struct{}{}
	return m
}

func txnVisitHashMap(m map[common.Hash]VisitTxnID, hash common.Hash, txnID int64) map[common.Hash]VisitTxnID {
	v, ok := m[hash]
	if ok {
		v[txnID] = struct{}{}
		return m
	}
	m[hash] = make(VisitTxnID)
	m[hash][txnID] = struct{}{}
	return m
}
