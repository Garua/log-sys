package main

// 从kafka消费日志，发往es
import (
	"fmt"
	"gopkg.in/ini.v1"
	"log_transfer/es"
	"log_transfer/kafka"
	"log_transfer/model"
)

func main() {
	var cfg = new(model.Config)
	err := ini.MapTo(cfg, "./config/logtransfer.ini")
	if err != nil {
		fmt.Printf("load config failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("load config success!")
	// 3.连接es
	err = es.Init(cfg.ESConf.Address, cfg.ESConf.GoRoutineNum, cfg.ESConf.MaxChanSize)
	if err != nil {
		fmt.Printf("connect to es failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("connect to es success!")

	// 2.连接kafka
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Topic)
	if err != nil {
		fmt.Printf("connect to kafka failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("connect kafka success")


	//等
	select {}
}
