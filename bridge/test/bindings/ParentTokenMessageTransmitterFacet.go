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

// ParentTokenMessageTransmitterFacetMetaData contains all meta data concerning the ParentTokenMessageTransmitterFacet contract.
var ParentTokenMessageTransmitterFacetMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"depositERC1155Token\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"depositERC20Token\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"depositERC721Token\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"depositETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"depositRED\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC1155BatchReceived\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC1155Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ParentTokenMessageTransmitterFacetABI is the input ABI used to generate the binding from.
// Deprecated: Use ParentTokenMessageTransmitterFacetMetaData.ABI instead.
var ParentTokenMessageTransmitterFacetABI = ParentTokenMessageTransmitterFacetMetaData.ABI

// ParentTokenMessageTransmitterFacet is an auto generated Go binding around an Ethereum contract.
type ParentTokenMessageTransmitterFacet struct {
	ParentTokenMessageTransmitterFacetCaller     // Read-only binding to the contract
	ParentTokenMessageTransmitterFacetTransactor // Write-only binding to the contract
	ParentTokenMessageTransmitterFacetFilterer   // Log filterer for contract events
}

// ParentTokenMessageTransmitterFacetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ParentTokenMessageTransmitterFacetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ParentTokenMessageTransmitterFacetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ParentTokenMessageTransmitterFacetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ParentTokenMessageTransmitterFacetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ParentTokenMessageTransmitterFacetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ParentTokenMessageTransmitterFacetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ParentTokenMessageTransmitterFacetSession struct {
	Contract     *ParentTokenMessageTransmitterFacet // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                       // Call options to use throughout this session
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// ParentTokenMessageTransmitterFacetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ParentTokenMessageTransmitterFacetCallerSession struct {
	Contract *ParentTokenMessageTransmitterFacetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                             // Call options to use throughout this session
}

// ParentTokenMessageTransmitterFacetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ParentTokenMessageTransmitterFacetTransactorSession struct {
	Contract     *ParentTokenMessageTransmitterFacetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                             // Transaction auth options to use throughout this session
}

// ParentTokenMessageTransmitterFacetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ParentTokenMessageTransmitterFacetRaw struct {
	Contract *ParentTokenMessageTransmitterFacet // Generic contract binding to access the raw methods on
}

// ParentTokenMessageTransmitterFacetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ParentTokenMessageTransmitterFacetCallerRaw struct {
	Contract *ParentTokenMessageTransmitterFacetCaller // Generic read-only contract binding to access the raw methods on
}

// ParentTokenMessageTransmitterFacetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ParentTokenMessageTransmitterFacetTransactorRaw struct {
	Contract *ParentTokenMessageTransmitterFacetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewParentTokenMessageTransmitterFacet creates a new instance of ParentTokenMessageTransmitterFacet, bound to a specific deployed contract.
func NewParentTokenMessageTransmitterFacet(address common.Address, backend bind.ContractBackend) (*ParentTokenMessageTransmitterFacet, error) {
	contract, err := bindParentTokenMessageTransmitterFacet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ParentTokenMessageTransmitterFacet{ParentTokenMessageTransmitterFacetCaller: ParentTokenMessageTransmitterFacetCaller{contract: contract}, ParentTokenMessageTransmitterFacetTransactor: ParentTokenMessageTransmitterFacetTransactor{contract: contract}, ParentTokenMessageTransmitterFacetFilterer: ParentTokenMessageTransmitterFacetFilterer{contract: contract}}, nil
}

// NewParentTokenMessageTransmitterFacetCaller creates a new read-only instance of ParentTokenMessageTransmitterFacet, bound to a specific deployed contract.
func NewParentTokenMessageTransmitterFacetCaller(address common.Address, caller bind.ContractCaller) (*ParentTokenMessageTransmitterFacetCaller, error) {
	contract, err := bindParentTokenMessageTransmitterFacet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ParentTokenMessageTransmitterFacetCaller{contract: contract}, nil
}

