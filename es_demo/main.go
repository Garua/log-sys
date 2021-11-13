package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func main() {
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	fmt.Println("connect to es success")
	exists, err := client.IndexExists("user").Do(context.Background())
	if !exists {
		//elasticsearch7默认不在支持指定索引类型，默认索引类型是_doc，
		//如果想改变，则配置include_type_name: true 即可(这个没有测试
		//官方文档说的，无论是否可行，建议不要这么做，因为elasticsearch8后就不在提供该字段)
		mapping := `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
			"properties":{
				"name":{
					"type":"keyword"
				},
				"age":{
					"type":"text",
					"store": true
				}
			}
	}
}
`
		_, err := client.CreateIndex("user").Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}

	}
	p1 := Person{Name: "loocc", Age: 30, Married: false}
	put1, err := client.Index().Index("user").
		BodyJson(p1).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Index user %s to index %s,type %s\n", put1.Id, put1.Index, put1.Type)

}
