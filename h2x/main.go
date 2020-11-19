package main

import (
	"fmt"
	"time"
)

func main() {
	{
		start := time.Now()
		var i = 0
		for i < 1e9 {
			x2dm("1A5C")
			i++
		}
		fmt.Printf("time %s: %f sec.\n", "timeA", time.Now().Sub(start).Seconds())
	}

	{
		start := time.Now()
		var i = 0
		for i < 1e9 {
			x2dmod("1A5C")
			i++
		}
		fmt.Printf("time %s: %f sec.\n", "timeB", time.Now().Sub(start).Seconds())
	}
}

func x2dm(s string) (d int32) {
	for _, x := range s {
		if x > '9' {
			d = d<<4 | (x - ('A' - 10))
		} else {
			d = d<<4 | (x - '0')
		}
	}
	return d
}

func x2dmod(s string) (d int32) {
	for _, x := range s {
		if x > '9' {
			d = d<<4 | x%('A'-10)%'0'
		} else {
			d = d<<4 | x%'0'
		}
	}
	return d
}
