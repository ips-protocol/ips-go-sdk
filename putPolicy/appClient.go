package putPolicy

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

// 第三方应用密钥对象
type AppClient struct {
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
	Description string `json:"description,omitempty"`
}

// 解码后的上传策略字符串内容
type DecodedPutPolicy struct {
	AccessKey string
	Sign      string
	AppClient AppClient
	PutPolicy PutPolicy
}

//
// 使用 PutPolicy 对象生成上传策略字符串
//
func (client *AppClient) MakePolicyWithPutPolicy(policy PutPolicy) string {
	jsonString, err := policy.ToJSON()
	if err != nil {
		return ""
	}

	return client.MakeWithJsonString(jsonString)
}

//
// 使用 JSON 字符串生成上传策略字符串
//
func (client *AppClient) MakeWithJsonString(json string) string {
	encodedPutPolicy := UrlSafeBase64Encode(json)
	sign := client.SignContent(encodedPutPolicy)

	return client.AccessKey + ":" + sign + ":" + encodedPutPolicy
}

//
// 解码上传策略字符串，并得到解码后的对象。解码失败或签名错误时返回 error
//
func DecodePutPolicyString(policy string) (DecodedPutPolicy, error) {

	// 1. 拆分策略字符串，分离出 AccessKey、签名和编码后的策略内容
	splitPolicy := strings.SplitN(policy, ":", 3)

	if len(splitPolicy) != 3 {
		return DecodedPutPolicy{}, errors.New("invalid policy string")
	}

	accessKey := splitPolicy[0]
	sign := splitPolicy[1]
	encodedPutPolicy := splitPolicy[2]

	// 2. 根据 AccessKey 查找 AppClient。找不到时抛出错误
	appClient, err := GetClientByAccessKey(accessKey)
	if err != nil {
		return DecodedPutPolicy{}, err
	}

	// 3. 验证签名
	if sign != appClient.SignContent(encodedPutPolicy) {
		return DecodedPutPolicy{}, errors.New("invalid signature")
	}

	// 4. 解码策略字符串
	jsonString := UrlSafeBase64Decode(encodedPutPolicy)
	putPolicy := PutPolicy{}

	err = json.Unmarshal([]byte(jsonString), &putPolicy)
	if err != nil {
		return DecodedPutPolicy{}, err
	}

	return DecodedPutPolicy{
		AccessKey: accessKey,
		Sign:      sign,
		AppClient: appClient,
		PutPolicy: putPolicy,
	}, nil
}

//
// 签名策略字符串
//
func (client *AppClient) SignContent(encodedPutPolicy string) string {
	h := hmac.New(sha1.New, []byte(client.SecretKey))
	h.Write([]byte(encodedPutPolicy))

	return hex.EncodeToString(h.Sum(nil))
}
