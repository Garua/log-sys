package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

//初始化全局的KafkaClient
func Init(address []string, chanSize int64) (err error) {
	//kafka生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll //all
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true //确认
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka:producer closed,err:%v\n", err)
		return
	}
	//初始化MsgChan
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	//
	go sendMsg()
	return
}

//从通道中读取消息，发送到kafka
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Info("send msg failed,err:", err)
				continue
			}
			logrus.Infof("send msg to kafka success,pid:%d offset:%v", pid, offset)
		}
	}

}

//定义一个函数，向外暴露msgChan,单向通道
func ToMsgChan(msg *sarama.ProducerMessage){
	 msgChan <- msg
}
