package conf

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const walletKey = "5BEE78415C36E7DC7C7652957157C3E74011E1E8A8A344BD738A17E64DE37988"

type ECConfig struct {
	DataShards int `json:"data_shards"`
	ParShards  int `json:"par_shards"`
}

type ContractConfig struct {
	ClientKeyHex     string `json:"client_key_hex"`
	ContractNodeAddr string `json:"contract_node_addr"`
}

func (cfg ContractConfig) GetClientKey() *ecdsa.PrivateKey {
	clientKey, err := crypto.HexToECDSA(cfg.ClientKeyHex)
	if err != nil {
		panic(err)
	}
	return clientKey
}

func (cfg ContractConfig) GetClientAddress() common.Address {
	pubKey := cfg.GetClientKey().PublicKey
	return crypto.PubkeyToAddress(pubKey)
}

type Config struct {
	NodesRefreshIntervalInSecond int            `json:"nodes_refresh_interval_in_second"`
	NodeRequestTimeoutInSecond   int            `json:"node_request_timeout_in_second"`
	BlockUpWorkerCount           int            `json:"block_up_worker_count"`
	BlockDownloadWorkerCount     int            `json:"block_download_worker_count"`
	ContractConf                 ContractConfig `json:"contract_conf"`
	ECConfig
}

func LoadConf(cfg interface{}, cfgPath string) error {
	f, err := os.Open(cfgPath)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, cfg)
}

func GetWalletPubKey() (pubKey string, err error) {
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
