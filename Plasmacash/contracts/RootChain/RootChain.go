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

// RootChainABI is the input ABI used to generate the binding from.
const RootChainABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"depositIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getNextExit\",\"outputs\":[{\"name\":\"depID\",\"type\":\"uint64\"},{\"name\":\"tokenID\",\"type\":\"uint64\"},{\"name\":\"exitableTS\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"},{\"name\":\"blk\",\"type\":\"uint64\"}],\"name\":\"challenge\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentDepositIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_depositIndex\",\"type\":\"uint64\"}],\"name\":\"depositExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"childChain\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"prevTxBytes\",\"type\":\"bytes\"},{\"name\":\"prevProof\",\"type\":\"bytes\"},{\"name\":\"prevBlk\",\"type\":\"uint64\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"},{\"name\":\"blk\",\"type\":\"uint64\"}],\"name\":\"startExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"},{\"name\":\"blk\",\"type\":\"uint64\"},{\"name\":\"faultyTxBytes\",\"type\":\"bytes\"},{\"name\":\"faultyProof\",\"type\":\"bytes\"},{\"name\":\"faultyBlk\",\"type\":\"uint64\"}],\"name\":\"challengeBefore\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentBlkNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeExits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"exits\",\"outputs\":[{\"name\":\"prevBlk\",\"type\":\"uint64\"},{\"name\":\"exitBlk\",\"type\":\"uint64\"},{\"name\":\"exitor\",\"type\":\"address\"},{\"name\":\"exitableTS\",\"type\":\"uint256\"},{\"name\":\"bond\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_blkRoot\",\"type\":\"bytes32\"},{\"name\":\"_blknum\",\"type\":\"uint64\"}],\"name\":\"submitBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"depositBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_depositor\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_depositIndex\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"_denomination\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_tokenID\",\"type\":\"uint64\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_exiter\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_depositIndex\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"_denomination\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_tokenID\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"StartExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_rootHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"_blknum\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_currentDepositIndex\",\"type\":\"uint64\"}],\"name\":\"PublishedBlock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_exiter\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_depositIndex\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"_denomination\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_tokenID\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"FinalizedExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_challenger\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenID\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Challenge\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_priority\",\"type\":\"uint256\"}],\"name\":\"ExitStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_priority\",\"type\":\"uint256\"}],\"name\":\"ExitCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"depID\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"tokenID\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"exitableTS\",\"type\":\"uint256\"}],\"name\":\"CurrtentExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"exitableTS\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"cuurrentTS\",\"type\":\"uint256\"}],\"name\":\"ExitTime\",\"type\":\"event\"}]"

// RootChain is an auto generated Go binding around an Ethereum contract.
type RootChain struct {
	RootChainCaller     // Read-only binding to the contract
	RootChainTransactor // Write-only binding to the contract
	RootChainFilterer   // Log filterer for contract events
}

// RootChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type RootChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RootChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RootChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RootChainSession struct {
	Contract     *RootChain        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RootChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RootChainCallerSession struct {
	Contract *RootChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// RootChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RootChainTransactorSession struct {
	Contract     *RootChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RootChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type RootChainRaw struct {
	Contract *RootChain // Generic contract binding to access the raw methods on
}

// RootChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RootChainCallerRaw struct {
	Contract *RootChainCaller // Generic read-only contract binding to access the raw methods on
}

// RootChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RootChainTransactorRaw struct {
	Contract *RootChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRootChain creates a new instance of RootChain, bound to a specific deployed contract.
func NewRootChain(address common.Address, backend bind.ContractBackend) (*RootChain, error) {
	contract, err := bindRootChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RootChain{RootChainCaller: RootChainCaller{contract: contract}, RootChainTransactor: RootChainTransactor{contract: contract}, RootChainFilterer: RootChainFilterer{contract: contract}}, nil
}

// NewRootChainCaller creates a new read-only instance of RootChain, bound to a specific deployed contract.
func NewRootChainCaller(address common.Address, caller bind.ContractCaller) (*RootChainCaller, error) {
	contract, err := bindRootChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RootChainCaller{contract: contract}, nil
}

// NewRootChainTransactor creates a new write-only instance of RootChain, bound to a specific deployed contract.
func NewRootChainTransactor(address common.Address, transactor bind.ContractTransactor) (*RootChainTransactor, error) {
	contract, err := bindRootChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RootChainTransactor{contract: contract}, nil
}

// NewRootChainFilterer creates a new log filterer instance of RootChain, bound to a specific deployed contract.
func NewRootChainFilterer(address common.Address, filterer bind.ContractFilterer) (*RootChainFilterer, error) {
	contract, err := bindRootChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RootChainFilterer{contract: contract}, nil
}

// bindRootChain binds a generic wrapper to an already deployed contract.
func bindRootChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RootChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootChain *RootChainRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootChain.Contract.RootChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootChain *RootChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.Contract.RootChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootChain *RootChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootChain.Contract.RootChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootChain *RootChainCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootChain *RootChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootChain *RootChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootChain.Contract.contract.Transact(opts, method, params...)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() constant returns(address)
func (_RootChain *RootChainCaller) Authority(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "authority")
	return *ret0, err
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() constant returns(address)
func (_RootChain *RootChainSession) Authority() (common.Address, error) {
	return _RootChain.Contract.Authority(&_RootChain.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() constant returns(address)
func (_RootChain *RootChainCallerSession) Authority() (common.Address, error) {
	return _RootChain.Contract.Authority(&_RootChain.CallOpts)
}

// ChildChain is a free data retrieval call binding the contract method 0x6434693a.
//
// Solidity: function childChain( uint64) constant returns(bytes32)
func (_RootChain *RootChainCaller) ChildChain(opts *bind.CallOpts, arg0 uint64) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "childChain", arg0)
	return *ret0, err
}

// ChildChain is a free data retrieval call binding the contract method 0x6434693a.
//
// Solidity: function childChain( uint64) constant returns(bytes32)
func (_RootChain *RootChainSession) ChildChain(arg0 uint64) ([32]byte, error) {
	return _RootChain.Contract.ChildChain(&_RootChain.CallOpts, arg0)
}

// ChildChain is a free data retrieval call binding the contract method 0x6434693a.
//
// Solidity: function childChain( uint64) constant returns(bytes32)
func (_RootChain *RootChainCallerSession) ChildChain(arg0 uint64) ([32]byte, error) {
	return _RootChain.Contract.ChildChain(&_RootChain.CallOpts, arg0)
}

// CurrentBlkNum is a free data retrieval call binding the contract method 0x8c1cdd5d.
//
// Solidity: function currentBlkNum() constant returns(uint64)
func (_RootChain *RootChainCaller) CurrentBlkNum(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "currentBlkNum")
	return *ret0, err
}

