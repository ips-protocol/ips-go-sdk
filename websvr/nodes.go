package websvr

import "github.com/kataras/iris"

func (s *Service) NodesList(ctx iris.Context) {
	ctx.JSON(s.Node.Nodes)
	return
}
