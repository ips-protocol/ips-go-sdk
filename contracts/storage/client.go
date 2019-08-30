package storage

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/storage/contract"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipweb-group/go-sdk/conf"
)

var storageDepositContractAddr = common.HexToAddress("0x0000000000000000000000000000000000000010")

type Client struct {
	conf.ContractConfig
	*ethclient.Client
	nonce    uint64
	nonceMux sync.RWMutex
}

func NewClient(cfg conf.ContractConfig) (cli *Client, err error) {
	client, err := ethclient.Dial(cfg.ContractNodeAddr)
	if err != nil {
		return
	}
	cli = &Client{Client: client, ContractConfig: cfg}
	//nonce, err := cli.NonceAt(context.Background(), cli.GetClientAddress(), nil)
	//if err != nil {
	//	return
	//}
	//cli.nonce = nonce
	return
}
func (c *Client) NewKeyedTransactor(clientKey string) *bind.TransactOpts {
	transactor := bind.NewKeyedTransactor(c.PrivateKey(clientKey))
	return transactor
}

func (c *Client) GetStorageAccount(fileHash string) (stgAccountAddr common.Address, err error) {
	storageDeposit, err := contract.NewStorageDeposit(storageDepositContractAddr, c)
	if err != nil {
		return
	}

	fileAddress := common.BytesToAddress(crypto.Keccak256([]byte(fileHash)))
	stgAccountAddr, err = storageDeposit.GetStorageAccount(nil, fileAddress)
	return
}

func (c *Client) NewUploadJob(fileHash string, fsize int64, shards int, shardSize int64) (job *contract.StorageDepositNewUploadJob, err error) {
	return c.NewUploadJobByClientKey(c.GetClientKey(), fileHash, fsize, shards, shardSize)
}

func (c *Client) NewUploadJobByClientKey(clientKey string, fileHash string, fsize int64, shards int, shardSize int64) (job *contract.StorageDepositNewUploadJob, err error) {
	storageDeposit, err := contract.NewStorageDeposit(storageDepositContractAddr, c)
	if err != nil {
		return
	}
	transactor := c.NewKeyedTransactor(clientKey)

	bytePrice, err := storageDeposit.GetBytePrice(nil)
	if err != nil {
		return
	}

	c.ContractConfig.PrivateKey(clientKey)
	// transactor.GasPrice = big.NewInt(0)

	//c.nonceMux.Lock()
	//c.nonce++
	//transactor.Nonce = big.NewInt(int64(c.nonce))
	//c.nonceMux.Unlock()

	fmt.Printf("client.NewUploadJob hash: %s bytePrice: %d, shardSize: %d, shards: %d\n", fileHash, bytePrice, shardSize, shards)
	value := big.NewInt(bytePrice.Int64())
	value.Mul(value, big.NewInt(shardSize))
	value.Mul(value, big.NewInt(int64(shards)))
	transactor.Value = value

	// pubKey := c.ContractConfig.PrivateKey(clientKey).PublicKey
	// fromAddress := crypto.PubkeyToAddress(pubKey)
	// nonce, err := c.Client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	return
	// }
	// fmt.Println("client.NewUploadJobByClientKey nonce=", nonce)

	fileAddress := common.BytesToAddress(crypto.Keccak256([]byte(fileHash)))
	tx, err := storageDeposit.NewUploadJob(transactor, fileAddress, big.NewInt(fsize), big.NewInt(int64(shards)), big.NewInt(shardSize))
	if err != nil {
		return
	}

	ctx := context.Background()
	receipt, err := c.waitTransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return
	}

	if len(receipt.Logs) <= 0 {
		err = fmt.Errorf("no receipt logs, status: %d, GasUsed: %d, CumulativeGasUsed: %d", receipt.Status, receipt.GasUsed, receipt.CumulativeGasUsed)
		return
	}

	storageABI, err := abi.JSON(bytes.NewReader([]byte(contract.StorageDepositABI)))
	if err != nil {
		return
	}

	job = &contract.StorageDepositNewUploadJob{}
	err = storageABI.Unpack(job, "NewUploadJob", receipt.Logs[0].Data)
	return
}

func (c *Client) DeleteFile(fileHash string) error {
	return c.DeleteFileByClientKey(c.ContractConfig.GetClientKey(), fileHash)
}

func (c *Client) DeleteFileByClientKey(clientKey string, fileHash string) error {
	storageDeposit, err := contract.NewStorageDeposit(storageDepositContractAddr, c)
	if err != nil {
		return err
	}

	fileAddress := common.BytesToAddress(crypto.Keccak256([]byte(fileHash)))
	transactor := c.NewKeyedTransactor(clientKey)
	tx, err := storageDeposit.DeleteFile(transactor, fileAddress)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = c.waitTransactionReceipt(ctx, tx.Hash())
	return err
}

type BlockInfo struct {
	BlockHash string
	PeerId    string
}

func (c *Client) GetBlocksInfo(fileHash string) (blocksInfo []BlockInfo, err error) {
	stgAccountAddr, err := c.GetStorageAccount(fileHash)
	if err != nil {
		return
	}

	stgAccount, err := contract.NewStorageAccount(stgAccountAddr, c)
	if err != nil {
		return
	}

	fInfo, err := stgAccount.GetFileInfo(nil)
	if err != nil {
		return
	}
	if fInfo.BlockNums.Cmp(fInfo.UploadedBlockNums) != 0 {
		err = fmt.Errorf("incomplete file, block number: %d, upload block number: %s", fInfo.BlockNums, fInfo.UploadedBlockNums)
		return
	}

	blockNums := int(fInfo.BlockNums.Int64())
	for i := 0; i < blockNums; i++ {
		blkInfo, err := stgAccount.GetBlockInfo(nil, big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}

		blocksInfo = append(blocksInfo, BlockInfo{string(blkInfo.BlockHash), string(blkInfo.PeerInfo)})
	}
	return
}

func (c *Client) waitTransactionReceipt(ctx context.Context, tx common.Hash) (receipt *types.Receipt, err error) {
	for {
		receipt, err = c.TransactionReceipt(ctx, tx)
		if err != nil {
			if err != ethereum.NotFound {
				return
			}
			time.Sleep(time.Second)
			continue
		}

		if receipt.Status != 1 {
			err = fmt.Errorf("tx %s status is failed", tx.String())
		}
		return
	}
	return
}
