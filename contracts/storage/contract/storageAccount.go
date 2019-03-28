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
const StorageAccountABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getAllBlocksInfo\",\"outputs\":[{\"name\":\"blocksHash\",\"type\":\"bytes32[]\"},{\"name\":\"peersInfo\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"downloadTotal\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"uploadedBlockNums\",\"outputs\":[{\"name\":\"\",\"type\":\"uint128\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fileAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_fileAddress\",\"type\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\"},{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"peerInfo\",\"type\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"string\"}],\"name\":\"commitBlockInfo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"downloadSuccess\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFileInfo\",\"outputs\":[{\"name\":\"_fileAddress\",\"type\":\"address\"},{\"name\":\"_blockNums\",\"type\":\"uint256\"},{\"name\":\"_uploadedBlockNums\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getBlockInfo\",\"outputs\":[{\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"name\":\"peerInfo\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"initBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"blockNums\",\"outputs\":[{\"name\":\"\",\"type\":\"uint128\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_fileAddress\",\"type\":\"address\"},{\"name\":\"_block_nums\",\"type\":\"uint128\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"}]"

// StorageAccountBin is the compiled bytecode used for deploying new contracts.
const StorageAccountBin = `6080604052604051604080610bea8339810180604052810190808051906020019092919080519060200190929190505050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600360106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550346005819055505050610b29806100c16000396000f3006080604052600436106100a4576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806310b0a7c9146100a9578063267ef2f61461015d578063637a4ef91461018857806373460e52146101d75780639006799d1461022e57806396ab02df146102dd578063adc0861d146102f4578063bb141cf414610359578063d59a2ad6146103b1578063edc8110d146103dc575b600080fd5b3480156100b557600080fd5b506100be61042b565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156101055780820151818401526020810190506100ea565b50505050905001838103825284818151815260200191508051906020019060200280838360005b8381101561014757808201518184015260208101905061012c565b5050505090500194505050505060405180910390f35b34801561016957600080fd5b50610172610432565b6040518082815260200191505060405180910390f35b34801561019457600080fd5b5061019d610438565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101e357600080fd5b506101ec61045a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561023a57600080fd5b506102db600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080356000191690602001909291908035600019169060200190929190803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919291929050505061047f565b005b3480156102e957600080fd5b506102f261080f565b005b34801561030057600080fd5b506103096109cf565b604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828152602001935050505060405180910390f35b34801561036557600080fd5b5061038460048036038101908080359060200190929190505050610a63565b60405180836000191660001916815260200182600019166000191681526020019250505060405180910390f35b3480156103bd57600080fd5b506103c6610ad5565b6040518082815260200191505060405180910390f35b3480156103e857600080fd5b506103f1610adb565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6060809091565b60045481565b600360009054906101000a90046fffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16851015156104bf57600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff1614151561051a57600080fd5b6003600160008781526020019081526020016000208054905010151561053f57600080fd5b600090505b600160008681526020019081526020016000208054905081101561063d578260001916600160008781526020019081526020016000208281548110151561058757fe5b9060005260206000209060030201600101546000191614801561062657503373ffffffffffffffffffffffffffffffffffffffff1660016000878152602001908152602001600020828154811015156105dc57fe5b906000526020600020906003020160020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b1561063057600080fd5b8080600101915050610544565b6003600081819054906101000a90046fffffffffffffffffffffffffffffffff168092919060010191906101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550506001600086815260200190815260200160002060606040519081016040528086600019168152602001856000191681526020013373ffffffffffffffffffffffffffffffffffffffff1681525090806001815401808255809150509060018203906000526020600020906003020160009091929091909150600082015181600001906000191690556020820151816001019060001916905560408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050503373ffffffffffffffffffffffffffffffffffffffff166108fc60038060009054906101000a90046fffffffffffffffffffffffffffffffff16600502026fffffffffffffffffffffffffffffffff166005548115156107da57fe5b049081150290604051600060405180830381858888f19350505050158015610806573d6000803e3d6000fd5b50505050505050565b600080600080600460008154809291906001019190505550600a6004541115610837576109c9565b60009350600092505b600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168310156108a0576001600084815260200190815260200160002080549050840193508280600101935050610840565b836032026004600554028115156108b357fe5b049150600092505b600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168310156109c857600090505b60016000848152602001908152602001600020805490508110156109bb57600160008481526020019081526020016000208181548110151561093657fe5b906000526020600020906003020160020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f193505050501580156109ad573d6000803e3d6000fd5b5080806001019150506108f8565b82806001019350506108bb565b5b50505050565b60008060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169250600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff169150600360009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff169050909192565b600080600160008481526020019081526020016000206000815481101515610a8757fe5b9060005260206000209060030201600001549150600160008481526020019081526020016000206000815481101515610abc57fe5b9060005260206000209060030201600101549050915091565b60055481565b600360109054906101000a90046fffffffffffffffffffffffffffffffff16815600a165627a7a72305820d6039eed7a0cc63ca152336d17eac3aaa55b704afd4296d66c4376a2c04c336d0029`

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

