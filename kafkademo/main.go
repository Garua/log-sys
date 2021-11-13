package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

//kafka client demo
func main() {
	config := sarama.NewConfig()
	//发送完需要leader和follow都确认 ,ack
	config.Producer.RequiredAcks = sarama.WaitForAll
	//选择分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 确认，成功交付的信息将在success channel返回
	config.Producer.Return.Successes = true

	//连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer closed,err:",err)
		return
	}
	defer client.Close()

	//构造消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "shopping"
	msg.Value = sarama.StringEncoder("this is a test log")

	//发送消息
	pid,offset,err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed,err",err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n",pid,offset)

}
