package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//连接etcd

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Println("connect to etcd failed,err:%v", err)
		return
	}
	defer client.Close()
	//put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	str := `[{"path":"d:\\logs\\s4.log","topic":"s4_log"},{"path":"e:\\logs\\web.log","topic":"web_log"}]`
	//str := `[{"path":"d:\\logs\\s4.log","topic":"s4_log"},{"path":"e:\\logs\\web.log","topic":"web_log"},{"path":"f:\\logs\\nazha.log","topic":"nazha_log"}]`
	_, err = client.Put(ctx, "collect_log_192.168.18.185_conf", str)
	if err != nil {
		fmt.Printf("put to etcd failed,err:%v", err)
		return
	}
	cancel()
	//get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	gr, err := client.Get(ctx, "collect_log_conf")
	if err != nil {
		fmt.Printf("get from etcd failed,err:%v", err)
		return
	}
	for _, ev := range gr.Kvs {
		fmt.Printf("key:%s,value:%s\n",ev.Key,ev.Value)
	}
	cancel()
}
