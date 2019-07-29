package persistent

import "encoding/json"

// 视频转换任务
// 用于 Redis 中存储和视频转换相关的任务细节
type VideoTask struct {
	Cid                 string `json:"cid"`
	FilePath            string `json:"filePath"` // 临时文件保存的路径
	PersistentOps       string `json:"persistentOps"`
	PersistentNotifyUrl string `json:"persistentNotifyUrl"`
}

func (v *VideoTask) ToJSON() string {
	j, _ := json.Marshal(v)
	return string(j)
}

func UnmarshalVideoTask(str string) *VideoTask {
	ret := &VideoTask{}

	err := json.Unmarshal([]byte(str), ret)
	if err != nil {
		return nil
	}

	return ret
}
