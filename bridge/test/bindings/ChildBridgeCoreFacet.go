// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ChildBridgeCoreFacetMetaData contains all meta data concerning the ChildBridgeCoreFacet contract.
var ChildBridgeCoreFacetMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"name\":\"AppendMessageEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"xDomainCalldataHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"SentMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"UpwardMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"erc1155Address\",\"type\":\"address\"}],\"name\":\"getBridgedERC1155TokenChild\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"erc20Address\",\"type\":\"address\"}],\"name\":\"getBridgedERC20TokenChild\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"erc721Address\",\"type\":\"address\"}],\"name\":\"getBridgedERC721TokenChild\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridgedERC1155TokenAddress\",\"type\":\"address\"}],\"name\":\"getERC1155TokenChild\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridgedERC20Address\",\"type\":\"address\"}],\"name\":\"getERC20TokenChild\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gettL1RedTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseBridge\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseStatusBridge\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"sendUpwardMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"l1RedTokenAddress\",\"type\":\"address\"}],\"name\":\"setRedTokenAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpauseBridge\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ChildBridgeCoreFacetABI is the input ABI used to generate the binding from.
// Deprecated: Use ChildBridgeCoreFacetMetaData.ABI instead.
var ChildBridgeCoreFacetABI = ChildBridgeCoreFacetMetaData.ABI

// ChildBridgeCoreFacet is an auto generated Go binding around an Ethereum contract.
type ChildBridgeCoreFacet struct {
	ChildBridgeCoreFacetCaller     // Read-only binding to the contract
	ChildBridgeCoreFacetTransactor // Write-only binding to the contract
	ChildBridgeCoreFacetFilterer   // Log filterer for contract events
}

// ChildBridgeCoreFacetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChildBridgeCoreFacetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChildBridgeCoreFacetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChildBridgeCoreFacetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChildBridgeCoreFacetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChildBridgeCoreFacetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChildBridgeCoreFacetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChildBridgeCoreFacetSession struct {
	Contract     *ChildBridgeCoreFacet // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ChildBridgeCoreFacetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChildBridgeCoreFacetCallerSession struct {
	Contract *ChildBridgeCoreFacetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ChildBridgeCoreFacetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChildBridgeCoreFacetTransactorSession struct {
	Contract     *ChildBridgeCoreFacetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ChildBridgeCoreFacetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChildBridgeCoreFacetRaw struct {
	Contract *ChildBridgeCoreFacet // Generic contract binding to access the raw methods on
}

// ChildBridgeCoreFacetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChildBridgeCoreFacetCallerRaw struct {
	Contract *ChildBridgeCoreFacetCaller // Generic read-only contract binding to access the raw methods on
}

// ChildBridgeCoreFacetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChildBridgeCoreFacetTransactorRaw struct {
	Contract *ChildBridgeCoreFacetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChildBridgeCoreFacet creates a new instance of ChildBridgeCoreFacet, bound to a specific deployed contract.
func NewChildBridgeCoreFacet(address common.Address, backend bind.ContractBackend) (*ChildBridgeCoreFacet, error) {
	contract, err := bindChildBridgeCoreFacet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChildBridgeCoreFacet{ChildBridgeCoreFacetCaller: ChildBridgeCoreFacetCaller{contract: contract}, ChildBridgeCoreFacetTransactor: ChildBridgeCoreFacetTransactor{contract: contract}, ChildBridgeCoreFacetFilterer: ChildBridgeCoreFacetFilterer{contract: contract}}, nil
}

// NewChildBridgeCoreFacetCaller creates a new read-only instance of ChildBridgeCoreFacet, bound to a specific deployed contract.
func NewChildBridgeCoreFacetCaller(address common.Address, caller bind.ContractCaller) (*ChildBridgeCoreFacetCaller, error) {
	contract, err := bindChildBridgeCoreFacet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChildBridgeCoreFacetCaller{contract: contract}, nil
}

// NewChildBridgeCoreFacetTransactor creates a new write-only instance of ChildBridgeCoreFacet, bound to a specific deployed contract.
func NewChildBridgeCoreFacetTransactor(address common.Address, transactor bind.ContractTransactor) (*ChildBridgeCoreFacetTransactor, error) {
	contract, err := bindChildBridgeCoreFacet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChildBridgeCoreFacetTransactor{contract: contract}, nil
}

