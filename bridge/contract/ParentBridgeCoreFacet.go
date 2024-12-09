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

// ParentBridgeCoreFacetMetaData contains all meta data concerning the ParentBridgeCoreFacet contract.
var ParentBridgeCoreFacetMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"DownwardMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"queueIndex\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"QueueTransaction\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"name\":\"RelayedMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"}],\"name\":\"estimateCrossMessageFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gettL1RedTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextCrossDomainMessageIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseBridge\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseStatusBridge\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"ethAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"sendDownwardMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"l1RedTokenAddress\",\"type\":\"address\"}],\"name\":\"setL1RedTokenAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpauseBridge\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ParentBridgeCoreFacetABI is the input ABI used to generate the binding from.
// Deprecated: Use ParentBridgeCoreFacetMetaData.ABI instead.
var ParentBridgeCoreFacetABI = ParentBridgeCoreFacetMetaData.ABI

// ParentBridgeCoreFacet is an auto generated Go binding around an Ethereum contract.
type ParentBridgeCoreFacet struct {
	ParentBridgeCoreFacetCaller     // Read-only binding to the contract
	ParentBridgeCoreFacetTransactor // Write-only binding to the contract
	ParentBridgeCoreFacetFilterer   // Log filterer for contract events
}

// ParentBridgeCoreFacetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ParentBridgeCoreFacetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ParentBridgeCoreFacetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ParentBridgeCoreFacetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ParentBridgeCoreFacetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ParentBridgeCoreFacetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ParentBridgeCoreFacetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ParentBridgeCoreFacetSession struct {
	Contract     *ParentBridgeCoreFacet // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ParentBridgeCoreFacetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ParentBridgeCoreFacetCallerSession struct {
	Contract *ParentBridgeCoreFacetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ParentBridgeCoreFacetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ParentBridgeCoreFacetTransactorSession struct {
	Contract     *ParentBridgeCoreFacetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ParentBridgeCoreFacetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ParentBridgeCoreFacetRaw struct {
	Contract *ParentBridgeCoreFacet // Generic contract binding to access the raw methods on
}

// ParentBridgeCoreFacetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ParentBridgeCoreFacetCallerRaw struct {
	Contract *ParentBridgeCoreFacetCaller // Generic read-only contract binding to access the raw methods on
}

// ParentBridgeCoreFacetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ParentBridgeCoreFacetTransactorRaw struct {
	Contract *ParentBridgeCoreFacetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewParentBridgeCoreFacet creates a new instance of ParentBridgeCoreFacet, bound to a specific deployed contract.
func NewParentBridgeCoreFacet(address common.Address, backend bind.ContractBackend) (*ParentBridgeCoreFacet, error) {
	contract, err := bindParentBridgeCoreFacet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ParentBridgeCoreFacet{ParentBridgeCoreFacetCaller: ParentBridgeCoreFacetCaller{contract: contract}, ParentBridgeCoreFacetTransactor: ParentBridgeCoreFacetTransactor{contract: contract}, ParentBridgeCoreFacetFilterer: ParentBridgeCoreFacetFilterer{contract: contract}}, nil
}

// NewParentBridgeCoreFacetCaller creates a new read-only instance of ParentBridgeCoreFacet, bound to a specific deployed contract.
func NewParentBridgeCoreFacetCaller(address common.Address, caller bind.ContractCaller) (*ParentBridgeCoreFacetCaller, error) {
	contract, err := bindParentBridgeCoreFacet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ParentBridgeCoreFacetCaller{contract: contract}, nil
}

// NewParentBridgeCoreFacetTransactor creates a new write-only instance of ParentBridgeCoreFacet, bound to a specific deployed contract.
func NewParentBridgeCoreFacetTransactor(address common.Address, transactor bind.ContractTransactor) (*ParentBridgeCoreFacetTransactor, error) {
	contract, err := bindParentBridgeCoreFacet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ParentBridgeCoreFacetTransactor{contract: contract}, nil
}

