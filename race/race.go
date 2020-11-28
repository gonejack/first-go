package main

import (
	"fmt"
	"time"
)

func main() {
	i := 0

	go func() {
		for {
			i++
			time.Sleep(time.Nanosecond)
		}
	}()

	for {
		time.Sleep(time.Millisecond)
		fmt.Println(i)
	}
}
