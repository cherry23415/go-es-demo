package controller

import (
	"github.com/kataras/iris"
	"go-es-demo/src/server/constant"
	"go-es-demo/src/server/entity"
	"go-es-demo/src/server/errcode"
	"go-es-demo/src/server/service"
)

type BlogController struct {
}

func (*BlogController) Save(ctx iris.Context) {
	var blog []entity.Blog
	blogService := service.BlogService{}
	err := ctx.ReadJSON(&blog)
	if err != nil {
		ctx.JSON(errcode.PARAM_ERROR.Result())
	}
	blogService.Save(constant.CHERRY_INDEX, constant.CHERRY_INDEX_BLOG_TYPE, blog)
	ctx.JSON(errcode.SUCCESS.Result())
	return
}

func (*BlogController) Search(ctx iris.Context) {
	blogService := service.BlogService{}
	query := ctx.URLParam("query")
	page, e1 := ctx.URLParamInt("page")
	if e1 != nil {
		ctx.JSON(errcode.PARAM_ERROR.Result())
	}
	size, e2 := ctx.URLParamInt("size")
	if e2 != nil {
		ctx.JSON(errcode.PARAM_ERROR.Result())
	}
	returnJson, err := blogService.Search(constant.CHERRY_INDEX, constant.CHERRY_INDEX_BLOG_TYPE, query, page, size)
	if err != nil {
		ctx.JSON(errcode.SYSTEM_ERROR.ResultWithMsg(err.Error()))
	}
	ctx.JSON(errcode.SUCCESS.ResultWithData(returnJson))
	return
}
