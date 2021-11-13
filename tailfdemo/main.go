package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)
func main() {
	filename := `./xx.log`
	config := tail.Config{
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{Offset: 0,Whence: 2},
		MustExist: false,
		Poll: true,
	}
	//打开文件，开始读取数据
	file, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Println("tailfile %s failed,err:%v\n",filename,err)
	}
	var (
		msg *tail.Line
		ok bool
	)
	for {
		msg,ok = <-file.Lines
		if !ok {
			fmt.Printf("tailfile file close reopen,filename:%s\n",file.Filename)
			time.Sleep(time.Second)//读取出错，等一秒
			continue
		}
		fmt.Println("msg:",msg.Text)
	}
}