// CurrentBlkNum is a free data retrieval call binding the contract method 0x8c1cdd5d.
//
// Solidity: function currentBlkNum() constant returns(uint64)
func (_RootChain *RootChainSession) CurrentBlkNum() (uint64, error) {
	return _RootChain.Contract.CurrentBlkNum(&_RootChain.CallOpts)
}

// CurrentBlkNum is a free data retrieval call binding the contract method 0x8c1cdd5d.
//
// Solidity: function currentBlkNum() constant returns(uint64)
func (_RootChain *RootChainCallerSession) CurrentBlkNum() (uint64, error) {
	return _RootChain.Contract.CurrentBlkNum(&_RootChain.CallOpts)
}

// CurrentDepositIndex is a free data retrieval call binding the contract method 0x36e5df06.
//
// Solidity: function currentDepositIndex() constant returns(uint64)
func (_RootChain *RootChainCaller) CurrentDepositIndex(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "currentDepositIndex")
	return *ret0, err
}

// CurrentDepositIndex is a free data retrieval call binding the contract method 0x36e5df06.
//
// Solidity: function currentDepositIndex() constant returns(uint64)
func (_RootChain *RootChainSession) CurrentDepositIndex() (uint64, error) {
	return _RootChain.Contract.CurrentDepositIndex(&_RootChain.CallOpts)
}

// CurrentDepositIndex is a free data retrieval call binding the contract method 0x36e5df06.
//
// Solidity: function currentDepositIndex() constant returns(uint64)
func (_RootChain *RootChainCallerSession) CurrentDepositIndex() (uint64, error) {
	return _RootChain.Contract.CurrentDepositIndex(&_RootChain.CallOpts)
}

// DepositBalance is a free data retrieval call binding the contract method 0xf34272ab.
//
// Solidity: function depositBalance( uint64) constant returns(uint64)
func (_RootChain *RootChainCaller) DepositBalance(opts *bind.CallOpts, arg0 uint64) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "depositBalance", arg0)
	return *ret0, err
}

// DepositBalance is a free data retrieval call binding the contract method 0xf34272ab.
//
// Solidity: function depositBalance( uint64) constant returns(uint64)
func (_RootChain *RootChainSession) DepositBalance(arg0 uint64) (uint64, error) {
	return _RootChain.Contract.DepositBalance(&_RootChain.CallOpts, arg0)
}

// DepositBalance is a free data retrieval call binding the contract method 0xf34272ab.
//
// Solidity: function depositBalance( uint64) constant returns(uint64)
func (_RootChain *RootChainCallerSession) DepositBalance(arg0 uint64) (uint64, error) {
	return _RootChain.Contract.DepositBalance(&_RootChain.CallOpts, arg0)
}

// DepositIndex is a free data retrieval call binding the contract method 0x0f8592ac.
//
// Solidity: function depositIndex( uint64) constant returns(uint64)
func (_RootChain *RootChainCaller) DepositIndex(opts *bind.CallOpts, arg0 uint64) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "depositIndex", arg0)
	return *ret0, err
}

// DepositIndex is a free data retrieval call binding the contract method 0x0f8592ac.
//
// Solidity: function depositIndex( uint64) constant returns(uint64)
func (_RootChain *RootChainSession) DepositIndex(arg0 uint64) (uint64, error) {
	return _RootChain.Contract.DepositIndex(&_RootChain.CallOpts, arg0)
}

