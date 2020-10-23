package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.SetFlags(0)
	log.Printf("Serving http://%s:3000", ipAddr())
	_ = http.ListenAndServe(":3000", http.FileServer(http.Dir(".")))
}

func ipAddr() string {
	cn, err := net.DialTimeout("udp", "8.8.8.8:80", time.Second)
	if err == nil {
		ip, _, err := net.SplitHostPort(cn.LocalAddr().String())
		if err == nil {
			return ip
		}
	}
	fmt.Printf("解析网络IP地址失败: %s", err)
	return ""
}