// BlockNums is a free data retrieval call binding the contract method 0xedc8110d.
//
// Solidity: function blockNums() constant returns(uint128)
func (_StorageAccount *StorageAccountCaller) BlockNums(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StorageAccount.contract.Call(opts, out, "blockNums")
	return *ret0, err
}

// BlockNums is a free data retrieval call binding the contract method 0xedc8110d.
//
// Solidity: function blockNums() constant returns(uint128)
func (_StorageAccount *StorageAccountSession) BlockNums() (*big.Int, error) {
	return _StorageAccount.Contract.BlockNums(&_StorageAccount.CallOpts)
}

// BlockNums is a free data retrieval call binding the contract method 0xedc8110d.
//
// Solidity: function blockNums() constant returns(uint128)
func (_StorageAccount *StorageAccountCallerSession) BlockNums() (*big.Int, error) {
	return _StorageAccount.Contract.BlockNums(&_StorageAccount.CallOpts)
}

// DownloadTotal is a free data retrieval call binding the contract method 0x267ef2f6.
//
// Solidity: function downloadTotal() constant returns(uint256)
func (_StorageAccount *StorageAccountCaller) DownloadTotal(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StorageAccount.contract.Call(opts, out, "downloadTotal")
	return *ret0, err
}

// DownloadTotal is a free data retrieval call binding the contract method 0x267ef2f6.
//
// Solidity: function downloadTotal() constant returns(uint256)
func (_StorageAccount *StorageAccountSession) DownloadTotal() (*big.Int, error) {
	return _StorageAccount.Contract.DownloadTotal(&_StorageAccount.CallOpts)
}

// DownloadTotal is a free data retrieval call binding the contract method 0x267ef2f6.
//
// Solidity: function downloadTotal() constant returns(uint256)
func (_StorageAccount *StorageAccountCallerSession) DownloadTotal() (*big.Int, error) {
	return _StorageAccount.Contract.DownloadTotal(&_StorageAccount.CallOpts)
}

// FileAddress is a free data retrieval call binding the contract method 0x73460e52.
//
// Solidity: function fileAddress() constant returns(address)
func (_StorageAccount *StorageAccountCaller) FileAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StorageAccount.contract.Call(opts, out, "fileAddress")
	return *ret0, err
}

// FileAddress is a free data retrieval call binding the contract method 0x73460e52.
//
// Solidity: function fileAddress() constant returns(address)
func (_StorageAccount *StorageAccountSession) FileAddress() (common.Address, error) {
	return _StorageAccount.Contract.FileAddress(&_StorageAccount.CallOpts)
}

// FileAddress is a free data retrieval call binding the contract method 0x73460e52.
//
// Solidity: function fileAddress() constant returns(address)
func (_StorageAccount *StorageAccountCallerSession) FileAddress() (common.Address, error) {
	return _StorageAccount.Contract.FileAddress(&_StorageAccount.CallOpts)
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

// GetFileInfo is a free data retrieval call binding the contract method 0xadc0861d.
//
// Solidity: function getFileInfo() constant returns(address _fileAddress, uint256 _blockNums, uint256 _uploadedBlockNums)
func (_StorageAccount *StorageAccountCaller) GetFileInfo(opts *bind.CallOpts) (struct {
	FileAddress       common.Address
	BlockNums         *big.Int
	UploadedBlockNums *big.Int
}, error) {
	ret := new(struct {
		FileAddress       common.Address
		BlockNums         *big.Int
		UploadedBlockNums *big.Int
	})
	out := ret
	err := _StorageAccount.contract.Call(opts, out, "getFileInfo")
	return *ret, err
}

// GetFileInfo is a free data retrieval call binding the contract method 0xadc0861d.
//
// Solidity: function getFileInfo() constant returns(address _fileAddress, uint256 _blockNums, uint256 _uploadedBlockNums)
func (_StorageAccount *StorageAccountSession) GetFileInfo() (struct {
	FileAddress       common.Address
	BlockNums         *big.Int
	UploadedBlockNums *big.Int
}, error) {
	return _StorageAccount.Contract.GetFileInfo(&_StorageAccount.CallOpts)
}

// GetFileInfo is a free data retrieval call binding the contract method 0xadc0861d.
//
// Solidity: function getFileInfo() constant returns(address _fileAddress, uint256 _blockNums, uint256 _uploadedBlockNums)
func (_StorageAccount *StorageAccountCallerSession) GetFileInfo() (struct {
	FileAddress       common.Address
	BlockNums         *big.Int
	UploadedBlockNums *big.Int
}, error) {
	return _StorageAccount.Contract.GetFileInfo(&_StorageAccount.CallOpts)
}

// InitBalance is a free data retrieval call binding the contract method 0xd59a2ad6.
//
// Solidity: function initBalance() constant returns(uint256)
func (_StorageAccount *StorageAccountCaller) InitBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StorageAccount.contract.Call(opts, out, "initBalance")
	return *ret0, err
}

