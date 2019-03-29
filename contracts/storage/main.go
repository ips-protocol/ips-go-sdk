package storage

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-sdk/contracts/storage/contract"
	"math/big"
	"time"
)

//go:generate abigen --abi contract/storageAccount.abi --bin contract/storageAccount.bin --type StorageAccount --pkg contract --out contract/storageAccount.go
//go:generate abigen --abi contract/storageDeposit.abi --bin contract/storageDeposit.bin --type StorageDeposit --pkg contract --out contract/storageDeposit.go

var (
	//storageDepositContractAddr = common.HexToAddress("0x0000000000000000000000000000000000000010")

	// client private Key, make sure its balance is enough
	testClientKey, _  = crypto.HexToECDSA("92D38B6F671F575EC9E47102364F53CA7F75B706A43606AA570E53917CBE2F9C")
	testClientAddress = crypto.PubkeyToAddress(testClientKey.PublicKey)

	// storage private key
	testStorageKey, _  = crypto.HexToECDSA("CED8FF231B09B14F09D8FF977C5C6C079EF4B485FC2A0D3B2955182B77310A04")
	testStorageAddress = crypto.PubkeyToAddress(testStorageKey.PublicKey)
)

var (
	rpcUrl = "http://127.0.0.1:8545" // node rpc address
)

func Proof(data []byte, prv *ecdsa.PrivateKey) ([]byte, error) {
	hash := crypto.Keccak256(data)
	return crypto.Sign(hash, prv)
}

func main() {

	// show balance of testClientAddress and testStorageAddress
	Status()

	//DeployStorageDepositContract()

	// client side new a upload job
	job, err := NewUploadJob(storageDepositContractAddr)
	if err != nil {
		panic(err)
	}

	// storage side commit blocks for a new upload job
	if err := CommitBlocks(job); err != nil {
		panic(err)
	}

	// client side get block info
	if err := GetCommitBlockInfo(job.StorageAccount); err != nil {
		panic(err)
	}

	// The client should invoke the DownloadSuccess method after download.
	if err := DownloadSuccess(job.StorageAccount); err != nil {
		panic(err)
	}
}

func DeployStorageDepositContract() {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}
	auth := bind.NewKeyedTransactor(testClientKey)
	addr, tx, _, err := contract.DeployStorageDeposit(auth, client)
	if err != nil {
		panic(err)
	}
	wait(client, tx.Hash())
	fmt.Println("storage deposit addr is", addr.Hex())
	storageDepositContractAddr = addr
}

// Make sure balance is enough
func Status() {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	balance, err := client.BalanceAt(ctx, testClientAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("testClientAccount address is", testClientAddress.Hex(), "balance", balance)

	balance, err = client.BalanceAt(ctx, testStorageAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("testStorageAddress address is", testStorageAddress.Hex(), "balance", balance)
}

func NewUploadJob(addr common.Address) (*contract.StorageDepositNewUploadJob, error) {

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}
	auth := bind.NewKeyedTransactor(testClientKey)
	storageDeposit, err := contract.NewStorageDeposit(addr, client)
	if err != nil {
		panic(err)
	}

	auth.Value = big.NewInt(1e6)

	fileAddress := common.BytesToAddress(crypto.Keccak256([]byte("file content" + time.Now().String()))[0:20])
	tx, err := storageDeposit.NewUploadJob(auth, fileAddress, big.NewInt(10240), big.NewInt(3))
	if err != nil {
		panic(err)
	}
	fmt.Println("invoking newUploadJob, tx hash is", tx.Hash().Hex())

	if err := wait(client, tx.Hash()); err != nil {
		return nil, err
	}

	ctx := context.Background()
	receipt, err := client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		panic(err)
	}
	if len(receipt.Logs) <= 0 {
		fmt.Println("no receipt logs", receipt.Status, receipt.GasUsed, receipt.CumulativeGasUsed)
		return nil, errors.New("unknown")
	}

	storageABI, err := abi.JSON(bytes.NewReader([]byte(contract.StorageDepositABI)))
	if err != nil {
		panic(err)
	}
	newUploadJob := contract.StorageDepositNewUploadJob{}
	if err := storageABI.Unpack(&newUploadJob, "NewUploadJob", receipt.Logs[0].Data); err != nil {
		panic(err)
	}

	fmt.Printf(`newUploadJob found:
	fileAddress : %s
	deposit : %d
	storageAccount : %s
`, newUploadJob.FileAddress.Hex(), newUploadJob.Deposit, newUploadJob.StorageAccount.Hex())
	return &newUploadJob, nil

}

