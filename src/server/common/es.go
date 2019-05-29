package common

import (
	"github.com/elastic/go-elasticsearch/v6"
	"log"
)

func SetEsClient() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("create es client failed, err: %v", err)
	}
	log.Println(es.Info())
	ES = es
}

var ES *elasticsearch.Client