// DepositIndex is a free data retrieval call binding the contract method 0x0f8592ac.
//
// Solidity: function depositIndex( uint64) constant returns(uint64)
func (_RootChain *RootChainCallerSession) DepositIndex(arg0 uint64) (uint64, error) {
	return _RootChain.Contract.DepositIndex(&_RootChain.CallOpts, arg0)
}

// Exits is a free data retrieval call binding the contract method 0xd6463d40.
//
// Solidity: function exits( uint64) constant returns(prevBlk uint64, exitBlk uint64, exitor address, exitableTS uint256, bond uint256)
func (_RootChain *RootChainCaller) Exits(opts *bind.CallOpts, arg0 uint64) (struct {
	PrevBlk    uint64
	ExitBlk    uint64
	Exitor     common.Address
	ExitableTS *big.Int
	Bond       *big.Int
}, error) {
	ret := new(struct {
		PrevBlk    uint64
		ExitBlk    uint64
		Exitor     common.Address
		ExitableTS *big.Int
		Bond       *big.Int
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "exits", arg0)
	return *ret, err
}

// Exits is a free data retrieval call binding the contract method 0xd6463d40.
//
// Solidity: function exits( uint64) constant returns(prevBlk uint64, exitBlk uint64, exitor address, exitableTS uint256, bond uint256)
func (_RootChain *RootChainSession) Exits(arg0 uint64) (struct {
	PrevBlk    uint64
	ExitBlk    uint64
	Exitor     common.Address
	ExitableTS *big.Int
	Bond       *big.Int
}, error) {
	return _RootChain.Contract.Exits(&_RootChain.CallOpts, arg0)
}

// Exits is a free data retrieval call binding the contract method 0xd6463d40.
//
// Solidity: function exits( uint64) constant returns(prevBlk uint64, exitBlk uint64, exitor address, exitableTS uint256, bond uint256)
func (_RootChain *RootChainCallerSession) Exits(arg0 uint64) (struct {
	PrevBlk    uint64
	ExitBlk    uint64
	Exitor     common.Address
	ExitableTS *big.Int
	Bond       *big.Int
}, error) {
	return _RootChain.Contract.Exits(&_RootChain.CallOpts, arg0)
}

// GetNextExit is a free data retrieval call binding the contract method 0x2fd56a25.
//
// Solidity: function getNextExit() constant returns(depID uint64, tokenID uint64, exitableTS uint256)
func (_RootChain *RootChainCaller) GetNextExit(opts *bind.CallOpts) (struct {
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
	err := _RootChain.contract.Call(opts, out, "getNextExit")
	return *ret, err
}

// GetNextExit is a free data retrieval call binding the contract method 0x2fd56a25.
//
// Solidity: function getNextExit() constant returns(depID uint64, tokenID uint64, exitableTS uint256)
func (_RootChain *RootChainSession) GetNextExit() (struct {
	DepID      uint64
	TokenID    uint64
	ExitableTS *big.Int
}, error) {
	return _RootChain.Contract.GetNextExit(&_RootChain.CallOpts)
}

// GetNextExit is a free data retrieval call binding the contract method 0x2fd56a25.
//
// Solidity: function getNextExit() constant returns(depID uint64, tokenID uint64, exitableTS uint256)
func (_RootChain *RootChainCallerSession) GetNextExit() (struct {
	DepID      uint64
	TokenID    uint64
	ExitableTS *big.Int
}, error) {
	return _RootChain.Contract.GetNextExit(&_RootChain.CallOpts)
}

// Challenge is a paid mutator transaction binding the contract method 0x31ddf645.
//
// Solidity: function challenge(txBytes bytes, proof bytes, blk uint64) returns()
func (_RootChain *RootChainTransactor) Challenge(opts *bind.TransactOpts, txBytes []byte, proof []byte, blk uint64) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "challenge", txBytes, proof, blk)
}

// Challenge is a paid mutator transaction binding the contract method 0x31ddf645.
//
// Solidity: function challenge(txBytes bytes, proof bytes, blk uint64) returns()
func (_RootChain *RootChainSession) Challenge(txBytes []byte, proof []byte, blk uint64) (*types.Transaction, error) {
	return _RootChain.Contract.Challenge(&_RootChain.TransactOpts, txBytes, proof, blk)
}

// Challenge is a paid mutator transaction binding the contract method 0x31ddf645.
//
// Solidity: function challenge(txBytes bytes, proof bytes, blk uint64) returns()
func (_RootChain *RootChainTransactorSession) Challenge(txBytes []byte, proof []byte, blk uint64) (*types.Transaction, error) {
	return _RootChain.Contract.Challenge(&_RootChain.TransactOpts, txBytes, proof, blk)
}

