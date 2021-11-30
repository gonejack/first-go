package main

import (
	"context"
	"github.com/shenghui0779/yiigo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

func main() {
	yiigo.Init(
		yiigo.WithDB(yiigo.Default, yiigo.MySQL, "dsn"),
		yiigo.WithDB("other", yiigo.MySQL, "dsn"),
		yiigo.WithMongo("abc", "dsn"),
	)

	// create pool
	pool := yiigo.NewGRPCPool(
		func() (*grpc.ClientConn, error) {
			return grpc.DialContext(context.Background(), "target",
				grpc.WithInsecure(),
				grpc.WithBlock(),
				grpc.WithKeepaliveParams(keepalive.ClientParameters{
					Time:    time.Second * 30,
					Timeout: time.Second * 10,
				}),
			)
		},
		yiigo.WithPoolSize(10),
		yiigo.WithPoolLimit(20),
		yiigo.WithPoolIdleTimeout(600*time.Second),
	)

	// use pool
	conn, _ := pool.Get(context.Background())
	defer pool.Put(conn)

	yiigo.LoadEnvFromFile("abc")
	yiigo.Mongo("abc").Database("").Collection("abc")
}
