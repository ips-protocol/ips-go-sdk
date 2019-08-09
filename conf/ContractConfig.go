package conf

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type ContractConfig struct {
	ClientKeyHex     string `json:"client_key_hex" yaml:"client_key_hex"`
	ContractNodeAddr string `json:"contract_node_addr" yaml:"contract_node_addr"`
}

func (cfg ContractConfig) GetClientKey() string {
	return cfg.ClientKeyHex
}

func (cfg ContractConfig) PrivateKey(clientKeyHex string) *ecdsa.PrivateKey {
	clientKey, err := crypto.HexToECDSA(clientKeyHex)
	if err != nil {
		panic(err)
	}
	return clientKey
}

func (cfg ContractConfig) GetClientAddress() common.Address {
	pubKey := cfg.PrivateKey(cfg.GetClientKey()).PublicKey
	return crypto.PubkeyToAddress(pubKey)
}

func (cfg ContractConfig) PublicKey(clientKeyHex string) common.Address {
	pubKey := cfg.PrivateKey(clientKeyHex).PublicKey
	return crypto.PubkeyToAddress(pubKey)
}
