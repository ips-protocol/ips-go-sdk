package putPolicy

import (
	"encoding/json"
)

type PutPolicy struct {
	Deadline int32 `json:"deadline"`
	//ReturnBody   string `json:"returnBody"`
	EndUser      string `json:"endUser,omitempty"`
	CallbackUrl  string `json:"callbackUrl,omitempty"`
	CallbackBody string `json:"callbackBody,omitempty"`
	FSizeLimit   int32  `json:"fSizeLimit,omitempty"`
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
