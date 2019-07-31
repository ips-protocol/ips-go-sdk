package rpc

import "github.com/ipweb-group/go-sdk/conf"

var clientInstance *Client

// 获取 RPC Client 的单例，此方法将在全局维护一个 client 实例，
// 避免在多个位置重复初始化 RPC Client，并简化代码结构
func GetClientInstance() (client *Client, err error) {
	if clientInstance == nil {
		clientInstance, err = NewClient(conf.GetConfig().NodeConf)
		if err != nil {
			return
		}
	}

	client = clientInstance
	return
}
