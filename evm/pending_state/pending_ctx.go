package pending_state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

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

func (sctx *StateContext) GetReadState() map[common.Address]map[common.Hash]struct{} {
	return sctx.Read.State
}

func (sctx *StateContext) GetWriteState() map[common.Address]map[common.Hash]struct{} {
	return sctx.Write.State
}

func (sctx *StateContext) ReadAddress() map[common.Address]struct{} {
	return sctx.Read.Address
}

func (sctx *StateContext) GetWriteAddress() map[common.Address]struct{} {
	return sctx.Write.Address
}

func (sctx *StateContext) GetReadAddress() map[common.Address]struct{} {
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

func (sctx *StateContext) SelfDestruct(addr common.Address) {
	sctx.Write.VisitDestruct(addr)
}

func (sctx *StateContext) WriteAccount(addr common.Address) {
	sctx.Write.VisitAccount(addr)
}

func (sctx *StateContext) WriteBalance(addr common.Address) {
	sctx.Write.VisitBalance(addr)
}

func (sctx *StateContext) ReadBalance(addr common.Address) {
	sctx.Read.VisitBalance(addr)
}

func (sctx *StateContext) WriteCode(addr common.Address) {
	sctx.Write.VisitCode(addr)
}

func (sctx *StateContext) ReadCode(addr common.Address) {
	sctx.Read.VisitCode(addr)
}

func (sctx *StateContext) WriteState(addr common.Address, key common.Hash) {
	sctx.Write.VisitState(addr, key)
}

func (sctx *StateContext) ReadState(addr common.Address, key common.Hash) {
	sctx.Read.VisitState(addr, key)
}

type VisitedAddress struct {
	Address  map[common.Address]struct{}
	Account  map[common.Address]struct{}
	Destruct map[common.Address]struct{}
	Balance  map[common.Address]struct{}
	Code     map[common.Address]struct{}
	State    map[common.Address]map[common.Hash]struct{}
}

func NewVisitedAddress() *VisitedAddress {
	return &VisitedAddress{
		Address:  make(map[common.Address]struct{}),
		Account:  make(map[common.Address]struct{}),
		Destruct: make(map[common.Address]struct{}),
		Balance:  make(map[common.Address]struct{}),
		Code:     make(map[common.Address]struct{}),
		State:    make(map[common.Address]map[common.Hash]struct{}),
	}
}

func (v *VisitedAddress) VisitAccount(addr common.Address) {
	v.Address[addr] = struct{}{}
	v.Account[addr] = struct{}{}
}

func (v *VisitedAddress) VisitBalance(addr common.Address) {
	v.Address[addr] = struct{}{}
	v.Balance[addr] = struct{}{}
}

func (v *VisitedAddress) VisitCode(addr common.Address) {
	v.Address[addr] = struct{}{}
	v.Code[addr] = struct{}{}
}

func (v *VisitedAddress) VisitState(addr common.Address, key common.Hash) {
	v1, ok := v.State[addr]
	if !ok {
		v1 = make(map[common.Hash]struct{})
	}
	v1[key] = struct{}{}
	v.State[addr] = v1
}

func (v *VisitedAddress) VisitDestruct(addr common.Address) {
	v.Address[addr] = struct{}{}
	v.Destruct[addr] = struct{}{}
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

func IsAddressConflict(a1, a2 map[common.Address]struct{}) bool {
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

func IsStateConflict(a1, a2 map[common.Address]map[common.Hash]struct{}) bool {
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
