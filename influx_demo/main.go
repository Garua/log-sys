package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"time"
)



func main() {
	conn := connInflux()
	if conn != nil {
		fmt.Println("连接成功！")
	}
	//ticker := time.Tick(time.Second)
	//select {
	//case <-ticker:
	//	cpuInfo()
	//
	//}
	writePoints()
	queryDB(org)
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

var client influxdb2.Client

const bucket = "stat"
const org = "own"

func connInflux() influxdb2.Client {
	// You can generate a Token from the "Tokens Tab" in the UI
	const token = "P5AnFpfNP3DVbd-ZH-S40HGg5b91Xsj1FCOCEpX-kcV28t7-rPJGrzX3riEQl8RijfND-_djFyNJtqNo8p4QwQ=="

	client = influxdb2.NewClient("http://localhost:8086", token)
	// always close client at the end
	return client
}
func queryDB(org string) {
	query := fmt.Sprintf("from(bucket:\"%v\")|> range(start: -1h) |> filter(fn: (r) => r._measurement == \"stat\")", bucket)
	fmt.Println(query)
	// Get query client
	queryAPI := client.QueryAPI(org)
	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			fmt.Printf("value: %v field:%v\n", result.Record().Value(),result.Record().Field())
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
}

func writePoints() {
	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45},
		time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// create point using fluent style
	p = influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()


}
