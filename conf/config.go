package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ECConfig struct {
	DataShards int `json:"data_shards"`
	ParShards  int `json:"par_shards"`
}

type Config struct {
	NodeRefreshIntervalInSecond int            `json:"node_refresh_interval_in_second"`
	NodeRefreshWorkers          int            `json:"node_refresh_workers"`
	NodeRequestTimeoutInSecond  int            `json:"node_request_timeout_in_second"`
	NodeCloseIntervalInSecond   int            `json:"node_close_interval_in_second"`
	ConnQuotaPerNode            int            `json:"conn_quota_per_node"`
	BlockUploadWorkers          int            `json:"block_upload_workers"`
	BlockDownloadWorkers        int            `json:"block_download_workers"`
	ContractConf                ContractConfig `json:"contract_conf"`
	ECConfig
}

type ExternalConfig struct {
	Ffmpeg  string `json:"ffmpeg"`
	Ffprobe string `json:"ffprobe"`
}

/**
 * 服务器配置
 */
type ServerConfig struct {
	ServerWriteTimeoutInSecond int            `json:"server_write_timeout_in_second"`
	ServerReadTimeoutInSecond  int            `json:"server_read_timeout_in_second"`
	ServerHost                 string         `json:"server_host"`
	NodeConf                   Config         `json:"node_conf"`
	RedisConfig                RedisConfig    `json:"redis_config"`
	ExternalConfig             ExternalConfig `json:"external"`
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
