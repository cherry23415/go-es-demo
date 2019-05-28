package initialize

import (
	"fmt"
	"github.com/go-errors/errors"
	"go-es-demo/src/server/errcode"
	"gopkg.in/kataras/iris.v5"
)

func SetErrorDeal() {
	iris.Use(iris.HandlerFunc(func(ctx *iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("发生panic异常: %v\n", errors.Wrap(err, 2).ErrorStack())
				fmt.Println(msg)
				ctx.JSON(iris.StatusOK, errcode.SYSTEM_ERROR.Result())
			}
		}()
		ctx.Next()
	}))

}
