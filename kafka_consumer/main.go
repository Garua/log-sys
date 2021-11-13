package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer,err:%s\n",err)
		return
	}
	defer consumer.Close()
	//根据topic取所有分区
	partitions, err := consumer.Partitions("web_log")
	if err != nil {
		fmt.Printf("fail to get list of partition,err:%s\n",err)
		return
	}
	fmt.Println(partitions)
	var wg sync.WaitGroup
	//遍历所有分区
	for partition := range partitions{
		//针对每个分区创建一个对应的消费者
		pc,err := consumer.ConsumePartition("web_log",int32(partition),sarama.OffsetNewest)
		if err !=nil {
			fmt.Printf("failed to start consumer for partition %d,err:%s\n",partition,err)
			return
		}
		defer pc.AsyncClose()
		//异步从每个分区获取消息
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages(){
				fmt.Printf("Partition:%d Offset:%d key:%s Value:%s\n",msg.Partition,msg.Offset,
				msg.Key,msg.Value,)
			}
		}(pc)

	}
	wg.Wait()
}
