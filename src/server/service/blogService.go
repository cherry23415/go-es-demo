package service

import (
	"context"
	"fmt"
	"go-es-demo/src/server/common"
	"gopkg.in/olivere/elastic.v6"
	"log"
	"strconv"
)

type BlogService struct {
}

func (*BlogService) CreateIndex(index string) bool {
	//查询索引是否存在
	exists, err := common.EsClient.IndexExists(index).Do(context.Background())
	if err != nil {
		// Handle error
		log.Printf("check index failed, err: %v\n", err)
	}
	//不存在就创建索引
	if !exists {
		// Index does not exist yet.
		result, err := common.EsClient.CreateIndex(index).Do(context.Background())
		if err != nil {
			log.Printf("create index failed, err: %v\n", err)
		}
		return result.Acknowledged
	}
	return true
}

//批量插入
func (*BlogService) Batch(index string, _type string, datas ...interface{}) {

	bulkRequest := common.EsClient.Bulk()
	for i, data := range datas {
		doc := elastic.NewBulkIndexRequest().Index(index).Type(_type).Id(strconv.Itoa(i)).Doc(data)
		bulkRequest = bulkRequest.Add(doc)
	}

	response, err := bulkRequest.Do(context.TODO())
	if err != nil {
		panic(err)
	}
	failed := response.Failed()
	iter := len(failed)
	fmt.Printf("error: %v, %v\n", response.Errors, iter)
}
