// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ticketex

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

// TicketexMetaData contains all meta data concerning the Ticketex contract.
var TicketexMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"minted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_tbaAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_ticketAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_ticketAddress\",\"type\":\"address\"}],\"name\":\"buyTicket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_ticketExtendedAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"ticketName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ticketSymbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ticketUri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"ticketAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ticketPrice\",\"type\":\"uint256\"}],\"name\":\"issueTicket\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ownedTicket\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newTicketAddress\",\"type\":\"address\"}],\"name\":\"updateTicketAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_ticketAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"updateWhiteList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TicketexABI is the input ABI used to generate the binding from.
// Deprecated: Use TicketexMetaData.ABI instead.
var TicketexABI = TicketexMetaData.ABI

// Ticketex is an auto generated Go binding around an Ethereum contract.
type Ticketex struct {
	TicketexCaller     // Read-only binding to the contract
	TicketexTransactor // Write-only binding to the contract
	TicketexFilterer   // Log filterer for contract events
}

// TicketexCaller is an auto generated read-only Go binding around an Ethereum contract.
type TicketexCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TicketexTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TicketexTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TicketexFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TicketexFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TicketexSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TicketexSession struct {
	Contract     *Ticketex         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TicketexCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TicketexCallerSession struct {
	Contract *TicketexCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TicketexTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TicketexTransactorSession struct {
	Contract     *TicketexTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TicketexRaw is an auto generated low-level Go binding around an Ethereum contract.
type TicketexRaw struct {
	Contract *Ticketex // Generic contract binding to access the raw methods on
}

// TicketexCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TicketexCallerRaw struct {
	Contract *TicketexCaller // Generic read-only contract binding to access the raw methods on
}

// TicketexTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TicketexTransactorRaw struct {
	Contract *TicketexTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTicketex creates a new instance of Ticketex, bound to a specific deployed contract.
func NewTicketex(address common.Address, backend bind.ContractBackend) (*Ticketex, error) {
	contract, err := bindTicketex(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ticketex{TicketexCaller: TicketexCaller{contract: contract}, TicketexTransactor: TicketexTransactor{contract: contract}, TicketexFilterer: TicketexFilterer{contract: contract}}, nil
}

// NewTicketexCaller creates a new read-only instance of Ticketex, bound to a specific deployed contract.
func NewTicketexCaller(address common.Address, caller bind.ContractCaller) (*TicketexCaller, error) {
	contract, err := bindTicketex(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TicketexCaller{contract: contract}, nil
}

// NewTicketexTransactor creates a new write-only instance of Ticketex, bound to a specific deployed contract.
func NewTicketexTransactor(address common.Address, transactor bind.ContractTransactor) (*TicketexTransactor, error) {
	contract, err := bindTicketex(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TicketexTransactor{contract: contract}, nil
}

// NewTicketexFilterer creates a new log filterer instance of Ticketex, bound to a specific deployed contract.
func NewTicketexFilterer(address common.Address, filterer bind.ContractFilterer) (*TicketexFilterer, error) {
	contract, err := bindTicketex(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TicketexFilterer{contract: contract}, nil
}

// bindTicketex binds a generic wrapper to an already deployed contract.
func bindTicketex(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TicketexMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ticketex *TicketexRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ticketex.Contract.TicketexCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ticketex *TicketexRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ticketex.Contract.TicketexTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ticketex *TicketexRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ticketex.Contract.TicketexTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ticketex *TicketexCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ticketex.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ticketex *TicketexTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ticketex.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ticketex *TicketexTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ticketex.Contract.contract.Transact(opts, method, params...)
}

// TbaAddress is a free data retrieval call binding the contract method 0xea4a693a.
//
// Solidity: function _tbaAddress(address ) view returns(address)
func (_Ticketex *TicketexCaller) TbaAddress(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _Ticketex.contract.Call(opts, &out, "_tbaAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TbaAddress is a free data retrieval call binding the contract method 0xea4a693a.
//
// Solidity: function _tbaAddress(address ) view returns(address)
func (_Ticketex *TicketexSession) TbaAddress(arg0 common.Address) (common.Address, error) {
	return _Ticketex.Contract.TbaAddress(&_Ticketex.CallOpts, arg0)
}

// TbaAddress is a free data retrieval call binding the contract method 0xea4a693a.
//
// Solidity: function _tbaAddress(address ) view returns(address)
func (_Ticketex *TicketexCallerSession) TbaAddress(arg0 common.Address) (common.Address, error) {
	return _Ticketex.Contract.TbaAddress(&_Ticketex.CallOpts, arg0)
}

// TicketAddresses is a free data retrieval call binding the contract method 0x8161b878.
//
// Solidity: function _ticketAddresses(address , uint256 ) view returns(address)
func (_Ticketex *TicketexCaller) TicketAddresses(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Ticketex.contract.Call(opts, &out, "_ticketAddresses", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TicketAddresses is a free data retrieval call binding the contract method 0x8161b878.
//
// Solidity: function _ticketAddresses(address , uint256 ) view returns(address)
func (_Ticketex *TicketexSession) TicketAddresses(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _Ticketex.Contract.TicketAddresses(&_Ticketex.CallOpts, arg0, arg1)
}

// TicketAddresses is a free data retrieval call binding the contract method 0x8161b878.
//
// Solidity: function _ticketAddresses(address , uint256 ) view returns(address)
func (_Ticketex *TicketexCallerSession) TicketAddresses(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _Ticketex.Contract.TicketAddresses(&_Ticketex.CallOpts, arg0, arg1)
}

// OwnedTicket is a free data retrieval call binding the contract method 0x5ce0e13c.
//
// Solidity: function ownedTicket() view returns(address[])
func (_Ticketex *TicketexCaller) OwnedTicket(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Ticketex.contract.Call(opts, &out, "ownedTicket")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// OwnedTicket is a free data retrieval call binding the contract method 0x5ce0e13c.
//
// Solidity: function ownedTicket() view returns(address[])
func (_Ticketex *TicketexSession) OwnedTicket() ([]common.Address, error) {
	return _Ticketex.Contract.OwnedTicket(&_Ticketex.CallOpts)
}

// OwnedTicket is a free data retrieval call binding the contract method 0x5ce0e13c.
//
// Solidity: function ownedTicket() view returns(address[])
func (_Ticketex *TicketexCallerSession) OwnedTicket() ([]common.Address, error) {
	return _Ticketex.Contract.OwnedTicket(&_Ticketex.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ticketex *TicketexCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ticketex.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ticketex *TicketexSession) Owner() (common.Address, error) {
	return _Ticketex.Contract.Owner(&_Ticketex.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ticketex *TicketexCallerSession) Owner() (common.Address, error) {
	return _Ticketex.Contract.Owner(&_Ticketex.CallOpts)
}

// BuyTicket is a paid mutator transaction binding the contract method 0xa5f8cdbb.
//
// Solidity: function buyTicket(address _ticketAddress) payable returns(uint256)
func (_Ticketex *TicketexTransactor) BuyTicket(opts *bind.TransactOpts, _ticketAddress common.Address) (*types.Transaction, error) {
	return _Ticketex.contract.Transact(opts, "buyTicket", _ticketAddress)
}

// BuyTicket is a paid mutator transaction binding the contract method 0xa5f8cdbb.
//
// Solidity: function buyTicket(address _ticketAddress) payable returns(uint256)
func (_Ticketex *TicketexSession) BuyTicket(_ticketAddress common.Address) (*types.Transaction, error) {
	return _Ticketex.Contract.BuyTicket(&_Ticketex.TransactOpts, _ticketAddress)
}

// BuyTicket is a paid mutator transaction binding the contract method 0xa5f8cdbb.
//
// Solidity: function buyTicket(address _ticketAddress) payable returns(uint256)
func (_Ticketex *TicketexTransactorSession) BuyTicket(_ticketAddress common.Address) (*types.Transaction, error) {
	return _Ticketex.Contract.BuyTicket(&_Ticketex.TransactOpts, _ticketAddress)
}

// IssueTicket is a paid mutator transaction binding the contract method 0x5b851fde.
//
// Solidity: function issueTicket(address _ticketExtendedAddress, address _tokenAddress, string ticketName, string ticketSymbol, string ticketUri, address initialOwner, uint256 ticketAmount, uint256 ticketPrice) returns(address)
func (_Ticketex *TicketexTransactor) IssueTicket(opts *bind.TransactOpts, _ticketExtendedAddress common.Address, _tokenAddress common.Address, ticketName string, ticketSymbol string, ticketUri string, initialOwner common.Address, ticketAmount *big.Int, ticketPrice *big.Int) (*types.Transaction, error) {
	return _Ticketex.contract.Transact(opts, "issueTicket", _ticketExtendedAddress, _tokenAddress, ticketName, ticketSymbol, ticketUri, initialOwner, ticketAmount, ticketPrice)
}

// IssueTicket is a paid mutator transaction binding the contract method 0x5b851fde.
//
// Solidity: function issueTicket(address _ticketExtendedAddress, address _tokenAddress, string ticketName, string ticketSymbol, string ticketUri, address initialOwner, uint256 ticketAmount, uint256 ticketPrice) returns(address)
func (_Ticketex *TicketexSession) IssueTicket(_ticketExtendedAddress common.Address, _tokenAddress common.Address, ticketName string, ticketSymbol string, ticketUri string, initialOwner common.Address, ticketAmount *big.Int, ticketPrice *big.Int) (*types.Transaction, error) {
	return _Ticketex.Contract.IssueTicket(&_Ticketex.TransactOpts, _ticketExtendedAddress, _tokenAddress, ticketName, ticketSymbol, ticketUri, initialOwner, ticketAmount, ticketPrice)
}

// IssueTicket is a paid mutator transaction binding the contract method 0x5b851fde.
//
// Solidity: function issueTicket(address _ticketExtendedAddress, address _tokenAddress, string ticketName, string ticketSymbol, string ticketUri, address initialOwner, uint256 ticketAmount, uint256 ticketPrice) returns(address)
func (_Ticketex *TicketexTransactorSession) IssueTicket(_ticketExtendedAddress common.Address, _tokenAddress common.Address, ticketName string, ticketSymbol string, ticketUri string, initialOwner common.Address, ticketAmount *big.Int, ticketPrice *big.Int) (*types.Transaction, error) {
	return _Ticketex.Contract.IssueTicket(&_Ticketex.TransactOpts, _ticketExtendedAddress, _tokenAddress, ticketName, ticketSymbol, ticketUri, initialOwner, ticketAmount, ticketPrice)
}

// Mint is a paid mutator transaction binding the contract method 0xd0def521.
//
// Solidity: function mint(address to, string tokenURI) payable returns(uint256, address)
func (_Ticketex *TicketexTransactor) Mint(opts *bind.TransactOpts, to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Ticketex.contract.Transact(opts, "mint", to, tokenURI)
}

// Mint is a paid mutator transaction binding the contract method 0xd0def521.
//
// Solidity: function mint(address to, string tokenURI) payable returns(uint256, address)
func (_Ticketex *TicketexSession) Mint(to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Ticketex.Contract.Mint(&_Ticketex.TransactOpts, to, tokenURI)
}

// Mint is a paid mutator transaction binding the contract method 0xd0def521.
//
// Solidity: function mint(address to, string tokenURI) payable returns(uint256, address)
func (_Ticketex *TicketexTransactorSession) Mint(to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Ticketex.Contract.Mint(&_Ticketex.TransactOpts, to, tokenURI)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ticketex *TicketexTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ticketex.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ticketex *TicketexSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ticketex.Contract.RenounceOwnership(&_Ticketex.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ticketex *TicketexTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ticketex.Contract.RenounceOwnership(&_Ticketex.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ticketex *TicketexTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ticketex.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ticketex *TicketexSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ticketex.Contract.TransferOwnership(&_Ticketex.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ticketex *TicketexTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ticketex.Contract.TransferOwnership(&_Ticketex.TransactOpts, newOwner)
}

// UpdateTicketAddresses is a paid mutator transaction binding the contract method 0x57dca1fa.
//
// Solidity: function updateTicketAddresses(address buyer, address newTicketAddress) returns()
func (_Ticketex *TicketexTransactor) UpdateTicketAddresses(opts *bind.TransactOpts, buyer common.Address, newTicketAddress common.Address) (*types.Transaction, error) {
	return _Ticketex.contract.Transact(opts, "updateTicketAddresses", buyer, newTicketAddress)
}

// UpdateTicketAddresses is a paid mutator transaction binding the contract method 0x57dca1fa.
//
// Solidity: function updateTicketAddresses(address buyer, address newTicketAddress) returns()
func (_Ticketex *TicketexSession) UpdateTicketAddresses(buyer common.Address, newTicketAddress common.Address) (*types.Transaction, error) {
	return _Ticketex.Contract.UpdateTicketAddresses(&_Ticketex.TransactOpts, buyer, newTicketAddress)
}

// UpdateTicketAddresses is a paid mutator transaction binding the contract method 0x57dca1fa.
//
// Solidity: function updateTicketAddresses(address buyer, address newTicketAddress) returns()
func (_Ticketex *TicketexTransactorSession) UpdateTicketAddresses(buyer common.Address, newTicketAddress common.Address) (*types.Transaction, error) {
	return _Ticketex.Contract.UpdateTicketAddresses(&_Ticketex.TransactOpts, buyer, newTicketAddress)
}

// UpdateWhiteList is a paid mutator transaction binding the contract method 0xe3bdcfd3.
//
// Solidity: function updateWhiteList(address _ticketAddress, address to) returns()
func (_Ticketex *TicketexTransactor) UpdateWhiteList(opts *bind.TransactOpts, _ticketAddress common.Address, to common.Address) (*types.Transaction, error) {
	return _Ticketex.contract.Transact(opts, "updateWhiteList", _ticketAddress, to)
}

// UpdateWhiteList is a paid mutator transaction binding the contract method 0xe3bdcfd3.
//
// Solidity: function updateWhiteList(address _ticketAddress, address to) returns()
func (_Ticketex *TicketexSession) UpdateWhiteList(_ticketAddress common.Address, to common.Address) (*types.Transaction, error) {
	return _Ticketex.Contract.UpdateWhiteList(&_Ticketex.TransactOpts, _ticketAddress, to)
}

// UpdateWhiteList is a paid mutator transaction binding the contract method 0xe3bdcfd3.
//
// Solidity: function updateWhiteList(address _ticketAddress, address to) returns()
func (_Ticketex *TicketexTransactorSession) UpdateWhiteList(_ticketAddress common.Address, to common.Address) (*types.Transaction, error) {
	return _Ticketex.Contract.UpdateWhiteList(&_Ticketex.TransactOpts, _ticketAddress, to)
}

// TicketexOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ticketex contract.
type TicketexOwnershipTransferredIterator struct {
	Event *TicketexOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TicketexOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TicketexOwnershipTransferred)
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
		it.Event = new(TicketexOwnershipTransferred)
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
func (it *TicketexOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TicketexOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TicketexOwnershipTransferred represents a OwnershipTransferred event raised by the Ticketex contract.
type TicketexOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ticketex *TicketexFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TicketexOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ticketex.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TicketexOwnershipTransferredIterator{contract: _Ticketex.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ticketex *TicketexFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TicketexOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ticketex.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TicketexOwnershipTransferred)
				if err := _Ticketex.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ticketex *TicketexFilterer) ParseOwnershipTransferred(log types.Log) (*TicketexOwnershipTransferred, error) {
	event := new(TicketexOwnershipTransferred)
	if err := _Ticketex.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TicketexMintedIterator is returned from FilterMinted and is used to iterate over the raw logs and unpacked data for Minted events raised by the Ticketex contract.
type TicketexMintedIterator struct {
	Event *TicketexMinted // Event containing the contract specifics and raw log

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
func (it *TicketexMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TicketexMinted)
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
		it.Event = new(TicketexMinted)
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
func (it *TicketexMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TicketexMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TicketexMinted represents a Minted event raised by the Ticketex contract.
type TicketexMinted struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinted is a free log retrieval operation binding the contract event 0x7dc0bf3ff15656545da2c5f0567962839fe379f74aacdfc4e8025bb24e0c082d.
//
// Solidity: event minted(uint256 tokenId)
func (_Ticketex *TicketexFilterer) FilterMinted(opts *bind.FilterOpts) (*TicketexMintedIterator, error) {

	logs, sub, err := _Ticketex.contract.FilterLogs(opts, "minted")
	if err != nil {
		return nil, err
	}
	return &TicketexMintedIterator{contract: _Ticketex.contract, event: "minted", logs: logs, sub: sub}, nil
}

// WatchMinted is a free log subscription operation binding the contract event 0x7dc0bf3ff15656545da2c5f0567962839fe379f74aacdfc4e8025bb24e0c082d.
//
// Solidity: event minted(uint256 tokenId)
func (_Ticketex *TicketexFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *TicketexMinted) (event.Subscription, error) {

	logs, sub, err := _Ticketex.contract.WatchLogs(opts, "minted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TicketexMinted)
				if err := _Ticketex.contract.UnpackLog(event, "minted", log); err != nil {
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

// ParseMinted is a log parse operation binding the contract event 0x7dc0bf3ff15656545da2c5f0567962839fe379f74aacdfc4e8025bb24e0c082d.
//
// Solidity: event minted(uint256 tokenId)
func (_Ticketex *TicketexFilterer) ParseMinted(log types.Log) (*TicketexMinted, error) {
	event := new(TicketexMinted)
	if err := _Ticketex.contract.UnpackLog(event, "minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
