package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"logagent/common"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"
)

//类似的开源项目-->filebeat
//收集指定目录下的日志文件，发送到kafka

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"` //section
}
type KafkaConfig struct {
	Address  string `ini:"address"` //和配置文件对应
	ChanSize int64  `ini:"chan_size"`
}
type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}
type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

func main() {
	//获取本机Ip,为后续去etcd中取配置文件打下坚实基础
	ip, err := common.GetLocalIp()
	if err != nil {
		logrus.Errorf("get ip failed,err:%v\n",err)
		return
	}
	var configObj = new(Config)

	err = ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Errorf("load config failed,err:%v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)
	//连接kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed,err:%v", err)
		return
	}
	logrus.Info("init kafka success!")
	//初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed,err:%v\n", err)
		return
	}
	logrus.Info("etcd init success!")
	//从etcd中拉取要收集的日志配置项
	collectKey := fmt.Sprintf(configObj.EtcdConfig.CollectKey,ip)
	allConf, err := etcd.GetConf(collectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed,err:%v\n", err)
		return
	}
	fmt.Println("11  ",allConf)
	//监控etcd中对应的配置key对应的值是否发生变化

	go etcd.WatchConf(collectKey)
	//初始化tail
	err = tailfile.Init(allConf) //把从etcd中获取的配置项传到Init中
	if err != nil {
		logrus.Errorf("init tailfile failed,err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")
	for  {
		select {

		}
	}
}

//真正的业务逻辑,单个
//func run() (err error) {
//	for {
//		line, ok := <-tailfile.TailObj.Lines
//		if !ok {
//			logrus.Error("tail file close reopen,filename:%s\n", tailfile.TailObj.Filename)
//			time.Sleep(time.Second)
//			continue
//		}
//		//去掉换行符(win),不然golang客户端接收到的值有问题
//		if len(strings.Trim(line.Text, "\r")) == 0 {
//			continue
//		}
//
//		//利用通道将同步代码改为异步的
//		//把读出的数据包装成kafka里面的Message类型，放到通道中
//		msg := &sarama.ProducerMessage{}
//		msg.Topic = "web_log"
//		msg.Value = sarama.StringEncoder(line.Text)
//		kafka.ToMsgChan(msg)
//	}
//}
