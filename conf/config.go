package conf

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type ECConfig struct {
	DataShards int `json:"data_shards"`
	ParShards  int `json:"par_shards"`
}

type ContractConfig struct {
	ClientKeyHex     string `json:"client_key_hex"`
	StorageKeyHex    string `json:"storage_key_hex"`
	ContractNodeAddr string `json:"contract_node_addr"`

	TransactorGasLimit uint64 `json:"transactor_gas_limit"`
	TransactorGasPrice int64  `json:"transactor_gas_price"`
	TransactorValue    int64  `json:"transactor_value"`
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

func (cfg ContractConfig) GetStorageKey() *ecdsa.PrivateKey {
	storageKey, err := crypto.HexToECDSA(cfg.StorageKeyHex)
	if err != nil {
		panic(err)
	}
	return storageKey
}

func (cfg ContractConfig) GetStorageAddress() common.Address {
	pubKey := cfg.GetStorageKey().PublicKey
	return crypto.PubkeyToAddress(pubKey)
}

type Config struct {
	ECConfig
	ContractConfig
}
