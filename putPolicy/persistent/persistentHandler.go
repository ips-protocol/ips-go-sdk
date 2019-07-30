package persistent

import (
	"encoding/json"
	"fmt"
	"github.com/ipweb-group/go-sdk/putPolicy/mediaHandler"
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

// 检查持久化参数及输入文件，如果文件符合持久化操作的条件，就将任务加入到持久化队列中，
// 否则将不会做任何处理，直接返回 "允许删除临时文件" 的标识
func (h *Task) CheckAndQueue() (shouldRemoveTmpFile bool, err error) {
	// 是否应该在完成后删除临时文件，默认为 true。当需要持久化操作时返回该值为 false
	shouldRemoveTmpFile = true

	// 检查持久化参数的值，如果无需持久化操作，直接返回即可
	if h.PersistentOps == "" {
		return
	}

	// 如果是要求转换视频，并且上传文件为视频类型时，启动视频转换任务
	if h.PersistentOps == "convertVideo" {
		if match, _ := regexp.MatchString("video/.*", h.MediaInfo.MimeType); match {
			h.queueVideo()
			shouldRemoveTmpFile = false
			fmt.Printf("[INFO] Video file is queued to redis: %s \n", h.FilePath)
		}
	}

	return
}

// 将任务转换为 JSON 字符串（用于保存到 Redis）
func (h *Task) ToJSON() string {
	j, _ := json.Marshal(h)
	return string(j)
}

// 添加视频任务到队列中
func (h *Task) queueVideo() {
	AddTaskToUnprocessedQueue(h)
}
