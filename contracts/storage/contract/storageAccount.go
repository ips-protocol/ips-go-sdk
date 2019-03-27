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

// StorageAccountABI is the input ABI used to generate the binding from.
const StorageAccountABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getAllBlocksInfo\",\"outputs\":[{\"name\":\"blocksHash\",\"type\":\"bytes32[]\"},{\"name\":\"peersInfo\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"DownloadSuccess\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_fileAddress\",\"type\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\"},{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"peerInfo\",\"type\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"string\"}],\"name\":\"CommitBlockInfo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getBlockInfo\",\"outputs\":[{\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"name\":\"peerInfo\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_fileAddress\",\"type\":\"address\"},{\"name\":\"_block_nums\",\"type\":\"uint128\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"}]"

// StorageAccountBin is the compiled bytecode used for deploying new contracts.
const StorageAccountBin = `0x60806040526040516040806108ef8339810180604052810190808051906020019092919080519060200190929190505050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600360106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555034600581905550505061082e806100c16000396000f300608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806310b0a7c91461006757806364d8aff61461011b5780636d40eb5214610132578063bb141cf4146101e1575b600080fd5b34801561007357600080fd5b5061007c610239565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156100c35780820151818401526020810190506100a8565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156101055780820151818401526020810190506100ea565b5050505090500194505050505060405180910390f35b34801561012757600080fd5b50610130610240565b005b34801561013e57600080fd5b506101df600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080356000191690602001909291908035600019169060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610400565b005b3480156101ed57600080fd5b5061020c60048036038101908080359060200190929190505050610790565b60405180836000191660001916815260200182600019166000191681526020019250505060405180910390f35b6060809091565b600080600080600460008154809291906001019190505550600a6004541115610268576103fa565b60009350600092505b600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168310156102d1576001600084815260200190815260200160002080549050840193508280600101935050610271565b836032026004600554028115156102e457fe5b049150600092505b600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168310156103f957600090505b60016000848152602001908152602001600020805490508110156103ec57600160008481526020019081526020016000208181548110151561036757fe5b906000526020600020906003020160020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f193505050501580156103de573d6000803e3d6000fd5b508080600101915050610329565b82806001019350506102ec565b5b50505050565b6000600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168510151561044057600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff1614151561049b57600080fd5b600360016000878152602001908152602001600020805490501015156104c057600080fd5b600090505b60016000868152602001908152602001600020805490508110156105be578260001916600160008781526020019081526020016000208281548110151561050857fe5b906000526020600020906003020160010154600019161480156105a757503373ffffffffffffffffffffffffffffffffffffffff16600160008781526020019081526020016000208281548110151561055d57fe5b906000526020600020906003020160020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b156105b157600080fd5b80806001019150506104c5565b6003600081819054906101000a90046fffffffffffffffffffffffffffffffff168092919060010191906101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550506001600086815260200190815260200160002060606040519081016040528086600019168152602001856000191681526020013373ffffffffffffffffffffffffffffffffffffffff1681525090806001815401808255809150509060018203906000526020600020906003020160009091929091909150600082015181600001906000191690556020820151816001019060001916905560408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050503373ffffffffffffffffffffffffffffffffffffffff166108fc60038060009054906101000a90046fffffffffffffffffffffffffffffffff16600502026fffffffffffffffffffffffffffffffff1660055481151561075b57fe5b049081150290604051600060405180830381858888f19350505050158015610787573d6000803e3d6000fd5b50505050505050565b6000806001600084815260200190815260200160002060008154811015156107b457fe5b90600052602060002090600302016000015491506001600084815260200190815260200160002060008154811015156107e957fe5b90600052602060002090600302016001015490509150915600a165627a7a723058209df59819afdf3634b303084d1a73afb4c29ea35a20cd2fd31428251eb80a33770029`

// DeployStorageAccount deploys a new Ethereum contract, binding an instance of StorageAccount to it.
func DeployStorageAccount(auth *bind.TransactOpts, backend bind.ContractBackend, _fileAddress common.Address, _block_nums *big.Int) (common.Address, *types.Transaction, *StorageAccount, error) {
	parsed, err := abi.JSON(strings.NewReader(StorageAccountABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StorageAccountBin), backend, _fileAddress, _block_nums)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StorageAccount{StorageAccountCaller: StorageAccountCaller{contract: contract}, StorageAccountTransactor: StorageAccountTransactor{contract: contract}, StorageAccountFilterer: StorageAccountFilterer{contract: contract}}, nil
}

