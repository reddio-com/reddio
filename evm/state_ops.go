package evm

type StateFunc int

const (
	CreateAccount = iota

	SubBalance
	AddBalance
	GetBalance

	GetNonce
	SetNonce

	GetCodeHash
	GetCode
	SetCode
	GetCodeSize

	GetCommittedState
	GetState
	SetState
	GetStorageRoot

	GetTransientState
	SetTransientState

	SelfDestruct
	HasSelfDestructed
	Selfdestruct6780

	Exist
	Empty

	AddressInAccessList
	SlotInAccessList
	AddAddressToAccessList
	AddSlotToAccessList
	Prepare

	RevertToSnapshot
	Snapshot
	AddLog
	AddPreimage
)

