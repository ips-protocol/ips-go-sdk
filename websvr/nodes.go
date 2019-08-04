package websvr

import "github.com/kataras/iris"

func (s *Service) NodesList(ctx iris.Context) {
	var ids []string
	for id := range s.Node.Nodes {
		ids = append(ids, id)
	}

	ctx.JSON(iris.Map{"ids": ids})
	return
}
