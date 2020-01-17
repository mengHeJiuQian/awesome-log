package main

import (
	"github.com/astaxie/beego/logs"
	"net"
)

/**
 * Create Time : 2020/1/15 下午3:08
 * Update Time :
 * Author : sheldon
 * Decription :
 */

var (
	ipArrays []string
)

func getLocalIP() (ips []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logs.Error("get ip arr failed,err:%v", err)
		return
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return
}
