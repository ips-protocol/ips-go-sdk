package websvr

import (
	"io"

	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/rpc"
	"github.com/kataras/iris"
)

type Service struct {
	Node *rpc.Client
}

func NewService(cfg conf.Config) (svr Service, err error) {
	cli, err := rpc.NewClient(cfg)
	if err != nil {
		return
	}

	svr = Service{Node: cli}
	return
}

func (s *Service) FileUpload(ctx iris.Context) {
	lg := ctx.Application().Logger()

	file, fi, err := ctx.FormFile("file")
	if err != nil {
		lg.Error("FormFile failed:", err)
		return
	}
	defer file.Close()

	cid, err := s.Node.Upload(file, fi.Filename, fi.Size)
	if err != nil {
		lg.Error("upload to ipfs failed:", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"err": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"cid": cid})
	return
}

func (s *Service) FileDownload(ctx iris.Context) {
	lg := ctx.Application().Logger()

	cid := ctx.Params().Get("cid")
	if cid == "" {
		lg.Warn("cid is null")
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	lg.Info("file download cid:", cid)

	rd, _, err := s.Node.Download(cid)
	if err != nil {
		lg.Errorf("file cid: %s download failed err: %s", cid, err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"err": err.Error()})
		return
	}
	defer rd.Close()

	_, err = io.Copy(ctx, rd)
	if err != nil {
		lg.Error("file copy failed:", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"err": err.Error()})
		return
	}

	return
}
