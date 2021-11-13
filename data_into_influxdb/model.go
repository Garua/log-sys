package main

import (
	"github.com/shirou/gopsutil/v3/disk"
)

const (
	CpuInfoType  = "cpu"
	MemInfoType  = "mem"
	DiskInfoType = "disk"
	NetInfoType  = "net"
)

type SysInfo struct {
	InfoType string
	IP       string
	Data     interface{}
}
type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}

type MemInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
}

type UsageStat struct {
	Path              string  `json:"path"`
	FsType            string  `json:"fs_type"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"used_percent"`
	InodesTotal       uint64  `json:"inodes_total"`
	InodesUsed        uint64  `json:"inodes_used"`
	InodesFree        uint64  `json:"inodes_free"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
}


type DiskInfo struct {
	PartitionUsageStat map[string]*disk.UsageStat
}

type NetInfo struct {
	NetIOCounterStat map[string]*IOStat
}
// 网上速率
type IOStat struct {
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
	BytesSentRate   float64 `json:"bytes_sent_rate"`
	BytesRecvRate   float64 `json:"bytes_recv_rate"`
	PacketsSentRate float64 `json:"packets_sent_rate"`
	PacketsRecvRate float64 `json:"packets_recv_rate"`
}