// NewChildBridgeCoreFacetFilterer creates a new log filterer instance of ChildBridgeCoreFacet, bound to a specific deployed contract.
func NewChildBridgeCoreFacetFilterer(address common.Address, filterer bind.ContractFilterer) (*ChildBridgeCoreFacetFilterer, error) {
	contract, err := bindChildBridgeCoreFacet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChildBridgeCoreFacetFilterer{contract: contract}, nil
}

// bindChildBridgeCoreFacet binds a generic wrapper to an already deployed contract.
func bindChildBridgeCoreFacet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChildBridgeCoreFacetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChildBridgeCoreFacet.Contract.ChildBridgeCoreFacetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.ChildBridgeCoreFacetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.ChildBridgeCoreFacetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChildBridgeCoreFacet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.contract.Transact(opts, method, params...)
}

// GetBridgedERC1155TokenChild is a free data retrieval call binding the contract method 0x8544e8a2.
//
// Solidity: function getBridgedERC1155TokenChild(address erc1155Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCaller) GetBridgedERC1155TokenChild(opts *bind.CallOpts, erc1155Address common.Address) (common.Address, error) {
	var out []interface{}
	err := _ChildBridgeCoreFacet.contract.Call(opts, &out, "getBridgedERC1155TokenChild", erc1155Address)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBridgedERC1155TokenChild is a free data retrieval call binding the contract method 0x8544e8a2.
//
// Solidity: function getBridgedERC1155TokenChild(address erc1155Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) GetBridgedERC1155TokenChild(erc1155Address common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetBridgedERC1155TokenChild(&_ChildBridgeCoreFacet.CallOpts, erc1155Address)
}

// GetBridgedERC1155TokenChild is a free data retrieval call binding the contract method 0x8544e8a2.
//
// Solidity: function getBridgedERC1155TokenChild(address erc1155Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCallerSession) GetBridgedERC1155TokenChild(erc1155Address common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetBridgedERC1155TokenChild(&_ChildBridgeCoreFacet.CallOpts, erc1155Address)
}

// GetBridgedERC20TokenChild is a free data retrieval call binding the contract method 0x630ac212.
//
// Solidity: function getBridgedERC20TokenChild(address erc20Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCaller) GetBridgedERC20TokenChild(opts *bind.CallOpts, erc20Address common.Address) (common.Address, error) {
	var out []interface{}
	err := _ChildBridgeCoreFacet.contract.Call(opts, &out, "getBridgedERC20TokenChild", erc20Address)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBridgedERC20TokenChild is a free data retrieval call binding the contract method 0x630ac212.
//
// Solidity: function getBridgedERC20TokenChild(address erc20Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) GetBridgedERC20TokenChild(erc20Address common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetBridgedERC20TokenChild(&_ChildBridgeCoreFacet.CallOpts, erc20Address)
}

// GetBridgedERC20TokenChild is a free data retrieval call binding the contract method 0x630ac212.
//
// Solidity: function getBridgedERC20TokenChild(address erc20Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCallerSession) GetBridgedERC20TokenChild(erc20Address common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetBridgedERC20TokenChild(&_ChildBridgeCoreFacet.CallOpts, erc20Address)
}

// GetBridgedERC721TokenChild is a free data retrieval call binding the contract method 0x828a10b3.
//
// Solidity: function getBridgedERC721TokenChild(address erc721Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCaller) GetBridgedERC721TokenChild(opts *bind.CallOpts, erc721Address common.Address) (common.Address, error) {
	var out []interface{}
	err := _ChildBridgeCoreFacet.contract.Call(opts, &out, "getBridgedERC721TokenChild", erc721Address)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBridgedERC721TokenChild is a free data retrieval call binding the contract method 0x828a10b3.
//
// Solidity: function getBridgedERC721TokenChild(address erc721Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) GetBridgedERC721TokenChild(erc721Address common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetBridgedERC721TokenChild(&_ChildBridgeCoreFacet.CallOpts, erc721Address)
}

// GetBridgedERC721TokenChild is a free data retrieval call binding the contract method 0x828a10b3.
//
// Solidity: function getBridgedERC721TokenChild(address erc721Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCallerSession) GetBridgedERC721TokenChild(erc721Address common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetBridgedERC721TokenChild(&_ChildBridgeCoreFacet.CallOpts, erc721Address)
}

