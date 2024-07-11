package evm

type StateOpcode int

const (
	CreateAccount StateOpcode = iota

	SubBalance
	AddBalance
	GetBalance

	GetNonce
	SetNonce

	GetCodeHash
	GetCode
	SetCode
	GetCodeSize

	AddRefund
	SubRefund
	GetRefund

	GetCommittedState
	GetState
	SetState
	GetStorageRoot

	// local variables
	GetTransientState
	SetTransientState

	SelfDestruct
	HasSelfDestructed
	Selfdestruct6780

	Exist
	Empty

	// local variables
	AddressInAccessList
	SlotInAccessList
	AddAddressToAccessList
	AddSlotToAccessList
	Prepare

	// Apply: snapshot + snapshot
	RevertToSnapshot
	Snapshot

	// merge
	AddLog
	AddPreimage
)

var (
	Reads  = append(AccountReads, StateReads...)
	Writes = append(AccountWrites, StateWrites...)

	AccountReads  = append(EOAReads, ContractReads...)
	AccountWrites = append(EOAWrites, ContractWrites...)

	EOAReads = []StateOpcode{
		GetBalance,
	}
	EOAWrites = []StateOpcode{
		CreateAccount,
		SubBalance, AddBalance,
	}

	ContractReads = []StateOpcode{
		GetNonce,
		GetCodeHash, GetCode, GetCodeSize,
		GetStorageRoot,
		HasSelfDestructed,
		Exist, Empty,
	}
	ContractWrites = []StateOpcode{
		SetNonce,
		SetCode,
		SelfDestruct, Selfdestruct6780,
	}

	StateReads = []StateOpcode{
		GetCommittedState, GetState,
	}
	StateWrites = []StateOpcode{SetState}

	//TransientReads  = []StateOpcode{GetTransientState}
	//TransientWrites = []StateOpcode{SetTransientState}
	//
	//AccessListReads  = []StateOpcode{AddressInAccessList, SlotInAccessList}
	//AccessListWrites = []StateOpcode{AddAddressToAccessList, AddSlotToAccessList}
	//
	//PreimageWrites = []StateOpcode{AddPreimage}
)
