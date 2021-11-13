package common

import (
	"fmt"
	"net"
	"strings"
)

const (
	CanNotGetIp = "get.ip.failed"
)

// CollectEntry 要收集的日志的配置
type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// GetLocalIp 获取本机Ip,利用
func GetLocalIp() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}
