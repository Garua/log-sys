package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"sync"
	"time"
)

func main() {
	//getNet()
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0;i < 100;i++{
		go fib(10000)
	}
	wg.Wait()
}
func fib(n int64) int64 {
	if n <= 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}




// 网络信息
func getNet()  {
	counters, _ := net.IOCounters(true)
	for _,count := range counters{
		fmt.Println(count)
	}
}



func getDiskInfo() {
	partitions, err := disk.Partitions(true)
	if err != nil {
		fmt.Errorf("err:%s\n", err)
	}
	for _, partition := range partitions {
		//磁盘使用情况
		pi, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			fmt.Printf("get part stat failed,err:%v\n",err)
			return
		}
		fmt.Println(pi)
	}
	// 磁盘io
	counters, _ := disk.IOCounters()
	for k,v := range counters{
		fmt.Printf("key:%v  value:%v\n",k,v)
	}
}

// 主机信息
func getHostInfo() {
	info, err := host.Info()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info)
}

// getLoad,win没有？？？
func getLoad() {
	info, err := load.Avg()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info)

}

// getMemInfo 获取内存信息
func getMemInfo() {
	memory, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(memory)
}

// cpuInfo cpu使用率
func cpuInfo() {
	infos, err := cpu.Info()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, info := range infos {
		fmt.Println(info)
	}
	for {
		percent, _ := cpu.Percent(time.Second, false)
		fmt.Println(percent)
	}
}
