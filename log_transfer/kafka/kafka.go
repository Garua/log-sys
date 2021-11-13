package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log_transfer/es"
)

// 初始化kafka连接
// 从kafka里面取出日志




func Init(address []string, topic string) error {
	consumer, err := sarama.NewConsumer(address, nil)
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		panic(err)
		return err
	}
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Printf("fail to get list of partition err:%v\n", err)
		return err
	}
	for partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			panic(err)
			return err
		}
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				// 异步写入es
				fmt.Printf("partition:%d Offset:%d key :%s value :%s",msg.Partition,msg.Offset,msg.Key,msg.Value)
				var m1 map[string]interface{}
				err = json.Unmarshal(msg.Value,&m1)
				if err != nil {
					fmt.Printf("unmarshal msg failed,err:%v\n",err)
					continue
				}
				es.PutLodDaa(m1)
			}
		}(pc)
	}
	return nil
}
