// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StorageDepositABI is the input ABI used to generate the binding from.
const StorageDepositABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"fileAddress\",\"type\":\"address\"}],\"name\":\"getStorageAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"fileAddress\",\"type\":\"address\"},{\"name\":\"fsize\",\"type\":\"uint256\"},{\"name\":\"block_nums\",\"type\":\"uint128\"}],\"name\":\"newUploadJob\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"fileAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"storageAccount\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"fsize\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"}],\"name\":\"NewUploadJob\",\"type\":\"event\"}]"

// StorageDepositBin is the compiled bytecode used for deploying new contracts.
const StorageDepositBin = `608060405234801561001057600080fd5b50610d34806100206000396000f30060806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680635023d4cd14610051578063fd5e0b08146100d4575b600080fd5b34801561005d57600080fd5b50610092600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610130565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61012e600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080356fffffffffffffffffffffffffffffffff169060200190929190505050610198565b005b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b60006101a3836103ff565b341115156101b057600080fd5b600073ffffffffffffffffffffffffffffffffffffffff166000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151561024957600080fd5b348483610254610409565b808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168152602001925050506040518091039082f0801580156102d0573d6000803e3d6000fd5b5090509050806000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507ff66ba18204f3e105f07d7014bed66a25c8a59cacca752d1874229a653b86f1ee84828534604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182815260200194505050505060405180910390a150505050565b6000809050919050565b6040516108ef8061041a83390190560060806040526040516040806108ef8339810180604052810190808051906020019092919080519060200190929190505050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600360106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555034600581905550505061082e806100c16000396000f300608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806310b0a7c91461006757806364d8aff61461011b5780636d40eb5214610132578063bb141cf4146101e1575b600080fd5b34801561007357600080fd5b5061007c610239565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156100c35780820151818401526020810190506100a8565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156101055780820151818401526020810190506100ea565b5050505090500194505050505060405180910390f35b34801561012757600080fd5b50610130610240565b005b34801561013e57600080fd5b506101df600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080356000191690602001909291908035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610400565b005b3480156101ed57600080fd5b5061020c60048036038101908080359060200190929190505050610790565b60405180836000191660001916815260200182600019166000191681526020019250505060405180910390f35b6060809091565b600080600080600460008154809291906001019190505550600a6004541115610268576103fa565b60009350600092505b600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168310156102d1576001600084815260200190815260200160002080549050840193508280600101935050610271565b836032026004600554028115156102e457fe5b049150600092505b600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168310156103f957600090505b60016000848152602001908152602001600020805490508110156103ec57600160008481526020019081526020016000208181548110151561036757fe5b906000526020600020906003020160020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f193505050501580156103de573d6000803e3d6000fd5b508080600101915050610329565b82806001019350506102ec565b5b50505050565b6000600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168510151561044057600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff1614151561049b57600080fd5b600360016000878152602001908152602001600020805490501015156104c057600080fd5b600090505b60016000868152602001908152602001600020805490508110156105be578260001916600160008781526020019081526020016000208281548110151561050857fe5b906000526020600020906003020160010154600019161480156105a757503373ffffffffffffffffffffffffffffffffffffffff16600160008781526020019081526020016000208281548110151561055d57fe5b906000526020600020906003020160020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b156105b157600080fd5b80806001019150506104c5565b6003600081819054906101000a90046fffffffffffffffffffffffffffffffff168092919060010191906101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550506001600086815260200190815260200160002060606040519081016040528086600019168152602001856000191681526020013373ffffffffffffffffffffffffffffffffffffffff1681525090806001815401808255809150509060018203906000526020600020906003020160009091929091909150600082015181600001906000191690556020820151816001019060001916905560408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050503373ffffffffffffffffffffffffffffffffffffffff166108fc60038060009054906101000a90046fffffffffffffffffffffffffffffffff16600502026fffffffffffffffffffffffffffffffff1660055481151561075b57fe5b049081150290604051600060405180830381858888f19350505050158015610787573d6000803e3d6000fd5b50505050505050565b6000806001600084815260200190815260200160002060008154811015156107b457fe5b90600052602060002090600302016000015491506001600084815260200190815260200160002060008154811015156107e957fe5b90600052602060002090600302016001015490509150915600a165627a7a72305820ab2a8e28e63788ac89b28f8843b79129a882a265f523b2a9fefd0d0465e144d80029a165627a7a72305820e93d7e56e7bf42a7de3859ff8df962af8c54ba6a613da399cc14730f90c452f90029`

