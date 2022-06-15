// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package weth9

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
)

// Weth9MetaData contains all meta data concerning the Weth9 contract.
var Weth9MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"guy\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040526040518060400160405280600d81526020017f57726170706564204574686572000000000000000000000000000000000000008152506000908051906020019062000051929190620000d0565b506040518060400160405280600481526020017f5745544800000000000000000000000000000000000000000000000000000000815250600190805190602001906200009f929190620000d0565b506012600260006101000a81548160ff021916908360ff160217905550348015620000c957600080fd5b50620001e5565b828054620000de9062000180565b90600052602060002090601f0160209004810192826200010257600085556200014e565b82601f106200011d57805160ff19168380011785556200014e565b828001600101855582156200014e579182015b828111156200014d57825182559160200191906001019062000130565b5b5090506200015d919062000161565b5090565b5b808211156200017c57600081600090555060010162000162565b5090565b600060028204905060018216806200019957607f821691505b60208210811415620001b057620001af620001b6565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b610ede80620001f56000396000f3fe6080604052600436106100a05760003560e01c8063313ce56711610064578063313ce567146101a057806370a08231146101cb57806395d89b4114610208578063a9059cbb14610233578063d0e30db014610270578063dd62ed3e1461027a576100af565b806306fdde03146100b4578063095ea7b3146100df57806318160ddd1461011c57806323b872dd146101475780632e1a7d4d14610184576100af565b366100af576100ad6102b7565b005b600080fd5b3480156100c057600080fd5b506100c961035d565b6040516100d69190610c4e565b60405180910390f35b3480156100eb57600080fd5b5061010660048036038101906101019190610b60565b6103eb565b6040516101139190610c33565b60405180910390f35b34801561012857600080fd5b506101316104dd565b60405161013e9190610c70565b60405180910390f35b34801561015357600080fd5b5061016e60048036038101906101699190610b0d565b6104e5565b60405161017b9190610c33565b60405180910390f35b61019e60048036038101906101999190610ba0565b610849565b005b3480156101ac57600080fd5b506101b5610983565b6040516101c29190610c8b565b60405180910390f35b3480156101d757600080fd5b506101f260048036038101906101ed9190610aa0565b610996565b6040516101ff9190610c70565b60405180910390f35b34801561021457600080fd5b5061021d6109ae565b60405161022a9190610c4e565b60405180910390f35b34801561023f57600080fd5b5061025a60048036038101906102559190610b60565b610a3c565b6040516102679190610c33565b60405180910390f35b6102786102b7565b005b34801561028657600080fd5b506102a1600480360381019061029c9190610acd565b610a51565b6040516102ae9190610c70565b60405180910390f35b34600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546103069190610cc2565b925050819055503373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c346040516103539190610c70565b60405180910390a2565b6000805461036a90610dd4565b80601f016020809104026020016040519081016040528092919081815260200182805461039690610dd4565b80156103e35780601f106103b8576101008083540402835291602001916103e3565b820191906000526020600020905b8154815290600101906020018083116103c657829003601f168201915b505050505081565b600081600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516104cb9190610c70565b60405180910390a36001905092915050565b600047905090565b600081600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561053357600080fd5b3373ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415801561060b57507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414155b1561072d5781600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561069957600080fd5b81600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546107259190610d18565b925050819055505b81600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461077c9190610d18565b9250508190555081600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546107d29190610cc2565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516108369190610c70565b60405180910390a3600190509392505050565b80600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561089557600080fd5b80600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546108e49190610d18565b925050819055503373ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610931573d6000803e3d6000fd5b503373ffffffffffffffffffffffffffffffffffffffff167f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65826040516109789190610c70565b60405180910390a250565b600260009054906101000a900460ff1681565b60036020528060005260406000206000915090505481565b600180546109bb90610dd4565b80601f01602080910402602001604051908101604052809291908181526020018280546109e790610dd4565b8015610a345780601f10610a0957610100808354040283529160200191610a34565b820191906000526020600020905b815481529060010190602001808311610a1757829003601f168201915b505050505081565b6000610a493384846104e5565b905092915050565b6004602052816000526040600020602052806000526040600020600091509150505481565b600081359050610a8581610e7a565b92915050565b600081359050610a9a81610e91565b92915050565b600060208284031215610ab657610ab5610e64565b5b6000610ac484828501610a76565b91505092915050565b60008060408385031215610ae457610ae3610e64565b5b6000610af285828601610a76565b9250506020610b0385828601610a76565b9150509250929050565b600080600060608486031215610b2657610b25610e64565b5b6000610b3486828701610a76565b9350506020610b4586828701610a76565b9250506040610b5686828701610a8b565b9150509250925092565b60008060408385031215610b7757610b76610e64565b5b6000610b8585828601610a76565b9250506020610b9685828601610a8b565b9150509250929050565b600060208284031215610bb657610bb5610e64565b5b6000610bc484828501610a8b565b91505092915050565b610bd681610d5e565b82525050565b6000610be782610ca6565b610bf18185610cb1565b9350610c01818560208601610da1565b610c0a81610e69565b840191505092915050565b610c1e81610d8a565b82525050565b610c2d81610d94565b82525050565b6000602082019050610c486000830184610bcd565b92915050565b60006020820190508181036000830152610c688184610bdc565b905092915050565b6000602082019050610c856000830184610c15565b92915050565b6000602082019050610ca06000830184610c24565b92915050565b600081519050919050565b600082825260208201905092915050565b6000610ccd82610d8a565b9150610cd883610d8a565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115610d0d57610d0c610e06565b5b828201905092915050565b6000610d2382610d8a565b9150610d2e83610d8a565b925082821015610d4157610d40610e06565b5b828203905092915050565b6000610d5782610d6a565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600060ff82169050919050565b60005b83811015610dbf578082015181840152602081019050610da4565b83811115610dce576000848401525b50505050565b60006002820490506001821680610dec57607f821691505b60208210811415610e0057610dff610e35565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600080fd5b6000601f19601f8301169050919050565b610e8381610d4c565b8114610e8e57600080fd5b50565b610e9a81610d8a565b8114610ea557600080fd5b5056fea2646970667358221220767df344a75c531737e8ccce0754840a667a8cafcb8795f0a846ba7794f0eeed64736f6c63430008070033",
}

