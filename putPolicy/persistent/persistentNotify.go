package persistent

type NotifyRequestBody struct {
	Hash          string `json:"hash"`          // 原始文件的 CID
	Code          int    `json:"code"`          // 状态码。0 表示成功；1 表示失败
	Desc          string `json:"desc"`          // 状态对应的描述
	PersistentOps string `json:"persistentOps"` // 持久化操作类型
	DstHash       string `json:"dstHash"`       // 持久化操作生成的目标文件的 CID
}

// 成功
const CodeSuccess = 0

// 失败
const CodeFailed = 1
