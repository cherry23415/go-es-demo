package initialize

import (
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
	"time"
)

func SetContext() {
	iris.UseFunc(func(ctx *iris.Context) {
		startAt := time.Now().UnixNano() / 1000000
		ctx.Set("startAt", startAt)
		ctx.Response.Header.Set("X-Powered-By", "soda-manager/v"+viper.GetString("version"))
		ctx.Next()
	})
}