// Weth9ABI is the input ABI used to generate the binding from.
// Deprecated: Use Weth9MetaData.ABI instead.
var Weth9ABI = Weth9MetaData.ABI

// Weth9Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Weth9MetaData.Bin instead.
var Weth9Bin = Weth9MetaData.Bin

// DeployWeth9 deploys a new Ethereum contract, binding an instance of Weth9 to it.
func DeployWeth9(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Weth9, error) {
	parsed, err := Weth9MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Weth9Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Weth9{Weth9Caller: Weth9Caller{contract: contract}, Weth9Transactor: Weth9Transactor{contract: contract}, Weth9Filterer: Weth9Filterer{contract: contract}}, nil
}

// Weth9 is an auto generated Go binding around an Ethereum contract.
type Weth9 struct {
	Weth9Caller     // Read-only binding to the contract
	Weth9Transactor // Write-only binding to the contract
	Weth9Filterer   // Log filterer for contract events
}

// Weth9Caller is an auto generated read-only Go binding around an Ethereum contract.
type Weth9Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Weth9Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Weth9Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Weth9Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Weth9Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Weth9Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Weth9Session struct {
	Contract     *Weth9            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Weth9CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Weth9CallerSession struct {
	Contract *Weth9Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Weth9TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Weth9TransactorSession struct {
	Contract     *Weth9Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Weth9Raw is an auto generated low-level Go binding around an Ethereum contract.
type Weth9Raw struct {
	Contract *Weth9 // Generic contract binding to access the raw methods on
}

// Weth9CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Weth9CallerRaw struct {
	Contract *Weth9Caller // Generic read-only contract binding to access the raw methods on
}

// Weth9TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Weth9TransactorRaw struct {
	Contract *Weth9Transactor // Generic write-only contract binding to access the raw methods on
}

// NewWeth9 creates a new instance of Weth9, bound to a specific deployed contract.
func NewWeth9(address common.Address, backend bind.ContractBackend) (*Weth9, error) {
	contract, err := bindWeth9(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Weth9{Weth9Caller: Weth9Caller{contract: contract}, Weth9Transactor: Weth9Transactor{contract: contract}, Weth9Filterer: Weth9Filterer{contract: contract}}, nil
}

// NewWeth9Caller creates a new read-only instance of Weth9, bound to a specific deployed contract.
func NewWeth9Caller(address common.Address, caller bind.ContractCaller) (*Weth9Caller, error) {
	contract, err := bindWeth9(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Weth9Caller{contract: contract}, nil
}

// NewWeth9Transactor creates a new write-only instance of Weth9, bound to a specific deployed contract.
func NewWeth9Transactor(address common.Address, transactor bind.ContractTransactor) (*Weth9Transactor, error) {
	contract, err := bindWeth9(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Weth9Transactor{contract: contract}, nil
}

// NewWeth9Filterer creates a new log filterer instance of Weth9, bound to a specific deployed contract.
func NewWeth9Filterer(address common.Address, filterer bind.ContractFilterer) (*Weth9Filterer, error) {
	contract, err := bindWeth9(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Weth9Filterer{contract: contract}, nil
}

// bindWeth9 binds a generic wrapper to an already deployed contract.
func bindWeth9(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Weth9ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Weth9 *Weth9Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Weth9.Contract.Weth9Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Weth9 *Weth9Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Weth9.Contract.Weth9Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Weth9 *Weth9Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Weth9.Contract.Weth9Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Weth9 *Weth9CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Weth9.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Weth9 *Weth9TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Weth9.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Weth9 *Weth9TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Weth9.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_Weth9 *Weth9Caller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Weth9.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_Weth9 *Weth9Session) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Weth9.Contract.Allowance(&_Weth9.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_Weth9 *Weth9CallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Weth9.Contract.Allowance(&_Weth9.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_Weth9 *Weth9Caller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Weth9.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_Weth9 *Weth9Session) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _Weth9.Contract.BalanceOf(&_Weth9.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_Weth9 *Weth9CallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _Weth9.Contract.BalanceOf(&_Weth9.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Weth9 *Weth9Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Weth9.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Weth9 *Weth9Session) Decimals() (uint8, error) {
	return _Weth9.Contract.Decimals(&_Weth9.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Weth9 *Weth9CallerSession) Decimals() (uint8, error) {
	return _Weth9.Contract.Decimals(&_Weth9.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Weth9 *Weth9Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Weth9.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Weth9 *Weth9Session) Name() (string, error) {
	return _Weth9.Contract.Name(&_Weth9.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Weth9 *Weth9CallerSession) Name() (string, error) {
	return _Weth9.Contract.Name(&_Weth9.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Weth9 *Weth9Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Weth9.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Weth9 *Weth9Session) Symbol() (string, error) {
	return _Weth9.Contract.Symbol(&_Weth9.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Weth9 *Weth9CallerSession) Symbol() (string, error) {
	return _Weth9.Contract.Symbol(&_Weth9.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Weth9 *Weth9Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Weth9.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Weth9 *Weth9Session) TotalSupply() (*big.Int, error) {
	return _Weth9.Contract.TotalSupply(&_Weth9.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Weth9 *Weth9CallerSession) TotalSupply() (*big.Int, error) {
	return _Weth9.Contract.TotalSupply(&_Weth9.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_Weth9 *Weth9Transactor) Approve(opts *bind.TransactOpts, guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.contract.Transact(opts, "approve", guy, wad)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_Weth9 *Weth9Session) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.Contract.Approve(&_Weth9.TransactOpts, guy, wad)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_Weth9 *Weth9TransactorSession) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.Contract.Approve(&_Weth9.TransactOpts, guy, wad)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Weth9 *Weth9Transactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Weth9.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Weth9 *Weth9Session) Deposit() (*types.Transaction, error) {
	return _Weth9.Contract.Deposit(&_Weth9.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Weth9 *Weth9TransactorSession) Deposit() (*types.Transaction, error) {
	return _Weth9.Contract.Deposit(&_Weth9.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_Weth9 *Weth9Transactor) Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.contract.Transact(opts, "transfer", dst, wad)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_Weth9 *Weth9Session) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.Contract.Transfer(&_Weth9.TransactOpts, dst, wad)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_Weth9 *Weth9TransactorSession) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.Contract.Transfer(&_Weth9.TransactOpts, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_Weth9 *Weth9Transactor) TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.contract.Transact(opts, "transferFrom", src, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_Weth9 *Weth9Session) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.Contract.TransferFrom(&_Weth9.TransactOpts, src, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_Weth9 *Weth9TransactorSession) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.Contract.TransferFrom(&_Weth9.TransactOpts, src, dst, wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) payable returns()
func (_Weth9 *Weth9Transactor) Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _Weth9.contract.Transact(opts, "withdraw", wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) payable returns()
func (_Weth9 *Weth9Session) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _Weth9.Contract.Withdraw(&_Weth9.TransactOpts, wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) payable returns()
func (_Weth9 *Weth9TransactorSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _Weth9.Contract.Withdraw(&_Weth9.TransactOpts, wad)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Weth9 *Weth9Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Weth9.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Weth9 *Weth9Session) Receive() (*types.Transaction, error) {
	return _Weth9.Contract.Receive(&_Weth9.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Weth9 *Weth9TransactorSession) Receive() (*types.Transaction, error) {
	return _Weth9.Contract.Receive(&_Weth9.TransactOpts)
}

// Weth9ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Weth9 contract.
type Weth9ApprovalIterator struct {
	Event *Weth9Approval // Event containing the contract specifics and raw log

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
func (it *Weth9ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Weth9Approval)
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
		it.Event = new(Weth9Approval)
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
func (it *Weth9ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Weth9ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Weth9Approval represents a Approval event raised by the Weth9 contract.
type Weth9Approval struct {
	Src common.Address
	Guy common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_Weth9 *Weth9Filterer) FilterApproval(opts *bind.FilterOpts, src []common.Address, guy []common.Address) (*Weth9ApprovalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _Weth9.contract.FilterLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return &Weth9ApprovalIterator{contract: _Weth9.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_Weth9 *Weth9Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Weth9Approval, src []common.Address, guy []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _Weth9.contract.WatchLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Weth9Approval)
				if err := _Weth9.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_Weth9 *Weth9Filterer) ParseApproval(log types.Log) (*Weth9Approval, error) {
	event := new(Weth9Approval)
	if err := _Weth9.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Weth9DepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Weth9 contract.
type Weth9DepositIterator struct {
	Event *Weth9Deposit // Event containing the contract specifics and raw log

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
func (it *Weth9DepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Weth9Deposit)
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
		it.Event = new(Weth9Deposit)
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
func (it *Weth9DepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Weth9DepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Weth9Deposit represents a Deposit event raised by the Weth9 contract.
type Weth9Deposit struct {
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_Weth9 *Weth9Filterer) FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*Weth9DepositIterator, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _Weth9.contract.FilterLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return &Weth9DepositIterator{contract: _Weth9.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_Weth9 *Weth9Filterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *Weth9Deposit, dst []common.Address) (event.Subscription, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _Weth9.contract.WatchLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Weth9Deposit)
				if err := _Weth9.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_Weth9 *Weth9Filterer) ParseDeposit(log types.Log) (*Weth9Deposit, error) {
	event := new(Weth9Deposit)
	if err := _Weth9.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Weth9TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Weth9 contract.
type Weth9TransferIterator struct {
	Event *Weth9Transfer // Event containing the contract specifics and raw log

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
func (it *Weth9TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Weth9Transfer)
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
		it.Event = new(Weth9Transfer)
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
func (it *Weth9TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Weth9TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Weth9Transfer represents a Transfer event raised by the Weth9 contract.
type Weth9Transfer struct {
	Src common.Address
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_Weth9 *Weth9Filterer) FilterTransfer(opts *bind.FilterOpts, src []common.Address, dst []common.Address) (*Weth9TransferIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _Weth9.contract.FilterLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return &Weth9TransferIterator{contract: _Weth9.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_Weth9 *Weth9Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Weth9Transfer, src []common.Address, dst []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _Weth9.contract.WatchLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Weth9Transfer)
				if err := _Weth9.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_Weth9 *Weth9Filterer) ParseTransfer(log types.Log) (*Weth9Transfer, error) {
	event := new(Weth9Transfer)
	if err := _Weth9.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Weth9WithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the Weth9 contract.
type Weth9WithdrawalIterator struct {
	Event *Weth9Withdrawal // Event containing the contract specifics and raw log

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
func (it *Weth9WithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Weth9Withdrawal)
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
		it.Event = new(Weth9Withdrawal)
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
func (it *Weth9WithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Weth9WithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Weth9Withdrawal represents a Withdrawal event raised by the Weth9 contract.
type Weth9Withdrawal struct {
	Src common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_Weth9 *Weth9Filterer) FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*Weth9WithdrawalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _Weth9.contract.FilterLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return &Weth9WithdrawalIterator{contract: _Weth9.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_Weth9 *Weth9Filterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *Weth9Withdrawal, src []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _Weth9.contract.WatchLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Weth9Withdrawal)
				if err := _Weth9.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_Weth9 *Weth9Filterer) ParseWithdrawal(log types.Log) (*Weth9Withdrawal, error) {
	event := new(Weth9Withdrawal)
	if err := _Weth9.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