// NewParentBridgeCoreFacetFilterer creates a new log filterer instance of ParentBridgeCoreFacet, bound to a specific deployed contract.
func NewParentBridgeCoreFacetFilterer(address common.Address, filterer bind.ContractFilterer) (*ParentBridgeCoreFacetFilterer, error) {
	contract, err := bindParentBridgeCoreFacet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ParentBridgeCoreFacetFilterer{contract: contract}, nil
}

// bindParentBridgeCoreFacet binds a generic wrapper to an already deployed contract.
func bindParentBridgeCoreFacet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ParentBridgeCoreFacetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ParentBridgeCoreFacet.Contract.ParentBridgeCoreFacetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.ParentBridgeCoreFacetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.ParentBridgeCoreFacetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ParentBridgeCoreFacet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.contract.Transact(opts, method, params...)
}

// EstimateCrossMessageFee is a free data retrieval call binding the contract method 0x5bf6b67e.
//
// Solidity: function estimateCrossMessageFee(uint256 _gasLimit) view returns(uint256)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetCaller) EstimateCrossMessageFee(opts *bind.CallOpts, _gasLimit *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ParentBridgeCoreFacet.contract.Call(opts, &out, "estimateCrossMessageFee", _gasLimit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EstimateCrossMessageFee is a free data retrieval call binding the contract method 0x5bf6b67e.
//
// Solidity: function estimateCrossMessageFee(uint256 _gasLimit) view returns(uint256)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetSession) EstimateCrossMessageFee(_gasLimit *big.Int) (*big.Int, error) {
	return _ParentBridgeCoreFacet.Contract.EstimateCrossMessageFee(&_ParentBridgeCoreFacet.CallOpts, _gasLimit)
}

// EstimateCrossMessageFee is a free data retrieval call binding the contract method 0x5bf6b67e.
//
// Solidity: function estimateCrossMessageFee(uint256 _gasLimit) view returns(uint256)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetCallerSession) EstimateCrossMessageFee(_gasLimit *big.Int) (*big.Int, error) {
	return _ParentBridgeCoreFacet.Contract.EstimateCrossMessageFee(&_ParentBridgeCoreFacet.CallOpts, _gasLimit)
}

// GettL1RedTokenAddress is a free data retrieval call binding the contract method 0xdada7d76.
//
// Solidity: function gettL1RedTokenAddress() view returns(address)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetCaller) GettL1RedTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ParentBridgeCoreFacet.contract.Call(opts, &out, "gettL1RedTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GettL1RedTokenAddress is a free data retrieval call binding the contract method 0xdada7d76.
//
// Solidity: function gettL1RedTokenAddress() view returns(address)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetSession) GettL1RedTokenAddress() (common.Address, error) {
	return _ParentBridgeCoreFacet.Contract.GettL1RedTokenAddress(&_ParentBridgeCoreFacet.CallOpts)
}

// GettL1RedTokenAddress is a free data retrieval call binding the contract method 0xdada7d76.
//
// Solidity: function gettL1RedTokenAddress() view returns(address)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetCallerSession) GettL1RedTokenAddress() (common.Address, error) {
	return _ParentBridgeCoreFacet.Contract.GettL1RedTokenAddress(&_ParentBridgeCoreFacet.CallOpts)
}

// NextCrossDomainMessageIndex is a free data retrieval call binding the contract method 0xfd0ad31e.
//
// Solidity: function nextCrossDomainMessageIndex() view returns(uint256)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetCaller) NextCrossDomainMessageIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ParentBridgeCoreFacet.contract.Call(opts, &out, "nextCrossDomainMessageIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextCrossDomainMessageIndex is a free data retrieval call binding the contract method 0xfd0ad31e.
//
// Solidity: function nextCrossDomainMessageIndex() view returns(uint256)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetSession) NextCrossDomainMessageIndex() (*big.Int, error) {
	return _ParentBridgeCoreFacet.Contract.NextCrossDomainMessageIndex(&_ParentBridgeCoreFacet.CallOpts)
}