// ChallengeBefore is a paid mutator transaction binding the contract method 0x743098a7.
//
// Solidity: function challengeBefore(txBytes bytes, proof bytes, blk uint64, faultyTxBytes bytes, faultyProof bytes, faultyBlk uint64) returns()
func (_RootChain *RootChainTransactor) ChallengeBefore(opts *bind.TransactOpts, txBytes []byte, proof []byte, blk uint64, faultyTxBytes []byte, faultyProof []byte, faultyBlk uint64) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "challengeBefore", txBytes, proof, blk, faultyTxBytes, faultyProof, faultyBlk)
}

// ChallengeBefore is a paid mutator transaction binding the contract method 0x743098a7.
//
// Solidity: function challengeBefore(txBytes bytes, proof bytes, blk uint64, faultyTxBytes bytes, faultyProof bytes, faultyBlk uint64) returns()
func (_RootChain *RootChainSession) ChallengeBefore(txBytes []byte, proof []byte, blk uint64, faultyTxBytes []byte, faultyProof []byte, faultyBlk uint64) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeBefore(&_RootChain.TransactOpts, txBytes, proof, blk, faultyTxBytes, faultyProof, faultyBlk)
}

// ChallengeBefore is a paid mutator transaction binding the contract method 0x743098a7.
//
// Solidity: function challengeBefore(txBytes bytes, proof bytes, blk uint64, faultyTxBytes bytes, faultyProof bytes, faultyBlk uint64) returns()
func (_RootChain *RootChainTransactorSession) ChallengeBefore(txBytes []byte, proof []byte, blk uint64, faultyTxBytes []byte, faultyProof []byte, faultyBlk uint64) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeBefore(&_RootChain.TransactOpts, txBytes, proof, blk, faultyTxBytes, faultyProof, faultyBlk)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_RootChain *RootChainTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_RootChain *RootChainSession) Deposit() (*types.Transaction, error) {
	return _RootChain.Contract.Deposit(&_RootChain.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_RootChain *RootChainTransactorSession) Deposit() (*types.Transaction, error) {
	return _RootChain.Contract.Deposit(&_RootChain.TransactOpts)
}

// DepositExit is a paid mutator transaction binding the contract method 0x3d2e2590.
//
// Solidity: function depositExit(_depositIndex uint64) returns()
func (_RootChain *RootChainTransactor) DepositExit(opts *bind.TransactOpts, _depositIndex uint64) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "depositExit", _depositIndex)
}

// DepositExit is a paid mutator transaction binding the contract method 0x3d2e2590.
//
// Solidity: function depositExit(_depositIndex uint64) returns()
func (_RootChain *RootChainSession) DepositExit(_depositIndex uint64) (*types.Transaction, error) {
	return _RootChain.Contract.DepositExit(&_RootChain.TransactOpts, _depositIndex)
}

// DepositExit is a paid mutator transaction binding the contract method 0x3d2e2590.
//
// Solidity: function depositExit(_depositIndex uint64) returns()
func (_RootChain *RootChainTransactorSession) DepositExit(_depositIndex uint64) (*types.Transaction, error) {
	return _RootChain.Contract.DepositExit(&_RootChain.TransactOpts, _depositIndex)
}

// FinalizeExits is a paid mutator transaction binding the contract method 0xc6ab44cd.
//
// Solidity: function finalizeExits() returns()
func (_RootChain *RootChainTransactor) FinalizeExits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "finalizeExits")
}

// FinalizeExits is a paid mutator transaction binding the contract method 0xc6ab44cd.
//
// Solidity: function finalizeExits() returns()
func (_RootChain *RootChainSession) FinalizeExits() (*types.Transaction, error) {
	return _RootChain.Contract.FinalizeExits(&_RootChain.TransactOpts)
}

// FinalizeExits is a paid mutator transaction binding the contract method 0xc6ab44cd.
//
// Solidity: function finalizeExits() returns()
func (_RootChain *RootChainTransactorSession) FinalizeExits() (*types.Transaction, error) {
	return _RootChain.Contract.FinalizeExits(&_RootChain.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_RootChain *RootChainTransactor) Kill(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "kill")
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_RootChain *RootChainSession) Kill() (*types.Transaction, error) {
	return _RootChain.Contract.Kill(&_RootChain.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_RootChain *RootChainTransactorSession) Kill() (*types.Transaction, error) {
	return _RootChain.Contract.Kill(&_RootChain.TransactOpts)
}

// StartExit is a paid mutator transaction binding the contract method 0x6a1473ee.
//
// Solidity: function startExit(prevTxBytes bytes, prevProof bytes, prevBlk uint64, txBytes bytes, proof bytes, blk uint64) returns()
func (_RootChain *RootChainTransactor) StartExit(opts *bind.TransactOpts, prevTxBytes []byte, prevProof []byte, prevBlk uint64, txBytes []byte, proof []byte, blk uint64) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "startExit", prevTxBytes, prevProof, prevBlk, txBytes, proof, blk)
}

// StartExit is a paid mutator transaction binding the contract method 0x6a1473ee.
//
// Solidity: function startExit(prevTxBytes bytes, prevProof bytes, prevBlk uint64, txBytes bytes, proof bytes, blk uint64) returns()
func (_RootChain *RootChainSession) StartExit(prevTxBytes []byte, prevProof []byte, prevBlk uint64, txBytes []byte, proof []byte, blk uint64) (*types.Transaction, error) {
	return _RootChain.Contract.StartExit(&_RootChain.TransactOpts, prevTxBytes, prevProof, prevBlk, txBytes, proof, blk)
}

