package service

import (
	"bytes"
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

	log.Println(strings.Repeat("-", 40))
}

func (*BlogService) Search(index string, _type string, searchStr string) {
	var (
		r      map[string]interface{}
		buf    bytes.Buffer
		matchs [2]map[string]interface{}
	)
	matchs[0] = map[string]interface{}{
		"title": searchStr,
	}
	matchs[1] = map[string]interface{}{
		"content": searchStr,
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": matchs,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := common.ES.Search(
		common.ES.Search.WithContext(context.Background()),
		common.ES.Search.WithIndex(index),
		//common.ES.Search.WithSearchType(_type),
		common.ES.Search.WithBody(&buf),
		common.ES.Search.WithTrackTotalHits(true),
		common.ES.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 40))
}