// DeployStorageDeposit deploys a new Ethereum contract, binding an instance of StorageDeposit to it.
func DeployStorageDeposit(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StorageDeposit, error) {
	parsed, err := abi.JSON(strings.NewReader(StorageDepositABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StorageDepositBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StorageDeposit{StorageDepositCaller: StorageDepositCaller{contract: contract}, StorageDepositTransactor: StorageDepositTransactor{contract: contract}, StorageDepositFilterer: StorageDepositFilterer{contract: contract}}, nil
}

// StorageDeposit is an auto generated Go binding around an Ethereum contract.
type StorageDeposit struct {
	StorageDepositCaller     // Read-only binding to the contract
	StorageDepositTransactor // Write-only binding to the contract
	StorageDepositFilterer   // Log filterer for contract events
}

// StorageDepositCaller is an auto generated read-only Go binding around an Ethereum contract.
type StorageDepositCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageDepositTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StorageDepositTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageDepositFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StorageDepositFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageDepositSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StorageDepositSession struct {
	Contract     *StorageDeposit   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageDepositCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StorageDepositCallerSession struct {
	Contract *StorageDepositCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// StorageDepositTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StorageDepositTransactorSession struct {
	Contract     *StorageDepositTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// StorageDepositRaw is an auto generated low-level Go binding around an Ethereum contract.
type StorageDepositRaw struct {
	Contract *StorageDeposit // Generic contract binding to access the raw methods on
}

// StorageDepositCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StorageDepositCallerRaw struct {
	Contract *StorageDepositCaller // Generic read-only contract binding to access the raw methods on
}

// StorageDepositTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StorageDepositTransactorRaw struct {
	Contract *StorageDepositTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageDeposit creates a new instance of StorageDeposit, bound to a specific deployed contract.
func NewStorageDeposit(address common.Address, backend bind.ContractBackend) (*StorageDeposit, error) {
	contract, err := bindStorageDeposit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageDeposit{StorageDepositCaller: StorageDepositCaller{contract: contract}, StorageDepositTransactor: StorageDepositTransactor{contract: contract}, StorageDepositFilterer: StorageDepositFilterer{contract: contract}}, nil
}

// NewStorageDepositCaller creates a new read-only instance of StorageDeposit, bound to a specific deployed contract.
func NewStorageDepositCaller(address common.Address, caller bind.ContractCaller) (*StorageDepositCaller, error) {
	contract, err := bindStorageDeposit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageDepositCaller{contract: contract}, nil
}

// NewStorageDepositTransactor creates a new write-only instance of StorageDeposit, bound to a specific deployed contract.
func NewStorageDepositTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageDepositTransactor, error) {
	contract, err := bindStorageDeposit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageDepositTransactor{contract: contract}, nil
}

// NewStorageDepositFilterer creates a new log filterer instance of StorageDeposit, bound to a specific deployed contract.
func NewStorageDepositFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageDepositFilterer, error) {
	contract, err := bindStorageDeposit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageDepositFilterer{contract: contract}, nil
}

// bindStorageDeposit binds a generic wrapper to an already deployed contract.
func bindStorageDeposit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StorageDepositABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageDeposit *StorageDepositRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StorageDeposit.Contract.StorageDepositCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageDeposit *StorageDepositRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageDeposit.Contract.StorageDepositTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageDeposit *StorageDepositRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageDeposit.Contract.StorageDepositTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageDeposit *StorageDepositCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StorageDeposit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageDeposit *StorageDepositTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageDeposit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageDeposit *StorageDepositTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageDeposit.Contract.contract.Transact(opts, method, params...)
}

// GetStorageAccount is a free data retrieval call binding the contract method 0x5023d4cd.
//
// Solidity: function getStorageAccount(address fileAddress) constant returns(address)
func (_StorageDeposit *StorageDepositCaller) GetStorageAccount(opts *bind.CallOpts, fileAddress common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StorageDeposit.contract.Call(opts, out, "getStorageAccount", fileAddress)
	return *ret0, err
}

// GetStorageAccount is a free data retrieval call binding the contract method 0x5023d4cd.
//
// Solidity: function getStorageAccount(address fileAddress) constant returns(address)
func (_StorageDeposit *StorageDepositSession) GetStorageAccount(fileAddress common.Address) (common.Address, error) {
	return _StorageDeposit.Contract.GetStorageAccount(&_StorageDeposit.CallOpts, fileAddress)
}

