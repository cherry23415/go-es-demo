package initialize

import (
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"go-es-demo/src/server/route/api"
)

/**
iris版本11
*/
func SetupServer() {

	port := viper.GetString("server.port")
	host := viper.GetString("server.host")
	app := iris.Default()
	api.Api(app)
	app.Run(iris.Addr(host + ":" + port))
}
