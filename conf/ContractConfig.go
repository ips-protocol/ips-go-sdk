package conf

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

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
