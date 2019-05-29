package service

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"go-es-demo/src/server/common"
	"go-es-demo/src/server/entity"
	"log"
	"strings"
	"sync"
)

type BlogService struct {
}

func (*BlogService) Save(index string, _type string, datas []entity.Blog) {
	var wg sync.WaitGroup
	for i, data := range datas {
		wg.Add(1)

		go func(i int, data entity.Blog) {
			defer wg.Done()
			d, _ := json.Marshal(data)
			// Set up the request object directly.
			req := esapi.IndexRequest{
				Index:        index,
				DocumentType: _type,
				DocumentID:   data.Id,
				Body:         strings.NewReader(string(d)),
				Refresh:      "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), common.ES)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(i, data)
	}
	wg.Wait()

	log.Println(strings.Repeat("-", 37))
}
