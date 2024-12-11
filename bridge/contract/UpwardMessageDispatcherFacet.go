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
	PayloadType uint32
	Payload     []byte
	Nonce       *big.Int
}

// UpwardMessageDispatcherFacetMetaData contains all meta data concerning the UpwardMessageDispatcherFacet contract.
var UpwardMessageDispatcherFacetMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"RelayedMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"isL2MessageExecuted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structUpwardMessage[]\",\"name\":\"upwardMessages\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signaturesArray\",\"type\":\"bytes[]\"}],\"name\":\"receiveUpwardMessages\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"relayMessageWithProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// IsL2MessageExecuted is a free data retrieval call binding the contract method 0x088681a7.
//
// Solidity: function isL2MessageExecuted(bytes32 hash) view returns(bool)
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetCaller) IsL2MessageExecuted(opts *bind.CallOpts, hash [32]byte) (bool, error) {
	var out []interface{}
	err := _UpwardMessageDispatcherFacet.contract.Call(opts, &out, "isL2MessageExecuted", hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsL2MessageExecuted is a free data retrieval call binding the contract method 0x088681a7.
//
// Solidity: function isL2MessageExecuted(bytes32 hash) view returns(bool)
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetSession) IsL2MessageExecuted(hash [32]byte) (bool, error) {
	return _UpwardMessageDispatcherFacet.Contract.IsL2MessageExecuted(&_UpwardMessageDispatcherFacet.CallOpts, hash)
}

// IsL2MessageExecuted is a free data retrieval call binding the contract method 0x088681a7.
//
// Solidity: function isL2MessageExecuted(bytes32 hash) view returns(bool)
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetCallerSession) IsL2MessageExecuted(hash [32]byte) (bool, error) {
	return _UpwardMessageDispatcherFacet.Contract.IsL2MessageExecuted(&_UpwardMessageDispatcherFacet.CallOpts, hash)
}

// ReceiveUpwardMessages is a paid mutator transaction binding the contract method 0xb3ce8ad0.
//
// Solidity: function receiveUpwardMessages((uint32,bytes,uint256)[] upwardMessages, bytes[] signaturesArray) returns()
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetTransactor) ReceiveUpwardMessages(opts *bind.TransactOpts, upwardMessages []UpwardMessage, signaturesArray [][]byte) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.contract.Transact(opts, "receiveUpwardMessages", upwardMessages, signaturesArray)
}

// ReceiveUpwardMessages is a paid mutator transaction binding the contract method 0xb3ce8ad0.
//
// Solidity: function receiveUpwardMessages((uint32,bytes,uint256)[] upwardMessages, bytes[] signaturesArray) returns()
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetSession) ReceiveUpwardMessages(upwardMessages []UpwardMessage, signaturesArray [][]byte) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.ReceiveUpwardMessages(&_UpwardMessageDispatcherFacet.TransactOpts, upwardMessages, signaturesArray)
}

// ReceiveUpwardMessages is a paid mutator transaction binding the contract method 0xb3ce8ad0.
//
// Solidity: function receiveUpwardMessages((uint32,bytes,uint256)[] upwardMessages, bytes[] signaturesArray) returns()
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetTransactorSession) ReceiveUpwardMessages(upwardMessages []UpwardMessage, signaturesArray [][]byte) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.ReceiveUpwardMessages(&_UpwardMessageDispatcherFacet.TransactOpts, upwardMessages, signaturesArray)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x9b1c8de6.
//
// Solidity: function relayMessageWithProof(uint32 payloadType, bytes payload, uint256 nonce) returns()
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetTransactor) RelayMessageWithProof(opts *bind.TransactOpts, payloadType uint32, payload []byte, nonce *big.Int) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.contract.Transact(opts, "relayMessageWithProof", payloadType, payload, nonce)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x9b1c8de6.
//
// Solidity: function relayMessageWithProof(uint32 payloadType, bytes payload, uint256 nonce) returns()
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetSession) RelayMessageWithProof(payloadType uint32, payload []byte, nonce *big.Int) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.RelayMessageWithProof(&_UpwardMessageDispatcherFacet.TransactOpts, payloadType, payload, nonce)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x9b1c8de6.
//
// Solidity: function relayMessageWithProof(uint32 payloadType, bytes payload, uint256 nonce) returns()
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetTransactorSession) RelayMessageWithProof(payloadType uint32, payload []byte, nonce *big.Int) (*types.Transaction, error) {
	return _UpwardMessageDispatcherFacet.Contract.RelayMessageWithProof(&_UpwardMessageDispatcherFacet.TransactOpts, payloadType, payload, nonce)
}

