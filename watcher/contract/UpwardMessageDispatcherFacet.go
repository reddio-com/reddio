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

// UpwardMessage is an auto generated low-level Go binding around an user-defined struct.
type UpwardMessage struct {
	Sequence    *big.Int
	PayloadType uint32
	Payload     []byte
}

// UpwardMessageDispatcherFacetMetaData contains all meta data concerning the UpwardMessageDispatcherFacet contract.
var UpwardMessageDispatcherFacetMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structUpwardMessage[]\",\"name\":\"upwardMessages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signaturesArray\",\"type\":\"bytes[]\"}],\"name\":\"receiveUpwardMessages\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// UpwardMessageDispatcherFacetABI is the input ABI used to generate the binding from.
// Deprecated: Use UpwardMessageDispatcherFacetMetaData.ABI instead.
var UpwardMessageDispatcherFacetABI = UpwardMessageDispatcherFacetMetaData.ABI

// UpwardMessageDispatcherFacet is an auto generated Go binding around an Ethereum contract.
type UpwardMessageDispatcherFacet struct {
	UpwardMessageDispatcherFacetCaller     // Read-only binding to the contract
	UpwardMessageDispatcherFacetTransactor // Write-only binding to the contract
	UpwardMessageDispatcherFacetFilterer   // Log filterer for contract events
}

// UpwardMessageDispatcherFacetCaller is an auto generated read-only Go binding around an Ethereum contract.
type UpwardMessageDispatcherFacetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpwardMessageDispatcherFacetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpwardMessageDispatcherFacetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpwardMessageDispatcherFacetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpwardMessageDispatcherFacetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpwardMessageDispatcherFacetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpwardMessageDispatcherFacetSession struct {
	Contract     *UpwardMessageDispatcherFacet // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// UpwardMessageDispatcherFacetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpwardMessageDispatcherFacetCallerSession struct {
	Contract *UpwardMessageDispatcherFacetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// UpwardMessageDispatcherFacetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpwardMessageDispatcherFacetTransactorSession struct {
	Contract     *UpwardMessageDispatcherFacetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// UpwardMessageDispatcherFacetRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpwardMessageDispatcherFacetRaw struct {
	Contract *UpwardMessageDispatcherFacet // Generic contract binding to access the raw methods on
}

// UpwardMessageDispatcherFacetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpwardMessageDispatcherFacetCallerRaw struct {
	Contract *UpwardMessageDispatcherFacetCaller // Generic read-only contract binding to access the raw methods on
}

// UpwardMessageDispatcherFacetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpwardMessageDispatcherFacetTransactorRaw struct {
	Contract *UpwardMessageDispatcherFacetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpwardMessageDispatcherFacet creates a new instance of UpwardMessageDispatcherFacet, bound to a specific deployed contract.
func NewUpwardMessageDispatcherFacet(address common.Address, backend bind.ContractBackend) (*UpwardMessageDispatcherFacet, error) {
	contract, err := bindUpwardMessageDispatcherFacet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpwardMessageDispatcherFacet{UpwardMessageDispatcherFacetCaller: UpwardMessageDispatcherFacetCaller{contract: contract}, UpwardMessageDispatcherFacetTransactor: UpwardMessageDispatcherFacetTransactor{contract: contract}, UpwardMessageDispatcherFacetFilterer: UpwardMessageDispatcherFacetFilterer{contract: contract}}, nil
}

// NewUpwardMessageDispatcherFacetCaller creates a new read-only instance of UpwardMessageDispatcherFacet, bound to a specific deployed contract.
func NewUpwardMessageDispatcherFacetCaller(address common.Address, caller bind.ContractCaller) (*UpwardMessageDispatcherFacetCaller, error) {
	contract, err := bindUpwardMessageDispatcherFacet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpwardMessageDispatcherFacetCaller{contract: contract}, nil
}

// NewUpwardMessageDispatcherFacetTransactor creates a new write-only instance of UpwardMessageDispatcherFacet, bound to a specific deployed contract.
func NewUpwardMessageDispatcherFacetTransactor(address common.Address, transactor bind.ContractTransactor) (*UpwardMessageDispatcherFacetTransactor, error) {
	contract, err := bindUpwardMessageDispatcherFacet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpwardMessageDispatcherFacetTransactor{contract: contract}, nil
}

// NewUpwardMessageDispatcherFacetFilterer creates a new log filterer instance of UpwardMessageDispatcherFacet, bound to a specific deployed contract.
func NewUpwardMessageDispatcherFacetFilterer(address common.Address, filterer bind.ContractFilterer) (*UpwardMessageDispatcherFacetFilterer, error) {
	contract, err := bindUpwardMessageDispatcherFacet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpwardMessageDispatcherFacetFilterer{contract: contract}, nil
}

// bindUpwardMessageDispatcherFacet binds a generic wrapper to an already deployed contract.
func bindUpwardMessageDispatcherFacet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UpwardMessageDispatcherFacetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpwardMessageDispatcherFacet.Contract.UpwardMessageDispatcherFacetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.UpwardMessageDispatcherFacetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.UpwardMessageDispatcherFacetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpwardMessageDispatcherFacet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.contract.Transact(opts, method, params...)
}

// ReceiveUpwardMessages is a paid mutator transaction binding the contract method 0x37a611ba.
//
// Solidity: function receiveUpwardMessages((uint256,uint32,bytes)[] upwardMessages, bytes[] signaturesArray) returns()
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetTransactor) ReceiveUpwardMessages(opts *bind.TransactOpts, upwardMessages []UpwardMessage, signaturesArray [][]byte) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.contract.Transact(opts, "receiveUpwardMessages", upwardMessages, signaturesArray)
}

// ReceiveUpwardMessages is a paid mutator transaction binding the contract method 0x37a611ba.
//
// Solidity: function receiveUpwardMessages((uint256,uint32,bytes)[] upwardMessages, bytes[] signaturesArray) returns()
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetSession) ReceiveUpwardMessages(upwardMessages []UpwardMessage, signaturesArray [][]byte) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.ReceiveUpwardMessages(&_UpwardMessageDispatcherFacet.TransactOpts, upwardMessages, signaturesArray)
}

// ReceiveUpwardMessages is a paid mutator transaction binding the contract method 0x37a611ba.
//
// Solidity: function receiveUpwardMessages((uint256,uint32,bytes)[] upwardMessages, bytes[] signaturesArray) returns()
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetTransactorSession) ReceiveUpwardMessages(upwardMessages []UpwardMessage, signaturesArray [][]byte) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.ReceiveUpwardMessages(&_UpwardMessageDispatcherFacet.TransactOpts, upwardMessages, signaturesArray)
}