// NextCrossDomainMessageIndex is a free data retrieval call binding the contract method 0xfd0ad31e.
//
// Solidity: function nextCrossDomainMessageIndex() view returns(uint256)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetCallerSession) NextCrossDomainMessageIndex() (*big.Int, error) {
	return _ParentBridgeCoreFacet.Contract.NextCrossDomainMessageIndex(&_ParentBridgeCoreFacet.CallOpts)
}

// PauseStatusBridge is a free data retrieval call binding the contract method 0xe7616cb0.
//
// Solidity: function pauseStatusBridge() view returns(bool)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetCaller) PauseStatusBridge(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ParentBridgeCoreFacet.contract.Call(opts, &out, "pauseStatusBridge")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PauseStatusBridge is a free data retrieval call binding the contract method 0xe7616cb0.
//
// Solidity: function pauseStatusBridge() view returns(bool)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetSession) PauseStatusBridge() (bool, error) {
	return _ParentBridgeCoreFacet.Contract.PauseStatusBridge(&_ParentBridgeCoreFacet.CallOpts)
}

// PauseStatusBridge is a free data retrieval call binding the contract method 0xe7616cb0.
//
// Solidity: function pauseStatusBridge() view returns(bool)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetCallerSession) PauseStatusBridge() (bool, error) {
	return _ParentBridgeCoreFacet.Contract.PauseStatusBridge(&_ParentBridgeCoreFacet.CallOpts)
}

// PauseBridge is a paid mutator transaction binding the contract method 0x7dd0480f.
//
// Solidity: function pauseBridge() returns(bool)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactor) PauseBridge(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.contract.Transact(opts, "pauseBridge")
}

// PauseBridge is a paid mutator transaction binding the contract method 0x7dd0480f.
//
// Solidity: function pauseBridge() returns(bool)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetSession) PauseBridge() (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.PauseBridge(&_ParentBridgeCoreFacet.TransactOpts)
}

// PauseBridge is a paid mutator transaction binding the contract method 0x7dd0480f.
//
// Solidity: function pauseBridge() returns(bool)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactorSession) PauseBridge() (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.PauseBridge(&_ParentBridgeCoreFacet.TransactOpts)
}

// SendDownwardMessage is a paid mutator transaction binding the contract method 0xd479ceaf.
//
// Solidity: function sendDownwardMessage(uint32 payloadType, bytes payload, uint256 ethAmount, uint256 gasLimit, uint256 value) returns()
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactor) SendDownwardMessage(opts *bind.TransactOpts, payloadType uint32, payload []byte, ethAmount *big.Int, gasLimit *big.Int, value *big.Int) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.contract.Transact(opts, "sendDownwardMessage", payloadType, payload, ethAmount, gasLimit, value)
}

// SendDownwardMessage is a paid mutator transaction binding the contract method 0xd479ceaf.
//
// Solidity: function sendDownwardMessage(uint32 payloadType, bytes payload, uint256 ethAmount, uint256 gasLimit, uint256 value) returns()
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetSession) SendDownwardMessage(payloadType uint32, payload []byte, ethAmount *big.Int, gasLimit *big.Int, value *big.Int) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.SendDownwardMessage(&_ParentBridgeCoreFacet.TransactOpts, payloadType, payload, ethAmount, gasLimit, value)
}

// SendDownwardMessage is a paid mutator transaction binding the contract method 0xd479ceaf.
//
// Solidity: function sendDownwardMessage(uint32 payloadType, bytes payload, uint256 ethAmount, uint256 gasLimit, uint256 value) returns()
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactorSession) SendDownwardMessage(payloadType uint32, payload []byte, ethAmount *big.Int, gasLimit *big.Int, value *big.Int) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.SendDownwardMessage(&_ParentBridgeCoreFacet.TransactOpts, payloadType, payload, ethAmount, gasLimit, value)
}

