// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package heroticket

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

// HeroticketMetaData contains all meta data concerning the Heroticket contract.
var HeroticketMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"minted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_tbaAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_ticketAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_ticketAddress\",\"type\":\"address\"}],\"name\":\"buyTicket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_ticketExtendedAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"ticketName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ticketSymbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ticketUri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"ticketAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ticketPrice\",\"type\":\"uint256\"}],\"name\":\"issueTicket\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ownedTicket\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newTicketAddress\",\"type\":\"address\"}],\"name\":\"updateTicketAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_ticketAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"updateWhiteList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// HeroticketABI is the input ABI used to generate the binding from.
// Deprecated: Use HeroticketMetaData.ABI instead.
var HeroticketABI = HeroticketMetaData.ABI

// Heroticket is an auto generated Go binding around an Ethereum contract.
type Heroticket struct {
	HeroticketCaller     // Read-only binding to the contract
	HeroticketTransactor // Write-only binding to the contract
	HeroticketFilterer   // Log filterer for contract events
}

// HeroticketCaller is an auto generated read-only Go binding around an Ethereum contract.
type HeroticketCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HeroticketTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HeroticketTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HeroticketFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HeroticketFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HeroticketSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HeroticketSession struct {
	Contract     *Heroticket       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HeroticketCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HeroticketCallerSession struct {
	Contract *HeroticketCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// HeroticketTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HeroticketTransactorSession struct {
	Contract     *HeroticketTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// HeroticketRaw is an auto generated low-level Go binding around an Ethereum contract.
type HeroticketRaw struct {
	Contract *Heroticket // Generic contract binding to access the raw methods on
}

// HeroticketCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HeroticketCallerRaw struct {
	Contract *HeroticketCaller // Generic read-only contract binding to access the raw methods on
}

// HeroticketTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HeroticketTransactorRaw struct {
	Contract *HeroticketTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHeroticket creates a new instance of Heroticket, bound to a specific deployed contract.
func NewHeroticket(address common.Address, backend bind.ContractBackend) (*Heroticket, error) {
	contract, err := bindHeroticket(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Heroticket{HeroticketCaller: HeroticketCaller{contract: contract}, HeroticketTransactor: HeroticketTransactor{contract: contract}, HeroticketFilterer: HeroticketFilterer{contract: contract}}, nil
}

// NewHeroticketCaller creates a new read-only instance of Heroticket, bound to a specific deployed contract.
func NewHeroticketCaller(address common.Address, caller bind.ContractCaller) (*HeroticketCaller, error) {
	contract, err := bindHeroticket(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HeroticketCaller{contract: contract}, nil
}

// NewHeroticketTransactor creates a new write-only instance of Heroticket, bound to a specific deployed contract.
func NewHeroticketTransactor(address common.Address, transactor bind.ContractTransactor) (*HeroticketTransactor, error) {
	contract, err := bindHeroticket(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HeroticketTransactor{contract: contract}, nil
}

// NewHeroticketFilterer creates a new log filterer instance of Heroticket, bound to a specific deployed contract.
func NewHeroticketFilterer(address common.Address, filterer bind.ContractFilterer) (*HeroticketFilterer, error) {
	contract, err := bindHeroticket(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HeroticketFilterer{contract: contract}, nil
}

// bindHeroticket binds a generic wrapper to an already deployed contract.
func bindHeroticket(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HeroticketMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Heroticket *HeroticketRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Heroticket.Contract.HeroticketCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Heroticket *HeroticketRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Heroticket.Contract.HeroticketTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Heroticket *HeroticketRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Heroticket.Contract.HeroticketTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Heroticket *HeroticketCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Heroticket.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Heroticket *HeroticketTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Heroticket.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Heroticket *HeroticketTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Heroticket.Contract.contract.Transact(opts, method, params...)
}

// TbaAddress is a free data retrieval call binding the contract method 0xea4a693a.
//
// Solidity: function _tbaAddress(address ) view returns(address)
func (_Heroticket *HeroticketCaller) TbaAddress(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _Heroticket.contract.Call(opts, &out, "_tbaAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TbaAddress is a free data retrieval call binding the contract method 0xea4a693a.
//
// Solidity: function _tbaAddress(address ) view returns(address)
func (_Heroticket *HeroticketSession) TbaAddress(arg0 common.Address) (common.Address, error) {
	return _Heroticket.Contract.TbaAddress(&_Heroticket.CallOpts, arg0)
}

// TbaAddress is a free data retrieval call binding the contract method 0xea4a693a.
//
// Solidity: function _tbaAddress(address ) view returns(address)
func (_Heroticket *HeroticketCallerSession) TbaAddress(arg0 common.Address) (common.Address, error) {
	return _Heroticket.Contract.TbaAddress(&_Heroticket.CallOpts, arg0)
}

// TicketAddresses is a free data retrieval call binding the contract method 0x8161b878.
//
// Solidity: function _ticketAddresses(address , uint256 ) view returns(address)
func (_Heroticket *HeroticketCaller) TicketAddresses(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Heroticket.contract.Call(opts, &out, "_ticketAddresses", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TicketAddresses is a free data retrieval call binding the contract method 0x8161b878.
//
// Solidity: function _ticketAddresses(address , uint256 ) view returns(address)
func (_Heroticket *HeroticketSession) TicketAddresses(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _Heroticket.Contract.TicketAddresses(&_Heroticket.CallOpts, arg0, arg1)
}

// TicketAddresses is a free data retrieval call binding the contract method 0x8161b878.
//
// Solidity: function _ticketAddresses(address , uint256 ) view returns(address)
func (_Heroticket *HeroticketCallerSession) TicketAddresses(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _Heroticket.Contract.TicketAddresses(&_Heroticket.CallOpts, arg0, arg1)
}

// OwnedTicket is a free data retrieval call binding the contract method 0x5ce0e13c.
//
// Solidity: function ownedTicket() view returns(address[])
func (_Heroticket *HeroticketCaller) OwnedTicket(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Heroticket.contract.Call(opts, &out, "ownedTicket")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// OwnedTicket is a free data retrieval call binding the contract method 0x5ce0e13c.
//
// Solidity: function ownedTicket() view returns(address[])
func (_Heroticket *HeroticketSession) OwnedTicket() ([]common.Address, error) {
	return _Heroticket.Contract.OwnedTicket(&_Heroticket.CallOpts)
}

// OwnedTicket is a free data retrieval call binding the contract method 0x5ce0e13c.
//
// Solidity: function ownedTicket() view returns(address[])
func (_Heroticket *HeroticketCallerSession) OwnedTicket() ([]common.Address, error) {
	return _Heroticket.Contract.OwnedTicket(&_Heroticket.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Heroticket *HeroticketCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Heroticket.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Heroticket *HeroticketSession) Owner() (common.Address, error) {
	return _Heroticket.Contract.Owner(&_Heroticket.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Heroticket *HeroticketCallerSession) Owner() (common.Address, error) {
	return _Heroticket.Contract.Owner(&_Heroticket.CallOpts)
}

// BuyTicket is a paid mutator transaction binding the contract method 0xa5f8cdbb.
//
// Solidity: function buyTicket(address _ticketAddress) payable returns(uint256)
func (_Heroticket *HeroticketTransactor) BuyTicket(opts *bind.TransactOpts, _ticketAddress common.Address) (*types.Transaction, error) {
	return _Heroticket.contract.Transact(opts, "buyTicket", _ticketAddress)
}

// BuyTicket is a paid mutator transaction binding the contract method 0xa5f8cdbb.
//
// Solidity: function buyTicket(address _ticketAddress) payable returns(uint256)
func (_Heroticket *HeroticketSession) BuyTicket(_ticketAddress common.Address) (*types.Transaction, error) {
	return _Heroticket.Contract.BuyTicket(&_Heroticket.TransactOpts, _ticketAddress)
}

// BuyTicket is a paid mutator transaction binding the contract method 0xa5f8cdbb.
//
// Solidity: function buyTicket(address _ticketAddress) payable returns(uint256)
func (_Heroticket *HeroticketTransactorSession) BuyTicket(_ticketAddress common.Address) (*types.Transaction, error) {
	return _Heroticket.Contract.BuyTicket(&_Heroticket.TransactOpts, _ticketAddress)
}

// IssueTicket is a paid mutator transaction binding the contract method 0x5b851fde.
//
// Solidity: function issueTicket(address _ticketExtendedAddress, address _tokenAddress, string ticketName, string ticketSymbol, string ticketUri, address initialOwner, uint256 ticketAmount, uint256 ticketPrice) returns(address)
func (_Heroticket *HeroticketTransactor) IssueTicket(opts *bind.TransactOpts, _ticketExtendedAddress common.Address, _tokenAddress common.Address, ticketName string, ticketSymbol string, ticketUri string, initialOwner common.Address, ticketAmount *big.Int, ticketPrice *big.Int) (*types.Transaction, error) {
	return _Heroticket.contract.Transact(opts, "issueTicket", _ticketExtendedAddress, _tokenAddress, ticketName, ticketSymbol, ticketUri, initialOwner, ticketAmount, ticketPrice)
}

// IssueTicket is a paid mutator transaction binding the contract method 0x5b851fde.
//
// Solidity: function issueTicket(address _ticketExtendedAddress, address _tokenAddress, string ticketName, string ticketSymbol, string ticketUri, address initialOwner, uint256 ticketAmount, uint256 ticketPrice) returns(address)
func (_Heroticket *HeroticketSession) IssueTicket(_ticketExtendedAddress common.Address, _tokenAddress common.Address, ticketName string, ticketSymbol string, ticketUri string, initialOwner common.Address, ticketAmount *big.Int, ticketPrice *big.Int) (*types.Transaction, error) {
	return _Heroticket.Contract.IssueTicket(&_Heroticket.TransactOpts, _ticketExtendedAddress, _tokenAddress, ticketName, ticketSymbol, ticketUri, initialOwner, ticketAmount, ticketPrice)
}

// IssueTicket is a paid mutator transaction binding the contract method 0x5b851fde.
//
// Solidity: function issueTicket(address _ticketExtendedAddress, address _tokenAddress, string ticketName, string ticketSymbol, string ticketUri, address initialOwner, uint256 ticketAmount, uint256 ticketPrice) returns(address)
func (_Heroticket *HeroticketTransactorSession) IssueTicket(_ticketExtendedAddress common.Address, _tokenAddress common.Address, ticketName string, ticketSymbol string, ticketUri string, initialOwner common.Address, ticketAmount *big.Int, ticketPrice *big.Int) (*types.Transaction, error) {
	return _Heroticket.Contract.IssueTicket(&_Heroticket.TransactOpts, _ticketExtendedAddress, _tokenAddress, ticketName, ticketSymbol, ticketUri, initialOwner, ticketAmount, ticketPrice)
}

// Mint is a paid mutator transaction binding the contract method 0xd0def521.
//
// Solidity: function mint(address to, string tokenURI) payable returns(uint256, address)
func (_Heroticket *HeroticketTransactor) Mint(opts *bind.TransactOpts, to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Heroticket.contract.Transact(opts, "mint", to, tokenURI)
}

// Mint is a paid mutator transaction binding the contract method 0xd0def521.
//
// Solidity: function mint(address to, string tokenURI) payable returns(uint256, address)
func (_Heroticket *HeroticketSession) Mint(to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Heroticket.Contract.Mint(&_Heroticket.TransactOpts, to, tokenURI)
}

// Mint is a paid mutator transaction binding the contract method 0xd0def521.
//
// Solidity: function mint(address to, string tokenURI) payable returns(uint256, address)
func (_Heroticket *HeroticketTransactorSession) Mint(to common.Address, tokenURI string) (*types.Transaction, error) {
	return _Heroticket.Contract.Mint(&_Heroticket.TransactOpts, to, tokenURI)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Heroticket *HeroticketTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Heroticket.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Heroticket *HeroticketSession) RenounceOwnership() (*types.Transaction, error) {
	return _Heroticket.Contract.RenounceOwnership(&_Heroticket.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Heroticket *HeroticketTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Heroticket.Contract.RenounceOwnership(&_Heroticket.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Heroticket *HeroticketTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Heroticket.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Heroticket *HeroticketSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Heroticket.Contract.TransferOwnership(&_Heroticket.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Heroticket *HeroticketTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Heroticket.Contract.TransferOwnership(&_Heroticket.TransactOpts, newOwner)
}

// UpdateTicketAddresses is a paid mutator transaction binding the contract method 0x57dca1fa.
//
// Solidity: function updateTicketAddresses(address buyer, address newTicketAddress) returns()
func (_Heroticket *HeroticketTransactor) UpdateTicketAddresses(opts *bind.TransactOpts, buyer common.Address, newTicketAddress common.Address) (*types.Transaction, error) {
	return _Heroticket.contract.Transact(opts, "updateTicketAddresses", buyer, newTicketAddress)
}

// UpdateTicketAddresses is a paid mutator transaction binding the contract method 0x57dca1fa.
//
// Solidity: function updateTicketAddresses(address buyer, address newTicketAddress) returns()
func (_Heroticket *HeroticketSession) UpdateTicketAddresses(buyer common.Address, newTicketAddress common.Address) (*types.Transaction, error) {
	return _Heroticket.Contract.UpdateTicketAddresses(&_Heroticket.TransactOpts, buyer, newTicketAddress)
}

// UpdateTicketAddresses is a paid mutator transaction binding the contract method 0x57dca1fa.
//
// Solidity: function updateTicketAddresses(address buyer, address newTicketAddress) returns()
func (_Heroticket *HeroticketTransactorSession) UpdateTicketAddresses(buyer common.Address, newTicketAddress common.Address) (*types.Transaction, error) {
	return _Heroticket.Contract.UpdateTicketAddresses(&_Heroticket.TransactOpts, buyer, newTicketAddress)
}

// UpdateWhiteList is a paid mutator transaction binding the contract method 0xe3bdcfd3.
//
// Solidity: function updateWhiteList(address _ticketAddress, address to) returns()
func (_Heroticket *HeroticketTransactor) UpdateWhiteList(opts *bind.TransactOpts, _ticketAddress common.Address, to common.Address) (*types.Transaction, error) {
	return _Heroticket.contract.Transact(opts, "updateWhiteList", _ticketAddress, to)
}

// UpdateWhiteList is a paid mutator transaction binding the contract method 0xe3bdcfd3.
//
// Solidity: function updateWhiteList(address _ticketAddress, address to) returns()
func (_Heroticket *HeroticketSession) UpdateWhiteList(_ticketAddress common.Address, to common.Address) (*types.Transaction, error) {
	return _Heroticket.Contract.UpdateWhiteList(&_Heroticket.TransactOpts, _ticketAddress, to)
}

// UpdateWhiteList is a paid mutator transaction binding the contract method 0xe3bdcfd3.
//
// Solidity: function updateWhiteList(address _ticketAddress, address to) returns()
func (_Heroticket *HeroticketTransactorSession) UpdateWhiteList(_ticketAddress common.Address, to common.Address) (*types.Transaction, error) {
	return _Heroticket.Contract.UpdateWhiteList(&_Heroticket.TransactOpts, _ticketAddress, to)
}

// HeroticketOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Heroticket contract.
type HeroticketOwnershipTransferredIterator struct {
	Event *HeroticketOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *HeroticketOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HeroticketOwnershipTransferred)
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
		it.Event = new(HeroticketOwnershipTransferred)
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
func (it *HeroticketOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HeroticketOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HeroticketOwnershipTransferred represents a OwnershipTransferred event raised by the Heroticket contract.
type HeroticketOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Heroticket *HeroticketFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*HeroticketOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Heroticket.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &HeroticketOwnershipTransferredIterator{contract: _Heroticket.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Heroticket *HeroticketFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *HeroticketOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Heroticket.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HeroticketOwnershipTransferred)
				if err := _Heroticket.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Heroticket *HeroticketFilterer) ParseOwnershipTransferred(log types.Log) (*HeroticketOwnershipTransferred, error) {
	event := new(HeroticketOwnershipTransferred)
	if err := _Heroticket.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HeroticketMintedIterator is returned from FilterMinted and is used to iterate over the raw logs and unpacked data for Minted events raised by the Heroticket contract.
type HeroticketMintedIterator struct {
	Event *HeroticketMinted // Event containing the contract specifics and raw log

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
func (it *HeroticketMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HeroticketMinted)
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
		it.Event = new(HeroticketMinted)
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
func (it *HeroticketMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HeroticketMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HeroticketMinted represents a Minted event raised by the Heroticket contract.
type HeroticketMinted struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinted is a free log retrieval operation binding the contract event 0x7dc0bf3ff15656545da2c5f0567962839fe379f74aacdfc4e8025bb24e0c082d.
//
// Solidity: event minted(uint256 tokenId)
func (_Heroticket *HeroticketFilterer) FilterMinted(opts *bind.FilterOpts) (*HeroticketMintedIterator, error) {

	logs, sub, err := _Heroticket.contract.FilterLogs(opts, "minted")
	if err != nil {
		return nil, err
	}
	return &HeroticketMintedIterator{contract: _Heroticket.contract, event: "minted", logs: logs, sub: sub}, nil
}

// WatchMinted is a free log subscription operation binding the contract event 0x7dc0bf3ff15656545da2c5f0567962839fe379f74aacdfc4e8025bb24e0c082d.
//
// Solidity: event minted(uint256 tokenId)
func (_Heroticket *HeroticketFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *HeroticketMinted) (event.Subscription, error) {

	logs, sub, err := _Heroticket.contract.WatchLogs(opts, "minted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HeroticketMinted)
				if err := _Heroticket.contract.UnpackLog(event, "minted", log); err != nil {
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
func (_Heroticket *HeroticketFilterer) ParseMinted(log types.Log) (*HeroticketMinted, error) {
	event := new(HeroticketMinted)
	if err := _Heroticket.contract.UnpackLog(event, "minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
