package uploadController

import (
	"fmt"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/putPolicy"
	"github.com/ipweb-group/go-sdk/rpc"
	"github.com/ipweb-group/go-sdk/utils"
	"github.com/kataras/iris"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"path"
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

	// 上传文件到 IPFS
	cid, err := s.Node.Upload(file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		throwError(iris.StatusInternalServerError, "Failed to Upload, "+err.Error(), ctx)
		return
	}

	mimeType := mime.TypeByExtension(path.Ext(fileHeader.Filename))
	var width int
	var height int

	if mimeType == "image/jpeg" || mimeType == "image/png" || mimeType == "image/gif" {
		// 获取图片宽高信息出错时不做任何处理，保持默认宽高为 0 即可
		width, height, err = getImageSize(file)
		if err != nil {
			ctx.Application().Logger().Warn(err.Error())
		}
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
	if err != nil {
	}

	fmt.Println(string(body))
	responseBody = string(body)
	return
}