// NewParentTokenMessageTransmitterFacetTransactor creates a new write-only instance of ParentTokenMessageTransmitterFacet, bound to a specific deployed contract.
func NewParentTokenMessageTransmitterFacetTransactor(address common.Address, transactor bind.ContractTransactor) (*ParentTokenMessageTransmitterFacetTransactor, error) {
	contract, err := bindParentTokenMessageTransmitterFacet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ParentTokenMessageTransmitterFacetTransactor{contract: contract}, nil
}

// NewParentTokenMessageTransmitterFacetFilterer creates a new log filterer instance of ParentTokenMessageTransmitterFacet, bound to a specific deployed contract.
func NewParentTokenMessageTransmitterFacetFilterer(address common.Address, filterer bind.ContractFilterer) (*ParentTokenMessageTransmitterFacetFilterer, error) {
	contract, err := bindParentTokenMessageTransmitterFacet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ParentTokenMessageTransmitterFacetFilterer{contract: contract}, nil
}

// bindParentTokenMessageTransmitterFacet binds a generic wrapper to an already deployed contract.
func bindParentTokenMessageTransmitterFacet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ParentTokenMessageTransmitterFacetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ParentTokenMessageTransmitterFacet.Contract.ParentTokenMessageTransmitterFacetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.ParentTokenMessageTransmitterFacetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.ParentTokenMessageTransmitterFacetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ParentTokenMessageTransmitterFacet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.contract.Transact(opts, method, params...)
}

// DepositERC1155Token is a paid mutator transaction binding the contract method 0xdeb5fcdc.
//
// Solidity: function depositERC1155Token(address tokenAddress, address recipient, uint256[] tokenIds, uint256[] amounts, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactor) DepositERC1155Token(opts *bind.TransactOpts, tokenAddress common.Address, recipient common.Address, tokenIds []*big.Int, amounts []*big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.contract.Transact(opts, "depositERC1155Token", tokenAddress, recipient, tokenIds, amounts, gasLimit)
}

// DepositERC1155Token is a paid mutator transaction binding the contract method 0xdeb5fcdc.
//
// Solidity: function depositERC1155Token(address tokenAddress, address recipient, uint256[] tokenIds, uint256[] amounts, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetSession) DepositERC1155Token(tokenAddress common.Address, recipient common.Address, tokenIds []*big.Int, amounts []*big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositERC1155Token(&_ParentTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, tokenIds, amounts, gasLimit)
}

// DepositERC1155Token is a paid mutator transaction binding the contract method 0xdeb5fcdc.
//
// Solidity: function depositERC1155Token(address tokenAddress, address recipient, uint256[] tokenIds, uint256[] amounts, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorSession) DepositERC1155Token(tokenAddress common.Address, recipient common.Address, tokenIds []*big.Int, amounts []*big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositERC1155Token(&_ParentTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, tokenIds, amounts, gasLimit)
}

// DepositERC20Token is a paid mutator transaction binding the contract method 0x156929f9.
//
// Solidity: function depositERC20Token(address tokenAddress, address recipient, uint256 amount, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactor) DepositERC20Token(opts *bind.TransactOpts, tokenAddress common.Address, recipient common.Address, amount *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.contract.Transact(opts, "depositERC20Token", tokenAddress, recipient, amount, gasLimit)
}

// DepositERC20Token is a paid mutator transaction binding the contract method 0x156929f9.
//
// Solidity: function depositERC20Token(address tokenAddress, address recipient, uint256 amount, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetSession) DepositERC20Token(tokenAddress common.Address, recipient common.Address, amount *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositERC20Token(&_ParentTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, amount, gasLimit)
}

// DepositERC20Token is a paid mutator transaction binding the contract method 0x156929f9.
//
// Solidity: function depositERC20Token(address tokenAddress, address recipient, uint256 amount, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorSession) DepositERC20Token(tokenAddress common.Address, recipient common.Address, amount *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositERC20Token(&_ParentTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, amount, gasLimit)
}

