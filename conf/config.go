package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ECConfig struct {
	DataShards int `json:"data_shards" yaml:"data_shards"`
	ParShards  int `json:"par_shards" yaml:"par_shards"`
}

type Config struct {
	NodeRefreshIntervalInSecond int            `json:"node_refresh_interval_in_second" yaml:"node_refresh_interval_in_second"`
	NodeRefreshWorkers          int            `json:"node_refresh_workers" yaml:"node_refresh_workers"`
	NodeRequestTimeoutInSecond  int            `json:"node_request_timeout_in_second" yaml:"node_request_timeout_in_second"`
	NodeCloseIntervalInSecond   int            `json:"node_close_interval_in_second" yaml:"node_close_interval_in_second"`
	ConnQuotaPerNode            int            `json:"conn_quota_per_node" yaml:"conn_quota_per_node"`
	BlockUploadWorkers          int            `json:"block_upload_workers" yaml:"block_upload_workers"`
	BlockDownloadWorkers        int            `json:"block_download_workers" yaml:"block_download_workers"`
	ContractConf                ContractConfig `json:"contract_conf" yaml:"contract_conf"`
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
