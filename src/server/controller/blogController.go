package controller

import (
	"encoding/json"
	"github.com/kataras/iris"
	"go-es-demo/src/server/constant"
	"go-es-demo/src/server/entity"
	"go-es-demo/src/server/errcode"
	"go-es-demo/src/server/service"
)

type BlogController struct {
}

func (*BlogController) Save(ctx iris.Context) {
	var blog entity.Blog
	blogService := service.BlogService{}
	err := ctx.ReadJSON(&blog)
	if err != nil {
		ctx.JSON(errcode.PARAM_ERROR.Result())
	}
	b, _ := json.Marshal(blog)
	blogService.Save(constant.CHERRY_INDEX, constant.CHERRY_INDEX_BLOG_TYPE, string(b), blog.Id)
	ctx.JSON(errcode.SUCCESS.Result())
	return
}
