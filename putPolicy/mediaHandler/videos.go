package mediaHandler

import (
	"encoding/json"
	"fmt"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/utils"
)

// ffprobe 返回的视频格式信息
type videoInfo struct {
	Streams []videoStream `json:"streams"`
}

// ffprobe 返回的音视频流数据信息
type videoStream struct {
	Index     int    `json:"index"`
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Duration  string `json:"duration"`
}

// 获取视频信息
func GetVideoInfo(filePath string, mediaInfo *MediaInfo) (err error) {
	// 调用 ffprobe 获取视频信息
	ffprobe := conf.GetConfig().ExternalConfig.Ffprobe
	fields := "stream=index,codec_name,codec_type,width,height,duration"
	command := fmt.Sprintf("%s -hide_banner -v quiet -print_format json -show_entries %s  -i %s", ffprobe, fields, filePath)

	result, err := utils.ExecCommand(command)
	if err != nil {
		fmt.Printf("[WARN] Get video properties failed [%v]", err)
		return
	}

	info := videoInfo{}
	err = json.Unmarshal([]byte(result), &info)
	if err != nil {
		fmt.Printf("[WARN] Get video properties failed, parse result failed, [%v]", err)
		return
	}

	// 解析返回的流信息
	isParsedVideo := false // 是否已经解析过视频流，避免重复解析
	for _, streamInfo := range info.Streams {
		if !isParsedVideo && streamInfo.CodecType == "video" {
			isParsedVideo = true

			mediaInfo.Type = streamInfo.CodecName
			mediaInfo.Width = streamInfo.Width
			mediaInfo.Height = streamInfo.Height
			mediaInfo.Duration = streamInfo.Duration
		}
	}

	return
}
