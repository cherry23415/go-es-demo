package common

import (
	"fmt"
	"gopkg.in/olivere/elastic.v6"
	"log"
)

var EsClient *elastic.Client

func SetEs() {
	//default 127.0.0.1:9200
	EsClient, err := elastic.NewClient()
	if err != nil {
		fmt.Printf("create es client failed, err: %v", err)
	}
	log.Println(EsClient.String())
}
