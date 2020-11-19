package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"time"
)

func main() {
	var s = semaphore.NewWeighted(10)

	for {
		err := s.Acquire(context.TODO(), 1)

		if err == nil {
			fmt.Println("get")

			go func() {
				time.Sleep(time.Second)
				fmt.Println("release")
				s.Release(1)
			}()
		}
	}
}
