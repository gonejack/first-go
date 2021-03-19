package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.TODO(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	server, err := newServer(ctx, Cloudflare())
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		log.Println("DNS server listening on :53")
		if err := server.ListenAndServe(); err != nil {
			log.Println("DNS server crashed: ", err)
		}
		log.Println("DNS server stopped")
	}()

	<-ctx.Done()

	timeout, _ := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	err = server.ShutdownContext(timeout)
	if err != nil {
		log.Println("DNS server shutdown error: ", err)
	}
}
