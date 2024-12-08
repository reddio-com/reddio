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

// ChildTokenMessageTransmitterFacetMetaData contains all meta data concerning the ChildTokenMessageTransmitterFacet contract.
var ChildTokenMessageTransmitterFacetMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"withdrawErc1155BatchToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawErc20Token\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"withdrawErc721Token\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawRED\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// ChildTokenMessageTransmitterFacetABI is the input ABI used to generate the binding from.
// Deprecated: Use ChildTokenMessageTransmitterFacetMetaData.ABI instead.
var ChildTokenMessageTransmitterFacetABI = ChildTokenMessageTransmitterFacetMetaData.ABI

// ChildTokenMessageTransmitterFacet is an auto generated Go binding around an Ethereum contract.
type ChildTokenMessageTransmitterFacet struct {
	ChildTokenMessageTransmitterFacetCaller     // Read-only binding to the contract
	ChildTokenMessageTransmitterFacetTransactor // Write-only binding to the contract
	ChildTokenMessageTransmitterFacetFilterer   // Log filterer for contract events
}

// ChildTokenMessageTransmitterFacetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChildTokenMessageTransmitterFacetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChildTokenMessageTransmitterFacetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChildTokenMessageTransmitterFacetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChildTokenMessageTransmitterFacetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChildTokenMessageTransmitterFacetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChildTokenMessageTransmitterFacetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChildTokenMessageTransmitterFacetSession struct {
	Contract     *ChildTokenMessageTransmitterFacet // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                      // Call options to use throughout this session
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// ChildTokenMessageTransmitterFacetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChildTokenMessageTransmitterFacetCallerSession struct {
	Contract *ChildTokenMessageTransmitterFacetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                            // Call options to use throughout this session
}

// ChildTokenMessageTransmitterFacetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChildTokenMessageTransmitterFacetTransactorSession struct {
	Contract     *ChildTokenMessageTransmitterFacetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                            // Transaction auth options to use throughout this session
}

// ChildTokenMessageTransmitterFacetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChildTokenMessageTransmitterFacetRaw struct {
	Contract *ChildTokenMessageTransmitterFacet // Generic contract binding to access the raw methods on
}

// ChildTokenMessageTransmitterFacetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChildTokenMessageTransmitterFacetCallerRaw struct {
	Contract *ChildTokenMessageTransmitterFacetCaller // Generic read-only contract binding to access the raw methods on
}

// ChildTokenMessageTransmitterFacetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChildTokenMessageTransmitterFacetTransactorRaw struct {
	Contract *ChildTokenMessageTransmitterFacetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChildTokenMessageTransmitterFacet creates a new instance of ChildTokenMessageTransmitterFacet, bound to a specific deployed contract.
func NewChildTokenMessageTransmitterFacet(address common.Address, backend bind.ContractBackend) (*ChildTokenMessageTransmitterFacet, error) {
	contract, err := bindChildTokenMessageTransmitterFacet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChildTokenMessageTransmitterFacet{ChildTokenMessageTransmitterFacetCaller: ChildTokenMessageTransmitterFacetCaller{contract: contract}, ChildTokenMessageTransmitterFacetTransactor: ChildTokenMessageTransmitterFacetTransactor{contract: contract}, ChildTokenMessageTransmitterFacetFilterer: ChildTokenMessageTransmitterFacetFilterer{contract: contract}}, nil
}

// NewChildTokenMessageTransmitterFacetCaller creates a new read-only instance of ChildTokenMessageTransmitterFacet, bound to a specific deployed contract.
func NewChildTokenMessageTransmitterFacetCaller(address common.Address, caller bind.ContractCaller) (*ChildTokenMessageTransmitterFacetCaller, error) {
	contract, err := bindChildTokenMessageTransmitterFacet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChildTokenMessageTransmitterFacetCaller{contract: contract}, nil
}

// NewChildTokenMessageTransmitterFacetTransactor creates a new write-only instance of ChildTokenMessageTransmitterFacet, bound to a specific deployed contract.
func NewChildTokenMessageTransmitterFacetTransactor(address common.Address, transactor bind.ContractTransactor) (*ChildTokenMessageTransmitterFacetTransactor, error) {
	contract, err := bindChildTokenMessageTransmitterFacet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChildTokenMessageTransmitterFacetTransactor{contract: contract}, nil
}

// NewChildTokenMessageTransmitterFacetFilterer creates a new log filterer instance of ChildTokenMessageTransmitterFacet, bound to a specific deployed contract.
func NewChildTokenMessageTransmitterFacetFilterer(address common.Address, filterer bind.ContractFilterer) (*ChildTokenMessageTransmitterFacetFilterer, error) {
	contract, err := bindChildTokenMessageTransmitterFacet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChildTokenMessageTransmitterFacetFilterer{contract: contract}, nil
}

