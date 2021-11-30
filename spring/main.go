package main

import (
	"context"
	"github.com/go-spring/spring-boot"
	_ "github.com/go-spring/starter-echo"
)

func init() {
	SpringBoot.RegisterBean(new(Service)).Init(func(s *Service) {
		SpringBoot.GetBinding("/", s.Echo)
	})
}

type Service struct {
	GoPath string `value:"${GOPATH}"`
}

type EchoRequest struct{}

func (s *Service) Echo(ctx context.Context, req *EchoRequest) interface{} {
	return s.GoPath
}

func main() {
	SpringBoot.RunApplication()
}