// DepositERC721Token is a paid mutator transaction binding the contract method 0x6c28888c.
//
// Solidity: function depositERC721Token(address tokenAddress, address recipient, uint256 tokenId, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactor) DepositERC721Token(opts *bind.TransactOpts, tokenAddress common.Address, recipient common.Address, tokenId *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.contract.Transact(opts, "depositERC721Token", tokenAddress, recipient, tokenId, gasLimit)
}

// DepositERC721Token is a paid mutator transaction binding the contract method 0x6c28888c.
//
// Solidity: function depositERC721Token(address tokenAddress, address recipient, uint256 tokenId, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetSession) DepositERC721Token(tokenAddress common.Address, recipient common.Address, tokenId *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositERC721Token(&_ParentTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, tokenId, gasLimit)
}

// DepositERC721Token is a paid mutator transaction binding the contract method 0x6c28888c.
//
// Solidity: function depositERC721Token(address tokenAddress, address recipient, uint256 tokenId, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorSession) DepositERC721Token(tokenAddress common.Address, recipient common.Address, tokenId *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositERC721Token(&_ParentTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, tokenId, gasLimit)
}

// DepositETH is a paid mutator transaction binding the contract method 0xce0b63ce.
//
// Solidity: function depositETH(address recipient, uint256 amount, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactor) DepositETH(opts *bind.TransactOpts, recipient common.Address, amount *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.contract.Transact(opts, "depositETH", recipient, amount, gasLimit)
}

// DepositETH is a paid mutator transaction binding the contract method 0xce0b63ce.
//
// Solidity: function depositETH(address recipient, uint256 amount, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetSession) DepositETH(recipient common.Address, amount *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositETH(&_ParentTokenMessageTransmitterFacet.TransactOpts, recipient, amount, gasLimit)
}

// DepositETH is a paid mutator transaction binding the contract method 0xce0b63ce.
//
// Solidity: function depositETH(address recipient, uint256 amount, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorSession) DepositETH(recipient common.Address, amount *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositETH(&_ParentTokenMessageTransmitterFacet.TransactOpts, recipient, amount, gasLimit)
}

// DepositRED is a paid mutator transaction binding the contract method 0x06e2aa65.
//
// Solidity: function depositRED(address recipient, uint256 amount, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactor) DepositRED(opts *bind.TransactOpts, recipient common.Address, amount *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.contract.Transact(opts, "depositRED", recipient, amount, gasLimit)
}

// DepositRED is a paid mutator transaction binding the contract method 0x06e2aa65.
//
// Solidity: function depositRED(address recipient, uint256 amount, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetSession) DepositRED(recipient common.Address, amount *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositRED(&_ParentTokenMessageTransmitterFacet.TransactOpts, recipient, amount, gasLimit)
}

// DepositRED is a paid mutator transaction binding the contract method 0x06e2aa65.
//
// Solidity: function depositRED(address recipient, uint256 amount, uint256 gasLimit) payable returns()
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorSession) DepositRED(recipient common.Address, amount *big.Int, gasLimit *big.Int) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.DepositRED(&_ParentTokenMessageTransmitterFacet.TransactOpts, recipient, amount, gasLimit)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactor) OnERC1155BatchReceived(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.contract.Transact(opts, "onERC1155BatchReceived", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.OnERC1155BatchReceived(&_ParentTokenMessageTransmitterFacet.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.OnERC1155BatchReceived(&_ParentTokenMessageTransmitterFacet.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactor) OnERC1155Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.contract.Transact(opts, "onERC1155Received", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.OnERC1155Received(&_ParentTokenMessageTransmitterFacet.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.OnERC1155Received(&_ParentTokenMessageTransmitterFacet.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactor) OnERC721Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.contract.Transact(opts, "onERC721Received", arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.OnERC721Received(&_ParentTokenMessageTransmitterFacet.TransactOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_ParentTokenMessageTransmitterFacet *ParentTokenMessageTransmitterFacetTransactorSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _ParentTokenMessageTransmitterFacet.Contract.OnERC721Received(&_ParentTokenMessageTransmitterFacet.TransactOpts, arg0, arg1, arg2, arg3)
}
