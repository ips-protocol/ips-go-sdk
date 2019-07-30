package putPolicy

import (
	"encoding/json"
	"github.com/ipweb-group/go-sdk/utils"
)

type PutPolicy struct {
	Deadline int32 `json:"deadline"`
	//ReturnBody          string `json:"returnBody"`
	EndUser             string `json:"endUser,omitempty"`
	CallbackUrl         string `json:"callbackUrl,omitempty"`
	CallbackBody        string `json:"callbackBody,omitempty"`
	FSizeLimit          int32  `json:"fSizeLimit,omitempty"`
	PersistentOps       string `json:"persistentOps,omitempty"`
	PersistentNotifyUrl string `json:"persistentNotifyUrl,omitempty"`
}

//
// 转换策略为 JSON 字符串
//
func (p *PutPolicy) ToJSON() (string, error) {
	str, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(str), nil
}

// 执行回调并返回回调响应内容
func (p *PutPolicy) ExecCallback(variable MagicVariable) (responseBody string, err error) {
	callbackBody := variable.ApplyMagicVariables(p.CallbackBody)

	responseBody, err = utils.RequestPost(p.CallbackUrl, callbackBody, utils.RequestContentTypeFormUrlencoded)
	return
}