// GetERC1155TokenChild is a free data retrieval call binding the contract method 0x3fb645aa.
//
// Solidity: function getERC1155TokenChild(address bridgedERC1155TokenAddress) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCaller) GetERC1155TokenChild(opts *bind.CallOpts, bridgedERC1155TokenAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _ChildBridgeCoreFacet.contract.Call(opts, &out, "getERC1155TokenChild", bridgedERC1155TokenAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetERC1155TokenChild is a free data retrieval call binding the contract method 0x3fb645aa.
//
// Solidity: function getERC1155TokenChild(address bridgedERC1155TokenAddress) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) GetERC1155TokenChild(bridgedERC1155TokenAddress common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetERC1155TokenChild(&_ChildBridgeCoreFacet.CallOpts, bridgedERC1155TokenAddress)
}

// GetERC1155TokenChild is a free data retrieval call binding the contract method 0x3fb645aa.
//
// Solidity: function getERC1155TokenChild(address bridgedERC1155TokenAddress) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCallerSession) GetERC1155TokenChild(bridgedERC1155TokenAddress common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetERC1155TokenChild(&_ChildBridgeCoreFacet.CallOpts, bridgedERC1155TokenAddress)
}

// GetERC20TokenChild is a free data retrieval call binding the contract method 0xa31ca56d.
//
// Solidity: function getERC20TokenChild(address bridgedERC20Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCaller) GetERC20TokenChild(opts *bind.CallOpts, bridgedERC20Address common.Address) (common.Address, error) {
	var out []interface{}
	err := _ChildBridgeCoreFacet.contract.Call(opts, &out, "getERC20TokenChild", bridgedERC20Address)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetERC20TokenChild is a free data retrieval call binding the contract method 0xa31ca56d.
//
// Solidity: function getERC20TokenChild(address bridgedERC20Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) GetERC20TokenChild(bridgedERC20Address common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetERC20TokenChild(&_ChildBridgeCoreFacet.CallOpts, bridgedERC20Address)
}

// GetERC20TokenChild is a free data retrieval call binding the contract method 0xa31ca56d.
//
// Solidity: function getERC20TokenChild(address bridgedERC20Address) view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCallerSession) GetERC20TokenChild(bridgedERC20Address common.Address) (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GetERC20TokenChild(&_ChildBridgeCoreFacet.CallOpts, bridgedERC20Address)
}

// GettL1RedTokenAddress is a free data retrieval call binding the contract method 0xdada7d76.
//
// Solidity: function gettL1RedTokenAddress() view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCaller) GettL1RedTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChildBridgeCoreFacet.contract.Call(opts, &out, "gettL1RedTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GettL1RedTokenAddress is a free data retrieval call binding the contract method 0xdada7d76.
//
// Solidity: function gettL1RedTokenAddress() view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) GettL1RedTokenAddress() (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GettL1RedTokenAddress(&_ChildBridgeCoreFacet.CallOpts)
}

// GettL1RedTokenAddress is a free data retrieval call binding the contract method 0xdada7d76.
//
// Solidity: function gettL1RedTokenAddress() view returns(address)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCallerSession) GettL1RedTokenAddress() (common.Address, error) {
	return _ChildBridgeCoreFacet.Contract.GettL1RedTokenAddress(&_ChildBridgeCoreFacet.CallOpts)
}

// PauseStatusBridge is a free data retrieval call binding the contract method 0xe7616cb0.
//
// Solidity: function pauseStatusBridge() view returns(bool)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCaller) PauseStatusBridge(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ChildBridgeCoreFacet.contract.Call(opts, &out, "pauseStatusBridge")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PauseStatusBridge is a free data retrieval call binding the contract method 0xe7616cb0.
//
// Solidity: function pauseStatusBridge() view returns(bool)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) PauseStatusBridge() (bool, error) {
	return _ChildBridgeCoreFacet.Contract.PauseStatusBridge(&_ChildBridgeCoreFacet.CallOpts)
}

// PauseStatusBridge is a free data retrieval call binding the contract method 0xe7616cb0.
//
// Solidity: function pauseStatusBridge() view returns(bool)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetCallerSession) PauseStatusBridge() (bool, error) {
	return _ChildBridgeCoreFacet.Contract.PauseStatusBridge(&_ChildBridgeCoreFacet.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) Initialize() (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.Initialize(&_ChildBridgeCoreFacet.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactorSession) Initialize() (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.Initialize(&_ChildBridgeCoreFacet.TransactOpts)
}

// PauseBridge is a paid mutator transaction binding the contract method 0x7dd0480f.
//
// Solidity: function pauseBridge() returns(bool)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactor) PauseBridge(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.contract.Transact(opts, "pauseBridge")
}

