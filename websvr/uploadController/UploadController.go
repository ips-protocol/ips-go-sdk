package uploadController

import (
	"fmt"
	"github.com/ipweb-group/go-sdk/putPolicy"
	"github.com/ipweb-group/go-sdk/putPolicy/mediaHandler"
	"github.com/ipweb-group/go-sdk/putPolicy/persistent"
	"github.com/ipweb-group/go-sdk/rpc"
	"github.com/ipweb-group/go-sdk/utils"
	"github.com/kataras/iris"
	"mime"
	"os"
	"path"
)

type UploadController struct {
	Node *rpc.Client
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

	// 写入临时文件
	fileExt := path.Ext(fileHeader.Filename)
	mimeType := mime.TypeByExtension(fileExt)
	tmpFilePath, err := mediaHandler.WriteTmpFile(file, cid, fileExt)

	// 初始化魔法变量对象
	magicVariable := putPolicy.MagicVariable{
		FName:    fileHeader.Filename,
		Hash:     cid,
		FSize:    fileHeader.Size,
		EndUser:  decodedPutPolicy.PutPolicy.EndUser,
		MimeType: mimeType,
	}

	// 检测媒体文件信息。当上传文件为图片或视频时，会检测文件的尺寸、时长等信息
	mediaInfo, err := mediaHandler.DetectMediaInfo(tmpFilePath, mimeType)
	if err == nil {
		magicVariable.Width = mediaInfo.Width
		magicVariable.Height = mediaInfo.Height
		magicVariable.Duration = mediaInfo.Duration
	} else {
		fmt.Printf("[WARN] Detect media info failed, [%v] \n", err)
	}

	// 处理持久化任务
	persistentTask := persistent.Task{
		Cid:                 cid,
		FilePath:            tmpFilePath,
		PersistentOps:       decodedPutPolicy.PutPolicy.PersistentOps,
		PersistentNotifyUrl: decodedPutPolicy.PutPolicy.PersistentNotifyUrl,
		MediaInfo:           mediaInfo,
	}
	shouldRemoveTmpFile, err := persistentTask.CheckAndQueue()

	// 删除临时文件
	if shouldRemoveTmpFile {
		_ = os.Remove(tmpFilePath)
	}

	// 如果上传策略中指定了回调地址，就异步去请求该地址
	if decodedPutPolicy.PutPolicy.CallbackUrl != "" {
		responseBody, err := decodedPutPolicy.PutPolicy.ExecCallback(magicVariable)
		if err != nil {
			fmt.Printf("[WARN] Callback to %s failed, %v \n", decodedPutPolicy.PutPolicy.CallbackUrl, err)
			throwError(utils.StatusCallbackFailed, "Callback Failed, "+err.Error(), ctx)
			return
		}
		fmt.Printf("[DEBUG] Callback to %s responds %s \n", decodedPutPolicy.PutPolicy.CallbackUrl, responseBody)

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
	ctx.Application().Logger().Error(error)
	ctx.StatusCode(statusCode)
	_, _ = ctx.JSON(iris.Map{
		"error": error,
	})
}
