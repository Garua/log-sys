package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

// 将日志写入es

var (
	esClient *EsClient
)

type EsClient struct {
	client *elastic.Client
	logDataChan chan interface{}
	index string
}

func Init(address string,gn,maxSize int) (err error) {
	client, err := elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		panic(err)
		return err
	}
	esClient = new(EsClient)
	esClient.client = client
	esClient.logDataChan = make(chan interface{},maxSize)
	fmt.Println("connect to es success!")
	//从通道中取数据
	for i :=0;i < gn;i++{
		go sendToES(esClient.index)
	}
	return nil
}
func sendToES(index string)  {
	for m1 := range esClient.logDataChan{
		//b,err := json.Marshal(m1)
		//if err != nil {
		//	fmt.Printf("Marshal m1 failed,m1 %s,err:%v\n",m1,err)
		//	continue
		//}
		p1,err := esClient.
			client.
			Index().
			Index(index).
			BodyJson(m1).
			Do(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Indexed user %s to index %s,type %s\n",p1.Id,p1.Index,p1.Type)
	}
}


// 通过函数从包外接收msg
func PutLodDaa(msg interface{})  {
	esClient.logDataChan<-msg
}