package main

import (
	"fmt"
	"golang.org/x/exp/rand"
	"time"
)

const n = 100
const size = 1e3

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))

	arr := make([]int, size, size)
	for i := range arr {
		arr[i] = i + 1
	}
	for i := 0; i < size; i++ {
		sw := i + rand.Intn(size-i)
		arr[i], arr[sw] = arr[sw], arr[i]
	}

	fmt.Println(arr[:n])
}
