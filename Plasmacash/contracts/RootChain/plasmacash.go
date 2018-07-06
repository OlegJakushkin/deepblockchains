/** Copyright 2018 Wolk Inc.
* This file is part of the Plasmacash library.
*
* The plasmacash library is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The Plasmacash library is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
*
* You should have received a copy of the GNU Lesser General Public License
* along with the plasmacash library. If not, see <http://www.gnu.org/licenses/>.
 */

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package plasmachain

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// PlasmaCashABI is the input ABI used to generate the binding from.
const PlasmaCashABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"depositIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getNextExit\",\"outputs\":[{\"name\":\"depID\",\"type\":\"uint64\"},{\"name\":\"tokenID\",\"type\":\"uint64\"},{\"name\":\"exitableTS\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentDepositIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_depositIndex\",\"type\":\"uint64\"}],\"name\":\"depositExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokenID\",\"type\":\"uint64\"},{\"name\":\"txBytes1\",\"type\":\"bytes\"},{\"name\":\"txBytes2\",\"type\":\"bytes\"},{\"name\":\"proof1\",\"type\":\"bytes\"},{\"name\":\"proof2\",\"type\":\"bytes\"},{\"name\":\"blk1\",\"type\":\"uint64\"},{\"name\":\"blk2\",\"type\":\"uint64\"}],\"name\":\"startExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"childChain\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokenID\",\"type\":\"uint64\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"},{\"name\":\"blk\",\"type\":\"uint64\"}],\"name\":\"challenge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentBlkNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeExits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"exits\",\"outputs\":[{\"name\":\"txblk1\",\"type\":\"uint64\"},{\"name\":\"txblk2\",\"type\":\"uint64\"},{\"name\":\"exitor\",\"type\":\"address\"},{\"name\":\"exitableTS\",\"type\":\"uint256\"},{\"name\":\"bond\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_blkRoot\",\"type\":\"bytes32\"},{\"name\":\"_blknum\",\"type\":\"uint64\"}],\"name\":\"submitBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"depositBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_depositor\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_depositIndex\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"_denomination\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_tokenID\",\"type\":\"uint64\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_exiter\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_depositIndex\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"_denomination\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_tokenID\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"StartExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_rootHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"_blknum\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_currentDepositIndex\",\"type\":\"uint64\"}],\"name\":\"PublishedBlock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_exiter\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_depositIndex\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"_denomination\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_tokenID\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"FinalizedExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_challenger\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenID\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Challenge\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_priority\",\"type\":\"uint256\"}],\"name\":\"ExitStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_priority\",\"type\":\"uint256\"}],\"name\":\"ExitCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"depID\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"tokenID\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"exitableTS\",\"type\":\"uint256\"}],\"name\":\"CurrtentExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"exitableTS\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"cuurrentTS\",\"type\":\"uint256\"}],\"name\":\"ExitTime\",\"type\":\"event\"}]"

// PlasmaCash is an auto generated Go binding around an Ethereum contract.
type PlasmaCash struct {
	PlasmaCashCaller     // Read-only binding to the contract
	PlasmaCashTransactor // Write-only binding to the contract
	PlasmaCashFilterer   // Log filterer for contract events
}

// PlasmaCashCaller is an auto generated read-only Go binding around an Ethereum contract.
type PlasmaCashCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlasmaCashTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PlasmaCashTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlasmaCashFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PlasmaCashFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlasmaCashSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PlasmaCashSession struct {
	Contract     *PlasmaCash       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PlasmaCashCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PlasmaCashCallerSession struct {
	Contract *PlasmaCashCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// PlasmaCashTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PlasmaCashTransactorSession struct {
	Contract     *PlasmaCashTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// PlasmaCashRaw is an auto generated low-level Go binding around an Ethereum contract.
type PlasmaCashRaw struct {
	Contract *PlasmaCash // Generic contract binding to access the raw methods on
}

// PlasmaCashCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PlasmaCashCallerRaw struct {
	Contract *PlasmaCashCaller // Generic read-only contract binding to access the raw methods on
}

// PlasmaCashTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PlasmaCashTransactorRaw struct {
	Contract *PlasmaCashTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPlasmaCash creates a new instance of PlasmaCash, bound to a specific deployed contract.
func NewPlasmaCash(address common.Address, backend bind.ContractBackend) (*PlasmaCash, error) {
	contract, err := bindPlasmaCash(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PlasmaCash{PlasmaCashCaller: PlasmaCashCaller{contract: contract}, PlasmaCashTransactor: PlasmaCashTransactor{contract: contract}, PlasmaCashFilterer: PlasmaCashFilterer{contract: contract}}, nil
}

// NewPlasmaCashCaller creates a new read-only instance of PlasmaCash, bound to a specific deployed contract.
func NewPlasmaCashCaller(address common.Address, caller bind.ContractCaller) (*PlasmaCashCaller, error) {
	contract, err := bindPlasmaCash(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashCaller{contract: contract}, nil
}

// NewPlasmaCashTransactor creates a new write-only instance of PlasmaCash, bound to a specific deployed contract.
func NewPlasmaCashTransactor(address common.Address, transactor bind.ContractTransactor) (*PlasmaCashTransactor, error) {
	contract, err := bindPlasmaCash(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashTransactor{contract: contract}, nil
}

// NewPlasmaCashFilterer creates a new log filterer instance of PlasmaCash, bound to a specific deployed contract.
func NewPlasmaCashFilterer(address common.Address, filterer bind.ContractFilterer) (*PlasmaCashFilterer, error) {
	contract, err := bindPlasmaCash(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashFilterer{contract: contract}, nil
}

// bindPlasmaCash binds a generic wrapper to an already deployed contract.
func bindPlasmaCash(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PlasmaCashABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PlasmaCash *PlasmaCashRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PlasmaCash.Contract.PlasmaCashCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PlasmaCash *PlasmaCashRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaCash.Contract.PlasmaCashTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PlasmaCash *PlasmaCashRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PlasmaCash.Contract.PlasmaCashTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PlasmaCash *PlasmaCashCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PlasmaCash.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PlasmaCash *PlasmaCashTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaCash.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PlasmaCash *PlasmaCashTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PlasmaCash.Contract.contract.Transact(opts, method, params...)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() constant returns(address)
func (_PlasmaCash *PlasmaCashCaller) Authority(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _PlasmaCash.contract.Call(opts, out, "authority")
	return *ret0, err
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() constant returns(address)
func (_PlasmaCash *PlasmaCashSession) Authority() (common.Address, error) {
	return _PlasmaCash.Contract.Authority(&_PlasmaCash.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() constant returns(address)
func (_PlasmaCash *PlasmaCashCallerSession) Authority() (common.Address, error) {
	return _PlasmaCash.Contract.Authority(&_PlasmaCash.CallOpts)
}

// ChildChain is a free data retrieval call binding the contract method 0x6434693a.
//
// Solidity: function childChain( uint64) constant returns(bytes32)
func (_PlasmaCash *PlasmaCashCaller) ChildChain(opts *bind.CallOpts, arg0 uint64) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _PlasmaCash.contract.Call(opts, out, "childChain", arg0)
	return *ret0, err
}

// ChildChain is a free data retrieval call binding the contract method 0x6434693a.
//
// Solidity: function childChain( uint64) constant returns(bytes32)
func (_PlasmaCash *PlasmaCashSession) ChildChain(arg0 uint64) ([32]byte, error) {
	return _PlasmaCash.Contract.ChildChain(&_PlasmaCash.CallOpts, arg0)
}

// ChildChain is a free data retrieval call binding the contract method 0x6434693a.
//
// Solidity: function childChain( uint64) constant returns(bytes32)
func (_PlasmaCash *PlasmaCashCallerSession) ChildChain(arg0 uint64) ([32]byte, error) {
	return _PlasmaCash.Contract.ChildChain(&_PlasmaCash.CallOpts, arg0)
}

// CurrentBlkNum is a free data retrieval call binding the contract method 0x8c1cdd5d.
//
// Solidity: function currentBlkNum() constant returns(uint64)
func (_PlasmaCash *PlasmaCashCaller) CurrentBlkNum(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _PlasmaCash.contract.Call(opts, out, "currentBlkNum")
	return *ret0, err
}

// CurrentBlkNum is a free data retrieval call binding the contract method 0x8c1cdd5d.
//
// Solidity: function currentBlkNum() constant returns(uint64)
func (_PlasmaCash *PlasmaCashSession) CurrentBlkNum() (uint64, error) {
	return _PlasmaCash.Contract.CurrentBlkNum(&_PlasmaCash.CallOpts)
}

// CurrentBlkNum is a free data retrieval call binding the contract method 0x8c1cdd5d.
//
// Solidity: function currentBlkNum() constant returns(uint64)
func (_PlasmaCash *PlasmaCashCallerSession) CurrentBlkNum() (uint64, error) {
	return _PlasmaCash.Contract.CurrentBlkNum(&_PlasmaCash.CallOpts)
}

// CurrentDepositIndex is a free data retrieval call binding the contract method 0x36e5df06.
//
// Solidity: function currentDepositIndex() constant returns(uint64)
func (_PlasmaCash *PlasmaCashCaller) CurrentDepositIndex(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _PlasmaCash.contract.Call(opts, out, "currentDepositIndex")
	return *ret0, err
}

// CurrentDepositIndex is a free data retrieval call binding the contract method 0x36e5df06.
//
// Solidity: function currentDepositIndex() constant returns(uint64)
func (_PlasmaCash *PlasmaCashSession) CurrentDepositIndex() (uint64, error) {
	return _PlasmaCash.Contract.CurrentDepositIndex(&_PlasmaCash.CallOpts)
}

// CurrentDepositIndex is a free data retrieval call binding the contract method 0x36e5df06.
//
// Solidity: function currentDepositIndex() constant returns(uint64)
func (_PlasmaCash *PlasmaCashCallerSession) CurrentDepositIndex() (uint64, error) {
	return _PlasmaCash.Contract.CurrentDepositIndex(&_PlasmaCash.CallOpts)
}

// DepositBalance is a free data retrieval call binding the contract method 0xf34272ab.
//
// Solidity: function depositBalance( uint64) constant returns(uint64)
func (_PlasmaCash *PlasmaCashCaller) DepositBalance(opts *bind.CallOpts, arg0 uint64) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _PlasmaCash.contract.Call(opts, out, "depositBalance", arg0)
	return *ret0, err
}

// DepositBalance is a free data retrieval call binding the contract method 0xf34272ab.
//
// Solidity: function depositBalance( uint64) constant returns(uint64)
func (_PlasmaCash *PlasmaCashSession) DepositBalance(arg0 uint64) (uint64, error) {
	return _PlasmaCash.Contract.DepositBalance(&_PlasmaCash.CallOpts, arg0)
}

// DepositBalance is a free data retrieval call binding the contract method 0xf34272ab.
//
// Solidity: function depositBalance( uint64) constant returns(uint64)
func (_PlasmaCash *PlasmaCashCallerSession) DepositBalance(arg0 uint64) (uint64, error) {
	return _PlasmaCash.Contract.DepositBalance(&_PlasmaCash.CallOpts, arg0)
}

// DepositIndex is a free data retrieval call binding the contract method 0x0f8592ac.
//
// Solidity: function depositIndex( uint64) constant returns(uint64)
func (_PlasmaCash *PlasmaCashCaller) DepositIndex(opts *bind.CallOpts, arg0 uint64) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _PlasmaCash.contract.Call(opts, out, "depositIndex", arg0)
	return *ret0, err
}

// DepositIndex is a free data retrieval call binding the contract method 0x0f8592ac.
//
// Solidity: function depositIndex( uint64) constant returns(uint64)
func (_PlasmaCash *PlasmaCashSession) DepositIndex(arg0 uint64) (uint64, error) {
	return _PlasmaCash.Contract.DepositIndex(&_PlasmaCash.CallOpts, arg0)
}

// DepositIndex is a free data retrieval call binding the contract method 0x0f8592ac.
//
// Solidity: function depositIndex( uint64) constant returns(uint64)
func (_PlasmaCash *PlasmaCashCallerSession) DepositIndex(arg0 uint64) (uint64, error) {
	return _PlasmaCash.Contract.DepositIndex(&_PlasmaCash.CallOpts, arg0)
}

// Exits is a free data retrieval call binding the contract method 0xd6463d40.
//
// Solidity: function exits( uint64) constant returns(txblk1 uint64, txblk2 uint64, exitor address, exitableTS uint256, bond uint256)
func (_PlasmaCash *PlasmaCashCaller) Exits(opts *bind.CallOpts, arg0 uint64) (struct {
	Txblk1     uint64
	Txblk2     uint64
	Exitor     common.Address
	ExitableTS *big.Int
	Bond       *big.Int
}, error) {
	ret := new(struct {
		Txblk1     uint64
		Txblk2     uint64
		Exitor     common.Address
		ExitableTS *big.Int
		Bond       *big.Int
	})
	out := ret
	err := _PlasmaCash.contract.Call(opts, out, "exits", arg0)
	return *ret, err
}

// Exits is a free data retrieval call binding the contract method 0xd6463d40.
//
// Solidity: function exits( uint64) constant returns(txblk1 uint64, txblk2 uint64, exitor address, exitableTS uint256, bond uint256)
func (_PlasmaCash *PlasmaCashSession) Exits(arg0 uint64) (struct {
	Txblk1     uint64
	Txblk2     uint64
	Exitor     common.Address
	ExitableTS *big.Int
	Bond       *big.Int
}, error) {
	return _PlasmaCash.Contract.Exits(&_PlasmaCash.CallOpts, arg0)
}

// Exits is a free data retrieval call binding the contract method 0xd6463d40.
//
// Solidity: function exits( uint64) constant returns(txblk1 uint64, txblk2 uint64, exitor address, exitableTS uint256, bond uint256)
func (_PlasmaCash *PlasmaCashCallerSession) Exits(arg0 uint64) (struct {
	Txblk1     uint64
	Txblk2     uint64
	Exitor     common.Address
	ExitableTS *big.Int
	Bond       *big.Int
}, error) {
	return _PlasmaCash.Contract.Exits(&_PlasmaCash.CallOpts, arg0)
}

// GetNextExit is a free data retrieval call binding the contract method 0x2fd56a25.
//
// Solidity: function getNextExit() constant returns(depID uint64, tokenID uint64, exitableTS uint256)
func (_PlasmaCash *PlasmaCashCaller) GetNextExit(opts *bind.CallOpts) (struct {
	DepID      uint64
	TokenID    uint64
	ExitableTS *big.Int
}, error) {
	ret := new(struct {
		DepID      uint64
		TokenID    uint64
		ExitableTS *big.Int
	})
	out := ret
	err := _PlasmaCash.contract.Call(opts, out, "getNextExit")
	return *ret, err
}

// GetNextExit is a free data retrieval call binding the contract method 0x2fd56a25.
//
// Solidity: function getNextExit() constant returns(depID uint64, tokenID uint64, exitableTS uint256)
func (_PlasmaCash *PlasmaCashSession) GetNextExit() (struct {
	DepID      uint64
	TokenID    uint64
	ExitableTS *big.Int
}, error) {
	return _PlasmaCash.Contract.GetNextExit(&_PlasmaCash.CallOpts)
}

// GetNextExit is a free data retrieval call binding the contract method 0x2fd56a25.
//
// Solidity: function getNextExit() constant returns(depID uint64, tokenID uint64, exitableTS uint256)
func (_PlasmaCash *PlasmaCashCallerSession) GetNextExit() (struct {
	DepID      uint64
	TokenID    uint64
	ExitableTS *big.Int
}, error) {
	return _PlasmaCash.Contract.GetNextExit(&_PlasmaCash.CallOpts)
}

// Challenge is a paid mutator transaction binding the contract method 0x6b68497e.
//
// Solidity: function challenge(tokenID uint64, txBytes bytes, proof bytes, blk uint64) returns()
func (_PlasmaCash *PlasmaCashTransactor) Challenge(opts *bind.TransactOpts, tokenID uint64, txBytes []byte, proof []byte, blk uint64) (*types.Transaction, error) {
	return _PlasmaCash.contract.Transact(opts, "challenge", tokenID, txBytes, proof, blk)
}

// Challenge is a paid mutator transaction binding the contract method 0x6b68497e.
//
// Solidity: function challenge(tokenID uint64, txBytes bytes, proof bytes, blk uint64) returns()
func (_PlasmaCash *PlasmaCashSession) Challenge(tokenID uint64, txBytes []byte, proof []byte, blk uint64) (*types.Transaction, error) {
	return _PlasmaCash.Contract.Challenge(&_PlasmaCash.TransactOpts, tokenID, txBytes, proof, blk)
}

// Challenge is a paid mutator transaction binding the contract method 0x6b68497e.
//
// Solidity: function challenge(tokenID uint64, txBytes bytes, proof bytes, blk uint64) returns()
func (_PlasmaCash *PlasmaCashTransactorSession) Challenge(tokenID uint64, txBytes []byte, proof []byte, blk uint64) (*types.Transaction, error) {
	return _PlasmaCash.Contract.Challenge(&_PlasmaCash.TransactOpts, tokenID, txBytes, proof, blk)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_PlasmaCash *PlasmaCashTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaCash.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_PlasmaCash *PlasmaCashSession) Deposit() (*types.Transaction, error) {
	return _PlasmaCash.Contract.Deposit(&_PlasmaCash.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_PlasmaCash *PlasmaCashTransactorSession) Deposit() (*types.Transaction, error) {
	return _PlasmaCash.Contract.Deposit(&_PlasmaCash.TransactOpts)
}

// DepositExit is a paid mutator transaction binding the contract method 0x3d2e2590.
//
// Solidity: function depositExit(_depositIndex uint64) returns()
func (_PlasmaCash *PlasmaCashTransactor) DepositExit(opts *bind.TransactOpts, _depositIndex uint64) (*types.Transaction, error) {
	return _PlasmaCash.contract.Transact(opts, "depositExit", _depositIndex)
}

// DepositExit is a paid mutator transaction binding the contract method 0x3d2e2590.
//
// Solidity: function depositExit(_depositIndex uint64) returns()
func (_PlasmaCash *PlasmaCashSession) DepositExit(_depositIndex uint64) (*types.Transaction, error) {
	return _PlasmaCash.Contract.DepositExit(&_PlasmaCash.TransactOpts, _depositIndex)
}

// DepositExit is a paid mutator transaction binding the contract method 0x3d2e2590.
//
// Solidity: function depositExit(_depositIndex uint64) returns()
func (_PlasmaCash *PlasmaCashTransactorSession) DepositExit(_depositIndex uint64) (*types.Transaction, error) {
	return _PlasmaCash.Contract.DepositExit(&_PlasmaCash.TransactOpts, _depositIndex)
}

// FinalizeExits is a paid mutator transaction binding the contract method 0xc6ab44cd.
//
// Solidity: function finalizeExits() returns()
func (_PlasmaCash *PlasmaCashTransactor) FinalizeExits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaCash.contract.Transact(opts, "finalizeExits")
}

// FinalizeExits is a paid mutator transaction binding the contract method 0xc6ab44cd.
//
// Solidity: function finalizeExits() returns()
func (_PlasmaCash *PlasmaCashSession) FinalizeExits() (*types.Transaction, error) {
	return _PlasmaCash.Contract.FinalizeExits(&_PlasmaCash.TransactOpts)
}

// FinalizeExits is a paid mutator transaction binding the contract method 0xc6ab44cd.
//
// Solidity: function finalizeExits() returns()
func (_PlasmaCash *PlasmaCashTransactorSession) FinalizeExits() (*types.Transaction, error) {
	return _PlasmaCash.Contract.FinalizeExits(&_PlasmaCash.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_PlasmaCash *PlasmaCashTransactor) Kill(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaCash.contract.Transact(opts, "kill")
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_PlasmaCash *PlasmaCashSession) Kill() (*types.Transaction, error) {
	return _PlasmaCash.Contract.Kill(&_PlasmaCash.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_PlasmaCash *PlasmaCashTransactorSession) Kill() (*types.Transaction, error) {
	return _PlasmaCash.Contract.Kill(&_PlasmaCash.TransactOpts)
}

// StartExit is a paid mutator transaction binding the contract method 0x5f10a47c.
//
// Solidity: function startExit(tokenID uint64, txBytes1 bytes, txBytes2 bytes, proof1 bytes, proof2 bytes, blk1 uint64, blk2 uint64) returns()
func (_PlasmaCash *PlasmaCashTransactor) StartExit(opts *bind.TransactOpts, tokenID uint64, txBytes1 []byte, txBytes2 []byte, proof1 []byte, proof2 []byte, blk1 uint64, blk2 uint64) (*types.Transaction, error) {
	return _PlasmaCash.contract.Transact(opts, "startExit", tokenID, txBytes1, txBytes2, proof1, proof2, blk1, blk2)
}

// StartExit is a paid mutator transaction binding the contract method 0x5f10a47c.
//
// Solidity: function startExit(tokenID uint64, txBytes1 bytes, txBytes2 bytes, proof1 bytes, proof2 bytes, blk1 uint64, blk2 uint64) returns()
func (_PlasmaCash *PlasmaCashSession) StartExit(tokenID uint64, txBytes1 []byte, txBytes2 []byte, proof1 []byte, proof2 []byte, blk1 uint64, blk2 uint64) (*types.Transaction, error) {
	return _PlasmaCash.Contract.StartExit(&_PlasmaCash.TransactOpts, tokenID, txBytes1, txBytes2, proof1, proof2, blk1, blk2)
}

// StartExit is a paid mutator transaction binding the contract method 0x5f10a47c.
//
// Solidity: function startExit(tokenID uint64, txBytes1 bytes, txBytes2 bytes, proof1 bytes, proof2 bytes, blk1 uint64, blk2 uint64) returns()
func (_PlasmaCash *PlasmaCashTransactorSession) StartExit(tokenID uint64, txBytes1 []byte, txBytes2 []byte, proof1 []byte, proof2 []byte, blk1 uint64, blk2 uint64) (*types.Transaction, error) {
	return _PlasmaCash.Contract.StartExit(&_PlasmaCash.TransactOpts, tokenID, txBytes1, txBytes2, proof1, proof2, blk1, blk2)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xefcfd072.
//
// Solidity: function submitBlock(_blkRoot bytes32, _blknum uint64) returns()
func (_PlasmaCash *PlasmaCashTransactor) SubmitBlock(opts *bind.TransactOpts, _blkRoot [32]byte, _blknum uint64) (*types.Transaction, error) {
	return _PlasmaCash.contract.Transact(opts, "submitBlock", _blkRoot, _blknum)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xefcfd072.
//
// Solidity: function submitBlock(_blkRoot bytes32, _blknum uint64) returns()
func (_PlasmaCash *PlasmaCashSession) SubmitBlock(_blkRoot [32]byte, _blknum uint64) (*types.Transaction, error) {
	return _PlasmaCash.Contract.SubmitBlock(&_PlasmaCash.TransactOpts, _blkRoot, _blknum)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xefcfd072.
//
// Solidity: function submitBlock(_blkRoot bytes32, _blknum uint64) returns()
func (_PlasmaCash *PlasmaCashTransactorSession) SubmitBlock(_blkRoot [32]byte, _blknum uint64) (*types.Transaction, error) {
	return _PlasmaCash.Contract.SubmitBlock(&_PlasmaCash.TransactOpts, _blkRoot, _blknum)
}

// PlasmaCashChallengeIterator is returned from FilterChallenge and is used to iterate over the raw logs and unpacked data for Challenge events raised by the PlasmaCash contract.
type PlasmaCashChallengeIterator struct {
	Event *PlasmaCashChallenge // Event containing the contract specifics and raw log

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
func (it *PlasmaCashChallengeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaCashChallenge)
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
		it.Event = new(PlasmaCashChallenge)
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
func (it *PlasmaCashChallengeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaCashChallengeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaCashChallenge represents a Challenge event raised by the PlasmaCash contract.
type PlasmaCashChallenge struct {
	Challenger common.Address
	TokenID    uint64
	Timestamp  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterChallenge is a free log retrieval operation binding the contract event 0xf4a2799edd9b2af498a251fd0f07a43047189b7090f3a6290d4e28645ffd4cd9.
//
// Solidity: event Challenge(_challenger address, _tokenID indexed uint64, _timestamp indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) FilterChallenge(opts *bind.FilterOpts, _tokenID []uint64, _timestamp []*big.Int) (*PlasmaCashChallengeIterator, error) {

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}
	var _timestampRule []interface{}
	for _, _timestampItem := range _timestamp {
		_timestampRule = append(_timestampRule, _timestampItem)
	}

	logs, sub, err := _PlasmaCash.contract.FilterLogs(opts, "Challenge", _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashChallengeIterator{contract: _PlasmaCash.contract, event: "Challenge", logs: logs, sub: sub}, nil
}

// WatchChallenge is a free log subscription operation binding the contract event 0xf4a2799edd9b2af498a251fd0f07a43047189b7090f3a6290d4e28645ffd4cd9.
//
// Solidity: event Challenge(_challenger address, _tokenID indexed uint64, _timestamp indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) WatchChallenge(opts *bind.WatchOpts, sink chan<- *PlasmaCashChallenge, _tokenID []uint64, _timestamp []*big.Int) (event.Subscription, error) {

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}
	var _timestampRule []interface{}
	for _, _timestampItem := range _timestamp {
		_timestampRule = append(_timestampRule, _timestampItem)
	}

	logs, sub, err := _PlasmaCash.contract.WatchLogs(opts, "Challenge", _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaCashChallenge)
				if err := _PlasmaCash.contract.UnpackLog(event, "Challenge", log); err != nil {
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

// PlasmaCashCurrtentExitIterator is returned from FilterCurrtentExit and is used to iterate over the raw logs and unpacked data for CurrtentExit events raised by the PlasmaCash contract.
type PlasmaCashCurrtentExitIterator struct {
	Event *PlasmaCashCurrtentExit // Event containing the contract specifics and raw log

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
func (it *PlasmaCashCurrtentExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaCashCurrtentExit)
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
		it.Event = new(PlasmaCashCurrtentExit)
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
func (it *PlasmaCashCurrtentExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaCashCurrtentExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaCashCurrtentExit represents a CurrtentExit event raised by the PlasmaCash contract.
type PlasmaCashCurrtentExit struct {
	DepID      uint64
	TokenID    uint64
	ExitableTS *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCurrtentExit is a free log retrieval operation binding the contract event 0xd0c1c7697d3a5f8a4ce2442c9dfebc6c40c0bee361654b5f97169f97dad2723a.
//
// Solidity: event CurrtentExit(depID uint64, tokenID uint64, exitableTS uint256)
func (_PlasmaCash *PlasmaCashFilterer) FilterCurrtentExit(opts *bind.FilterOpts) (*PlasmaCashCurrtentExitIterator, error) {

	logs, sub, err := _PlasmaCash.contract.FilterLogs(opts, "CurrtentExit")
	if err != nil {
		return nil, err
	}
	return &PlasmaCashCurrtentExitIterator{contract: _PlasmaCash.contract, event: "CurrtentExit", logs: logs, sub: sub}, nil
}

// WatchCurrtentExit is a free log subscription operation binding the contract event 0xd0c1c7697d3a5f8a4ce2442c9dfebc6c40c0bee361654b5f97169f97dad2723a.
//
// Solidity: event CurrtentExit(depID uint64, tokenID uint64, exitableTS uint256)
func (_PlasmaCash *PlasmaCashFilterer) WatchCurrtentExit(opts *bind.WatchOpts, sink chan<- *PlasmaCashCurrtentExit) (event.Subscription, error) {

	logs, sub, err := _PlasmaCash.contract.WatchLogs(opts, "CurrtentExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaCashCurrtentExit)
				if err := _PlasmaCash.contract.UnpackLog(event, "CurrtentExit", log); err != nil {
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

// PlasmaCashDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the PlasmaCash contract.
type PlasmaCashDepositIterator struct {
	Event *PlasmaCashDeposit // Event containing the contract specifics and raw log

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
func (it *PlasmaCashDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaCashDeposit)
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
		it.Event = new(PlasmaCashDeposit)
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
func (it *PlasmaCashDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaCashDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaCashDeposit represents a Deposit event raised by the PlasmaCash contract.
type PlasmaCashDeposit struct {
	Depositor    common.Address
	DepositIndex uint64
	Denomination uint64
	TokenID      uint64
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x96d929db0520d785b9981429377486f9182e32225c9a8b9f1b371519644cc68a.
//
// Solidity: event Deposit(_depositor address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64)
func (_PlasmaCash *PlasmaCashFilterer) FilterDeposit(opts *bind.FilterOpts, _depositIndex []uint64, _tokenID []uint64) (*PlasmaCashDepositIterator, error) {

	var _depositIndexRule []interface{}
	for _, _depositIndexItem := range _depositIndex {
		_depositIndexRule = append(_depositIndexRule, _depositIndexItem)
	}

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}

	logs, sub, err := _PlasmaCash.contract.FilterLogs(opts, "Deposit", _depositIndexRule, _tokenIDRule)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashDepositIterator{contract: _PlasmaCash.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x96d929db0520d785b9981429377486f9182e32225c9a8b9f1b371519644cc68a.
//
// Solidity: event Deposit(_depositor address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64)
func (_PlasmaCash *PlasmaCashFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *PlasmaCashDeposit, _depositIndex []uint64, _tokenID []uint64) (event.Subscription, error) {

	var _depositIndexRule []interface{}
	for _, _depositIndexItem := range _depositIndex {
		_depositIndexRule = append(_depositIndexRule, _depositIndexItem)
	}

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}

	logs, sub, err := _PlasmaCash.contract.WatchLogs(opts, "Deposit", _depositIndexRule, _tokenIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaCashDeposit)
				if err := _PlasmaCash.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// PlasmaCashExitCompletedIterator is returned from FilterExitCompleted and is used to iterate over the raw logs and unpacked data for ExitCompleted events raised by the PlasmaCash contract.
type PlasmaCashExitCompletedIterator struct {
	Event *PlasmaCashExitCompleted // Event containing the contract specifics and raw log

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
func (it *PlasmaCashExitCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaCashExitCompleted)
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
		it.Event = new(PlasmaCashExitCompleted)
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
func (it *PlasmaCashExitCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaCashExitCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaCashExitCompleted represents a ExitCompleted event raised by the PlasmaCash contract.
type PlasmaCashExitCompleted struct {
	Priority *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterExitCompleted is a free log retrieval operation binding the contract event 0xd9e811ef520eaecd3949de532960c72eda128c01bc3f2981aa414b511b157849.
//
// Solidity: event ExitCompleted(_priority indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) FilterExitCompleted(opts *bind.FilterOpts, _priority []*big.Int) (*PlasmaCashExitCompletedIterator, error) {

	var _priorityRule []interface{}
	for _, _priorityItem := range _priority {
		_priorityRule = append(_priorityRule, _priorityItem)
	}

	logs, sub, err := _PlasmaCash.contract.FilterLogs(opts, "ExitCompleted", _priorityRule)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashExitCompletedIterator{contract: _PlasmaCash.contract, event: "ExitCompleted", logs: logs, sub: sub}, nil
}

// WatchExitCompleted is a free log subscription operation binding the contract event 0xd9e811ef520eaecd3949de532960c72eda128c01bc3f2981aa414b511b157849.
//
// Solidity: event ExitCompleted(_priority indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) WatchExitCompleted(opts *bind.WatchOpts, sink chan<- *PlasmaCashExitCompleted, _priority []*big.Int) (event.Subscription, error) {

	var _priorityRule []interface{}
	for _, _priorityItem := range _priority {
		_priorityRule = append(_priorityRule, _priorityItem)
	}

	logs, sub, err := _PlasmaCash.contract.WatchLogs(opts, "ExitCompleted", _priorityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaCashExitCompleted)
				if err := _PlasmaCash.contract.UnpackLog(event, "ExitCompleted", log); err != nil {
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

// PlasmaCashExitStartedIterator is returned from FilterExitStarted and is used to iterate over the raw logs and unpacked data for ExitStarted events raised by the PlasmaCash contract.
type PlasmaCashExitStartedIterator struct {
	Event *PlasmaCashExitStarted // Event containing the contract specifics and raw log

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
func (it *PlasmaCashExitStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaCashExitStarted)
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
		it.Event = new(PlasmaCashExitStarted)
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
func (it *PlasmaCashExitStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaCashExitStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaCashExitStarted represents a ExitStarted event raised by the PlasmaCash contract.
type PlasmaCashExitStarted struct {
	Priority *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterExitStarted is a free log retrieval operation binding the contract event 0x2e6c46902e4864b15d699069473ce23cfcd871fa89d450c6049ce12399fe2f92.
//
// Solidity: event ExitStarted(_priority indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) FilterExitStarted(opts *bind.FilterOpts, _priority []*big.Int) (*PlasmaCashExitStartedIterator, error) {

	var _priorityRule []interface{}
	for _, _priorityItem := range _priority {
		_priorityRule = append(_priorityRule, _priorityItem)
	}

	logs, sub, err := _PlasmaCash.contract.FilterLogs(opts, "ExitStarted", _priorityRule)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashExitStartedIterator{contract: _PlasmaCash.contract, event: "ExitStarted", logs: logs, sub: sub}, nil
}

// WatchExitStarted is a free log subscription operation binding the contract event 0x2e6c46902e4864b15d699069473ce23cfcd871fa89d450c6049ce12399fe2f92.
//
// Solidity: event ExitStarted(_priority indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) WatchExitStarted(opts *bind.WatchOpts, sink chan<- *PlasmaCashExitStarted, _priority []*big.Int) (event.Subscription, error) {

	var _priorityRule []interface{}
	for _, _priorityItem := range _priority {
		_priorityRule = append(_priorityRule, _priorityItem)
	}

	logs, sub, err := _PlasmaCash.contract.WatchLogs(opts, "ExitStarted", _priorityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaCashExitStarted)
				if err := _PlasmaCash.contract.UnpackLog(event, "ExitStarted", log); err != nil {
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

// PlasmaCashExitTimeIterator is returned from FilterExitTime and is used to iterate over the raw logs and unpacked data for ExitTime events raised by the PlasmaCash contract.
type PlasmaCashExitTimeIterator struct {
	Event *PlasmaCashExitTime // Event containing the contract specifics and raw log

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
func (it *PlasmaCashExitTimeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaCashExitTime)
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
		it.Event = new(PlasmaCashExitTime)
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
func (it *PlasmaCashExitTimeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaCashExitTimeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaCashExitTime represents a ExitTime event raised by the PlasmaCash contract.
type PlasmaCashExitTime struct {
	ExitableTS *big.Int
	CuurrentTS *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterExitTime is a free log retrieval operation binding the contract event 0x05b05d48f45d9984a4696f7e1cb5b91fc5cd1e4420b57db43b82a9b97dccfa6a.
//
// Solidity: event ExitTime(exitableTS uint256, cuurrentTS uint256)
func (_PlasmaCash *PlasmaCashFilterer) FilterExitTime(opts *bind.FilterOpts) (*PlasmaCashExitTimeIterator, error) {

	logs, sub, err := _PlasmaCash.contract.FilterLogs(opts, "ExitTime")
	if err != nil {
		return nil, err
	}
	return &PlasmaCashExitTimeIterator{contract: _PlasmaCash.contract, event: "ExitTime", logs: logs, sub: sub}, nil
}

// WatchExitTime is a free log subscription operation binding the contract event 0x05b05d48f45d9984a4696f7e1cb5b91fc5cd1e4420b57db43b82a9b97dccfa6a.
//
// Solidity: event ExitTime(exitableTS uint256, cuurrentTS uint256)
func (_PlasmaCash *PlasmaCashFilterer) WatchExitTime(opts *bind.WatchOpts, sink chan<- *PlasmaCashExitTime) (event.Subscription, error) {

	logs, sub, err := _PlasmaCash.contract.WatchLogs(opts, "ExitTime")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaCashExitTime)
				if err := _PlasmaCash.contract.UnpackLog(event, "ExitTime", log); err != nil {
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

// PlasmaCashFinalizedExitIterator is returned from FilterFinalizedExit and is used to iterate over the raw logs and unpacked data for FinalizedExit events raised by the PlasmaCash contract.
type PlasmaCashFinalizedExitIterator struct {
	Event *PlasmaCashFinalizedExit // Event containing the contract specifics and raw log

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
func (it *PlasmaCashFinalizedExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaCashFinalizedExit)
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
		it.Event = new(PlasmaCashFinalizedExit)
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
func (it *PlasmaCashFinalizedExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaCashFinalizedExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaCashFinalizedExit represents a FinalizedExit event raised by the PlasmaCash contract.
type PlasmaCashFinalizedExit struct {
	Exiter       common.Address
	DepositIndex uint64
	Denomination uint64
	TokenID      uint64
	Timestamp    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFinalizedExit is a free log retrieval operation binding the contract event 0xe109ec77e392bf51e903e88a71fb4158aac7a024428a47f4b7c9fdf334abb362.
//
// Solidity: event FinalizedExit(_exiter address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64, _timestamp indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) FilterFinalizedExit(opts *bind.FilterOpts, _depositIndex []uint64, _tokenID []uint64, _timestamp []*big.Int) (*PlasmaCashFinalizedExitIterator, error) {

	var _depositIndexRule []interface{}
	for _, _depositIndexItem := range _depositIndex {
		_depositIndexRule = append(_depositIndexRule, _depositIndexItem)
	}

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}
	var _timestampRule []interface{}
	for _, _timestampItem := range _timestamp {
		_timestampRule = append(_timestampRule, _timestampItem)
	}

	logs, sub, err := _PlasmaCash.contract.FilterLogs(opts, "FinalizedExit", _depositIndexRule, _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashFinalizedExitIterator{contract: _PlasmaCash.contract, event: "FinalizedExit", logs: logs, sub: sub}, nil
}

// WatchFinalizedExit is a free log subscription operation binding the contract event 0xe109ec77e392bf51e903e88a71fb4158aac7a024428a47f4b7c9fdf334abb362.
//
// Solidity: event FinalizedExit(_exiter address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64, _timestamp indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) WatchFinalizedExit(opts *bind.WatchOpts, sink chan<- *PlasmaCashFinalizedExit, _depositIndex []uint64, _tokenID []uint64, _timestamp []*big.Int) (event.Subscription, error) {

	var _depositIndexRule []interface{}
	for _, _depositIndexItem := range _depositIndex {
		_depositIndexRule = append(_depositIndexRule, _depositIndexItem)
	}

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}
	var _timestampRule []interface{}
	for _, _timestampItem := range _timestamp {
		_timestampRule = append(_timestampRule, _timestampItem)
	}

	logs, sub, err := _PlasmaCash.contract.WatchLogs(opts, "FinalizedExit", _depositIndexRule, _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaCashFinalizedExit)
				if err := _PlasmaCash.contract.UnpackLog(event, "FinalizedExit", log); err != nil {
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

// PlasmaCashPublishedBlockIterator is returned from FilterPublishedBlock and is used to iterate over the raw logs and unpacked data for PublishedBlock events raised by the PlasmaCash contract.
type PlasmaCashPublishedBlockIterator struct {
	Event *PlasmaCashPublishedBlock // Event containing the contract specifics and raw log

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
func (it *PlasmaCashPublishedBlockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaCashPublishedBlock)
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
		it.Event = new(PlasmaCashPublishedBlock)
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
func (it *PlasmaCashPublishedBlockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaCashPublishedBlockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaCashPublishedBlock represents a PublishedBlock event raised by the PlasmaCash contract.
type PlasmaCashPublishedBlock struct {
	RootHash            [32]byte
	Blknum              uint64
	CurrentDepositIndex uint64
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterPublishedBlock is a free log retrieval operation binding the contract event 0xafb4a42a2439df06ac89ceace483c1382c37f9574eb7d217050bc54a63a33811.
//
// Solidity: event PublishedBlock(_rootHash bytes32, _blknum indexed uint64, _currentDepositIndex indexed uint64)
func (_PlasmaCash *PlasmaCashFilterer) FilterPublishedBlock(opts *bind.FilterOpts, _blknum []uint64, _currentDepositIndex []uint64) (*PlasmaCashPublishedBlockIterator, error) {

	var _blknumRule []interface{}
	for _, _blknumItem := range _blknum {
		_blknumRule = append(_blknumRule, _blknumItem)
	}
	var _currentDepositIndexRule []interface{}
	for _, _currentDepositIndexItem := range _currentDepositIndex {
		_currentDepositIndexRule = append(_currentDepositIndexRule, _currentDepositIndexItem)
	}

	logs, sub, err := _PlasmaCash.contract.FilterLogs(opts, "PublishedBlock", _blknumRule, _currentDepositIndexRule)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashPublishedBlockIterator{contract: _PlasmaCash.contract, event: "PublishedBlock", logs: logs, sub: sub}, nil
}

// WatchPublishedBlock is a free log subscription operation binding the contract event 0xafb4a42a2439df06ac89ceace483c1382c37f9574eb7d217050bc54a63a33811.
//
// Solidity: event PublishedBlock(_rootHash bytes32, _blknum indexed uint64, _currentDepositIndex indexed uint64)
func (_PlasmaCash *PlasmaCashFilterer) WatchPublishedBlock(opts *bind.WatchOpts, sink chan<- *PlasmaCashPublishedBlock, _blknum []uint64, _currentDepositIndex []uint64) (event.Subscription, error) {

	var _blknumRule []interface{}
	for _, _blknumItem := range _blknum {
		_blknumRule = append(_blknumRule, _blknumItem)
	}
	var _currentDepositIndexRule []interface{}
	for _, _currentDepositIndexItem := range _currentDepositIndex {
		_currentDepositIndexRule = append(_currentDepositIndexRule, _currentDepositIndexItem)
	}

	logs, sub, err := _PlasmaCash.contract.WatchLogs(opts, "PublishedBlock", _blknumRule, _currentDepositIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaCashPublishedBlock)
				if err := _PlasmaCash.contract.UnpackLog(event, "PublishedBlock", log); err != nil {
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

// PlasmaCashStartExitIterator is returned from FilterStartExit and is used to iterate over the raw logs and unpacked data for StartExit events raised by the PlasmaCash contract.
type PlasmaCashStartExitIterator struct {
	Event *PlasmaCashStartExit // Event containing the contract specifics and raw log

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
func (it *PlasmaCashStartExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaCashStartExit)
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
		it.Event = new(PlasmaCashStartExit)
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
func (it *PlasmaCashStartExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaCashStartExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaCashStartExit represents a StartExit event raised by the PlasmaCash contract.
type PlasmaCashStartExit struct {
	Exiter       common.Address
	DepositIndex uint64
	Denomination uint64
	TokenID      uint64
	Timestamp    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStartExit is a free log retrieval operation binding the contract event 0x6b75ed42e8219410363fb31fa601dea7c8b3549bfc9428b1eeaf53a6226cb668.
//
// Solidity: event StartExit(_exiter address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64, _timestamp indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) FilterStartExit(opts *bind.FilterOpts, _depositIndex []uint64, _tokenID []uint64, _timestamp []*big.Int) (*PlasmaCashStartExitIterator, error) {

	var _depositIndexRule []interface{}
	for _, _depositIndexItem := range _depositIndex {
		_depositIndexRule = append(_depositIndexRule, _depositIndexItem)
	}

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}
	var _timestampRule []interface{}
	for _, _timestampItem := range _timestamp {
		_timestampRule = append(_timestampRule, _timestampItem)
	}

	logs, sub, err := _PlasmaCash.contract.FilterLogs(opts, "StartExit", _depositIndexRule, _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return &PlasmaCashStartExitIterator{contract: _PlasmaCash.contract, event: "StartExit", logs: logs, sub: sub}, nil
}

// WatchStartExit is a free log subscription operation binding the contract event 0x6b75ed42e8219410363fb31fa601dea7c8b3549bfc9428b1eeaf53a6226cb668.
//
// Solidity: event StartExit(_exiter address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64, _timestamp indexed uint256)
func (_PlasmaCash *PlasmaCashFilterer) WatchStartExit(opts *bind.WatchOpts, sink chan<- *PlasmaCashStartExit, _depositIndex []uint64, _tokenID []uint64, _timestamp []*big.Int) (event.Subscription, error) {

	var _depositIndexRule []interface{}
	for _, _depositIndexItem := range _depositIndex {
		_depositIndexRule = append(_depositIndexRule, _depositIndexItem)
	}

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}
	var _timestampRule []interface{}
	for _, _timestampItem := range _timestamp {
		_timestampRule = append(_timestampRule, _timestampItem)
	}

	logs, sub, err := _PlasmaCash.contract.WatchLogs(opts, "StartExit", _depositIndexRule, _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaCashStartExit)
				if err := _PlasmaCash.contract.UnpackLog(event, "StartExit", log); err != nil {
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