// StorageAccount is an auto generated Go binding around an Ethereum contract.
type StorageAccount struct {
	StorageAccountCaller     // Read-only binding to the contract
	StorageAccountTransactor // Write-only binding to the contract
	StorageAccountFilterer   // Log filterer for contract events
}

// StorageAccountCaller is an auto generated read-only Go binding around an Ethereum contract.
type StorageAccountCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageAccountTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StorageAccountTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageAccountFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StorageAccountFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageAccountSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StorageAccountSession struct {
	Contract     *StorageAccount   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageAccountCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StorageAccountCallerSession struct {
	Contract *StorageAccountCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// StorageAccountTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StorageAccountTransactorSession struct {
	Contract     *StorageAccountTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// StorageAccountRaw is an auto generated low-level Go binding around an Ethereum contract.
type StorageAccountRaw struct {
	Contract *StorageAccount // Generic contract binding to access the raw methods on
}

// StorageAccountCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StorageAccountCallerRaw struct {
	Contract *StorageAccountCaller // Generic read-only contract binding to access the raw methods on
}

// StorageAccountTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StorageAccountTransactorRaw struct {
	Contract *StorageAccountTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageAccount creates a new instance of StorageAccount, bound to a specific deployed contract.
func NewStorageAccount(address common.Address, backend bind.ContractBackend) (*StorageAccount, error) {
	contract, err := bindStorageAccount(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageAccount{StorageAccountCaller: StorageAccountCaller{contract: contract}, StorageAccountTransactor: StorageAccountTransactor{contract: contract}, StorageAccountFilterer: StorageAccountFilterer{contract: contract}}, nil
}

// NewStorageAccountCaller creates a new read-only instance of StorageAccount, bound to a specific deployed contract.
func NewStorageAccountCaller(address common.Address, caller bind.ContractCaller) (*StorageAccountCaller, error) {
	contract, err := bindStorageAccount(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageAccountCaller{contract: contract}, nil
}

// NewStorageAccountTransactor creates a new write-only instance of StorageAccount, bound to a specific deployed contract.
func NewStorageAccountTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageAccountTransactor, error) {
	contract, err := bindStorageAccount(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageAccountTransactor{contract: contract}, nil
}

// NewStorageAccountFilterer creates a new log filterer instance of StorageAccount, bound to a specific deployed contract.
func NewStorageAccountFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageAccountFilterer, error) {
	contract, err := bindStorageAccount(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageAccountFilterer{contract: contract}, nil
}

// bindStorageAccount binds a generic wrapper to an already deployed contract.
func bindStorageAccount(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StorageAccountABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageAccount *StorageAccountRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StorageAccount.Contract.StorageAccountCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageAccount *StorageAccountRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageAccount.Contract.StorageAccountTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageAccount *StorageAccountRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageAccount.Contract.StorageAccountTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageAccount *StorageAccountCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StorageAccount.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageAccount *StorageAccountTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageAccount.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageAccount *StorageAccountTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageAccount.Contract.contract.Transact(opts, method, params...)
}

// GetAllBlocksInfo is a free data retrieval call binding the contract method 0x10b0a7c9.
//
// Solidity: function getAllBlocksInfo() constant returns(bytes32[] blocksHash, bytes32[] peersInfo)
func (_StorageAccount *StorageAccountCaller) GetAllBlocksInfo(opts *bind.CallOpts) (struct {
	BlocksHash [][32]byte
	PeersInfo  [][32]byte
}, error) {
	ret := new(struct {
		BlocksHash [][32]byte
		PeersInfo  [][32]byte
	})
	out := ret
	err := _StorageAccount.contract.Call(opts, out, "getAllBlocksInfo")
	return *ret, err
}

// GetAllBlocksInfo is a free data retrieval call binding the contract method 0x10b0a7c9.
//
// Solidity: function getAllBlocksInfo() constant returns(bytes32[] blocksHash, bytes32[] peersInfo)
func (_StorageAccount *StorageAccountSession) GetAllBlocksInfo() (struct {
	BlocksHash [][32]byte
	PeersInfo  [][32]byte
}, error) {
	return _StorageAccount.Contract.GetAllBlocksInfo(&_StorageAccount.CallOpts)
}

// GetAllBlocksInfo is a free data retrieval call binding the contract method 0x10b0a7c9.
//
// Solidity: function getAllBlocksInfo() constant returns(bytes32[] blocksHash, bytes32[] peersInfo)
func (_StorageAccount *StorageAccountCallerSession) GetAllBlocksInfo() (struct {
	BlocksHash [][32]byte
	PeersInfo  [][32]byte
}, error) {
	return _StorageAccount.Contract.GetAllBlocksInfo(&_StorageAccount.CallOpts)
}

// GetBlockInfo is a free data retrieval call binding the contract method 0xbb141cf4.
//
// Solidity: function getBlockInfo(uint256 index) constant returns(bytes32 blockHash, bytes32 peerInfo)
func (_StorageAccount *StorageAccountCaller) GetBlockInfo(opts *bind.CallOpts, index *big.Int) (struct {
	BlockHash [32]byte
	PeerInfo  [32]byte
}, error) {
	ret := new(struct {
		BlockHash [32]byte
		PeerInfo  [32]byte
	})
	out := ret
	err := _StorageAccount.contract.Call(opts, out, "getBlockInfo", index)
	return *ret, err
}

// GetBlockInfo is a free data retrieval call binding the contract method 0xbb141cf4.
//
// Solidity: function getBlockInfo(uint256 index) constant returns(bytes32 blockHash, bytes32 peerInfo)
func (_StorageAccount *StorageAccountSession) GetBlockInfo(index *big.Int) (struct {
	BlockHash [32]byte
	PeerInfo  [32]byte
}, error) {
	return _StorageAccount.Contract.GetBlockInfo(&_StorageAccount.CallOpts, index)
}

// GetBlockInfo is a free data retrieval call binding the contract method 0xbb141cf4.
//
// Solidity: function getBlockInfo(uint256 index) constant returns(bytes32 blockHash, bytes32 peerInfo)
func (_StorageAccount *StorageAccountCallerSession) GetBlockInfo(index *big.Int) (struct {
	BlockHash [32]byte
	PeerInfo  [32]byte
}, error) {
	return _StorageAccount.Contract.GetBlockInfo(&_StorageAccount.CallOpts, index)
}

// CommitBlockInfo is a paid mutator transaction binding the contract method 0x6d40eb52.
//
// Solidity: function CommitBlockInfo(address _fileAddress, uint256 index, bytes32 hash, bytes32 peerInfo, string proof) returns()
func (_StorageAccount *StorageAccountTransactor) CommitBlockInfo(opts *bind.TransactOpts, _fileAddress common.Address, index *big.Int, hash [32]byte, peerInfo [32]byte, proof string) (*types.Transaction, error) {
	return _StorageAccount.contract.Transact(opts, "CommitBlockInfo", _fileAddress, index, hash, peerInfo, proof)
}

// CommitBlockInfo is a paid mutator transaction binding the contract method 0x6d40eb52.
//
// Solidity: function CommitBlockInfo(address _fileAddress, uint256 index, bytes32 hash, bytes32 peerInfo, string proof) returns()
func (_StorageAccount *StorageAccountSession) CommitBlockInfo(_fileAddress common.Address, index *big.Int, hash [32]byte, peerInfo [32]byte, proof string) (*types.Transaction, error) {
	return _StorageAccount.Contract.CommitBlockInfo(&_StorageAccount.TransactOpts, _fileAddress, index, hash, peerInfo, proof)
}

// CommitBlockInfo is a paid mutator transaction binding the contract method 0x6d40eb52.
//
// Solidity: function CommitBlockInfo(address _fileAddress, uint256 index, bytes32 hash, bytes32 peerInfo, string proof) returns()
func (_StorageAccount *StorageAccountTransactorSession) CommitBlockInfo(_fileAddress common.Address, index *big.Int, hash [32]byte, peerInfo [32]byte, proof string) (*types.Transaction, error) {
	return _StorageAccount.Contract.CommitBlockInfo(&_StorageAccount.TransactOpts, _fileAddress, index, hash, peerInfo, proof)
}

// DownloadSuccess is a paid mutator transaction binding the contract method 0x64d8aff6.
//
// Solidity: function DownloadSuccess() returns()
func (_StorageAccount *StorageAccountTransactor) DownloadSuccess(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageAccount.contract.Transact(opts, "DownloadSuccess")
}

// DownloadSuccess is a paid mutator transaction binding the contract method 0x64d8aff6.
//
// Solidity: function DownloadSuccess() returns()
func (_StorageAccount *StorageAccountSession) DownloadSuccess() (*types.Transaction, error) {
	return _StorageAccount.Contract.DownloadSuccess(&_StorageAccount.TransactOpts)
}

// DownloadSuccess is a paid mutator transaction binding the contract method 0x64d8aff6.
//
// Solidity: function DownloadSuccess() returns()
func (_StorageAccount *StorageAccountTransactorSession) DownloadSuccess() (*types.Transaction, error) {
	return _StorageAccount.Contract.DownloadSuccess(&_StorageAccount.TransactOpts)
}