// SetL1RedTokenAddress is a paid mutator transaction binding the contract method 0x1224811e.
//
// Solidity: function setL1RedTokenAddress(address l1RedTokenAddress) returns()
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactor) SetL1RedTokenAddress(opts *bind.TransactOpts, l1RedTokenAddress common.Address) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.contract.Transact(opts, "setL1RedTokenAddress", l1RedTokenAddress)
}

// SetL1RedTokenAddress is a paid mutator transaction binding the contract method 0x1224811e.
//
// Solidity: function setL1RedTokenAddress(address l1RedTokenAddress) returns()
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetSession) SetL1RedTokenAddress(l1RedTokenAddress common.Address) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.SetL1RedTokenAddress(&_ParentBridgeCoreFacet.TransactOpts, l1RedTokenAddress)
}

// SetL1RedTokenAddress is a paid mutator transaction binding the contract method 0x1224811e.
//
// Solidity: function setL1RedTokenAddress(address l1RedTokenAddress) returns()
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactorSession) SetL1RedTokenAddress(l1RedTokenAddress common.Address) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.SetL1RedTokenAddress(&_ParentBridgeCoreFacet.TransactOpts, l1RedTokenAddress)
}

// UnpauseBridge is a paid mutator transaction binding the contract method 0xa82f143c.
//
// Solidity: function unpauseBridge() returns(bool)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactor) UnpauseBridge(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.contract.Transact(opts, "unpauseBridge")
}

// UnpauseBridge is a paid mutator transaction binding the contract method 0xa82f143c.
//
// Solidity: function unpauseBridge() returns(bool)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetSession) UnpauseBridge() (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.UnpauseBridge(&_ParentBridgeCoreFacet.TransactOpts)
}

// UnpauseBridge is a paid mutator transaction binding the contract method 0xa82f143c.
//
// Solidity: function unpauseBridge() returns(bool)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactorSession) UnpauseBridge() (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.UnpauseBridge(&_ParentBridgeCoreFacet.TransactOpts)
}

