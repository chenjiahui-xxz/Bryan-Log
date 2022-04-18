package common

import (
	"fmt"
	"net"
	"strings"
)

// CollectEntity etcd 收集的日志的配置
type CollectEntity struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// GetIp 获取本机ip
func GetIp() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
