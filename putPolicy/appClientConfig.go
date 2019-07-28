package putPolicy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type appClientConfig struct {
	Clients []AppClient `json:"clients"`
}

// 以对象的形式保存所有的可用 AppClient 列表
var appClients appClientConfig

// 以 Map 的形式保存 AppClients 中 AccessKey 与 AppClient 的对应关系，方便查找
var appClientsMap map[string]AppClient

//
// 加载 app clients 配置文件，并解析其内容
//
func LoadAppClients(configFilePath string) {
	f, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &appClients)
	if err != nil {
		panic(err)
	}

	appClientsMap = make(map[string]AppClient)

	for _, client := range appClients.Clients {
		appClientsMap[client.AccessKey] = client
	}
}

//
// 根据 AccessKey 获取 AppClient 对象
//
func GetClientByAccessKey(accessKey string) (AppClient, error) {
	client, ok := appClientsMap[accessKey]
	if ok == false {
		return AppClient{}, errors.New("client does not exist")
	}

	return client, nil
}