// ParentBridgeCoreFacetDownwardMessageIterator is returned from FilterDownwardMessage and is used to iterate over the raw logs and unpacked data for DownwardMessage events raised by the ParentBridgeCoreFacet contract.
type ParentBridgeCoreFacetDownwardMessageIterator struct {
	Event *ParentBridgeCoreFacetDownwardMessage // Event containing the contract specifics and raw log

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
func (it *ParentBridgeCoreFacetDownwardMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ParentBridgeCoreFacetDownwardMessage)
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
		it.Event = new(ParentBridgeCoreFacetDownwardMessage)
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
func (it *ParentBridgeCoreFacetDownwardMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ParentBridgeCoreFacetDownwardMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ParentBridgeCoreFacetDownwardMessage represents a DownwardMessage event raised by the ParentBridgeCoreFacet contract.
type ParentBridgeCoreFacetDownwardMessage struct {
	PayloadType uint32
	Payload     []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDownwardMessage is a free log retrieval operation binding the contract event 0xd52319e1be2e700973de64f36eae406db3c6106f7d0961b897ce537d23373532.
//
// Solidity: event DownwardMessage(uint32 payloadType, bytes payload)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) FilterDownwardMessage(opts *bind.FilterOpts) (*ParentBridgeCoreFacetDownwardMessageIterator, error) {

	logs, sub, err := _ParentBridgeCoreFacet.contract.FilterLogs(opts, "DownwardMessage")
	if err != nil {
		return nil, err
	}
	return &ParentBridgeCoreFacetDownwardMessageIterator{contract: _ParentBridgeCoreFacet.contract, event: "DownwardMessage", logs: logs, sub: sub}, nil
}

// WatchDownwardMessage is a free log subscription operation binding the contract event 0xd52319e1be2e700973de64f36eae406db3c6106f7d0961b897ce537d23373532.
//
// Solidity: event DownwardMessage(uint32 payloadType, bytes payload)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) WatchDownwardMessage(opts *bind.WatchOpts, sink chan<- *ParentBridgeCoreFacetDownwardMessage) (event.Subscription, error) {

	logs, sub, err := _ParentBridgeCoreFacet.contract.WatchLogs(opts, "DownwardMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ParentBridgeCoreFacetDownwardMessage)
				if err := _ParentBridgeCoreFacet.contract.UnpackLog(event, "DownwardMessage", log); err != nil {
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

// ParseDownwardMessage is a log parse operation binding the contract event 0xd52319e1be2e700973de64f36eae406db3c6106f7d0961b897ce537d23373532.
//
// Solidity: event DownwardMessage(uint32 payloadType, bytes payload)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) ParseDownwardMessage(log types.Log) (*ParentBridgeCoreFacetDownwardMessage, error) {
	event := new(ParentBridgeCoreFacetDownwardMessage)
	if err := _ParentBridgeCoreFacet.contract.UnpackLog(event, "DownwardMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ParentBridgeCoreFacetQueueTransactionIterator is returned from FilterQueueTransaction and is used to iterate over the raw logs and unpacked data for QueueTransaction events raised by the ParentBridgeCoreFacet contract.
type ParentBridgeCoreFacetQueueTransactionIterator struct {
	Event *ParentBridgeCoreFacetQueueTransaction // Event containing the contract specifics and raw log

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
func (it *ParentBridgeCoreFacetQueueTransactionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ParentBridgeCoreFacetQueueTransaction)
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
		it.Event = new(ParentBridgeCoreFacetQueueTransaction)
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
func (it *ParentBridgeCoreFacetQueueTransactionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ParentBridgeCoreFacetQueueTransactionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ParentBridgeCoreFacetQueueTransaction represents a QueueTransaction event raised by the ParentBridgeCoreFacet contract.
type ParentBridgeCoreFacetQueueTransaction struct {
	Hash        [32]byte
	QueueIndex  uint64
	PayloadType uint32
	Payload     []byte
	GasLimit    *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterQueueTransaction is a free log retrieval operation binding the contract event 0xd7c6020d9703629cf2a1e3d85ec55cdc7e0b3516cffc7382a3c96f63a62e1970.
//
// Solidity: event QueueTransaction(bytes32 indexed hash, uint64 indexed queueIndex, uint32 payloadType, bytes payload, uint256 gasLimit)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) FilterQueueTransaction(opts *bind.FilterOpts, hash [][32]byte, queueIndex []uint64) (*ParentBridgeCoreFacetQueueTransactionIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}
	var queueIndexRule []interface{}
	for _, queueIndexItem := range queueIndex {
		queueIndexRule = append(queueIndexRule, queueIndexItem)
	}

	logs, sub, err := _ParentBridgeCoreFacet.contract.FilterLogs(opts, "QueueTransaction", hashRule, queueIndexRule)
	if err != nil {
		return nil, err
	}
	return &ParentBridgeCoreFacetQueueTransactionIterator{contract: _ParentBridgeCoreFacet.contract, event: "QueueTransaction", logs: logs, sub: sub}, nil
}

// WatchQueueTransaction is a free log subscription operation binding the contract event 0xd7c6020d9703629cf2a1e3d85ec55cdc7e0b3516cffc7382a3c96f63a62e1970.
//
// Solidity: event QueueTransaction(bytes32 indexed hash, uint64 indexed queueIndex, uint32 payloadType, bytes payload, uint256 gasLimit)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) WatchQueueTransaction(opts *bind.WatchOpts, sink chan<- *ParentBridgeCoreFacetQueueTransaction, hash [][32]byte, queueIndex []uint64) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}
	var queueIndexRule []interface{}
	for _, queueIndexItem := range queueIndex {
		queueIndexRule = append(queueIndexRule, queueIndexItem)
	}

	logs, sub, err := _ParentBridgeCoreFacet.contract.WatchLogs(opts, "QueueTransaction", hashRule, queueIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ParentBridgeCoreFacetQueueTransaction)
				if err := _ParentBridgeCoreFacet.contract.UnpackLog(event, "QueueTransaction", log); err != nil {
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

// ParseQueueTransaction is a log parse operation binding the contract event 0xd7c6020d9703629cf2a1e3d85ec55cdc7e0b3516cffc7382a3c96f63a62e1970.
//
// Solidity: event QueueTransaction(bytes32 indexed hash, uint64 indexed queueIndex, uint32 payloadType, bytes payload, uint256 gasLimit)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) ParseQueueTransaction(log types.Log) (*ParentBridgeCoreFacetQueueTransaction, error) {
	event := new(ParentBridgeCoreFacetQueueTransaction)
	if err := _ParentBridgeCoreFacet.contract.UnpackLog(event, "QueueTransaction", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ParentBridgeCoreFacetRelayedMessageIterator is returned from FilterRelayedMessage and is used to iterate over the raw logs and unpacked data for RelayedMessage events raised by the ParentBridgeCoreFacet contract.
type ParentBridgeCoreFacetRelayedMessageIterator struct {
	Event *ParentBridgeCoreFacetRelayedMessage // Event containing the contract specifics and raw log

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
func (it *ParentBridgeCoreFacetRelayedMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ParentBridgeCoreFacetRelayedMessage)
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
		it.Event = new(ParentBridgeCoreFacetRelayedMessage)
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
func (it *ParentBridgeCoreFacetRelayedMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ParentBridgeCoreFacetRelayedMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ParentBridgeCoreFacetRelayedMessage represents a RelayedMessage event raised by the ParentBridgeCoreFacet contract.
type ParentBridgeCoreFacetRelayedMessage struct {
	MessageHash [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRelayedMessage is a free log retrieval operation binding the contract event 0x4641df4a962071e12719d8c8c8e5ac7fc4d97b927346a3d7a335b1f7517e133c.
//
// Solidity: event RelayedMessage(bytes32 indexed messageHash)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) FilterRelayedMessage(opts *bind.FilterOpts, messageHash [][32]byte) (*ParentBridgeCoreFacetRelayedMessageIterator, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _ParentBridgeCoreFacet.contract.FilterLogs(opts, "RelayedMessage", messageHashRule)
	if err != nil {
		return nil, err
	}
	return &ParentBridgeCoreFacetRelayedMessageIterator{contract: _ParentBridgeCoreFacet.contract, event: "RelayedMessage", logs: logs, sub: sub}, nil
}

// WatchRelayedMessage is a free log subscription operation binding the contract event 0x4641df4a962071e12719d8c8c8e5ac7fc4d97b927346a3d7a335b1f7517e133c.
//
// Solidity: event RelayedMessage(bytes32 indexed messageHash)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) WatchRelayedMessage(opts *bind.WatchOpts, sink chan<- *ParentBridgeCoreFacetRelayedMessage, messageHash [][32]byte) (event.Subscription, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _ParentBridgeCoreFacet.contract.WatchLogs(opts, "RelayedMessage", messageHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ParentBridgeCoreFacetRelayedMessage)
				if err := _ParentBridgeCoreFacet.contract.UnpackLog(event, "RelayedMessage", log); err != nil {
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

// ParseRelayedMessage is a log parse operation binding the contract event 0x4641df4a962071e12719d8c8c8e5ac7fc4d97b927346a3d7a335b1f7517e133c.
//
// Solidity: event RelayedMessage(bytes32 indexed messageHash)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) ParseRelayedMessage(log types.Log) (*ParentBridgeCoreFacetRelayedMessage, error) {
	event := new(ParentBridgeCoreFacetRelayedMessage)
	if err := _ParentBridgeCoreFacet.contract.UnpackLog(event, "RelayedMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
