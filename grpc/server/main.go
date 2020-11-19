package main

import (
	"context"
	helloworld "first-go/grpc/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

// server is used to implement helloworld.GreeterServer.
type service struct {
	server *grpc.Server
}

// SayHello implements helloworld.GreeterServer
func (s *service) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", request.GetName())

	return &helloworld.HelloReply{Message: "Hello " + request.GetName()}, nil
}

func main() {
	var service = &service{grpc.NewServer()}

	reg := newConsulRegister()
	reg.Port = 8079
	reg.Name = "testservice"
	reg.Tag = []string{"testservice"}
	if err := reg.Register(); err != nil {
		panic(err)
	}
	grpc_health_v1.RegisterHealthServer(service.server, &health{Status: grpc_health_v1.HealthCheckResponse_SERVING})

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", reg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	helloworld.RegisterGreeterServer(service.server, service)
	if err := service.server.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
