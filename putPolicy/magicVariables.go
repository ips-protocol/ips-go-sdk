package putPolicy

import (
	"strconv"
	"strings"
)

/**
 * 魔法变量
 */
type MagicVariable struct {
	FName       string `json:"fname"`
	FSize       int64  `json:"fsize"`
	MimeType    string `json:"mimeType"`
	EndUser     string `json:"endUser"`
	Hash        string `json:"hash"`
	ImageWidth  int    `json:"imageWidth"`
	ImageHeight int    `json:"imageHeight"`
}

/**
 * 应用魔法变量
 *
 * 替换 returnBody 中的所有魔法变量文本并返回替换后的文本
 */
func (variables *MagicVariable) ApplyMagicVariables(returnBody string) (ret string) {
	ret = strings.Replace(returnBody, "$(fname)", variables.FName, -1)
	ret = strings.Replace(ret, "$(fsize)", strconv.FormatInt(variables.FSize, 10), -1)
	ret = strings.Replace(ret, "$(mimeType)", variables.MimeType, -1)
	ret = strings.Replace(ret, "$(endUser)", variables.EndUser, -1)
	ret = strings.Replace(ret, "$(hash)", variables.Hash, -1)
	ret = strings.Replace(ret, "$(imageWidth)", strconv.Itoa(variables.ImageWidth), -1)
	ret = strings.Replace(ret, "$(imageHeight)", strconv.Itoa(variables.ImageHeight), -1)

	return
}
