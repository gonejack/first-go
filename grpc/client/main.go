package main

import (
	"context"
	helloworld "first-go/grpc/pb"
	"google.golang.org/grpc"
	"log"
	"time"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
)

func main() {
	conn, err := grpc.Dial(
		"consul://127.0.0.1:8500/testservice?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err == nil {
		defer conn.Close()
	} else {
		log.Fatalf("did not connect: %v", err)
	}

	client := helloworld.NewGreeterClient(conn)
	{
		var n = 100
		for n > 0 {
			n -= 1
			timeout, _ := context.WithTimeout(context.Background(), time.Second)
			request := &helloworld.HelloRequest{Name: "request"}
			reply, err := client.SayHello(timeout, request)
			if err == nil {
				log.Printf("resonse: %s", reply.GetMessage())
			} else {
				log.Fatalf("could not greet: %v", err)
			}
		}
	}
}