// GetStorageAccount is a free data retrieval call binding the contract method 0x5023d4cd.
//
// Solidity: function getStorageAccount(address fileAddress) constant returns(address)
func (_StorageDeposit *StorageDepositCallerSession) GetStorageAccount(fileAddress common.Address) (common.Address, error) {
	return _StorageDeposit.Contract.GetStorageAccount(&_StorageDeposit.CallOpts, fileAddress)
}

// NewUploadJob is a paid mutator transaction binding the contract method 0xfd5e0b08.
//
// Solidity: function newUploadJob(address fileAddress, uint256 fsize, uint128 block_nums) returns()
func (_StorageDeposit *StorageDepositTransactor) NewUploadJob(opts *bind.TransactOpts, fileAddress common.Address, fsize *big.Int, block_nums *big.Int) (*types.Transaction, error) {
	return _StorageDeposit.contract.Transact(opts, "newUploadJob", fileAddress, fsize, block_nums)
}

// NewUploadJob is a paid mutator transaction binding the contract method 0xfd5e0b08.
//
// Solidity: function newUploadJob(address fileAddress, uint256 fsize, uint128 block_nums) returns()
func (_StorageDeposit *StorageDepositSession) NewUploadJob(fileAddress common.Address, fsize *big.Int, block_nums *big.Int) (*types.Transaction, error) {
	return _StorageDeposit.Contract.NewUploadJob(&_StorageDeposit.TransactOpts, fileAddress, fsize, block_nums)
}

// NewUploadJob is a paid mutator transaction binding the contract method 0xfd5e0b08.
//
// Solidity: function newUploadJob(address fileAddress, uint256 fsize, uint128 block_nums) returns()
func (_StorageDeposit *StorageDepositTransactorSession) NewUploadJob(fileAddress common.Address, fsize *big.Int, block_nums *big.Int) (*types.Transaction, error) {
	return _StorageDeposit.Contract.NewUploadJob(&_StorageDeposit.TransactOpts, fileAddress, fsize, block_nums)
}

// StorageDepositNewUploadJobIterator is returned from FilterNewUploadJob and is used to iterate over the raw logs and unpacked data for NewUploadJob events raised by the StorageDeposit contract.
type StorageDepositNewUploadJobIterator struct {
	Event *StorageDepositNewUploadJob // Event containing the contract specifics and raw log

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
func (it *StorageDepositNewUploadJobIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StorageDepositNewUploadJob)
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
		it.Event = new(StorageDepositNewUploadJob)
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
func (it *StorageDepositNewUploadJobIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StorageDepositNewUploadJobIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StorageDepositNewUploadJob represents a NewUploadJob event raised by the StorageDeposit contract.
type StorageDepositNewUploadJob struct {
	FileAddress    common.Address
	StorageAccount common.Address
	Fsize          *big.Int
	Deposit        *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterNewUploadJob is a free log retrieval operation binding the contract event 0xf66ba18204f3e105f07d7014bed66a25c8a59cacca752d1874229a653b86f1ee.
//
// Solidity: event NewUploadJob(address fileAddress, address storageAccount, uint256 fsize, uint256 deposit)
func (_StorageDeposit *StorageDepositFilterer) FilterNewUploadJob(opts *bind.FilterOpts) (*StorageDepositNewUploadJobIterator, error) {

	logs, sub, err := _StorageDeposit.contract.FilterLogs(opts, "NewUploadJob")
	if err != nil {
		return nil, err
	}
	return &StorageDepositNewUploadJobIterator{contract: _StorageDeposit.contract, event: "NewUploadJob", logs: logs, sub: sub}, nil
}

// WatchNewUploadJob is a free log subscription operation binding the contract event 0xf66ba18204f3e105f07d7014bed66a25c8a59cacca752d1874229a653b86f1ee.
//
// Solidity: event NewUploadJob(address fileAddress, address storageAccount, uint256 fsize, uint256 deposit)
func (_StorageDeposit *StorageDepositFilterer) WatchNewUploadJob(opts *bind.WatchOpts, sink chan<- *StorageDepositNewUploadJob) (event.Subscription, error) {

	logs, sub, err := _StorageDeposit.contract.WatchLogs(opts, "NewUploadJob")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StorageDepositNewUploadJob)
				if err := _StorageDeposit.contract.UnpackLog(event, "NewUploadJob", log); err != nil {
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
