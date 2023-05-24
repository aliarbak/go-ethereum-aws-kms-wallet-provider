// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bank_account

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

// BankAccountMetaData contains all meta data concerning the BankAccount contract.
var BankAccountMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"accountHolder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600080546001600160a01b0319163317905561068b806100326000396000f3fe60806040526004361061003f5760003560e01c8063030ba25d14610044578063ad7a672f14610066578063be5460791461008f578063d0e30db0146100c7575b600080fd5b34801561005057600080fd5b5061006461005f366004610548565b6100cf565b005b34801561007257600080fd5b5061007c60015481565b6040519081526020015b60405180910390f35b34801561009b57600080fd5b506000546100af906001600160a01b031681565b6040516001600160a01b039091168152602001610086565b61006461029d565b6040516bffffffffffffffffffffffff1930606090811b8216602084015233901b166034820152600090604801604051602081830303815290604052805190602001209050600061015783610151847f19457468657265756d205369676e6564204d6573736167653a0a3332000000006000908152601c91909152603c902090565b906102b6565b6000549091506001600160a01b038083169116146101ad5760405162461bcd60e51b815260206004820152600e60248201526d24b73b30b634b21039b4b3b732b960911b60448201526064015b60405180910390fd5b8360015410156101f55760405162461bcd60e51b8152602060048201526013602482015272496e737566666963656e742062616c616e636560681b60448201526064016101a4565b83600160008282546102079190610619565b9091555050604051600090339086908381818185875af1925050503d806000811461024e576040519150601f19603f3d011682016040523d82523d6000602084013e610253565b606091505b50509050806102965760405162461bcd60e51b815260206004820152600f60248201526e15da5d1a191c985dc819985a5b1959608a1b60448201526064016101a4565b5050505050565b34600160008282546102af919061062c565b9091555050565b60008060006102c585856102dc565b915091506102d281610321565b5090505b92915050565b60008082516041036103125760208301516040840151606085015160001a6103068782858561046e565b9450945050505061031a565b506000905060025b9250929050565b60008160048111156103355761033561063f565b0361033d5750565b60018160048111156103515761035161063f565b0361039e5760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e6174757265000000000000000060448201526064016101a4565b60028160048111156103b2576103b261063f565b036103ff5760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016101a4565b60038160048111156104135761041361063f565b0361046b5760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b60648201526084016101a4565b50565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08311156104a55750600090506003610529565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa1580156104f9573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b03811661052257600060019250925050610529565b9150600090505b94509492505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561055b57600080fd5b82359150602083013567ffffffffffffffff8082111561057a57600080fd5b818501915085601f83011261058e57600080fd5b8135818111156105a0576105a0610532565b604051601f8201601f19908116603f011681019083821181831017156105c8576105c8610532565b816040528281528860208487010111156105e157600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b634e487b7160e01b600052601160045260246000fd5b818103818111156102d6576102d6610603565b808201808211156102d6576102d6610603565b634e487b7160e01b600052602160045260246000fdfea2646970667358221220a73de4cbb8a5e27177307ac6e686d6feae6c5831842443aa6662a0d5dd6ec1aa64736f6c63430008110033",
}

// BankAccountABI is the input ABI used to generate the binding from.
// Deprecated: Use BankAccountMetaData.ABI instead.
var BankAccountABI = BankAccountMetaData.ABI

// BankAccountBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BankAccountMetaData.Bin instead.
var BankAccountBin = BankAccountMetaData.Bin

// DeployBankAccount deploys a new Ethereum contract, binding an instance of BankAccount to it.
func DeployBankAccount(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BankAccount, error) {
	parsed, err := BankAccountMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BankAccountBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BankAccount{BankAccountCaller: BankAccountCaller{contract: contract}, BankAccountTransactor: BankAccountTransactor{contract: contract}, BankAccountFilterer: BankAccountFilterer{contract: contract}}, nil
}

// BankAccount is an auto generated Go binding around an Ethereum contract.
type BankAccount struct {
	BankAccountCaller     // Read-only binding to the contract
	BankAccountTransactor // Write-only binding to the contract
	BankAccountFilterer   // Log filterer for contract events
}

