package service

import (
	"context"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	"github.com/go-errors/errors"
	"github.com/ricardolonga/jsongo"
	"go-es-demo/src/server/common"
	"go-es-demo/src/server/entity"
	"log"
	"strings"
	"sync"
)

type BlogService struct {
}

func (*BlogService) Save(index string, _type string, datas []entity.Blog) error {
	var wg sync.WaitGroup
	for i, data := range datas {
		wg.Add(1)
		go func(i int, data entity.Blog) {
			defer wg.Done()
			req := esapi.IndexRequest{
				Index:        index,
				DocumentType: _type,
				DocumentID:   data.Id,
				Body:         esutil.NewJSONReader(data),
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
	return nil
}

func (*BlogService) Search(index string, _type string, searchStr string, page int, size int) (*simplejson.Json, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	from := (page - 1) * size
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
	req := esapi.SearchRequest{
		Index:   []string{index},
		Body:    strings.NewReader(buf.String()),
		Size:    esapi.IntPtr(size),
		From:    esapi.IntPtr(from),
		Pretty:  true,
		Timeout: 100,
	}
	res, err := req.Do(context.Background(), common.ES)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e simplejson.Json
		json.NewDecoder(res.Body).Decode(&e)
		return nil, errors.New("status:" + res.Status() + ",error type:" +
			e.Get("error").Get("type").MustString() + ",error reason:" + e.Get("error").Get("reason").MustString())
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	returnJson := simplejson.New()
	hits := r.Get("hits")
	returnJson.Set("total", hits.Get("total").MustInt())

	datas := jsongo.Array()
	for _, hit := range r.Get("hits").Get("hits").MustArray() {
		h, _ := json.Marshal(hit)
		j, _ := simplejson.NewJson(h)
		datas.Put(j.Get("_source"))
	}
	returnJson.Set("datas", datas)
	returnJson.Set("page", page)
	returnJson.Set("size", size)
	return returnJson, nil
}
