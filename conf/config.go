package conf

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

const walletKey = "5BEE78415C36E7DC7C7652957157C3E74011E1E8A8A344BD738A17E64DE37988"

type ECConfig struct {
	DataShards int `json:"data_shards"`
	ParShards  int `json:"par_shards"`
}

type Config struct {
	NodeRefreshIntervalInSecond int            `json:"node_refresh_interval_in_second"`
	NodeRefreshWorkers          int            `json:"node_refresh_workers"`
	NodeRequestTimeoutInSecond  int            `json:"node_request_timeout_in_second"`
	BlockUploadWorkers          int            `json:"block_upload_workers"`
	BlockDownloadWorkers        int            `json:"block_download_workers"`
	ContractConf                ContractConfig `json:"contract_conf"`
	ECConfig
}

/**
 * 服务器配置
 */
type ServerConfig struct {
	ServerWriteTimeoutInSecond int    `json:"server_write_timeout_in_second"`
	ServerReadTimeoutInSecond  int    `json:"server_read_timeout_in_second"`
	ServerHost                 string `json:"server_host"`
	NodeConf                   Config `json:"node_conf"`
}

/**
 * 服务器配置缓存。服务启用后加载一次配置文件内容，之后将会缓存下来
 */
var configCache ServerConfig

func LoadConfig(cfgPath string) {
	f, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &configCache)
	if err != nil {
		panic(err)
	}
}

func GetConfig() ServerConfig {
	return configCache
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