// BankAccountCaller is an auto generated read-only Go binding around an Ethereum contract.
type BankAccountCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BankAccountTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BankAccountTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BankAccountFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BankAccountFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BankAccountSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BankAccountSession struct {
	Contract     *BankAccount      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BankAccountCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BankAccountCallerSession struct {
	Contract *BankAccountCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BankAccountTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BankAccountTransactorSession struct {
	Contract     *BankAccountTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BankAccountRaw is an auto generated low-level Go binding around an Ethereum contract.
type BankAccountRaw struct {
	Contract *BankAccount // Generic contract binding to access the raw methods on
}

// BankAccountCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BankAccountCallerRaw struct {
	Contract *BankAccountCaller // Generic read-only contract binding to access the raw methods on
}

// BankAccountTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BankAccountTransactorRaw struct {
	Contract *BankAccountTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBankAccount creates a new instance of BankAccount, bound to a specific deployed contract.
func NewBankAccount(address common.Address, backend bind.ContractBackend) (*BankAccount, error) {
	contract, err := bindBankAccount(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BankAccount{BankAccountCaller: BankAccountCaller{contract: contract}, BankAccountTransactor: BankAccountTransactor{contract: contract}, BankAccountFilterer: BankAccountFilterer{contract: contract}}, nil
}

// NewBankAccountCaller creates a new read-only instance of BankAccount, bound to a specific deployed contract.
func NewBankAccountCaller(address common.Address, caller bind.ContractCaller) (*BankAccountCaller, error) {
	contract, err := bindBankAccount(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BankAccountCaller{contract: contract}, nil
}

// NewBankAccountTransactor creates a new write-only instance of BankAccount, bound to a specific deployed contract.
func NewBankAccountTransactor(address common.Address, transactor bind.ContractTransactor) (*BankAccountTransactor, error) {
	contract, err := bindBankAccount(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BankAccountTransactor{contract: contract}, nil
}

// NewBankAccountFilterer creates a new log filterer instance of BankAccount, bound to a specific deployed contract.
func NewBankAccountFilterer(address common.Address, filterer bind.ContractFilterer) (*BankAccountFilterer, error) {
	contract, err := bindBankAccount(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BankAccountFilterer{contract: contract}, nil
}

// bindBankAccount binds a generic wrapper to an already deployed contract.
func bindBankAccount(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BankAccountABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BankAccount *BankAccountRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BankAccount.Contract.BankAccountCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BankAccount *BankAccountRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BankAccount.Contract.BankAccountTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BankAccount *BankAccountRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BankAccount.Contract.BankAccountTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BankAccount *BankAccountCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BankAccount.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BankAccount *BankAccountTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BankAccount.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BankAccount *BankAccountTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BankAccount.Contract.contract.Transact(opts, method, params...)
}

// AccountHolder is a free data retrieval call binding the contract method 0xbe546079.
//
// Solidity: function accountHolder() view returns(address)
func (_BankAccount *BankAccountCaller) AccountHolder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BankAccount.contract.Call(opts, &out, "accountHolder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccountHolder is a free data retrieval call binding the contract method 0xbe546079.
//
// Solidity: function accountHolder() view returns(address)
func (_BankAccount *BankAccountSession) AccountHolder() (common.Address, error) {
	return _BankAccount.Contract.AccountHolder(&_BankAccount.CallOpts)
}

// AccountHolder is a free data retrieval call binding the contract method 0xbe546079.
//
// Solidity: function accountHolder() view returns(address)
func (_BankAccount *BankAccountCallerSession) AccountHolder() (common.Address, error) {
	return _BankAccount.Contract.AccountHolder(&_BankAccount.CallOpts)
}

// TotalBalance is a free data retrieval call binding the contract method 0xad7a672f.
//
// Solidity: function totalBalance() view returns(uint256)
func (_BankAccount *BankAccountCaller) TotalBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BankAccount.contract.Call(opts, &out, "totalBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalBalance is a free data retrieval call binding the contract method 0xad7a672f.
//
// Solidity: function totalBalance() view returns(uint256)
func (_BankAccount *BankAccountSession) TotalBalance() (*big.Int, error) {
	return _BankAccount.Contract.TotalBalance(&_BankAccount.CallOpts)
}

// TotalBalance is a free data retrieval call binding the contract method 0xad7a672f.
//
// Solidity: function totalBalance() view returns(uint256)
func (_BankAccount *BankAccountCallerSession) TotalBalance() (*big.Int, error) {
	return _BankAccount.Contract.TotalBalance(&_BankAccount.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_BankAccount *BankAccountTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BankAccount.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_BankAccount *BankAccountSession) Deposit() (*types.Transaction, error) {
	return _BankAccount.Contract.Deposit(&_BankAccount.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_BankAccount *BankAccountTransactorSession) Deposit() (*types.Transaction, error) {
	return _BankAccount.Contract.Deposit(&_BankAccount.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x030ba25d.
//
// Solidity: function withdraw(uint256 amount, bytes signature) returns()
func (_BankAccount *BankAccountTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int, signature []byte) (*types.Transaction, error) {
	return _BankAccount.contract.Transact(opts, "withdraw", amount, signature)
}

// Withdraw is a paid mutator transaction binding the contract method 0x030ba25d.
//
// Solidity: function withdraw(uint256 amount, bytes signature) returns()
func (_BankAccount *BankAccountSession) Withdraw(amount *big.Int, signature []byte) (*types.Transaction, error) {
	return _BankAccount.Contract.Withdraw(&_BankAccount.TransactOpts, amount, signature)
}

// Withdraw is a paid mutator transaction binding the contract method 0x030ba25d.
//
// Solidity: function withdraw(uint256 amount, bytes signature) returns()
func (_BankAccount *BankAccountTransactorSession) Withdraw(amount *big.Int, signature []byte) (*types.Transaction, error) {
	return _BankAccount.Contract.Withdraw(&_BankAccount.TransactOpts, amount, signature)
}
