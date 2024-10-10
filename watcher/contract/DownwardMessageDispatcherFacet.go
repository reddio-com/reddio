// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// DownwardMessageDispatcherFacetDownwardMessage is an auto generated low-level Go binding around an user-defined struct.
type DownwardMessageDispatcherFacetDownwardMessage struct {
	Sequence    *big.Int
	PayloadType uint32
	Payload     []byte
}

// DownwardMessageDispatcherFacetMetaData contains all meta data concerning the DownwardMessageDispatcherFacet contract.
var DownwardMessageDispatcherFacetMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structDownwardMessageDispatcherFacet.DownwardMessage[]\",\"name\":\"downwardMessages\",\"type\":\"tuple[]\"}],\"name\":\"receiveDownwardMessages\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DownwardMessageDispatcherFacetABI is the input ABI used to generate the binding from.
// Deprecated: Use DownwardMessageDispatcherFacetMetaData.ABI instead.
var DownwardMessageDispatcherFacetABI = DownwardMessageDispatcherFacetMetaData.ABI

// DownwardMessageDispatcherFacet is an auto generated Go binding around an Ethereum contract.
type DownwardMessageDispatcherFacet struct {
	DownwardMessageDispatcherFacetCaller     // Read-only binding to the contract
	DownwardMessageDispatcherFacetTransactor // Write-only binding to the contract
	DownwardMessageDispatcherFacetFilterer   // Log filterer for contract events
}

// DownwardMessageDispatcherFacetCaller is an auto generated read-only Go binding around an Ethereum contract.
type DownwardMessageDispatcherFacetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DownwardMessageDispatcherFacetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DownwardMessageDispatcherFacetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DownwardMessageDispatcherFacetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DownwardMessageDispatcherFacetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DownwardMessageDispatcherFacetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DownwardMessageDispatcherFacetSession struct {
	Contract     *DownwardMessageDispatcherFacet // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// DownwardMessageDispatcherFacetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DownwardMessageDispatcherFacetCallerSession struct {
	Contract *DownwardMessageDispatcherFacetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// DownwardMessageDispatcherFacetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DownwardMessageDispatcherFacetTransactorSession struct {
	Contract     *DownwardMessageDispatcherFacetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// DownwardMessageDispatcherFacetRaw is an auto generated low-level Go binding around an Ethereum contract.
type DownwardMessageDispatcherFacetRaw struct {
	Contract *DownwardMessageDispatcherFacet // Generic contract binding to access the raw methods on
}

// DownwardMessageDispatcherFacetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DownwardMessageDispatcherFacetCallerRaw struct {
	Contract *DownwardMessageDispatcherFacetCaller // Generic read-only contract binding to access the raw methods on
}

// DownwardMessageDispatcherFacetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DownwardMessageDispatcherFacetTransactorRaw struct {
	Contract *DownwardMessageDispatcherFacetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDownwardMessageDispatcherFacet creates a new instance of DownwardMessageDispatcherFacet, bound to a specific deployed contract.
func NewDownwardMessageDispatcherFacet(address common.Address, backend bind.ContractBackend) (*DownwardMessageDispatcherFacet, error) {
	contract, err := bindDownwardMessageDispatcherFacet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DownwardMessageDispatcherFacet{DownwardMessageDispatcherFacetCaller: DownwardMessageDispatcherFacetCaller{contract: contract}, DownwardMessageDispatcherFacetTransactor: DownwardMessageDispatcherFacetTransactor{contract: contract}, DownwardMessageDispatcherFacetFilterer: DownwardMessageDispatcherFacetFilterer{contract: contract}}, nil
}

// NewDownwardMessageDispatcherFacetCaller creates a new read-only instance of DownwardMessageDispatcherFacet, bound to a specific deployed contract.
func NewDownwardMessageDispatcherFacetCaller(address common.Address, caller bind.ContractCaller) (*DownwardMessageDispatcherFacetCaller, error) {
	contract, err := bindDownwardMessageDispatcherFacet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DownwardMessageDispatcherFacetCaller{contract: contract}, nil
}

// NewDownwardMessageDispatcherFacetTransactor creates a new write-only instance of DownwardMessageDispatcherFacet, bound to a specific deployed contract.
func NewDownwardMessageDispatcherFacetTransactor(address common.Address, transactor bind.ContractTransactor) (*DownwardMessageDispatcherFacetTransactor, error) {
	contract, err := bindDownwardMessageDispatcherFacet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DownwardMessageDispatcherFacetTransactor{contract: contract}, nil
}

// NewDownwardMessageDispatcherFacetFilterer creates a new log filterer instance of DownwardMessageDispatcherFacet, bound to a specific deployed contract.
func NewDownwardMessageDispatcherFacetFilterer(address common.Address, filterer bind.ContractFilterer) (*DownwardMessageDispatcherFacetFilterer, error) {
	contract, err := bindDownwardMessageDispatcherFacet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DownwardMessageDispatcherFacetFilterer{contract: contract}, nil
}

// bindDownwardMessageDispatcherFacet binds a generic wrapper to an already deployed contract.
func bindDownwardMessageDispatcherFacet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DownwardMessageDispatcherFacetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DownwardMessageDispatcherFacet.Contract.DownwardMessageDispatcherFacetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.DownwardMessageDispatcherFacetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.DownwardMessageDispatcherFacetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DownwardMessageDispatcherFacet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.contract.Transact(opts, method, params...)
}

// ReceiveDownwardMessages is a paid mutator transaction binding the contract method 0x1f94d010.
//
// Solidity: function receiveDownwardMessages((uint256,uint32,bytes)[] downwardMessages) returns()
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetTransactor) ReceiveDownwardMessages(opts *bind.TransactOpts, downwardMessages []DownwardMessageDispatcherFacetDownwardMessage) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.contract.Transact(opts, "receiveDownwardMessages", downwardMessages)
}

// ReceiveDownwardMessages is a paid mutator transaction binding the contract method 0x1f94d010.
//
// Solidity: function receiveDownwardMessages((uint256,uint32,bytes)[] downwardMessages) returns()
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetSession) ReceiveDownwardMessages(downwardMessages []DownwardMessageDispatcherFacetDownwardMessage) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.ReceiveDownwardMessages(&_DownwardMessageDispatcherFacet.TransactOpts, downwardMessages)
}

// ReceiveDownwardMessages is a paid mutator transaction binding the contract method 0x1f94d010.
//
// Solidity: function receiveDownwardMessages((uint256,uint32,bytes)[] downwardMessages) returns()
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetTransactorSession) ReceiveDownwardMessages(downwardMessages []DownwardMessageDispatcherFacetDownwardMessage) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.ReceiveDownwardMessages(&_DownwardMessageDispatcherFacet.TransactOpts, downwardMessages)
}
