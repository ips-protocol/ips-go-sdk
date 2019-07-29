package uploadController

import (
	"fmt"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/putPolicy"
	"github.com/ipweb-group/go-sdk/putPolicy/persistent"
	"github.com/ipweb-group/go-sdk/rpc"
	"github.com/ipweb-group/go-sdk/utils"
	"github.com/kataras/iris"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type UploadController struct {
	Node *rpc.Client
}

func New() UploadController {
	config := conf.GetConfig().NodeConf
	cli, err := rpc.NewClient(config)
	if err != nil {
		panic(err)
	}

	return UploadController{
		Node: cli,
	}
}

/**
 * 文件上传
 */
func (s *UploadController) Upload(ctx iris.Context) {
	token := ctx.FormValue("token")
	if len(token) == 0 {
		throwError(iris.StatusUnprocessableEntity, "No Upload Token Specified", ctx)
		return
	}

	// 解码上传 Token
	decodedPutPolicy, err := putPolicy.DecodePutPolicyString(token)
	if err != nil {
		throwError(iris.StatusInternalServerError, err.Error(), ctx)
		return
	}

	// 获取表单上传的文件
	file, fileHeader, err := ctx.FormFile("file")
	if err != nil {
		throwError(iris.StatusBadRequest, "Invalid File", ctx)
		return
	}

	defer file.Close()

	// TODO 上传有效期的校验

	// TODO 文件大小限制

	// TODO 根据 EndUser 参数进行扣款。扣款操作直接记录在链上

	// 上传文件到 IPFS
	cid, err := s.Node.Upload(file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		throwError(iris.StatusInternalServerError, "Failed to Upload, "+err.Error(), ctx)
		return
	}

	fileExt := path.Ext(fileHeader.Filename)
	mimeType := mime.TypeByExtension(fileExt)
	var width int
	var height int

	if mimeType == "image/jpeg" || mimeType == "image/png" || mimeType == "image/gif" {
		ctx.Application().Logger().Info("File is of supported image type, will process image size detector")
		// 获取图片宽高信息出错时不做任何处理，保持默认宽高为 0 即可
		width, height, err = getImageSize(file)
		if err != nil {
			ctx.Application().Logger().Warn(err.Error())
		}
	}

	// 处理视频文件
	if match, _ := regexp.MatchString("video/.*", mimeType); match {
		// 如果启用了持久化功能，并且类型是转换视频格式，将会把当前文件写入临时目录并添加到 Redis 队列
		if decodedPutPolicy.PutPolicy.PersistentOps == "convertVideo" {
			ctx.Application().Logger().Info("File is of type video, will process video converter")
			tmpFilePath, err := writeTmpFile(file, cid, fileExt)
			// 写入临时文件出错时不做任何处理
			if err != nil {
				ctx.Application().Logger().Warnf("Write video to tmp dir failed [%v]", err)

			} else {
				// 写入视频任务到 Redis 队列
				persistent.AddTaskToUnprocessedQueue(&persistent.VideoTask{
					Cid:                 cid,
					FilePath:            tmpFilePath,
					PersistentOps:       decodedPutPolicy.PutPolicy.PersistentOps,
					PersistentNotifyUrl: decodedPutPolicy.PutPolicy.PersistentNotifyUrl,
				})

				ctx.Application().Logger().Infof("Save video to temp file and queue to redis: %s", tmpFilePath)
			}
		}

		// TODO 检测视频宽高及时长
	}

	// 根据文件内容生成魔法变量
	magicVariable := putPolicy.MagicVariable{
		FName:       fileHeader.Filename,
		Hash:        cid,
		FSize:       fileHeader.Size,
		EndUser:     decodedPutPolicy.PutPolicy.EndUser,
		MimeType:    mimeType,
		ImageWidth:  width,
		ImageHeight: height,
	}

	// 如果上传策略中指定了回调地址，就异步去请求该地址
	if decodedPutPolicy.PutPolicy.CallbackUrl != "" {
		callbackBody := magicVariable.ApplyMagicVariables(decodedPutPolicy.PutPolicy.CallbackBody)

		responseBody, err := requestCallback(decodedPutPolicy.PutPolicy.CallbackUrl, callbackBody)
		// TODO 需要处理 callbackUrl 端返回非 200 的情况
		if err != nil {
			ctx.Application().Logger().Warn(err.Error())
			throwError(utils.StatusCallbackFailed, "Callback Failed, "+err.Error(), ctx)
			return
		}

		ctx.Header("Content-Type", "application/json; charset=UTF-8")
		_, _ = ctx.WriteString(responseBody)
		return
	}

	// 未指定回调地址时，返回默认内容
	_, _ = ctx.JSON(iris.Map{
		"hash":   cid,
		"length": fileHeader.Size,
	})
}

func throwError(statusCode int, error string, ctx iris.Context) {
	ctx.StatusCode(statusCode)
	_, _ = ctx.JSON(iris.Map{
		"error": error,
	})
}

func getImageSize(file multipart.File) (width int, height int, err error) {
	// 重置 file reader 的读取位置
	_, err = file.Seek(0, 0)

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return
	}

	width = config.Width
	height = config.Height
	return
}

func requestCallback(callbackUrl string, callbackBody string) (responseBody string, err error) {
	client := &http.Client{
		Timeout: time.Second * 30, // 默认请求超时时间为 30 秒
	}

	req, err := http.NewRequest(iris.MethodPost, callbackUrl, strings.NewReader(callbackBody))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "IPWeb SDK")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	responseBody = string(body)
	return
}

// 写入上传文件到临时文件，并返回临时文件的绝对路径
func writeTmpFile(file multipart.File, cid string, ext string) (path string, err error) {
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