// UpwardMessageDispatcherFacetRelayedMessageIterator is returned from FilterRelayedMessage and is used to iterate over the raw logs and unpacked data for RelayedMessage events raised by the UpwardMessageDispatcherFacet contract.
type UpwardMessageDispatcherFacetRelayedMessageIterator struct {
	Event *UpwardMessageDispatcherFacetRelayedMessage // Event containing the contract specifics and raw log

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
func (it *UpwardMessageDispatcherFacetRelayedMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpwardMessageDispatcherFacetRelayedMessage)
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
		it.Event = new(UpwardMessageDispatcherFacetRelayedMessage)
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
func (it *UpwardMessageDispatcherFacetRelayedMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpwardMessageDispatcherFacetRelayedMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpwardMessageDispatcherFacetRelayedMessage represents a RelayedMessage event raised by the UpwardMessageDispatcherFacet contract.
type UpwardMessageDispatcherFacetRelayedMessage struct {
	MessageHash [32]byte
	PayloadType uint32
	Payload     []byte
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRelayedMessage is a free log retrieval operation binding the contract event 0x0b62c0b7f830f688170a35d8d74fac1ec72f7a47dacdda0c821e73629614ad08.
//
// Solidity: event RelayedMessage(bytes32 indexed messageHash, uint32 payloadType, bytes payload, uint256 nonce)
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetFilterer) FilterRelayedMessage(opts *bind.FilterOpts, messageHash [][32]byte) (*UpwardMessageDispatcherFacetRelayedMessageIterator, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _UpwardMessageDispatcherFacet.contract.FilterLogs(opts, "RelayedMessage", messageHashRule)
	if err != nil {
		return nil, err
	}
	return &UpwardMessageDispatcherFacetRelayedMessageIterator{contract: _UpwardMessageDispatcherFacet.contract, event: "RelayedMessage", logs: logs, sub: sub}, nil
}

// WatchRelayedMessage is a free log subscription operation binding the contract event 0x0b62c0b7f830f688170a35d8d74fac1ec72f7a47dacdda0c821e73629614ad08.
//
// Solidity: event RelayedMessage(bytes32 indexed messageHash, uint32 payloadType, bytes payload, uint256 nonce)
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetFilterer) WatchRelayedMessage(opts *bind.WatchOpts, sink chan<- *UpwardMessageDispatcherFacetRelayedMessage, messageHash [][32]byte) (event.Subscription, error) {

	var messageHashRule []interface{}
	for _, messageHashItem := range messageHash {
		messageHashRule = append(messageHashRule, messageHashItem)
	}

	logs, sub, err := _UpwardMessageDispatcherFacet.contract.WatchLogs(opts, "RelayedMessage", messageHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpwardMessageDispatcherFacetRelayedMessage)
				if err := _UpwardMessageDispatcherFacet.contract.UnpackLog(event, "RelayedMessage", log); err != nil {
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
func (_UpwardMessageDispatcherFacet *UpwardMessageDispatcherFacetFilterer) ParseRelayedMessage(log types.Log) (*UpwardMessageDispatcherFacetRelayedMessage, error) {
	event := new(UpwardMessageDispatcherFacetRelayedMessage)
	if err := _UpwardMessageDispatcherFacet.contract.UnpackLog(event, "RelayedMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
