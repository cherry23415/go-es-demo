package service

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"go-es-demo/src/server/common"
	"log"
	"strings"
)

type BlogService struct {
}

func (*BlogService) Save(index string, _type string, data string, id string) {
	// Set up the request object.
	req := esapi.IndexRequest{
		Index:        index,
		DocumentType: _type,
		DocumentID:   id,
		Body:         strings.NewReader(data),
		Refresh:      "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), common.Es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), id)
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
}