// bindChildTokenMessageTransmitterFacet binds a generic wrapper to an already deployed contract.
func bindChildTokenMessageTransmitterFacet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChildTokenMessageTransmitterFacetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChildTokenMessageTransmitterFacet.Contract.ChildTokenMessageTransmitterFacetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.ChildTokenMessageTransmitterFacetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.ChildTokenMessageTransmitterFacetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChildTokenMessageTransmitterFacet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.contract.Transact(opts, method, params...)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x4782f779.
//
// Solidity: function withdrawETH(address recipient, uint256 amount) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactor) WithdrawETH(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.contract.Transact(opts, "withdrawETH", recipient, amount)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x4782f779.
//
// Solidity: function withdrawETH(address recipient, uint256 amount) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetSession) WithdrawETH(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawETH(&_ChildTokenMessageTransmitterFacet.TransactOpts, recipient, amount)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0x4782f779.
//
// Solidity: function withdrawETH(address recipient, uint256 amount) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactorSession) WithdrawETH(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawETH(&_ChildTokenMessageTransmitterFacet.TransactOpts, recipient, amount)
}

// WithdrawErc1155BatchToken is a paid mutator transaction binding the contract method 0xcd35c3ce.
//
// Solidity: function withdrawErc1155BatchToken(address tokenAddress, address recipient, uint256[] tokenIds, uint256[] amounts) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactor) WithdrawErc1155BatchToken(opts *bind.TransactOpts, tokenAddress common.Address, recipient common.Address, tokenIds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.contract.Transact(opts, "withdrawErc1155BatchToken", tokenAddress, recipient, tokenIds, amounts)
}

// WithdrawErc1155BatchToken is a paid mutator transaction binding the contract method 0xcd35c3ce.
//
// Solidity: function withdrawErc1155BatchToken(address tokenAddress, address recipient, uint256[] tokenIds, uint256[] amounts) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetSession) WithdrawErc1155BatchToken(tokenAddress common.Address, recipient common.Address, tokenIds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawErc1155BatchToken(&_ChildTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, tokenIds, amounts)
}

// WithdrawErc1155BatchToken is a paid mutator transaction binding the contract method 0xcd35c3ce.
//
// Solidity: function withdrawErc1155BatchToken(address tokenAddress, address recipient, uint256[] tokenIds, uint256[] amounts) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactorSession) WithdrawErc1155BatchToken(tokenAddress common.Address, recipient common.Address, tokenIds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawErc1155BatchToken(&_ChildTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, tokenIds, amounts)
}

// WithdrawErc20Token is a paid mutator transaction binding the contract method 0x2d079734.
//
// Solidity: function withdrawErc20Token(address tokenAddress, address recipient, uint256 amount) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactor) WithdrawErc20Token(opts *bind.TransactOpts, tokenAddress common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.contract.Transact(opts, "withdrawErc20Token", tokenAddress, recipient, amount)
}

// WithdrawErc20Token is a paid mutator transaction binding the contract method 0x2d079734.
//
// Solidity: function withdrawErc20Token(address tokenAddress, address recipient, uint256 amount) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetSession) WithdrawErc20Token(tokenAddress common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawErc20Token(&_ChildTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, amount)
}

// WithdrawErc20Token is a paid mutator transaction binding the contract method 0x2d079734.
//
// Solidity: function withdrawErc20Token(address tokenAddress, address recipient, uint256 amount) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactorSession) WithdrawErc20Token(tokenAddress common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawErc20Token(&_ChildTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, amount)
}

// WithdrawErc721Token is a paid mutator transaction binding the contract method 0xa6f83281.
//
// Solidity: function withdrawErc721Token(address tokenAddress, address recipient, uint256 tokenId) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactor) WithdrawErc721Token(opts *bind.TransactOpts, tokenAddress common.Address, recipient common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.contract.Transact(opts, "withdrawErc721Token", tokenAddress, recipient, tokenId)
}

// WithdrawErc721Token is a paid mutator transaction binding the contract method 0xa6f83281.
//
// Solidity: function withdrawErc721Token(address tokenAddress, address recipient, uint256 tokenId) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetSession) WithdrawErc721Token(tokenAddress common.Address, recipient common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawErc721Token(&_ChildTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, tokenId)
}

// WithdrawErc721Token is a paid mutator transaction binding the contract method 0xa6f83281.
//
// Solidity: function withdrawErc721Token(address tokenAddress, address recipient, uint256 tokenId) returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactorSession) WithdrawErc721Token(tokenAddress common.Address, recipient common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawErc721Token(&_ChildTokenMessageTransmitterFacet.TransactOpts, tokenAddress, recipient, tokenId)
}

// WithdrawRED is a paid mutator transaction binding the contract method 0x71f5f10f.
//
// Solidity: function withdrawRED(address recipient) payable returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactor) WithdrawRED(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.contract.Transact(opts, "withdrawRED", recipient)
}

// WithdrawRED is a paid mutator transaction binding the contract method 0x71f5f10f.
//
// Solidity: function withdrawRED(address recipient) payable returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetSession) WithdrawRED(recipient common.Address) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawRED(&_ChildTokenMessageTransmitterFacet.TransactOpts, recipient)
}

// WithdrawRED is a paid mutator transaction binding the contract method 0x71f5f10f.
//
// Solidity: function withdrawRED(address recipient) payable returns()
func (_ChildTokenMessageTransmitterFacet *ChildTokenMessageTransmitterFacetTransactorSession) WithdrawRED(recipient common.Address) (*types.Transaction, error) {
	return _ChildTokenMessageTransmitterFacet.Contract.WithdrawRED(&_ChildTokenMessageTransmitterFacet.TransactOpts, recipient)
}
