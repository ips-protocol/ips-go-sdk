package controllers

import (
	"github.com/ipweb-group/go-sdk/rpc"
	"github.com/kataras/iris"
	"io"
	"mime"
	"path/filepath"
)

type DownloadController struct{}

func (d *DownloadController) StreamedDownload(ctx iris.Context) {
	lg := ctx.Application().Logger()

	cid := ctx.Params().Get("cid")
	if cid == "" {
		throwError(iris.StatusUnprocessableEntity, "Invalid file hash", ctx)
		return
	}

	lg.Info("File download: ", cid)

	// TODO 文件不存在时会发生什么？
	rpcClient, _ := rpc.GetClientInstance()
	stream, meta, err := rpcClient.StreamRead(cid)
	if err != nil {
		lg.Errorf("An error occurred while downloading %s, %v", cid, err)
		throwError(iris.StatusInternalServerError, "Download file failed, "+err.Error(), ctx)
		return
	}
	defer stream.Close()

	filename := meta.FName
	fileExt := filepath.Ext(filename)
	if fileExt != "" {
		mimeType := mime.TypeByExtension(fileExt)
		// iris 的 ContentType 方法会自动设置 charset，并不适合这里的应用场景，
		// 所以这里直接使用 Header 方法设置 ContentType
		ctx.Header("Content-Type", mimeType)
	}

	_, err = io.Copy(ctx, stream)
	if err != nil {
		lg.Errorf("Copy file stream to context failed, %v", err)
		throwError(iris.StatusInternalServerError, "Send file content failed", ctx)
		return
	}
}
