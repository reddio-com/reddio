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

// DownwardMessage is an auto generated low-level Go binding around an user-defined struct.
type DownwardMessage struct {
	PayloadType uint32
	Payload     []byte
	Nonce       *big.Int
}

// DownwardMessageDispatcherFacetMetaData contains all meta data concerning the DownwardMessageDispatcherFacet contract.
var DownwardMessageDispatcherFacetMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"RelayedMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"isL1MessageExecuted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structDownwardMessage[]\",\"name\":\"downwardMessages\",\"type\":\"tuple[]\"}],\"name\":\"receiveDownwardMessages\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"relayMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// IsL1MessageExecuted is a free data retrieval call binding the contract method 0x02345b50.
//
// Solidity: function isL1MessageExecuted(bytes32 hash) view returns(bool)
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetCaller) IsL1MessageExecuted(opts *bind.CallOpts, hash [32]byte) (bool, error) {
	var out []interface{}
	err := _DownwardMessageDispatcherFacet.contract.Call(opts, &out, "isL1MessageExecuted", hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsL1MessageExecuted is a free data retrieval call binding the contract method 0x02345b50.
//
// Solidity: function isL1MessageExecuted(bytes32 hash) view returns(bool)
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetSession) IsL1MessageExecuted(hash [32]byte) (bool, error) {
	return _DownwardMessageDispatcherFacet.Contract.IsL1MessageExecuted(&_DownwardMessageDispatcherFacet.CallOpts, hash)
}

// IsL1MessageExecuted is a free data retrieval call binding the contract method 0x02345b50.
//
// Solidity: function isL1MessageExecuted(bytes32 hash) view returns(bool)
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetCallerSession) IsL1MessageExecuted(hash [32]byte) (bool, error) {
	return _DownwardMessageDispatcherFacet.Contract.IsL1MessageExecuted(&_DownwardMessageDispatcherFacet.CallOpts, hash)
}

// ReceiveDownwardMessages is a paid mutator transaction binding the contract method 0x3f70ad6c.
//
// Solidity: function receiveDownwardMessages((uint32,bytes,uint256)[] downwardMessages) returns()
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetTransactor) ReceiveDownwardMessages(opts *bind.TransactOpts, downwardMessages []DownwardMessage) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.contract.Transact(opts, "receiveDownwardMessages", downwardMessages)
}

// ReceiveDownwardMessages is a paid mutator transaction binding the contract method 0x3f70ad6c.
//
// Solidity: function receiveDownwardMessages((uint32,bytes,uint256)[] downwardMessages) returns()
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetSession) ReceiveDownwardMessages(downwardMessages []DownwardMessage) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.ReceiveDownwardMessages(&_DownwardMessageDispatcherFacet.TransactOpts, downwardMessages)
}

// ReceiveDownwardMessages is a paid mutator transaction binding the contract method 0x3f70ad6c.
//
// Solidity: function receiveDownwardMessages((uint32,bytes,uint256)[] downwardMessages) returns()
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetTransactorSession) ReceiveDownwardMessages(downwardMessages []DownwardMessage) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.ReceiveDownwardMessages(&_DownwardMessageDispatcherFacet.TransactOpts, downwardMessages)
}

// RelayMessage is a paid mutator transaction binding the contract method 0xf49d811e.
//
// Solidity: function relayMessage(uint32 payloadType, bytes payload, uint256 nonce) returns()
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetTransactor) RelayMessage(opts *bind.TransactOpts, payloadType uint32, payload []byte, nonce *big.Int) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.contract.Transact(opts, "relayMessage", payloadType, payload, nonce)
}

// RelayMessage is a paid mutator transaction binding the contract method 0xf49d811e.
//
// Solidity: function relayMessage(uint32 payloadType, bytes payload, uint256 nonce) returns()
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetSession) RelayMessage(payloadType uint32, payload []byte, nonce *big.Int) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.RelayMessage(&_DownwardMessageDispatcherFacet.TransactOpts, payloadType, payload, nonce)
}

