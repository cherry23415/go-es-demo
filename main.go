package main

import (
	"go-es-demo/src/server/common"
	"go-es-demo/src/server/initialize"
)

func main() {

	initialize.SetupConfig()

	common.SetEs()

	initialize.SetErrorDeal()

	initialize.SetContext()

	initialize.SetupServer()
}
