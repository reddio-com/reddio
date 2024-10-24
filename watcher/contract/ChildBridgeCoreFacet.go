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

// ChildBridgeCoreFacetMetaData contains all meta data concerning the ChildBridgeCoreFacet contract.
var ChildBridgeCoreFacetMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"UpwardMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"sendUpwardMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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
	Sequence    *big.Int
	PayloadType uint32
	Payload     []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpwardMessage is a free log retrieval operation binding the contract event 0x1038383885a257d0e7b06035d5fde6388257f3d08e3db8ea01584b362768f4f4.
//
// Solidity: event UpwardMessage(uint256 indexed sequence, uint32 payloadType, bytes payload)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) FilterUpwardMessage(opts *bind.FilterOpts, sequence []*big.Int) (*ChildBridgeCoreFacetUpwardMessageIterator, error) {

	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ChildBridgeCoreFacet.contract.FilterLogs(opts, "UpwardMessage", sequenceRule)
	if err != nil {
		return nil, err
	}
	return &ChildBridgeCoreFacetUpwardMessageIterator{contract: _ChildBridgeCoreFacet.contract, event: "UpwardMessage", logs: logs, sub: sub}, nil
}

// WatchUpwardMessage is a free log subscription operation binding the contract event 0x1038383885a257d0e7b06035d5fde6388257f3d08e3db8ea01584b362768f4f4.
//
// Solidity: event UpwardMessage(uint256 indexed sequence, uint32 payloadType, bytes payload)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) WatchUpwardMessage(opts *bind.WatchOpts, sink chan<- *ChildBridgeCoreFacetUpwardMessage, sequence []*big.Int) (event.Subscription, error) {

	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ChildBridgeCoreFacet.contract.WatchLogs(opts, "UpwardMessage", sequenceRule)
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

// ParseUpwardMessage is a log parse operation binding the contract event 0x1038383885a257d0e7b06035d5fde6388257f3d08e3db8ea01584b362768f4f4.
//
// Solidity: event UpwardMessage(uint256 indexed sequence, uint32 payloadType, bytes payload)
func (_ChildBridgeCoreFacet *ChildBridgeCoreFacetFilterer) ParseUpwardMessage(log types.Log) (*ChildBridgeCoreFacetUpwardMessage, error) {
	event := new(ChildBridgeCoreFacetUpwardMessage)
	if err := _ChildBridgeCoreFacet.contract.UnpackLog(event, "UpwardMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
