package main

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"time"
)

var (
	client                 influxdb2.Client
	lastNetIOStatTimeStamp int64
	lastNetInfo            *NetInfo // 上一次的网络IO数据
)

const bucket = "monitor"
const org = "own"

func getCpuInfo() {
	percents, _ := cpu.Percent(time.Second, false)
	fmt.Printf("cpu percent %v\n", percents)
	//发送到influxdb
	for _, v := range percents {
		fmt.Printf("准备写CPU数据,%v\n", v)
		var cpuInfo = new(CpuInfo)
		cpuInfo.CpuPercent = percents[0]
		writeCpuPoints(cpuInfo)
	}
}

func initConnInflux() {
	// You can generate a Token from the "Tokens Tab" in the UI
	const token = "P5AnFpfNP3DVbd-ZH-S40HGg5b91Xsj1FCOCEpX-kcV28t7-rPJGrzX3riEQl8RijfND-_djFyNJtqNo8p4QwQ=="

	client = influxdb2.NewClient("http://localhost:8086", token)
	// always close client at the end
}

func writeCpuPoints(data *CpuInfo) {
	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)
	// 根据传入的数据类型插入数据 1
	//p := influxdb2.NewPoint("monitor", //表名
	//	map[string]string{"cpu": "percent"},            //tag
	//	map[string]interface{}{"cpu_percent": percent}, //field
	//	time.Now())
	// write point asynchronously
	//writeAPI.WritePoint(p)
	// create point using fluent style

	// 2 两种写法
	p := influxdb2.NewPointWithMeasurement("cpu").
		AddTag("cpu", "cpu0").
		AddField("cpu_percent", data.CpuPercent).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

}

func getMemInfo() {
	var memInfo = new(MemInfo)
	memory, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("get mem failed, err:%v\n", err)
		return
	}
	memInfo.Total = memory.Total
	memInfo.UsedPercent = memory.UsedPercent
	memInfo.Used = memory.Used
	memInfo.Buffers = memory.Buffers
	memInfo.Cached = memory.Cached
	memInfo.Available = memory.Available
	fmt.Println("准备写内存数据")
	writeMemPoints(memInfo)
}

func writeMemPoints(data *MemInfo) {
	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)
	// 根据传入的数据类型插入数据 1
	//p := influxdb2.NewPoint("monitor", //表名
	//	map[string]string{"cpu": "percent"},            //tag
	//	map[string]interface{}{"cpu_percent": percent}, //field
	//	time.Now())
	// write point asynchronously
	//writeAPI.WritePoint(p)
	// create point using fluent style

	// 2 两种写法
	p := influxdb2.NewPointWithMeasurement("mem").
		AddTag("cpu", "mem").
		AddField("total", data.Total).
		AddField("used", data.Used).
		AddField("cached", data.Cached).
		AddField("buffers", data.Buffers).
		AddField("usedPercent", data.UsedPercent).
		AddField("available", data.Available).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

}

func getDiskInfo() {
	var diskInfo = &DiskInfo{
		PartitionUsageStat: make(map[string]*disk.UsageStat, 16),
	}
	partitions, _ := disk.Partitions(true)
	for _, part := range partitions {
		// 拿到每个磁盘分区的信息
		usageInfo, _ := disk.Usage(part.Mountpoint) //传入挂载点
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageInfo
	}
	fmt.Println(diskInfo)
	writeDiskPoints(diskInfo)
}

func writeDiskPoints(data *DiskInfo) {
	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)

	for k, v := range data.PartitionUsageStat {
		tags := map[string]string{
			"path": k,
		}
		fields := map[string]interface{}{
			"total":               v.Total,
			"free":                v.Free,
			"used":                v.Used,
			"used_percent":        v.UsedPercent,
			"inodes_total":        v.InodesTotal,
			"inodes_used":         v.InodesUsed,
			"inodes_free":         v.InodesFree,
			"inodes_used_percent": v.InodesUsedPercent,
		}
		p := influxdb2.NewPoint("disk", tags, fields, time.Now())
		writeAPI.WritePoint(p)
		writeAPI.Flush()

	}

	// write point asynchronously
	// Flush writes

}

func getNetInfo() {
	var netInfo = &NetInfo{
		NetIOCounterStat: make(map[string]*IOStat, 8),
	}
	netIOs, _ := net.IOCounters(true)
	curTimeStamp := time.Now().Unix()
	for _, netIO := range netIOs {
		var ioStat = new(IOStat)
		// 当前的数据
		ioStat.BytesSent = netIO.BytesSent
		ioStat.BytesRecv = netIO.BytesRecv
		ioStat.PacketsSent = netIO.PacketsSent
		ioStat.PacketsRecv = netIO.PacketsRecv

		// 将具体的网卡数据添加到map
		netInfo.NetIOCounterStat[netIO.Name] = ioStat

		// 计算相关速率
		//第一次，不用计算
		if lastNetIOStatTimeStamp == 0 || lastNetInfo == nil {
			continue
		}
		// 计算时间间隔
		interval := curTimeStamp - lastNetIOStatTimeStamp

		ioStat.BytesSentRate = float64(ioStat.BytesSent - lastNetInfo.NetIOCounterStat[netIO.Name].BytesSent) / float64(interval)
		ioStat.BytesRecvRate = float64(ioStat.BytesRecv - lastNetInfo.NetIOCounterStat[netIO.Name].BytesRecv) / float64(interval)
		ioStat.PacketsSentRate = float64(ioStat.PacketsSent - lastNetInfo.NetIOCounterStat[netIO.Name].PacketsSent) / float64(interval)
		ioStat.PacketsRecvRate = float64(ioStat.PacketsRecv - lastNetInfo.NetIOCounterStat[netIO.Name].PacketsRecv) / float64(interval)


	}
	lastNetIOStatTimeStamp = curTimeStamp
	lastNetInfo = netInfo
	//发送到influxdb
	writeNetPoints(netInfo)
}

func writeNetPoints(data *NetInfo) {
	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)

	for k, v := range data.NetIOCounterStat {
		tags := map[string]string{
			"name": k,//网卡名索引
		}
		fields := map[string]interface{}{
			"bytes_sent_rate":v.BytesSentRate,
			"bytes_recv_rate":v.BytesRecvRate,
			"packets_sent_rate":v.PacketsSentRate,
			"packets_recv_rate":v.PacketsRecvRate,
		}
		fmt.Println(v,k)
		p := influxdb2.NewPoint("net", tags, fields, time.Now())
		writeAPI.WritePoint(p)
		writeAPI.Flush()

	}

	// write point asynchronously
	// Flush writes

}

func main() {
	initConnInflux()
	run(time.Second)
}

func run(interval time.Duration) {
	ticker := time.Tick(interval)
	for _ = range ticker {
		go getCpuInfo()
		go getMemInfo()
		go getDiskInfo()
		go getNetInfo()
	}
}
