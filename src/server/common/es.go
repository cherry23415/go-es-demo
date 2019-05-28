package common

import (
	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/estransport"
	"log"
)

var Es *elasticsearch.Client

/**
es用的v6版本
*/
func SetEs() {
	Es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("create es client failed, err: %v", err)
	}
	log.Println(elasticsearch.Version)
	res, err := Es.Info()
	if err != nil {
		log.Fatalf("Error getting the response: %s\n", err)
	}
	defer res.Body.Close()
	log.Println(res)
	log.Println(Es.Transport.(*estransport.Client).URLs())
}
