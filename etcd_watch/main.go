package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//watch
func main() {
	cli ,err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: time.Second*1,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed,err:%v",err)
		return
	}
	defer cli.Close()
	//watch,没有改动就会阻塞
	watchCh := cli.Watch(context.Background(), "s4")
	for wresp := range watchCh{
		for _,evt := range wresp.Events{
			fmt.Printf("type:%s key:%s,value:%s\n",evt.Type,evt.Kv.Key,evt.Kv.Value)
		}
	}
	
}