// StartExit is a paid mutator transaction binding the contract method 0x6a1473ee.
//
// Solidity: function startExit(prevTxBytes bytes, prevProof bytes, prevBlk uint64, txBytes bytes, proof bytes, blk uint64) returns()
func (_RootChain *RootChainTransactorSession) StartExit(prevTxBytes []byte, prevProof []byte, prevBlk uint64, txBytes []byte, proof []byte, blk uint64) (*types.Transaction, error) {
	return _RootChain.Contract.StartExit(&_RootChain.TransactOpts, prevTxBytes, prevProof, prevBlk, txBytes, proof, blk)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xefcfd072.
//
// Solidity: function submitBlock(_blkRoot bytes32, _blknum uint64) returns()
func (_RootChain *RootChainTransactor) SubmitBlock(opts *bind.TransactOpts, _blkRoot [32]byte, _blknum uint64) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "submitBlock", _blkRoot, _blknum)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xefcfd072.
//
// Solidity: function submitBlock(_blkRoot bytes32, _blknum uint64) returns()
func (_RootChain *RootChainSession) SubmitBlock(_blkRoot [32]byte, _blknum uint64) (*types.Transaction, error) {
	return _RootChain.Contract.SubmitBlock(&_RootChain.TransactOpts, _blkRoot, _blknum)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xefcfd072.
//
// Solidity: function submitBlock(_blkRoot bytes32, _blknum uint64) returns()
func (_RootChain *RootChainTransactorSession) SubmitBlock(_blkRoot [32]byte, _blknum uint64) (*types.Transaction, error) {
	return _RootChain.Contract.SubmitBlock(&_RootChain.TransactOpts, _blkRoot, _blknum)
}

// RootChainChallengeIterator is returned from FilterChallenge and is used to iterate over the raw logs and unpacked data for Challenge events raised by the RootChain contract.
type RootChainChallengeIterator struct {
	Event *RootChainChallenge // Event containing the contract specifics and raw log

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
func (it *RootChainChallengeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainChallenge)
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
		it.Event = new(RootChainChallenge)
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
func (it *RootChainChallengeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainChallengeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainChallenge represents a Challenge event raised by the RootChain contract.
type RootChainChallenge struct {
	Challenger common.Address
	TokenID    uint64
	Timestamp  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterChallenge is a free log retrieval operation binding the contract event 0xf4a2799edd9b2af498a251fd0f07a43047189b7090f3a6290d4e28645ffd4cd9.
//
// Solidity: event Challenge(_challenger address, _tokenID indexed uint64, _timestamp indexed uint256)
func (_RootChain *RootChainFilterer) FilterChallenge(opts *bind.FilterOpts, _tokenID []uint64, _timestamp []*big.Int) (*RootChainChallengeIterator, error) {

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}
	var _timestampRule []interface{}
	for _, _timestampItem := range _timestamp {
		_timestampRule = append(_timestampRule, _timestampItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "Challenge", _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return &RootChainChallengeIterator{contract: _RootChain.contract, event: "Challenge", logs: logs, sub: sub}, nil
}

// WatchChallenge is a free log subscription operation binding the contract event 0xf4a2799edd9b2af498a251fd0f07a43047189b7090f3a6290d4e28645ffd4cd9.
//
// Solidity: event Challenge(_challenger address, _tokenID indexed uint64, _timestamp indexed uint256)
func (_RootChain *RootChainFilterer) WatchChallenge(opts *bind.WatchOpts, sink chan<- *RootChainChallenge, _tokenID []uint64, _timestamp []*big.Int) (event.Subscription, error) {

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}
	var _timestampRule []interface{}
	for _, _timestampItem := range _timestamp {
		_timestampRule = append(_timestampRule, _timestampItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "Challenge", _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainChallenge)
				if err := _RootChain.contract.UnpackLog(event, "Challenge", log); err != nil {
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

// RootChainCurrtentExitIterator is returned from FilterCurrtentExit and is used to iterate over the raw logs and unpacked data for CurrtentExit events raised by the RootChain contract.
type RootChainCurrtentExitIterator struct {
	Event *RootChainCurrtentExit // Event containing the contract specifics and raw log

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
func (it *RootChainCurrtentExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainCurrtentExit)
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
		it.Event = new(RootChainCurrtentExit)
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
func (it *RootChainCurrtentExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainCurrtentExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainCurrtentExit represents a CurrtentExit event raised by the RootChain contract.
type RootChainCurrtentExit struct {
	DepID      uint64
	TokenID    uint64
	ExitableTS *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCurrtentExit is a free log retrieval operation binding the contract event 0xd0c1c7697d3a5f8a4ce2442c9dfebc6c40c0bee361654b5f97169f97dad2723a.
//
// Solidity: event CurrtentExit(depID uint64, tokenID uint64, exitableTS uint256)
func (_RootChain *RootChainFilterer) FilterCurrtentExit(opts *bind.FilterOpts) (*RootChainCurrtentExitIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "CurrtentExit")
	if err != nil {
		return nil, err
	}
	return &RootChainCurrtentExitIterator{contract: _RootChain.contract, event: "CurrtentExit", logs: logs, sub: sub}, nil
}

// WatchCurrtentExit is a free log subscription operation binding the contract event 0xd0c1c7697d3a5f8a4ce2442c9dfebc6c40c0bee361654b5f97169f97dad2723a.
//
// Solidity: event CurrtentExit(depID uint64, tokenID uint64, exitableTS uint256)
func (_RootChain *RootChainFilterer) WatchCurrtentExit(opts *bind.WatchOpts, sink chan<- *RootChainCurrtentExit) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "CurrtentExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainCurrtentExit)
				if err := _RootChain.contract.UnpackLog(event, "CurrtentExit", log); err != nil {
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

// RootChainDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the RootChain contract.
type RootChainDepositIterator struct {
	Event *RootChainDeposit // Event containing the contract specifics and raw log

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
func (it *RootChainDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainDeposit)
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
		it.Event = new(RootChainDeposit)
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
func (it *RootChainDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainDeposit represents a Deposit event raised by the RootChain contract.
type RootChainDeposit struct {
	Depositor    common.Address
	DepositIndex uint64
	Denomination uint64
	TokenID      uint64
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x96d929db0520d785b9981429377486f9182e32225c9a8b9f1b371519644cc68a.
//
// Solidity: event Deposit(_depositor address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64)
func (_RootChain *RootChainFilterer) FilterDeposit(opts *bind.FilterOpts, _depositIndex []uint64, _tokenID []uint64) (*RootChainDepositIterator, error) {

	var _depositIndexRule []interface{}
	for _, _depositIndexItem := range _depositIndex {
		_depositIndexRule = append(_depositIndexRule, _depositIndexItem)
	}

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "Deposit", _depositIndexRule, _tokenIDRule)
	if err != nil {
		return nil, err
	}
	return &RootChainDepositIterator{contract: _RootChain.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x96d929db0520d785b9981429377486f9182e32225c9a8b9f1b371519644cc68a.
//
// Solidity: event Deposit(_depositor address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64)
func (_RootChain *RootChainFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *RootChainDeposit, _depositIndex []uint64, _tokenID []uint64) (event.Subscription, error) {

	var _depositIndexRule []interface{}
	for _, _depositIndexItem := range _depositIndex {
		_depositIndexRule = append(_depositIndexRule, _depositIndexItem)
	}

	var _tokenIDRule []interface{}
	for _, _tokenIDItem := range _tokenID {
		_tokenIDRule = append(_tokenIDRule, _tokenIDItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "Deposit", _depositIndexRule, _tokenIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainDeposit)
				if err := _RootChain.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// RootChainExitCompletedIterator is returned from FilterExitCompleted and is used to iterate over the raw logs and unpacked data for ExitCompleted events raised by the RootChain contract.
type RootChainExitCompletedIterator struct {
	Event *RootChainExitCompleted // Event containing the contract specifics and raw log

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
func (it *RootChainExitCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainExitCompleted)
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
		it.Event = new(RootChainExitCompleted)
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
func (it *RootChainExitCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainExitCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainExitCompleted represents a ExitCompleted event raised by the RootChain contract.
type RootChainExitCompleted struct {
	Priority *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterExitCompleted is a free log retrieval operation binding the contract event 0xd9e811ef520eaecd3949de532960c72eda128c01bc3f2981aa414b511b157849.
//
// Solidity: event ExitCompleted(_priority indexed uint256)
func (_RootChain *RootChainFilterer) FilterExitCompleted(opts *bind.FilterOpts, _priority []*big.Int) (*RootChainExitCompletedIterator, error) {

	var _priorityRule []interface{}
	for _, _priorityItem := range _priority {
		_priorityRule = append(_priorityRule, _priorityItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "ExitCompleted", _priorityRule)
	if err != nil {
		return nil, err
	}
	return &RootChainExitCompletedIterator{contract: _RootChain.contract, event: "ExitCompleted", logs: logs, sub: sub}, nil
}

// WatchExitCompleted is a free log subscription operation binding the contract event 0xd9e811ef520eaecd3949de532960c72eda128c01bc3f2981aa414b511b157849.
//
// Solidity: event ExitCompleted(_priority indexed uint256)
func (_RootChain *RootChainFilterer) WatchExitCompleted(opts *bind.WatchOpts, sink chan<- *RootChainExitCompleted, _priority []*big.Int) (event.Subscription, error) {

	var _priorityRule []interface{}
	for _, _priorityItem := range _priority {
		_priorityRule = append(_priorityRule, _priorityItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "ExitCompleted", _priorityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainExitCompleted)
				if err := _RootChain.contract.UnpackLog(event, "ExitCompleted", log); err != nil {
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

// RootChainExitStartedIterator is returned from FilterExitStarted and is used to iterate over the raw logs and unpacked data for ExitStarted events raised by the RootChain contract.
type RootChainExitStartedIterator struct {
	Event *RootChainExitStarted // Event containing the contract specifics and raw log

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
func (it *RootChainExitStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainExitStarted)
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
		it.Event = new(RootChainExitStarted)
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
func (it *RootChainExitStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainExitStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainExitStarted represents a ExitStarted event raised by the RootChain contract.
type RootChainExitStarted struct {
	Priority *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterExitStarted is a free log retrieval operation binding the contract event 0x2e6c46902e4864b15d699069473ce23cfcd871fa89d450c6049ce12399fe2f92.
//
// Solidity: event ExitStarted(_priority indexed uint256)
func (_RootChain *RootChainFilterer) FilterExitStarted(opts *bind.FilterOpts, _priority []*big.Int) (*RootChainExitStartedIterator, error) {

	var _priorityRule []interface{}
	for _, _priorityItem := range _priority {
		_priorityRule = append(_priorityRule, _priorityItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "ExitStarted", _priorityRule)
	if err != nil {
		return nil, err
	}
	return &RootChainExitStartedIterator{contract: _RootChain.contract, event: "ExitStarted", logs: logs, sub: sub}, nil
}

// WatchExitStarted is a free log subscription operation binding the contract event 0x2e6c46902e4864b15d699069473ce23cfcd871fa89d450c6049ce12399fe2f92.
//
// Solidity: event ExitStarted(_priority indexed uint256)
func (_RootChain *RootChainFilterer) WatchExitStarted(opts *bind.WatchOpts, sink chan<- *RootChainExitStarted, _priority []*big.Int) (event.Subscription, error) {

	var _priorityRule []interface{}
	for _, _priorityItem := range _priority {
		_priorityRule = append(_priorityRule, _priorityItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "ExitStarted", _priorityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainExitStarted)
				if err := _RootChain.contract.UnpackLog(event, "ExitStarted", log); err != nil {
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

// RootChainExitTimeIterator is returned from FilterExitTime and is used to iterate over the raw logs and unpacked data for ExitTime events raised by the RootChain contract.
type RootChainExitTimeIterator struct {
	Event *RootChainExitTime // Event containing the contract specifics and raw log

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
func (it *RootChainExitTimeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainExitTime)
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
		it.Event = new(RootChainExitTime)
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
func (it *RootChainExitTimeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainExitTimeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainExitTime represents a ExitTime event raised by the RootChain contract.
type RootChainExitTime struct {
	ExitableTS *big.Int
	CuurrentTS *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterExitTime is a free log retrieval operation binding the contract event 0x05b05d48f45d9984a4696f7e1cb5b91fc5cd1e4420b57db43b82a9b97dccfa6a.
//
// Solidity: event ExitTime(exitableTS uint256, cuurrentTS uint256)
func (_RootChain *RootChainFilterer) FilterExitTime(opts *bind.FilterOpts) (*RootChainExitTimeIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "ExitTime")
	if err != nil {
		return nil, err
	}
	return &RootChainExitTimeIterator{contract: _RootChain.contract, event: "ExitTime", logs: logs, sub: sub}, nil
}

// WatchExitTime is a free log subscription operation binding the contract event 0x05b05d48f45d9984a4696f7e1cb5b91fc5cd1e4420b57db43b82a9b97dccfa6a.
//
// Solidity: event ExitTime(exitableTS uint256, cuurrentTS uint256)
func (_RootChain *RootChainFilterer) WatchExitTime(opts *bind.WatchOpts, sink chan<- *RootChainExitTime) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "ExitTime")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainExitTime)
				if err := _RootChain.contract.UnpackLog(event, "ExitTime", log); err != nil {
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

// RootChainFinalizedExitIterator is returned from FilterFinalizedExit and is used to iterate over the raw logs and unpacked data for FinalizedExit events raised by the RootChain contract.
type RootChainFinalizedExitIterator struct {
	Event *RootChainFinalizedExit // Event containing the contract specifics and raw log

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
func (it *RootChainFinalizedExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainFinalizedExit)
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
		it.Event = new(RootChainFinalizedExit)
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
func (it *RootChainFinalizedExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainFinalizedExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainFinalizedExit represents a FinalizedExit event raised by the RootChain contract.
type RootChainFinalizedExit struct {
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
func (_RootChain *RootChainFilterer) FilterFinalizedExit(opts *bind.FilterOpts, _depositIndex []uint64, _tokenID []uint64, _timestamp []*big.Int) (*RootChainFinalizedExitIterator, error) {

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

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "FinalizedExit", _depositIndexRule, _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return &RootChainFinalizedExitIterator{contract: _RootChain.contract, event: "FinalizedExit", logs: logs, sub: sub}, nil
}

// WatchFinalizedExit is a free log subscription operation binding the contract event 0xe109ec77e392bf51e903e88a71fb4158aac7a024428a47f4b7c9fdf334abb362.
//
// Solidity: event FinalizedExit(_exiter address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64, _timestamp indexed uint256)
func (_RootChain *RootChainFilterer) WatchFinalizedExit(opts *bind.WatchOpts, sink chan<- *RootChainFinalizedExit, _depositIndex []uint64, _tokenID []uint64, _timestamp []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "FinalizedExit", _depositIndexRule, _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainFinalizedExit)
				if err := _RootChain.contract.UnpackLog(event, "FinalizedExit", log); err != nil {
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

// RootChainPublishedBlockIterator is returned from FilterPublishedBlock and is used to iterate over the raw logs and unpacked data for PublishedBlock events raised by the RootChain contract.
type RootChainPublishedBlockIterator struct {
	Event *RootChainPublishedBlock // Event containing the contract specifics and raw log

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
func (it *RootChainPublishedBlockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainPublishedBlock)
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
		it.Event = new(RootChainPublishedBlock)
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
func (it *RootChainPublishedBlockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainPublishedBlockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainPublishedBlock represents a PublishedBlock event raised by the RootChain contract.
type RootChainPublishedBlock struct {
	RootHash            [32]byte
	Blknum              uint64
	CurrentDepositIndex uint64
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterPublishedBlock is a free log retrieval operation binding the contract event 0xafb4a42a2439df06ac89ceace483c1382c37f9574eb7d217050bc54a63a33811.
//
// Solidity: event PublishedBlock(_rootHash bytes32, _blknum indexed uint64, _currentDepositIndex indexed uint64)
func (_RootChain *RootChainFilterer) FilterPublishedBlock(opts *bind.FilterOpts, _blknum []uint64, _currentDepositIndex []uint64) (*RootChainPublishedBlockIterator, error) {

	var _blknumRule []interface{}
	for _, _blknumItem := range _blknum {
		_blknumRule = append(_blknumRule, _blknumItem)
	}
	var _currentDepositIndexRule []interface{}
	for _, _currentDepositIndexItem := range _currentDepositIndex {
		_currentDepositIndexRule = append(_currentDepositIndexRule, _currentDepositIndexItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "PublishedBlock", _blknumRule, _currentDepositIndexRule)
	if err != nil {
		return nil, err
	}
	return &RootChainPublishedBlockIterator{contract: _RootChain.contract, event: "PublishedBlock", logs: logs, sub: sub}, nil
}

// WatchPublishedBlock is a free log subscription operation binding the contract event 0xafb4a42a2439df06ac89ceace483c1382c37f9574eb7d217050bc54a63a33811.
//
// Solidity: event PublishedBlock(_rootHash bytes32, _blknum indexed uint64, _currentDepositIndex indexed uint64)
func (_RootChain *RootChainFilterer) WatchPublishedBlock(opts *bind.WatchOpts, sink chan<- *RootChainPublishedBlock, _blknum []uint64, _currentDepositIndex []uint64) (event.Subscription, error) {

	var _blknumRule []interface{}
	for _, _blknumItem := range _blknum {
		_blknumRule = append(_blknumRule, _blknumItem)
	}
	var _currentDepositIndexRule []interface{}
	for _, _currentDepositIndexItem := range _currentDepositIndex {
		_currentDepositIndexRule = append(_currentDepositIndexRule, _currentDepositIndexItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "PublishedBlock", _blknumRule, _currentDepositIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainPublishedBlock)
				if err := _RootChain.contract.UnpackLog(event, "PublishedBlock", log); err != nil {
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

// RootChainStartExitIterator is returned from FilterStartExit and is used to iterate over the raw logs and unpacked data for StartExit events raised by the RootChain contract.
type RootChainStartExitIterator struct {
	Event *RootChainStartExit // Event containing the contract specifics and raw log

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
func (it *RootChainStartExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainStartExit)
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
		it.Event = new(RootChainStartExit)
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
func (it *RootChainStartExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainStartExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainStartExit represents a StartExit event raised by the RootChain contract.
type RootChainStartExit struct {
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
func (_RootChain *RootChainFilterer) FilterStartExit(opts *bind.FilterOpts, _depositIndex []uint64, _tokenID []uint64, _timestamp []*big.Int) (*RootChainStartExitIterator, error) {

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

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "StartExit", _depositIndexRule, _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return &RootChainStartExitIterator{contract: _RootChain.contract, event: "StartExit", logs: logs, sub: sub}, nil
}

// WatchStartExit is a free log subscription operation binding the contract event 0x6b75ed42e8219410363fb31fa601dea7c8b3549bfc9428b1eeaf53a6226cb668.
//
// Solidity: event StartExit(_exiter address, _depositIndex indexed uint64, _denomination uint64, _tokenID indexed uint64, _timestamp indexed uint256)
func (_RootChain *RootChainFilterer) WatchStartExit(opts *bind.WatchOpts, sink chan<- *RootChainStartExit, _depositIndex []uint64, _tokenID []uint64, _timestamp []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "StartExit", _depositIndexRule, _tokenIDRule, _timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainStartExit)
				if err := _RootChain.contract.UnpackLog(event, "StartExit", log); err != nil {
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
