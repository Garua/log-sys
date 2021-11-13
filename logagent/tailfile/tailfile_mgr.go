package tailfile

import (
	"github.com/sirupsen/logrus"
	"logagent/common"
)

//tailTask 的管理者

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask       //所有的task任务
	collectEntryList []common.CollectEntry      //所有配置项
	confChan         chan []common.CollectEntry //等待新配置的通道
}

var (
	ttMgr *tailTaskMgr
)

// Init 在main中调用
func Init(allConf []common.CollectEntry) (err error) {
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		collectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry),
	}
	for _, conf := range allConf {
		task := newTailTask(conf.Path, conf.Topic) //创建日志收集任务
		err = task.Init()
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed,err:%s\n", conf.Path, err)
			continue
		}
		//去收集日志
		logrus.Infof("create a tail task for path:%s success\n ", conf.Path)
		//把创建的task存到Map,方便后续管理
		ttMgr.tailTaskMap[task.path] = task
		go task.run()

	}
	go ttMgr.watch() //等
	return
}
func (ttMgr *tailTaskMgr) watch() {
	for { //后台一直等，配置改变，就调整任务
		//初始化新配置的管道
		//一个阻塞的通道，没有新配置就等
		//监听etcd中配置项的改变
		newConf := <-ttMgr.confChan
		logrus.Infof("get new conf from etcd,conf:%v，start manager tail task\n", newConf)
		for _, conf := range newConf {
			//原来存在的任务不用动
			if ttMgr.isExist(conf) {
				continue
			}
			//新来的任务创建一个新的tailTask任务
			tt := newTailTask(conf.Path, conf.Topic)
			err := tt.Init()
			if err != nil {
				logrus.Errorf("create a tail task for path:%s ,top:%s,failed\n", conf.Path, conf.Topic)
				continue
			}
			ttMgr.tailTaskMap[tt.path] = tt
			go tt.run() //启动新任务

		}
		//原来有，现在没有，需要停掉对应的task
		//找出tailTaskMap中存在，但是newConf中不存在的tailTask,把他们都停掉
		for k, task := range ttMgr.tailTaskMap {
			var found bool
			for _, conf := range newConf {
				if k == conf.Path {
					found = true
					break
				}
			}
			if !found {
				logrus.Infof("the task collect path:%s need to stop.",task.path)
				//把需要停掉的服务从map中删除
				delete(ttMgr.tailTaskMap,k)
				task.cancel()
			}

		}
	}

}

func (t *tailTaskMgr) isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok

}

func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
