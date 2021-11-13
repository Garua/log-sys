package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"logagent/common"
	"logagent/tailfile"
	"time"
)

var (
	client *clientv3.Client
)

// Init
func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed,err:%v\n", err)
		return
	}
	return
}

// GetConf 拉取日志收集配置项的函数
func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by %s failed,err:%v\n", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len:0 conf from etcd by %s,err:%v\n", key)
		return
	}
	ret := resp.Kvs[0]
	fmt.Printf("value:%s\n", ret.Value)
	// ret.Value 存放的是json格式字符串
	// Unmarshal和marshal时，对应结构体字段应该是可导出的，不然值有问题
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("json unmarshal failed,err:%v\n", err)
		return nil, err
	}
	return
}

// WatchConf 监控etcd key对应的值变化的函数
func WatchConf(key string) {//后台一直监听
	for{
		//这里的context不要用带超时的，因为我们不知道什么时候，这个key对应的值会变
		watchCh :=client.Watch(context.Background(),key)
		for wresp := range watchCh{
			logrus.Info("get new conf from etcd!")
			for _,evt := range wresp.Events{
				fmt.Printf("type:%s key:%s valeu:%s\n",evt.Type,evt.Kv.Key,evt.Kv.Value)
				var newConf []common.CollectEntry

				//如果是删除事件,etcd中没有对应的配置key了
				if evt.Type == clientv3.EventTypeDelete{
					logrus.Warningf("FBI warning :etcd delete the key %s!!",evt.Kv.Key)
					//为nil不影响遍历(不会报错)
					tailfile.SendNewConf(newConf) //阻塞
					continue
				}

				err := json.Unmarshal(evt.Kv.Value,&newConf)
				if err != nil {
					logrus.Errorf("json Unmarshal failed,err:%v\n",err)
					continue
				}
				//告诉tailfile模块，应该启用新的配置
				tailfile.SendNewConf(newConf) //阻塞
			}
		}
	}
}

