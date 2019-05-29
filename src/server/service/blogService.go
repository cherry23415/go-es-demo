package service

import (
	"context"
	"encoding/json"
	"github.com/bitly/go-simplejson"
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
			req := esapi.IndexRequest{
				Index:        index,
				DocumentType: _type,
				DocumentID:   data.Id,
				Body:         strings.NewReader(string(d)),
				Refresh:      "true",
			}

			res, err := req.Do(context.Background(), common.ES)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				var r simplejson.Json
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					log.Printf("[%s] %s; version=%d", res.Status(), r.Get("result"), r.Get("_version").MustInt())
				}
			}
		}(i, data)
	}
	wg.Wait()
	log.Println(strings.Repeat("-", 40))
}

func (*BlogService) Search(index string, _type string, searchStr string) {
	var (
		r   simplejson.Json
		buf strings.Builder
	)
	buf.Reset()
	buf.WriteString(`{"query": {"bool": {"should": [{"match": {"title": "`)
	buf.WriteString(searchStr)
	buf.WriteString(`"}},{"match": {"content": "`)
	buf.WriteString(searchStr)
	buf.WriteString(`"}}]}}}`)
	// Perform the search request.
	res, err := common.ES.Search(
		common.ES.Search.WithContext(context.Background()),
		common.ES.Search.WithIndex(index),
		common.ES.Search.WithSearchType(_type),
		common.ES.Search.WithBody(strings.NewReader(buf.String())),
		common.ES.Search.WithTrackTotalHits(true),
		common.ES.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e simplejson.Json
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			log.Fatalf("[%s] %s: %s", res.Status(), e.Get("error").Get("type"), e.Get("error").Get("reason"))
		}
	} else {
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		log.Printf("[%s] %d hits; took: %dms", res.Status(), r.Get("hits").Get("total").MustInt(), r.Get("took").MustInt())
		for _, hit := range r.Get("hits").Get("hits").MustArray() {
			h, _ := simplejson.NewJson(hit.([]byte))
			log.Printf(" * ID=%s, %s", h.Get("_id"), h.Get("_source"))
		}
	}
	log.Println(strings.Repeat("=", 40))
}