// InitBalance is a free data retrieval call binding the contract method 0xd59a2ad6.
//
// Solidity: function initBalance() constant returns(uint256)
func (_StorageAccount *StorageAccountSession) InitBalance() (*big.Int, error) {
	return _StorageAccount.Contract.InitBalance(&_StorageAccount.CallOpts)
}

// InitBalance is a free data retrieval call binding the contract method 0xd59a2ad6.
//
// Solidity: function initBalance() constant returns(uint256)
func (_StorageAccount *StorageAccountCallerSession) InitBalance() (*big.Int, error) {
	return _StorageAccount.Contract.InitBalance(&_StorageAccount.CallOpts)
}

// UploadedBlockNums is a free data retrieval call binding the contract method 0x637a4ef9.
//
// Solidity: function uploadedBlockNums() constant returns(uint128)
func (_StorageAccount *StorageAccountCaller) UploadedBlockNums(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StorageAccount.contract.Call(opts, out, "uploadedBlockNums")
	return *ret0, err
}

// UploadedBlockNums is a free data retrieval call binding the contract method 0x637a4ef9.
//
// Solidity: function uploadedBlockNums() constant returns(uint128)
func (_StorageAccount *StorageAccountSession) UploadedBlockNums() (*big.Int, error) {
	return _StorageAccount.Contract.UploadedBlockNums(&_StorageAccount.CallOpts)
}

// UploadedBlockNums is a free data retrieval call binding the contract method 0x637a4ef9.
//
// Solidity: function uploadedBlockNums() constant returns(uint128)
func (_StorageAccount *StorageAccountCallerSession) UploadedBlockNums() (*big.Int, error) {
	return _StorageAccount.Contract.UploadedBlockNums(&_StorageAccount.CallOpts)
}

// CommitBlockInfo is a paid mutator transaction binding the contract method 0x9006799d.
//
// Solidity: function commitBlockInfo(address _fileAddress, uint256 index, bytes32 hash, bytes32 peerInfo, string proof) returns()
func (_StorageAccount *StorageAccountTransactor) CommitBlockInfo(opts *bind.TransactOpts, _fileAddress common.Address, index *big.Int, hash [32]byte, peerInfo [32]byte, proof string) (*types.Transaction, error) {
	return _StorageAccount.contract.Transact(opts, "commitBlockInfo", _fileAddress, index, hash, peerInfo, proof)
}

// CommitBlockInfo is a paid mutator transaction binding the contract method 0x9006799d.
//
// Solidity: function commitBlockInfo(address _fileAddress, uint256 index, bytes32 hash, bytes32 peerInfo, string proof) returns()
func (_StorageAccount *StorageAccountSession) CommitBlockInfo(_fileAddress common.Address, index *big.Int, hash [32]byte, peerInfo [32]byte, proof string) (*types.Transaction, error) {
	return _StorageAccount.Contract.CommitBlockInfo(&_StorageAccount.TransactOpts, _fileAddress, index, hash, peerInfo, proof)
}

// CommitBlockInfo is a paid mutator transaction binding the contract method 0x9006799d.
//
// Solidity: function commitBlockInfo(address _fileAddress, uint256 index, bytes32 hash, bytes32 peerInfo, string proof) returns()
func (_StorageAccount *StorageAccountTransactorSession) CommitBlockInfo(_fileAddress common.Address, index *big.Int, hash [32]byte, peerInfo [32]byte, proof string) (*types.Transaction, error) {
	return _StorageAccount.Contract.CommitBlockInfo(&_StorageAccount.TransactOpts, _fileAddress, index, hash, peerInfo, proof)
}

// DownloadSuccess is a paid mutator transaction binding the contract method 0x96ab02df.
//
// Solidity: function downloadSuccess() returns()
func (_StorageAccount *StorageAccountTransactor) DownloadSuccess(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageAccount.contract.Transact(opts, "downloadSuccess")
}

// DownloadSuccess is a paid mutator transaction binding the contract method 0x96ab02df.
//
// Solidity: function downloadSuccess() returns()
func (_StorageAccount *StorageAccountSession) DownloadSuccess() (*types.Transaction, error) {
	return _StorageAccount.Contract.DownloadSuccess(&_StorageAccount.TransactOpts)
}

// DownloadSuccess is a paid mutator transaction binding the contract method 0x96ab02df.
//
// Solidity: function downloadSuccess() returns()
func (_StorageAccount *StorageAccountTransactorSession) DownloadSuccess() (*types.Transaction, error) {
	return _StorageAccount.Contract.DownloadSuccess(&_StorageAccount.TransactOpts)
}
