package tailfile

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent/kafka"
	"strings"
	"time"
)

type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

func newTailTask(path, topic string) *tailTask {
	ctx, cancel := context.WithCancel(context.Background())
	tt := tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	return &tt
}
func (t *tailTask) Init() (err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tObj, err = tail.TailFile(t.path, cfg)
	return
}

func (task *tailTask) run() {
	//读取日志，发往kafka
	logrus.Infof("collect for path:%s is running.", task.path)
	for {
		select {
		case <-task.ctx.Done(): //接收取消信号(cancel(),)
			fmt.Printf("task path %s stop!",task.path)
			return
		case line, ok := <-task.tObj.Lines:
			//line, ok := <-task.tObj.Lines
			if !ok {
				logrus.Error("tail file close reopen,filename:%s\n", task.path)
				time.Sleep(time.Second)
				continue
			}
			//去掉换行符(win),不然golang客户端接收到的值有问题
			if len(strings.Trim(line.Text, "\r")) == 0 {
				continue
			}

			//利用通道将同步代码改为异步的
			//把读出的数据包装成kafka里面的Message类型，放到通道中
			msg := &sarama.ProducerMessage{}
			msg.Topic = task.topic
			msg.Value = sarama.StringEncoder(line.Text)
			kafka.ToMsgChan(msg)
		}

	}
}