// PauseBridge is a paid mutator transaction binding the contract method 0x7dd0480f.
//
// Solidity: function pauseBridge() returns(bool)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) PauseBridge() (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.PauseBridge(&_ChildBridgeCoreFacet.TransactOpts)
}

// PauseBridge is a paid mutator transaction binding the contract method 0x7dd0480f.
//
// Solidity: function pauseBridge() returns(bool)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactorSession) PauseBridge() (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.PauseBridge(&_ChildBridgeCoreFacet.TransactOpts)
}

// SendUpwardMessage is a paid mutator transaction binding the contract method 0xd3609a28.
//
// Solidity: function sendUpwardMessage(uint32 payloadType, bytes payload) returns()
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactor) SendUpwardMessage(opts *bind.TransactOpts, payloadType uint32, payload []byte) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.contract.Transact(opts, "sendUpwardMessage", payloadType, payload)
}

// SendUpwardMessage is a paid mutator transaction binding the contract method 0xd3609a28.
//
// Solidity: function sendUpwardMessage(uint32 payloadType, bytes payload) returns()
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) SendUpwardMessage(payloadType uint32, payload []byte) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.SendUpwardMessage(&_ChildBridgeCoreFacet.TransactOpts, payloadType, payload)
}

// SendUpwardMessage is a paid mutator transaction binding the contract method 0xd3609a28.
//
// Solidity: function sendUpwardMessage(uint32 payloadType, bytes payload) returns()
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactorSession) SendUpwardMessage(payloadType uint32, payload []byte) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.SendUpwardMessage(&_ChildBridgeCoreFacet.TransactOpts, payloadType, payload)
}

// SetRedTokenAddress is a paid mutator transaction binding the contract method 0xfd3d6a8d.
//
// Solidity: function setRedTokenAddress(address l1RedTokenAddress) returns()
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactor) SetRedTokenAddress(opts *bind.TransactOpts, l1RedTokenAddress common.Address) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.contract.Transact(opts, "setRedTokenAddress", l1RedTokenAddress)
}

// SetRedTokenAddress is a paid mutator transaction binding the contract method 0xfd3d6a8d.
//
// Solidity: function setRedTokenAddress(address l1RedTokenAddress) returns()
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) SetRedTokenAddress(l1RedTokenAddress common.Address) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.SetRedTokenAddress(&_ChildBridgeCoreFacet.TransactOpts, l1RedTokenAddress)
}

// SetRedTokenAddress is a paid mutator transaction binding the contract method 0xfd3d6a8d.
//
// Solidity: function setRedTokenAddress(address l1RedTokenAddress) returns()
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactorSession) SetRedTokenAddress(l1RedTokenAddress common.Address) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.SetRedTokenAddress(&_ChildBridgeCoreFacet.TransactOpts, l1RedTokenAddress)
}

// UnpauseBridge is a paid mutator transaction binding the contract method 0xa82f143c.
//
// Solidity: function unpauseBridge() returns(bool)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactor) UnpauseBridge(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.contract.Transact(opts, "unpauseBridge")
}

// UnpauseBridge is a paid mutator transaction binding the contract method 0xa82f143c.
//
// Solidity: function unpauseBridge() returns(bool)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetSession) UnpauseBridge() (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.UnpauseBridge(&_ChildBridgeCoreFacet.TransactOpts)
}

// UnpauseBridge is a paid mutator transaction binding the contract method 0xa82f143c.
//
// Solidity: function unpauseBridge() returns(bool)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetTransactorSession) UnpauseBridge() (*types.Transaction, error) {
	return _ChildBridgeCoreFacet.Contract.UnpauseBridge(&_ChildBridgeCoreFacet.TransactOpts)
}

