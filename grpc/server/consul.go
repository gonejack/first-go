package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"time"
)

func newConsulRegister() *ConsulRegister {
	return &ConsulRegister{
		Address:                        "localhost:8500", //consul address
		Name:                           "undefined",
		Tag:                            []string{},
		Port:                           3000,
		DeregisterCriticalServiceAfter: time.Duration(10) * time.Second,
		Interval:                       time.Duration(1) * time.Second,
	}
}

type ConsulRegister struct {
	Address                        string
	Name                           string
	Tag                            []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

func (r *ConsulRegister) Register() (err error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return
	}

	ip := getLocalIP()
	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", r.Name, ip, r.Port), // 服务节点的名称
		Name:    r.Name,                                      // 服务名称
		Tags:    r.Tag,                                       // tag，可以为空
		Port:    r.Port,                                      // 服务端口
		Address: ip,                                          // 服务 ip
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       r.Interval.String(),                         // 健康检查间隔
			GRPC:                           fmt.Sprintf("%v:%v/%v", ip, r.Port, r.Name), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),   // 注销时间，相当于过期时间
		},
	}

	return client.Agent().ServiceRegister(reg)
}

func getLocalIP() (s string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return
}
