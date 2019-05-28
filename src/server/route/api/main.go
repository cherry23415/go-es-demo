package api

import (
	"github.com/kataras/iris"
	"go-es-demo/src/server/controller"
)

func Api(app *iris.Application) {

	var (
		blog = &controller.BlogController{}
	)
	api := app.Party("/es")
	api.Post("/save", blog.Save)
}