// ChildBridgeCoreFacetAppendMessageEventIterator is returned from FilterAppendMessageEvent and is used to iterate over the raw logs and unpacked data for AppendMessageEvent events raised by the ChildBridgeCoreFacet contract.
type ChildBridgeCoreFacetAppendMessageEventIterator struct {
	Event *ChildBridgeCoreFacetAppendMessageEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChildBridgeCoreFacetAppendMessageEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChildBridgeCoreFacetAppendMessageEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChildBridgeCoreFacetAppendMessageEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChildBridgeCoreFacetAppendMessageEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChildBridgeCoreFacetAppendMessageEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChildBridgeCoreFacetAppendMessageEvent represents a AppendMessageEvent event raised by the ChildBridgeCoreFacet contract.
type ChildBridgeCoreFacetAppendMessageEvent struct {
	Index       *big.Int
	MessageHash [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAppendMessageEvent is a free log retrieval operation binding the contract event 0xbd09f2f4cb9fdb26ba06cc07aff189ef87d22a6ae5a6dfa0f1dad34b348d1f3d.
//
// Solidity: event AppendMessageEvent(uint256 index, bytes32 messageHash)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) FilterAppendMessageEvent(opts *bind.FilterOpts) (*ChildBridgeCoreFacetAppendMessageEventIterator, error) {

	logs, sub, err := _ChildBridgeCoreFacet.contract.FilterLogs(opts, "AppendMessageEvent")
	if err != nil {
		return nil, err
	}
	return &ChildBridgeCoreFacetAppendMessageEventIterator{contract: _ChildBridgeCoreFacet.contract, event: "AppendMessageEvent", logs: logs, sub: sub}, nil
}

// WatchAppendMessageEvent is a free log subscription operation binding the contract event 0xbd09f2f4cb9fdb26ba06cc07aff189ef87d22a6ae5a6dfa0f1dad34b348d1f3d.
//
// Solidity: event AppendMessageEvent(uint256 index, bytes32 messageHash)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) WatchAppendMessageEvent(opts *bind.WatchOpts, sink chan<- *ChildBridgeCoreFacetAppendMessageEvent) (event.Subscription, error) {

	logs, sub, err := _ChildBridgeCoreFacet.contract.WatchLogs(opts, "AppendMessageEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChildBridgeCoreFacetAppendMessageEvent)
				if err := _ChildBridgeCoreFacet.contract.UnpackLog(event, "AppendMessageEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAppendMessageEvent is a log parse operation binding the contract event 0xbd09f2f4cb9fdb26ba06cc07aff189ef87d22a6ae5a6dfa0f1dad34b348d1f3d.
//
// Solidity: event AppendMessageEvent(uint256 index, bytes32 messageHash)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) ParseAppendMessageEvent(log types.Log) (*ChildBridgeCoreFacetAppendMessageEvent, error) {
	event := new(ChildBridgeCoreFacetAppendMessageEvent)
	if err := _ChildBridgeCoreFacet.contract.UnpackLog(event, "AppendMessageEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChildBridgeCoreFacetSentMessageIterator is returned from FilterSentMessage and is used to iterate over the raw logs and unpacked data for SentMessage events raised by the ChildBridgeCoreFacet contract.
type ChildBridgeCoreFacetSentMessageIterator struct {
	Event *ChildBridgeCoreFacetSentMessage // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChildBridgeCoreFacetSentMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChildBridgeCoreFacetSentMessage)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChildBridgeCoreFacetSentMessage)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChildBridgeCoreFacetSentMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChildBridgeCoreFacetSentMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChildBridgeCoreFacetSentMessage represents a SentMessage event raised by the ChildBridgeCoreFacet contract.
type ChildBridgeCoreFacetSentMessage struct {
	XDomainCalldataHash [32]byte
	Nonce               *big.Int
	PayloadType         uint32
	Payload             []byte
	GasLimit            *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSentMessage is a free log retrieval operation binding the contract event 0x9717419fccdf63bc6c14a8f84640f92755acddeab6683f96607d2e422acf4676.
//
// Solidity: event SentMessage(bytes32 indexed xDomainCalldataHash, uint256 nonce, uint32 payloadType, bytes payload, uint256 gasLimit)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) FilterSentMessage(opts *bind.FilterOpts, xDomainCalldataHash [][32]byte) (*ChildBridgeCoreFacetSentMessageIterator, error) {

	var xDomainCalldataHashRule []interface{}
	for _, xDomainCalldataHashItem := range xDomainCalldataHash {
		xDomainCalldataHashRule = append(xDomainCalldataHashRule, xDomainCalldataHashItem)
	}

	logs, sub, err := _ChildBridgeCoreFacet.contract.FilterLogs(opts, "SentMessage", xDomainCalldataHashRule)
	if err != nil {
		return nil, err
	}
	return &ChildBridgeCoreFacetSentMessageIterator{contract: _ChildBridgeCoreFacet.contract, event: "SentMessage", logs: logs, sub: sub}, nil
}

// WatchSentMessage is a free log subscription operation binding the contract event 0x9717419fccdf63bc6c14a8f84640f92755acddeab6683f96607d2e422acf4676.
//
// Solidity: event SentMessage(bytes32 indexed xDomainCalldataHash, uint256 nonce, uint32 payloadType, bytes payload, uint256 gasLimit)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) WatchSentMessage(opts *bind.WatchOpts, sink chan<- *ChildBridgeCoreFacetSentMessage, xDomainCalldataHash [][32]byte) (event.Subscription, error) {

	var xDomainCalldataHashRule []interface{}
	for _, xDomainCalldataHashItem := range xDomainCalldataHash {
		xDomainCalldataHashRule = append(xDomainCalldataHashRule, xDomainCalldataHashItem)
	}

	logs, sub, err := _ChildBridgeCoreFacet.contract.WatchLogs(opts, "SentMessage", xDomainCalldataHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChildBridgeCoreFacetSentMessage)
				if err := _ChildBridgeCoreFacet.contract.UnpackLog(event, "SentMessage", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSentMessage is a log parse operation binding the contract event 0x9717419fccdf63bc6c14a8f84640f92755acddeab6683f96607d2e422acf4676.
//
// Solidity: event SentMessage(bytes32 indexed xDomainCalldataHash, uint256 nonce, uint32 payloadType, bytes payload, uint256 gasLimit)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) ParseSentMessage(log types.Log) (*ChildBridgeCoreFacetSentMessage, error) {
	event := new(ChildBridgeCoreFacetSentMessage)
	if err := _ChildBridgeCoreFacet.contract.UnpackLog(event, "SentMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChildBridgeCoreFacetUpwardMessageIterator is returned from FilterUpwardMessage and is used to iterate over the raw logs and unpacked data for UpwardMessage events raised by the ChildBridgeCoreFacet contract.
type ChildBridgeCoreFacetUpwardMessageIterator struct {
	Event *ChildBridgeCoreFacetUpwardMessage // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChildBridgeCoreFacetUpwardMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChildBridgeCoreFacetUpwardMessage)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChildBridgeCoreFacetUpwardMessage)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChildBridgeCoreFacetUpwardMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChildBridgeCoreFacetUpwardMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChildBridgeCoreFacetUpwardMessage represents a UpwardMessage event raised by the ChildBridgeCoreFacet contract.
type ChildBridgeCoreFacetUpwardMessage struct {
	PayloadType uint32
	Payload     []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpwardMessage is a free log retrieval operation binding the contract event 0x71682e85596cab934fc449074db8a6e222fd93232a40ff6bfa37745ab8bb085f.
//
// Solidity: event UpwardMessage(uint32 payloadType, bytes payload)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) FilterUpwardMessage(opts *bind.FilterOpts) (*ChildBridgeCoreFacetUpwardMessageIterator, error) {

	logs, sub, err := _ChildBridgeCoreFacet.contract.FilterLogs(opts, "UpwardMessage")
	if err != nil {
		return nil, err
	}
	return &ChildBridgeCoreFacetUpwardMessageIterator{contract: _ChildBridgeCoreFacet.contract, event: "UpwardMessage", logs: logs, sub: sub}, nil
}

// WatchUpwardMessage is a free log subscription operation binding the contract event 0x71682e85596cab934fc449074db8a6e222fd93232a40ff6bfa37745ab8bb085f.
//
// Solidity: event UpwardMessage(uint32 payloadType, bytes payload)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) WatchUpwardMessage(opts *bind.WatchOpts, sink chan<- *ChildBridgeCoreFacetUpwardMessage) (event.Subscription, error) {

	logs, sub, err := _ChildBridgeCoreFacet.contract.WatchLogs(opts, "UpwardMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChildBridgeCoreFacetUpwardMessage)
				if err := _ChildBridgeCoreFacet.contract.UnpackLog(event, "UpwardMessage", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpwardMessage is a log parse operation binding the contract event 0x71682e85596cab934fc449074db8a6e222fd93232a40ff6bfa37745ab8bb085f.
//
// Solidity: event UpwardMessage(uint32 payloadType, bytes payload)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) ParseUpwardMessage(log types.Log) (*ChildBridgeCoreFacetUpwardMessage, error) {
	event := new(ChildBridgeCoreFacetUpwardMessage)
	if err := _ChildBridgeCoreFacet.contract.UnpackLog(event, "UpwardMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
