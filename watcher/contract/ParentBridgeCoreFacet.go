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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"DownwardMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"sendDownwardMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// SendDownwardMessage is a paid mutator transaction binding the contract method 0xed5621b4.
//
// Solidity: function sendDownwardMessage(uint32 payloadType, bytes payload) returns()
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactor) SendDownwardMessage(opts *bind.TransactOpts, payloadType uint32, payload []byte) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.contract.Transact(opts, "sendDownwardMessage", payloadType, payload)
}

// SendDownwardMessage is a paid mutator transaction binding the contract method 0xed5621b4.
//
// Solidity: function sendDownwardMessage(uint32 payloadType, bytes payload) returns()
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetSession) SendDownwardMessage(payloadType uint32, payload []byte) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.SendDownwardMessage(&_ParentBridgeCoreFacet.TransactOpts, payloadType, payload)
}

// SendDownwardMessage is a paid mutator transaction binding the contract method 0xed5621b4.
//
// Solidity: function sendDownwardMessage(uint32 payloadType, bytes payload) returns()
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetTransactorSession) SendDownwardMessage(payloadType uint32, payload []byte) (*types.Transaction, error) {
	return _ParentBridgeCoreFacet.Contract.SendDownwardMessage(&_ParentBridgeCoreFacet.TransactOpts, payloadType, payload)
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
	Sequence    *big.Int
	PayloadType uint32
	Payload     []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDownwardMessage is a free log retrieval operation binding the contract event 0xff1cd53930bda8b9e19e21164146c24905d9fc6c928ef1f28c6e2f85c06ef095.
//
// Solidity: event DownwardMessage(uint256 indexed sequence, uint32 payloadType, bytes payload)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) FilterDownwardMessage(opts *bind.FilterOpts, sequence []*big.Int) (*ParentBridgeCoreFacetDownwardMessageIterator, error) {

	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ParentBridgeCoreFacet.contract.FilterLogs(opts, "DownwardMessage", sequenceRule)
	if err != nil {
		return nil, err
	}
	return &ParentBridgeCoreFacetDownwardMessageIterator{contract: _ParentBridgeCoreFacet.contract, event: "DownwardMessage", logs: logs, sub: sub}, nil
}

// WatchDownwardMessage is a free log subscription operation binding the contract event 0xff1cd53930bda8b9e19e21164146c24905d9fc6c928ef1f28c6e2f85c06ef095.
//
// Solidity: event DownwardMessage(uint256 indexed sequence, uint32 payloadType, bytes payload)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) WatchDownwardMessage(opts *bind.WatchOpts, sink chan<- *ParentBridgeCoreFacetDownwardMessage, sequence []*big.Int) (event.Subscription, error) {

	var sequenceRule []interface{}
	for _, sequenceItem := range sequence {
		sequenceRule = append(sequenceRule, sequenceItem)
	}

	logs, sub, err := _ParentBridgeCoreFacet.contract.WatchLogs(opts, "DownwardMessage", sequenceRule)
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

// ParseDownwardMessage is a log parse operation binding the contract event 0xff1cd53930bda8b9e19e21164146c24905d9fc6c928ef1f28c6e2f85c06ef095.
//
// Solidity: event DownwardMessage(uint256 indexed sequence, uint32 payloadType, bytes payload)
func (_ParentBridgeCoreFacet *ParentBridgeCoreFacetFilterer) ParseDownwardMessage(log types.Log) (*ParentBridgeCoreFacetDownwardMessage, error) {
	event := new(ParentBridgeCoreFacetDownwardMessage)
	if err := _ParentBridgeCoreFacet.contract.UnpackLog(event, "DownwardMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
