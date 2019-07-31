package persistent

import (
	"encoding/json"
	"fmt"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/rpc"
	"github.com/ipweb-group/go-sdk/utils"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// 处理多媒体文件的格式转换
func ConvertMediaJob() {
	fmt.Println("[INFO] Convert media job is started")

	// 每隔一段时间检查一次队列，并在有任务时执行任务
	for {
		time.Sleep(1 * time.Second)
		convertMedia()
	}
}

func convertMedia() {
	// 获取第一个转换任务，并添加到转换中的 Hash 表中
	task := GetFirstUnprocessedTask()
	if task == nil {
		return
	}

	fmt.Printf("[INFO] Convert task detected. Hash is %s, Ops is %s, NotifyUrl is %s \n", task.Cid, task.PersistentOps, task.PersistentNotifyUrl)

	// 是否需要转换，默认为需要。如果判断文件格式本身无需转换的话，将会设置此变量为 false
	needConvert := true
	var err error

	// 有任务时，添加任务到正在转换的 Hash 表中
	AddTaskToProcessingMap(task)

	dir, filename, ext := parseFilePath(task.FilePath)

	// 判断如果视频是 h264 格式的话，将不再进行转换，而是直接回调成功
	if task.MediaInfo.Type == "h264" {
		needConvert = false
	}

	// 启动转换任务
	inputFilePath := task.FilePath
	outputFilePath := dir + filename + "-converted" + ext

	if needConvert {
		ffmpeg := conf.GetConfig().ExternalConfig.Ffmpeg
		command := fmt.Sprintf("%s -stats -y -hide_banner -i %s -c:v libx264 -c:a libmp3lame %s", ffmpeg, inputFilePath, outputFilePath)
		fmt.Printf("[INFO] Start ffmpeg converter, command is %s \n", command)

		var result string // var 创建 result 变量，避免使用 := 时覆盖 err 变量的值
		result, err = utils.ExecCommand(command)
		if err != nil {
			// 转换失败后，把失败的任务插入到失败的 Hash 表中
			fmt.Printf("[INFO] Convert task failed [%v] \n", err)
			AddFailedTask(task, result)
		} else {
			fmt.Printf("%s \n", result)
			fmt.Println("[INFO] Convert completed")
		}
	}

	// 转换成功后，上传转换后的文件到 IPFS
	var dstCid string
	if err == nil && needConvert {
		fmt.Printf("[INFO] Uploading converted file to IPFS")

		file, err := os.Open(outputFilePath)
		if err != nil {
			fmt.Printf("[ERROR] Open file failed, %v \n", err)
			AddFailedTask(task, err.Error())

		} else {
			defer closeFile(file)
			dstFileInfo, _ := file.Stat()

			rpcClient, _ := rpc.GetClientInstance()
			dstCid, err = rpcClient.Upload(file, dstFileInfo.Name(), dstFileInfo.Size())
			if err != nil {
				fmt.Printf("[ERROR] Upload converted file failed, [%v] \n", err)
				AddFailedTask(task, err.Error())

			} else {
				fmt.Printf("[INFO] Upload converted file completed, cid is %s \n", dstCid)
			}
		}
	}

	// 删除临时文件（无论转换成功或者失败，都尝试删除文件）
	_ = os.Remove(inputFilePath)
	_ = os.Remove(outputFilePath)

	// 移除 Redis 缓存
	RemoveProcessingTask(task)
	if err == nil || !needConvert {
		RemoveTask(task)
	}

	// 完成后回调（无论成功还是失败）
	requestBody := NotifyRequestBody{
		Hash:          task.Cid,
		Code:          CodeSuccess,
		Desc:          "",
		PersistentOps: task.PersistentOps,
		DstHash:       dstCid,
	}

	if err != nil {
		requestBody.Code = CodeFailed
		requestBody.Desc = err.Error()
	}

	// 如果无需转换，直接设置转换后的 hash 为原文件的 CID
	if !needConvert {
		requestBody.DstHash = task.Cid
	}

	stringContent, _ := json.Marshal(requestBody)
	responseBody, err := utils.RequestPost(task.PersistentNotifyUrl, string(stringContent), utils.RequestContentTypeJson)
	if err != nil {
		fmt.Printf("[WARN] Callback failed in persistent process, %v \n", err)
	}
	fmt.Printf("[DEBUG] Callback in persistent process responds: %s", responseBody)
}

// 分割路径字符串为目录、文件名、文件后缀三部分
func parseFilePath(filePath string) (dir string, filename string, ext string) {
	dir = path.Dir(filePath) + string(filepath.Separator)
	base := path.Base(filePath)
	ext = path.Ext(base)
	filename = strings.TrimSuffix(base, ext)
	return
}

// 关闭文件。发生错误时不做任何处理
func closeFile(file io.Closer) {
	err := file.Close()
	if err != nil {
		fmt.Println("[WARN] An error occurred while closing file")
	}
}
