package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type ESClient struct {
	client      *elastic.Client
	index       string //index-- 数据库
	logDataChan chan interface{}
}

var esClient *ESClient

func Init(addr,index string, goroutineNum, maxSize int) (err error) {
	client, err := elastic.NewClient(elastic.SetURL("http://"+addr))
	if err != nil {
		// Handle error
		panic(err)
	}

	esClient = &ESClient{
		client : client,
		logDataChan: make(chan interface{}, maxSize),
		index: index,
	}

	//fmt.Println("connect to es success")

	for i := 0; i < goroutineNum; i++ {
		go sendToEs()
	}
	return
}

func sendToEs() {
	for m1 := range esClient.logDataChan {
		put1, err := esClient.client.Index().
			Index(esClient.index).
			BodyJson(m1).
			Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
}

func PutLogData(msg interface{}) {
	esClient.logDataChan <- msg
}
