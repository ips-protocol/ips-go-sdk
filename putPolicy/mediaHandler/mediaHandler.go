package mediaHandler

import (
	"github.com/ipweb-group/go-sdk/utils"
	"io"
	"mime/multipart"
	"os"
	"regexp"
)

type MediaInfo struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration string `json:"duration"` // 时长是个浮点数，这里直接用字符串保存
	Type     string `json:"type"`     // 类型。针对图片可能是 jpeg/png/gif；针对视频可能是 h264 等
	MimeType string `json:"mime_type"`
}

func DetectMediaInfo(filePath string, mimeType string) (info MediaInfo, err error) {
	info = MediaInfo{
		MimeType: mimeType,
	}

	lg := utils.GetLogger()

	// 如果文件是支持的图片类型，就调用图片处理器获取图片尺寸信息
	if mimeType == "image/jpeg" || mimeType == "image/png" || mimeType == "image/gif" {
		lg.Info("File is of supported image type, will process image size detector")
		err = GetImageInfo(filePath, &info)
		if err != nil {
			return
		}
	}

	// 如果文件是视频类型，就调用视频处理器获取视频的基本信息
	if match, _ := regexp.MatchString("video/.*", mimeType); match {
		lg.Info("File is of type video, will process video information detector")
		err = GetVideoInfo(filePath, &info)
		if err != nil {
			return
		}
	}

	return
}

// 写入上传文件到临时文件，并返回临时文件的绝对路径
func WriteTmpFile(file multipart.File, cid string, ext string) (path string, err error) {
	tmpDir := utils.GetTmpDir()
	path = tmpDir + "/" + cid + ext

	_, err = file.Seek(0, 0)
	if err != nil {
		return
	}

	dst, err := os.Create(path)
	if err != nil {
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return
}
