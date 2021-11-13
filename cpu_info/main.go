package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"time"
)



func getCpuInfo()  {
	cpuInfos,err := cpu.Info()
	if err != nil {
		fmt.Printf("get cpu info failed,err:%v\n",err)
	}
	//cpu使用率
	for _,ci := range cpuInfos{
		fmt.Println(ci)
	}
	for {
		percent,_ := cpu.Percent(time.Second,false)
		fmt.Printf("cpu percent :%v\n",percent)
	}
}
func main()  {
	getCpuInfo()
}