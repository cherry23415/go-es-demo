package common

import (
	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func SetEsClient() {
	//地址，多个逗号间隔
	addrs := strings.Split(viper.GetString("es.addr"), ",")
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: addrs,
	})
	if err != nil {
		log.Fatalf("create es client failed, err: %v", err)
	}
	log.Println(es.Info())
	ES = es
}

var ES *elasticsearch.Client
