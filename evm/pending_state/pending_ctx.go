package pending_state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

type StateContext struct {
	Read  *VisitedAddress
	Write *VisitedAddress

	addSlotToAddress    []slotToAddress
	addAddressToList    []common.Address
	prepareParams       *prepareParams
	HasRevertToSnapshot bool
}

func NewStateContext() *StateContext {
	sctx := &StateContext{
		Read:  NewVisitedAddress(),
		Write: NewVisitedAddress(),
	}
	return sctx
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

func (sctx *StateContext) WriteState(addr common.Address) {
	sctx.Write.VisitState(addr)
}

func (sctx *StateContext) ReadState(addr common.Address) {
	sctx.Read.VisitState(addr)
}

type VisitedAddress struct {
	Account  map[common.Address]struct{}
	Destruct map[common.Address]struct{}
	Balance  map[common.Address]struct{}
	Code     map[common.Address]struct{}
	State    map[common.Address]struct{}
}

func NewVisitedAddress() *VisitedAddress {
	return &VisitedAddress{
		Account:  make(map[common.Address]struct{}),
		Destruct: make(map[common.Address]struct{}),
		Balance:  make(map[common.Address]struct{}),
		Code:     make(map[common.Address]struct{}),
		State:    make(map[common.Address]struct{}),
	}
}

func (v *VisitedAddress) VisitAccount(addr common.Address) {
	v.Account[addr] = struct{}{}
}

func (v *VisitedAddress) VisitBalance(addr common.Address) {
	v.Balance[addr] = struct{}{}
}

func (v *VisitedAddress) VisitCode(addr common.Address) {
	v.Code[addr] = struct{}{}
}

func (v *VisitedAddress) VisitState(addr common.Address) {
	v.State[addr] = struct{}{}
}

func (v *VisitedAddress) VisitDestruct(addr common.Address) {
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

type preImage struct {
	hash common.Hash
	b    []byte
}

type slotToAddress struct {
	addr common.Address
	slot common.Hash
}
