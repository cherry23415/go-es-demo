package initialize

import (
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"go-es-demo/src/server/route/api"
)

/**
iris版本10
*/
func SetupServer() {

	port := viper.GetString("server.port")
	host := viper.GetString("server.host")
	app := iris.New()
	api.Api(app)
	app.Run(iris.Addr(host + ":" + port))
}
