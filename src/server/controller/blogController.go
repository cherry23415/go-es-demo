package controller

import (
	"github.com/kataras/iris"
	"go-es-demo/src/server/constant"
	"go-es-demo/src/server/errcode"
	"go-es-demo/src/server/service"
)

type BlogController struct {
}

func (*BlogController) CreateIndex(ctx iris.Context) {
	blogService := service.BlogService{}
	blogService.CreateIndex(constant.CHERRY_INDEX)
	ctx.JSON(errcode.SUCCESS.Result())
	return
}

//批量插入
func (*BlogController) Batch(ctx iris.Context) {
	blogService := service.BlogService{}
	param := ctx.Params()
	blogService.Batch(constant.CHERRY_INDEX, constant.CHERRY_INDEX_BLOG_TYPE, param)
	ctx.JSON(errcode.SUCCESS.Result())
	return
}
