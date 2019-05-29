package main

import (
	"go-es-demo/src/server/common"
	"go-es-demo/src/server/initialize"
)

func main() {

	initialize.SetupConfig()

	common.SetEsClient()

	initialize.SetErrorDeal()

	initialize.SetContext()

	initialize.SetupServer()
}
