// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// StorageAccountABI is the input ABI used to generate the binding from.
const StorageAccountABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"downloadTotal\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"uploadedBlockNums\",\"outputs\":[{\"name\":\"\",\"type\":\"uint128\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fileAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"downloadSuccess\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_fileAddress\",\"type\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\"},{\"name\":\"hash\",\"type\":\"bytes\"},{\"name\":\"peerInfo\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"commitBlockInfo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFileInfo\",\"outputs\":[{\"name\":\"_fileAddress\",\"type\":\"address\"},{\"name\":\"_blockNums\",\"type\":\"uint256\"},{\"name\":\"_uploadedBlockNums\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getBlockInfo\",\"outputs\":[{\"name\":\"blockHash\",\"type\":\"bytes\"},{\"name\":\"peerInfo\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"initBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"blockNums\",\"outputs\":[{\"name\":\"\",\"type\":\"uint128\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_fileAddress\",\"type\":\"address\"},{\"name\":\"_block_nums\",\"type\":\"uint128\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"}]"

// StorageAccountBin is the compiled bytecode used for deploying new contracts.
const StorageAccountBin = `60806040526040516040806110108339810180604052810190808051906020019092919080519060200190929190505050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600360106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550346005819055505050610f4f806100c16000396000f300608060405260043610610099576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063267ef2f61461009e578063637a4ef9146100c957806373460e521461011857806396ab02df1461016f57806399f4521114610186578063adc0861d146102a5578063bb141cf41461030a578063d59a2ad61461041c578063edc8110d14610447575b600080fd5b3480156100aa57600080fd5b506100b3610496565b6040518082815260200191505060405180910390f35b3480156100d557600080fd5b506100de61049c565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561012457600080fd5b5061012d6104be565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561017b57600080fd5b506101846104e3565b005b34801561019257600080fd5b506102a3600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506106a3565b005b3480156102b157600080fd5b506102ba610adc565b604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828152602001935050505060405180910390f35b34801561031657600080fd5b5061033560048036038101908080359060200190929190505050610b70565b604051808060200180602001838103835285818151815260200191508051906020019080838360005b8381101561037957808201518184015260208101905061035e565b50505050905090810190601f1680156103a65780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b838110156103df5780820151818401526020810190506103c4565b50505050905090810190601f16801561040c5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390f35b34801561042857600080fd5b50610431610d12565b6040518082815260200191505060405180910390f35b34801561045357600080fd5b5061045c610d18565b60405180826fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60045481565b600360009054906101000a90046fffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600080600080600460008154809291906001019190505550600a600454111561050b5761069d565b60009350600092505b600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16831015610574576001600084815260200190815260200160002080549050840193508280600101935050610514565b8360320260046005540281151561058757fe5b049150600092505b600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff1683101561069c57600090505b600160008481526020019081526020016000208054905081101561068f57600160008481526020019081526020016000208181548110151561060a57fe5b906000526020600020906003020160020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f19350505050158015610681573d6000803e3d6000fd5b5080806001019150506105cc565b828060010193505061058f565b5b50505050565b6000600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff16851015156106e357600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff1614151561073e57600080fd5b6003600160008781526020019081526020016000208054905010151561076357600080fd5b600090505b60016000868152602001908152602001600020805490508110156108f85761085960016000878152602001908152602001600020828154811015156107a957fe5b90600052602060002090600302016001018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561084e5780601f106108235761010080835404028352916020019161084e565b820191906000526020600020905b81548152906001019060200180831161083157829003601f168201915b505050505084610d3a565b80156108e157503373ffffffffffffffffffffffffffffffffffffffff16600160008781526020019081526020016000208281548110151561089757fe5b906000526020600020906003020160020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b156108eb57600080fd5b8080600101915050610768565b6003600081819054906101000a90046fffffffffffffffffffffffffffffffff168092919060010191906101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555050600160008681526020019081526020016000206060604051908101604052808681526020018581526020013373ffffffffffffffffffffffffffffffffffffffff168152509080600181540180825580915050906001820390600052602060002090600302016000909192909190915060008201518160000190805190602001906109e3929190610e7e565b506020820151816001019080519060200190610a00929190610e7e565b5060408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050503373ffffffffffffffffffffffffffffffffffffffff166108fc60038060009054906101000a90046fffffffffffffffffffffffffffffffff16600502026fffffffffffffffffffffffffffffffff16600554811515610aa757fe5b049081150290604051600060405180830381858888f19350505050158015610ad3573d6000803e3d6000fd5b50505050505050565b60008060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169250600360109054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff169150600360009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff169050909192565b606080600160008481526020019081526020016000206000815481101515610b9457fe5b90600052602060002090600302016000018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610c395780601f10610c0e57610100808354040283529160200191610c39565b820191906000526020600020905b815481529060010190602001808311610c1c57829003601f168201915b50505050509150600160008481526020019081526020016000206000815481101515610c6157fe5b90600052602060002090600302016001018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610d065780601f10610cdb57610100808354040283529160200191610d06565b820191906000526020600020905b815481529060010190602001808311610ce957829003601f168201915b50505050509050915091565b60055481565b600360109054906101000a90046fffffffffffffffffffffffffffffffff1681565b60008082518451141515610d515760009150610e77565b600090505b8351811015610e72578281815181101515610d6d57fe5b9060200101517f010000000000000000000000000000000000000000000000000000000000000090047f0100000000000000000000000000000000000000000000000000000000000000027effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168482815181101515610de857fe5b9060200101517f010000000000000000000000000000000000000000000000000000000000000090047f0100000000000000000000000000000000000000000000000000000000000000027effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916141515610e655760009150610e77565b8080600101915050610d56565b600191505b5092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610ebf57805160ff1916838001178555610eed565b82800160010185558215610eed579182015b82811115610eec578251825591602001919060010190610ed1565b5b509050610efa9190610efe565b5090565b610f2091905b80821115610f1c576000816000905550600101610f04565b5090565b905600a165627a7a723058205d60a4199659256e04206ba9e1014c8c12515ee95ff4656c2622e11dfa05c7b30029`

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

