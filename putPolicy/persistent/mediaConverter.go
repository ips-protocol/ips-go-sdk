package persistent

import (
	"encoding/json"
	"fmt"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/rpc"
	"github.com/ipweb-group/go-sdk/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// 处理多媒体文件的格式转换
func ConvertMediaJob() {
	fmt.Println("[INFO] Convert media job is started")

	rpcClient, err := rpc.NewClient(conf.GetConfig().NodeConf)
	if err != nil {
		panic(err)
	}

	// 每隔一段时间检查一次队列，并在有任务时执行任务
	for {
		time.Sleep(1 * time.Second)
		convertMedia(rpcClient)
	}
}

func convertMedia(rpcClient *rpc.Client) {
	// 获取第一个转换任务，并添加到转换中的 Hash 表中
	task := GetFirstUnprocessedTask()
	if task == nil {
		return
	}

	fmt.Printf("[INFO] Convert task detected. Cid is %s, Ops is %s, NotifyUrl is %s \n", task.Cid, task.PersistentOps, task.PersistentNotifyUrl)

	// 有任务时，添加任务到正在转换的 Hash 表中
	AddTaskToProcessingMap(task)

	dir, filename, ext := parseFilePath(task.FilePath)

	// TODO 判断如果视频是 h264 格式的话，将不再进行转换，而是直接回调成功

	// 启动转换任务
	inputFilePath := task.FilePath
	outputFilePath := dir + filename + "-converted" + ext
	ffmpeg := conf.GetConfig().ExternalConfig.Ffmpeg
	command := fmt.Sprintf("%s -stats -y -hide_banner -i %s -c:v libx264 -c:a libmp3lame %s", ffmpeg, inputFilePath, outputFilePath)

	fmt.Printf("[INFO] Start ffmpeg converter, command is %s \n", command)

	result, err := utils.ExecCommand(command)
	if err != nil {
		// 转换失败后，把失败的任务插入到失败的 Hash 表中
		fmt.Printf("[INFO] Convert task failed [%v] \n", err)
		AddFailedTask(task, result)
	} else {
		fmt.Printf("%s \n", result)
		fmt.Println("[INFO] Convert completed")
	}

	// 转换成功后，上传转换后的文件到 IPFS
	var dstCid string
	if err == nil {
		fmt.Printf("[INFO] Uploading converted file to IPFS")

		file, err := os.Open(outputFilePath)
		if err != nil {
			fmt.Printf("[ERROR] Open file failed, %v \n", err)
			AddFailedTask(task, err.Error())

		} else {
			dstFileInfo, _ := file.Stat()

			dstCid, err = rpcClient.Upload(file, dstFileInfo.Name(), dstFileInfo.Size())
			if err != nil {
				fmt.Printf("[ERROR] Upload converted file failed, [%v] \n", err)
				AddFailedTask(task, err.Error())

			} else {
				fmt.Printf("[INFO] Upload converted file completed, cid is %s \n", dstCid)
			}
		}
	}

	// 删除临时文件
	if err == nil {
		_ = os.Remove(inputFilePath)
		_ = os.Remove(outputFilePath)
	}

	// 移除 Redis 缓存
	RemoveProcessingTask(task)
	if err == nil {
		RemoveTask(task)
	}

	// 完成后回调（无论成功还是失败）
	requestBody := NotifyRequestBody{
		Cid:           task.Cid,
		Code:          CodeSuccess,
		Desc:          "",
		PersistentOps: task.PersistentOps,
		DstCid:        dstCid,
	}

	if err != nil {
		requestBody.Code = CodeFailed
		requestBody.Desc = err.Error()
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
