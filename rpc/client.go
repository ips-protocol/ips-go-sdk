package rpc

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/contracts/storage"
	"github.com/ipweb-group/go-sdk/p2p"
)

var ErrContractNotFound = errors.New("no contract code at given address")

const (
	P2pProtocl        = "/sys/http"
	NodeWeightInfoUrl = "http://teck.qiniudn.com/nodes.json?v=1"
)

type Client struct {
	Nodes                map[string]*Node
	NodesMux             map[string]*sync.RWMutex
	NodesAllocCond       *sync.Cond
	NodesWeightInfo      map[string]int
	NodeRefreshTime      time.Time
	NodeRefreshDuration  time.Duration
	NodeCloseDuration    time.Duration
	NodeRequestTimeout   time.Duration
	NodeRefreshWorkers   int
	ConnQuotaPerNode     int
	BlockUploadWorkers   int
	BlockDownloadWorkers int
	WalletPubKey         string
	*storage.Client
	*core.IpfsNode
}

func NewClient(cfg conf.Config) (cli *Client, err error) {
	ctx := context.Background()
	n, err := p2p.NewNode(ctx, cfg)
	if err != nil {
		return
	}

	cli = &Client{IpfsNode: n}
	cli.Nodes = make(map[string]*Node)
	cli.NodesMux = make(map[string]*sync.RWMutex)
	cli.NodesAllocCond = sync.NewCond(new(sync.Mutex))

	if cfg.NodeRefreshIntervalInSecond == 0 {
		cfg.NodeRefreshIntervalInSecond = 600
	}
	cli.NodeRefreshDuration = time.Second * time.Duration(cfg.NodeRefreshIntervalInSecond)

	if cfg.NodeRefreshWorkers == 0 {
		cfg.NodeRefreshWorkers = 10
	}
	cli.NodeRefreshWorkers = cfg.NodeRefreshWorkers

	if cfg.NodeCloseIntervalInSecond == 0 {
		cfg.NodeRefreshIntervalInSecond = 3600
	}
	cli.NodeCloseDuration = time.Second * time.Duration(cfg.NodeCloseIntervalInSecond)

	if cfg.NodeRequestTimeoutInSecond == 0 {
		cfg.NodeRequestTimeoutInSecond = 60
	}
	cli.NodeRequestTimeout = time.Second * time.Duration(cfg.NodeRequestTimeoutInSecond)

	if cfg.BlockUploadWorkers == 0 {
		cfg.BlockUploadWorkers = 5
	}
	cli.BlockUploadWorkers = cfg.BlockUploadWorkers

	if cfg.BlockDownloadWorkers == 0 {
		cfg.BlockDownloadWorkers = 5
	}
	cli.BlockDownloadWorkers = cfg.BlockDownloadWorkers

	if cfg.ConnQuotaPerNode == 0 {
		cfg.ConnQuotaPerNode = 2
	}
	cli.ConnQuotaPerNode = 2

	pubKey, err := GetWalletPubKey(cfg.ContractConf.ClientKeyHex)
	if err != nil {
		return
	}
	cli.WalletPubKey = pubKey

	resp, err := http.Get(NodeWeightInfoUrl)
	if err == nil {
		err = json.NewDecoder(resp.Body).Decode(&cli.NodesWeightInfo)
		if err != nil {
			return
		}
	}

	c, err := storage.NewClient(cfg.ContractConf)
	if err != nil {
		return
	}
	cli.Client = c
	go cli.refreshNodesTick()
	go cli.closeNodesTick()

	return
}

func GetWalletPubKey(walletKey string) (pubKey string, err error) {
	privateKey, err := crypto.HexToECDSA(walletKey)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("error casting public key to ECDSA")

	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fmt.Sprintf("%x", fromAddress), nil
}