// RelayMessage is a paid mutator transaction binding the contract method 0xf49d811e.
//
// Solidity: function relayMessage(uint32 payloadType, bytes payload, uint256 nonce) returns()
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetTransactorSession) RelayMessage(payloadType uint32, payload []byte, nonce *big.Int) (*types.Transaction, error) {
	return _DownwardMessageDispatcherFacet.Contract.RelayMessage(&_DownwardMessageDispatcherFacet.TransactOpts, payloadType, payload, nonce)
}

// DownwardMessageDispatcherFacetRelayedMessageIterator is returned from FilterRelayedMessage and is used to iterate over the raw logs and unpacked data for RelayedMessage events raised by the DownwardMessageDispatcherFacet contract.
type DownwardMessageDispatcherFacetRelayedMessageIterator struct {
	Event *DownwardMessageDispatcherFacetRelayedMessage // Event containing the contract specifics and raw log

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
func (it *DownwardMessageDispatcherFacetRelayedMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DownwardMessageDispatcherFacetRelayedMessage)
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
		it.Event = new(DownwardMessageDispatcherFacetRelayedMessage)
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
func (it *DownwardMessageDispatcherFacetRelayedMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DownwardMessageDispatcherFacetRelayedMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DownwardMessageDispatcherFacetRelayedMessage represents a RelayedMessage event raised by the DownwardMessageDispatcherFacet contract.
type DownwardMessageDispatcherFacetRelayedMessage struct {
	MessageHash [32]byte
	PayloadType uint32
	Payload     []byte
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRelayedMessage is a free log retrieval operation binding the contract event 0x0b62c0b7f830f688170a35d8d74fac1ec72f7a47dacdda0c821e73629614ad08.
//
// Solidity: event RelayedMessage(bytes32 indexed messageHash, uint32 payloadType, bytes payload, uint256 nonce)
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetFilterer) FilterRelayedMessage(opts *bind.FilterOpts, messageHash [][32]byte) (*DownwardMessageDispatcherFacetRelayedMessageIterator, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _DownwardMessageDispatcherFacet.contract.FilterLogs(opts, "RelayedMessage", messageHashRule)
	if err != nil {
		return nil, err
	}
	return &DownwardMessageDispatcherFacetRelayedMessageIterator{contract: _DownwardMessageDispatcherFacet.contract, event: "RelayedMessage", logs: logs, sub: sub}, nil
}

// WatchRelayedMessage is a free log subscription operation binding the contract event 0x0b62c0b7f830f688170a35d8d74fac1ec72f7a47dacdda0c821e73629614ad08.
//
// Solidity: event RelayedMessage(bytes32 indexed messageHash, uint32 payloadType, bytes payload, uint256 nonce)
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetFilterer) WatchRelayedMessage(opts *bind.WatchOpts, sink chan<- *DownwardMessageDispatcherFacetRelayedMessage, messageHash [][32]byte) (event.Subscription, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _DownwardMessageDispatcherFacet.contract.WatchLogs(opts, "RelayedMessage", messageHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DownwardMessageDispatcherFacetRelayedMessage)
				if err := _DownwardMessageDispatcherFacet.contract.UnpackLog(event, "RelayedMessage", log); err != nil {
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

// ParseRelayedMessage is a log parse operation binding the contract event 0x0b62c0b7f830f688170a35d8d74fac1ec72f7a47dacdda0c821e73629614ad08.
//
// Solidity: event RelayedMessage(bytes32 indexed messageHash, uint32 payloadType, bytes payload, uint256 nonce)
func (_DownwardMessageDispatcherFacet *DownwardMessageDispatcherFacetFilterer) ParseRelayedMessage(log types.Log) (*DownwardMessageDispatcherFacetRelayedMessage, error) {
	event := new(DownwardMessageDispatcherFacetRelayedMessage)
	if err := _DownwardMessageDispatcherFacet.contract.UnpackLog(event, "RelayedMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