func GetCommitBlockInfo(addr common.Address) error {

	fmt.Println("Getting committed Blocks")
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		return err
	}
	storageAccount, err := contract.NewStorageAccount(addr, client)
	if err != nil {
		return err
	}

	fileInfo, err := storageAccount.GetFileInfo(nil)
	fmt.Println("fileInfo", fileInfo.BlockNums, "\ttotal:", fileInfo.BlockNums.Int64(), "\tuploaded:", fileInfo.UploadedBlockNums.Int64())

	hash := append(crypto.Keccak256([]byte("data")), []byte("--")...)
	peer := append(crypto.Keccak256([]byte("peer")), []byte("--")...)
	for i := 0; i < 3; i++ {
		fmt.Println("GetBlockInfo", i)
		result, err := storageAccount.GetBlockInfo(nil, big.NewInt(int64(i)))
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !bytes.Equal(result.BlockHash, hash) {
			return fmt.Errorf("not expected hash, index is %d", i)
		}

		if !bytes.Equal(result.PeerInfo, peer) {
			return fmt.Errorf("not expected peer, index is %d", i)
		}

		fmt.Printf("block %d hash %x peer %x\n", i, result.BlockHash, result.PeerInfo)
	}
	return nil
}

func DownloadSuccess(storageAccountAddr common.Address) error {

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}
	auth := bind.NewKeyedTransactor(testClientKey)

	storageAccount, err := contract.NewStorageAccount(storageAccountAddr, client)
	if err != nil {
		return err
	}
	tx, err := storageAccount.DownloadSuccess(auth)
	if err != nil {
		return err
	}
	if err := wait(client, tx.Hash()); err != nil {
		return err
	}
	return nil
}

func CommitBlocks(job *contract.StorageDepositNewUploadJob) error {

	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		return err
	}
	storageAccount, err := contract.NewStorageAccount(job.StorageAccount, client)
	if err != nil {
		return err
	}

	auth := bind.NewKeyedTransactor(testStorageKey)

	hash := append(crypto.Keccak256([]byte("data")), []byte("--")...)
	peer := append(crypto.Keccak256([]byte("peer")), []byte("--")...)
	ctx := context.Background()

	fmt.Println("Committing Blocks")
	for i := 0; i < 3; i++ {

		tx, err := storageAccount.CommitBlockInfo(auth, job.FileAddress, big.NewInt(int64(i)), hash, peer, []byte("proof"))
		if err != nil {
			return err
		}
		if err := wait(client, tx.Hash()); err != nil {
			return err
		}
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			return err
		}
		if receipt.Status != 0 {
			return errors.New("tx failed")

		}
		fmt.Printf("block %d committed success!\n", i)
	}
	return nil
}

func wait(b *ethclient.Client, tx common.Hash) error {
	ctx := context.Background()
	for {
		receipt, err := b.TransactionReceipt(ctx, tx)
		if err != nil {
			if err != ethereum.NotFound {
				return err
			}
			time.Sleep(time.Second)
			continue
		}
		// TODO:update ethereum version
		if receipt.Status != 0 {
			return fmt.Errorf("tx %s status is failed", tx.String())
		}
		return nil
	}
	return nil
}
