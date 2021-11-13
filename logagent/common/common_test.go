package common

import (
	"fmt"
	"testing"
)

func TestGetLocalIp(t *testing.T) {
	ip,_ := GetLocalIp()
	fmt.Println(ip)
}
