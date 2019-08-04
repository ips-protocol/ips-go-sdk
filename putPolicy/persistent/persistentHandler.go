package persistent

import (
	"encoding/json"
	"github.com/ipweb-group/go-sdk/putPolicy/mediaHandler"
	"github.com/ipweb-group/go-sdk/utils"
	"regexp"
)

// 持久化任务处理器
type Task struct {
	Cid                 string                 `json:"cid"`
	FilePath            string                 `json:"filePath"` // 临时文件保存的路径
	PersistentOps       string                 `json:"persistentOps"`
	PersistentNotifyUrl string                 `json:"persistentNotifyUrl"`
	MediaInfo           mediaHandler.MediaInfo `json:"mediaInfo"`
}

// 解析 JSON 字符串为 Task 类型的数据
func UnmarshalTask(str string) *Task {
	ret := &Task{}

	err := json.Unmarshal([]byte(str), ret)
	if err != nil {
		return nil
	}

	return ret
}

// 检查任务是否需要被添加到转换队列中，如果不需要被转换，则会返回允许删除临时文件的标识
func (h *Task) CheckShouldQueueTask() bool {
	// 如果是要求转换视频，并且上传文件为视频类型时，启动视频转换任务
	if h.PersistentOps == "convertVideo" {
		if match, _ := regexp.MatchString("video/.*", h.MediaInfo.MimeType); match {
			return false
		}
	}

	return true
}

// 添加任务添加到未处理队列中
func (h *Task) Queue() {
	AddTaskToUnprocessedQueue(h)
	utils.GetLogger().Info("Video file is queued to redis: %s", h.FilePath)
}

// 将任务转换为 JSON 字符串（用于保存到 Redis）
func (h *Task) ToJSON() string {
	j, _ := json.Marshal(h)
	return string(j)
}