// GetBlockInfo is a free data retrieval call binding the contract method 0xbb141cf4.
//
// Solidity: function getBlockInfo(index uint256) constant returns(blockHash bytes, peerInfo bytes)
func (_StorageAccount *StorageAccountCaller) GetBlockInfo(opts *bind.CallOpts, index *big.Int) (struct {
	BlockHash []byte
	PeerInfo  []byte
}, error) {
	ret := new(struct {
		BlockHash []byte
		PeerInfo  []byte
	})
	out := ret
	err := _StorageAccount.contract.Call(opts, out, "getBlockInfo", index)
	return *ret, err
}

// GetBlockInfo is a free data retrieval call binding the contract method 0xbb141cf4.
//
// Solidity: function getBlockInfo(index uint256) constant returns(blockHash bytes, peerInfo bytes)
func (_StorageAccount *StorageAccountSession) GetBlockInfo(index *big.Int) (struct {
	BlockHash []byte
	PeerInfo  []byte
}, error) {
	return _StorageAccount.Contract.GetBlockInfo(&_StorageAccount.CallOpts, index)
}

// GetBlockInfo is a free data retrieval call binding the contract method 0xbb141cf4.
//
// Solidity: function getBlockInfo(index uint256) constant returns(blockHash bytes, peerInfo bytes)
func (_StorageAccount *StorageAccountCallerSession) GetBlockInfo(index *big.Int) (struct {
	BlockHash []byte
	PeerInfo  []byte
}, error) {
	return _StorageAccount.Contract.GetBlockInfo(&_StorageAccount.CallOpts, index)
}

// GetFileInfo is a free data retrieval call binding the contract method 0xadc0861d.
//
// Solidity: function getFileInfo() constant returns(_fileAddress address, _blockNums uint256, _uploadedBlockNums uint256)
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
// Solidity: function getFileInfo() constant returns(_fileAddress address, _blockNums uint256, _uploadedBlockNums uint256)
func (_StorageAccount *StorageAccountSession) GetFileInfo() (struct {
	FileAddress       common.Address
	BlockNums         *big.Int
	UploadedBlockNums *big.Int
}, error) {
	return _StorageAccount.Contract.GetFileInfo(&_StorageAccount.CallOpts)
}

// GetFileInfo is a free data retrieval call binding the contract method 0xadc0861d.
//
// Solidity: function getFileInfo() constant returns(_fileAddress address, _blockNums uint256, _uploadedBlockNums uint256)
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

// CommitBlockInfo is a paid mutator transaction binding the contract method 0x99f45211.
//
// Solidity: function commitBlockInfo(_fileAddress address, index uint256, hash bytes, peerInfo bytes, proof bytes) returns()
func (_StorageAccount *StorageAccountTransactor) CommitBlockInfo(opts *bind.TransactOpts, _fileAddress common.Address, index *big.Int, hash []byte, peerInfo []byte, proof []byte) (*types.Transaction, error) {
	return _StorageAccount.contract.Transact(opts, "commitBlockInfo", _fileAddress, index, hash, peerInfo, proof)
}

// CommitBlockInfo is a paid mutator transaction binding the contract method 0x99f45211.
//
// Solidity: function commitBlockInfo(_fileAddress address, index uint256, hash bytes, peerInfo bytes, proof bytes) returns()
func (_StorageAccount *StorageAccountSession) CommitBlockInfo(_fileAddress common.Address, index *big.Int, hash []byte, peerInfo []byte, proof []byte) (*types.Transaction, error) {
	return _StorageAccount.Contract.CommitBlockInfo(&_StorageAccount.TransactOpts, _fileAddress, index, hash, peerInfo, proof)
}

// CommitBlockInfo is a paid mutator transaction binding the contract method 0x99f45211.
//
// Solidity: function commitBlockInfo(_fileAddress address, index uint256, hash bytes, peerInfo bytes, proof bytes) returns()
func (_StorageAccount *StorageAccountTransactorSession) CommitBlockInfo(_fileAddress common.Address, index *big.Int, hash []byte, peerInfo []byte, proof []byte) (*types.Transaction, error) {
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
