package conf

type EC struct {
	DataShards int `json:"data_shards"`
	ParShards  int `json:"par_shards"`
}

type Contract struct {
	ClientKey                  []byte `json:"client_key"`
	StorageKey                 []byte `json:"storage_key"`
	StorageDepositContractAddr []byte `json:"storage_deposit_contract_addr"`
	ContractNodeAddr           string `json:"contract_node_addr"`
}

type Config struct {
	EC
	Contract
}
